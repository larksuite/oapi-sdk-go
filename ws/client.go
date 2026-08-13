package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcache "github.com/larksuite/oapi-sdk-go/v3/cache"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

type Client struct {
	appID                   string
	appSecret               string
	clientAssertionProvider larkcore.ClientAssertionProvider
	logLevel                larkcore.LogLevel
	logger                  larkcore.Logger
	eventHandler            *dispatcher.EventDispatcher
	cardHandler             *larkcard.CardActionHandler
	domain                  string
	headers                 http.Header
	source                  string
	conn                    *ws.Conn
	connUrl                 *url.URL
	serviceID               string
	connID                  string
	autoReconnect           bool          // 是否自动重连，默认开启
	reconnectNonce          int           // 首次重连抖动，单位秒
	reconnectCount          int           // 重连次数，负数无限次
	reconnectInterval       time.Duration // 重连间隔
	pingInterval            time.Duration // Ping间隔
	cache                   *larkcache.Cache
	lifecycleMu             sync.Mutex
	run                     *clientRun
	mu                      sync.Mutex
	writeMu                 sync.Mutex
	connRun                 *clientRun
	onReady                 func()
	onError                 func(err error)
	onReconnecting          func()
	onReconnected           func()
	onDisconnected          func()
}

type clientRun struct {
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
	once   sync.Once
	err    error
	wg     sync.WaitGroup

	eventsMu                 sync.Mutex
	connectionCallbackActive bool
	pendingDisconnected      func()
}

type clientCallbacks struct {
	onReady        func()
	onError        func(error)
	onReconnecting func()
	onReconnected  func()
}

func newClientRun(ctx context.Context) *clientRun {
	runCtx, cancel := context.WithCancel(ctx)
	return &clientRun{
		ctx:    runCtx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

func (r *clientRun) finish(err error) {
	r.once.Do(func() {
		r.err = err
		r.cancel()
		close(r.done)
	})
}

func (r *clientRun) result() error {
	<-r.done
	return r.err
}

var bootstrapHTTPClient = http.DefaultClient

type bootstrapErrorResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ClientOption func(cli *Client)

func WithEventHandler(handler *dispatcher.EventDispatcher) ClientOption {
	return func(cli *Client) {
		cli.eventHandler = handler
	}
}

//func WithCardHandler(handler *larkcard.CardActionHandler) ClientOption {
//	return func(cli *Client) {
//		cli.cardHandler = handler
//	}
//}

func WithLogLevel(level larkcore.LogLevel) ClientOption {
	return func(cli *Client) {
		cli.logLevel = level
	}
}

func WithLogger(logger larkcore.Logger) ClientOption {
	return func(cli *Client) {
		cli.logger = logger
	}
}

func WithAutoReconnect(b bool) ClientOption {
	return func(cli *Client) {
		cli.autoReconnect = b
	}
}

func WithDomain(domain string) ClientOption {
	return func(cli *Client) {
		cli.domain = domain
	}
}

func WithHeaders(header http.Header) ClientOption {
	return func(cli *Client) {
		cli.headers = header
	}
}

func WithSource(source string) ClientOption {
	return func(cli *Client) {
		cli.source = source
	}
}

func WithClientAssertionProvider(provider larkcore.ClientAssertionProvider) ClientOption {
	return func(cli *Client) {
		cli.clientAssertionProvider = provider
	}
}
func WithOnReady(f func()) ClientOption {
	return func(cli *Client) {
		cli.onReady = f
	}
}

func WithOnError(f func(err error)) ClientOption {
	return func(cli *Client) {
		cli.onError = f
	}
}

func WithOnReconnecting(f func()) ClientOption {
	return func(cli *Client) {
		cli.onReconnecting = f
	}
}

func WithOnReconnected(f func()) ClientOption {
	return func(cli *Client) {
		cli.onReconnected = f
	}
}

func WithOnDisconnected(f func()) ClientOption {
	return func(cli *Client) {
		cli.onDisconnected = f
	}
}

func (c *Client) SetOnReady(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReady = f
}

func (c *Client) SetOnReconnecting(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReconnecting = f
}

func (c *Client) SetOnReconnected(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onReconnected = f
}

func (c *Client) SetOnError(f func(err error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onError = f
}

func (c *Client) SetOnDisconnected(f func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onDisconnected = f
}

func (c *Client) Close() {
	run := c.currentRun()
	if run == nil {
		return
	}
	run.finish(context.Canceled)
	c.disconnectRun(context.Background(), run, nil)
}

func NewClient(appId, appSecret string, opts ...ClientOption) *Client {
	cli := &Client{
		appID:             appId,
		appSecret:         appSecret,
		autoReconnect:     true,
		reconnectNonce:    30,
		reconnectCount:    -1,
		reconnectInterval: 2 * time.Minute,
		pingInterval:      2 * time.Minute,
		cache:             larkcache.New(30 * time.Second),
		domain:            lark.FeishuBaseUrl,
	}

	for _, opt := range opts {
		opt(cli)
	}

	if cli.logger == nil {
		cli.logger = larkcore.NewDefaultLogger(cli.logLevel)
	}

	return cli
}

func (c *Client) beginRun(ctx context.Context) (*clientRun, error) {
	c.lifecycleMu.Lock()
	defer c.lifecycleMu.Unlock()
	if c.run != nil {
		return nil, ErrClientAlreadyStarted
	}
	run := newClientRun(ctx)
	c.run = run
	return run, nil
}

func (c *Client) currentRun() *clientRun {
	c.lifecycleMu.Lock()
	defer c.lifecycleMu.Unlock()
	return c.run
}

func (c *Client) endRun(run *clientRun) {
	c.lifecycleMu.Lock()
	defer c.lifecycleMu.Unlock()
	if c.run == run {
		c.run = nil
	}
}

func (c *Client) Start(ctx context.Context) (err error) {
	run, err := c.beginRun(ctx)
	if err != nil {
		return err
	}
	defer func() {
		run.cancel()
		c.disconnectRun(context.Background(), run, nil)
		run.wg.Wait()
		c.endRun(run)
	}()

	if err = run.ctx.Err(); err != nil {
		run.finish(err)
		return run.result()
	}

	err = c.connect(run)
	if err != nil {
		if ctxErr := run.ctx.Err(); ctxErr != nil {
			run.finish(ctxErr)
			return run.result()
		}
		c.logger.Error(run.ctx, c.fmtLog("connect failed, err: %v", err)...)
		c.callOnError(err)
		if _, ok := err.(*ClientError); ok {
			run.finish(err)
			return run.result()
		}
		if c.autoReconnectEnabled() {
			if err = c.reconnect(run); err != nil {
				if ctxErr := run.ctx.Err(); ctxErr != nil {
					run.finish(ctxErr)
				} else {
					run.finish(err)
				}
				return run.result()
			}
		} else {
			run.finish(err)
			return run.result()
		}
	} else {
		c.callConnectionCallback(run, c.callbacksSnapshot().onReady)
	}

	if ctxErr := run.ctx.Err(); ctxErr != nil {
		run.finish(ctxErr)
		return run.result()
	}

	run.wg.Add(2)
	go func() {
		defer run.wg.Done()
		c.pingLoop(run)
	}()
	go func() {
		defer run.wg.Done()
		c.receiveMessageLoop(run)
	}()

	select {
	case <-run.done:
	case <-run.ctx.Done():
		run.finish(run.ctx.Err())
	}
	return run.result()
}

func (c *Client) connect(run *clientRun) (err error) {
	ctx := run.ctx

	c.mu.Lock()
	if c.conn != nil {
		ownedByRun := c.connRun == run
		c.mu.Unlock()
		if ownedByRun {
			return nil
		}
		return ErrClientAlreadyStarted
	}
	c.mu.Unlock()

	// 获取建连URL
	connUrl, err := c.getConnURL(ctx)
	if err != nil {
		c.logger.Warn(ctx, c.fmtLog("get conn url failed, err: %v", err)...)
		return
	}

	// 验证URL
	u, err := url.Parse(connUrl)
	if err != nil {
		return
	}
	connID := u.Query().Get(DeviceID)
	serviceID := u.Query().Get(ServiceID)

	conn, resp, err := ws.DefaultDialer.DialContext(ctx, connUrl, nil)
	if err != nil {
		if resp == nil {
			return err
		}
		defer resp.Body.Close()
		return parseErr(resp)
	}
	if resp == nil || resp.StatusCode != http.StatusSwitchingProtocols {
		if conn != nil {
			_ = conn.Close()
		}
		if resp == nil {
			return fmt.Errorf("websocket handshake returned no response")
		}
		defer resp.Body.Close()
		// 连接失败
		return parseErr(resp)
	}

	c.mu.Lock()
	if err = ctx.Err(); err != nil {
		c.mu.Unlock()
		_ = conn.Close()
		return err
	}
	if c.conn != nil {
		ownedByRun := c.connRun == run
		c.mu.Unlock()
		_ = conn.Close()
		if ownedByRun {
			return nil
		}
		return ErrClientAlreadyStarted
	}
	c.conn = conn
	c.connRun = run
	c.connUrl = u
	c.connID = connID
	c.serviceID = serviceID
	c.mu.Unlock()

	c.logger.Info(ctx, c.fmtLog("connected to %s", u)...)
	return
}

func (c *Client) reconnect(run *clientRun) (err error) {
	ctx := run.ctx
	if err = ctx.Err(); err != nil {
		return err
	}
	if onReconnecting := c.callbacksSnapshot().onReconnecting; onReconnecting != nil {
		onReconnecting()
	}

	reconnectCount, _, reconnectNonce, _ := c.configSnapshot()

	// 首次重连随机抖动
	if reconnectNonce > 0 {
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(reconnectNonce * 1000)
		if err = sleepWithContext(ctx, time.Duration(num)*time.Millisecond); err != nil {
			return err
		}
	}

	if reconnectCount >= 0 {
		for i := 0; i < reconnectCount; i++ {
			if err = ctx.Err(); err != nil {
				return err
			}
			success, err := c.tryConnect(run, i)
			if success {
				if err = ctx.Err(); err != nil {
					return err
				}
				c.callConnectionCallback(run, c.callbacksSnapshot().onReconnected)
				return nil
			}
			if err != nil {
				return err
			}
			_, reconnectInterval, _, _ := c.configSnapshot()
			if err = sleepWithContext(ctx, reconnectInterval); err != nil {
				return err
			}
		}
		return fmt.Errorf("unable to connect to server after %d retries", reconnectCount)
	} else {
		i := 0
		for {
			if err = ctx.Err(); err != nil {
				return err
			}
			success, err := c.tryConnect(run, i)
			if success {
				if err = ctx.Err(); err != nil {
					return err
				}
				c.callConnectionCallback(run, c.callbacksSnapshot().onReconnected)
				return nil
			}
			if err != nil {
				return err
			}
			_, reconnectInterval, _, _ := c.configSnapshot()
			if err = sleepWithContext(ctx, reconnectInterval); err != nil {
				return err
			}
			i += 1
		}
	}
}

func (c *Client) tryConnect(run *clientRun, cnt int) (bool, error) {
	ctx := run.ctx
	c.logger.Info(ctx, c.fmtLog("trying to reconnect: %d", cnt+1)...)
	err := c.connect(run)
	if err == nil {
		return true, nil
	} else if ctx.Err() != nil {
		return false, ctx.Err()
	} else if _, ok := err.(*ClientError); ok {
		c.callOnError(err)
		return false, err
	} else {
		c.logger.Error(ctx, c.fmtLog("connect failed, err: %v", err)...)
		c.callOnError(err)
		return false, nil
	}
}

func (c *Client) disconnectRun(ctx context.Context, run *clientRun, expected *ws.Conn) {
	run.eventsMu.Lock()

	c.mu.Lock()
	if c.conn == nil || c.connRun != run || (expected != nil && c.conn != expected) {
		c.mu.Unlock()
		run.eventsMu.Unlock()
		return
	}

	conn := c.conn
	connURL := c.connUrl
	connID := c.connID
	onDisconnected := c.onDisconnected
	c.conn = nil
	c.connRun = nil
	c.connUrl = nil
	c.connID = ""
	c.serviceID = ""
	c.mu.Unlock()
	if run.connectionCallbackActive {
		run.pendingDisconnected = onDisconnected
		onDisconnected = nil
	}
	run.eventsMu.Unlock()

	_ = conn.Close()
	log := []interface{}{fmt.Sprintf("disconnected to %s", connURL)}
	if connID != "" {
		log = append(log, fmt.Sprintf("[conn_id=%s]", connID))
	}
	c.logger.Info(ctx, log...)

	if onDisconnected != nil {
		onDisconnected()
	}
}

func (c *Client) getConnURL(ctx context.Context) (url string, err error) {
	requestURL := strings.TrimRight(c.domain, "/") + GenEndpointUri
	body := &BootstrapRequest{AppID: c.appID}
	headers := make(http.Header)

	if c.clientAssertionProvider == nil && c.appSecret == "" {
		return "", NewClientError(larkcore.ErrCodeAppSecretAndClientAssertionEmpty, "appSecret and clientAssertionProvider cannot be nil")
	}

	if c.clientAssertionProvider != nil {
		aud, extractErr := extractAudFromWSURL(c.domain)
		if extractErr != nil {
			return "", extractErr
		}
		clientAssertionToken, retrieveErr := c.clientAssertionProvider.RetrieveToken(ctx, aud)
		if retrieveErr != nil {
			return "", retrieveErr
		}
		if clientAssertionToken == nil || clientAssertionToken.Value == "" {
			return "", NewClientError(larkcore.ErrCodeClientAssertionTokenEmpty, "client assertion token is empty")
		}
		body.ClientAssertion = clientAssertionToken.Value
		body.AppSecret = ""
		if clientAssertionToken.TargetInfo != nil {
			requestURL = buildWSProxyURL(clientAssertionToken.TargetInfo.TargetService, clientAssertionToken.TargetInfo.TargetPrefix, GenEndpointUri)
			headers.Set(larkcore.HeaderXTargetService, aud)
		}
	} else {
		body.AppSecret = c.appSecret
	}

	bs, err := json.Marshal(body)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(bs))
	if err != nil {
		return
	}

	req.Header.Add("locale", "zh")
	req.Header.Add("Content-Type", "application/json")
	for k, values := range c.headers {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}
	for k, values := range headers {
		req.Header.Del(k)
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}
	req.Header.Set("User-Agent", larkcore.UserAgent(c.source))
	resp, err := bootstrapHTTPClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.logger.Warn(ctx, "response status code %d", resp.StatusCode)
		serverMsg := "system busy"
		bootstrapErrResp := &bootstrapErrorResp{}
		if json.Unmarshal(respBody, bootstrapErrResp) == nil {
			if bootstrapErrResp.Msg != "" {
				serverMsg = bootstrapErrResp.Msg
			}
		}
		err = NewServerError(resp.StatusCode, serverMsg)
		return
	}

	endpointResp := &EndpointResp{}
	err = json.Unmarshal(respBody, endpointResp)
	if err != nil {
		return
	}

	switch endpointResp.Code {
	case OK:
	case SystemBusy:
		return "", NewServerError(endpointResp.Code, "system busy")
	case InternalError:
		return "", NewServerError(endpointResp.Code, endpointResp.Msg)
	default:
		return "", NewClientError(endpointResp.Code, endpointResp.Msg)
	}

	endpoint := endpointResp.Data
	if endpoint == nil || endpoint.Url == "" {
		err = NewServerError(http.StatusInternalServerError, "endpoint is null")
		return
	}

	url = endpoint.Url
	if endpoint.ClientConfig != nil {
		c.configure(endpoint.ClientConfig)
	}

	return
}

func extractAudFromWSURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Host != "" {
		return parsedURL.Host, nil
	}
	return "", fmt.Errorf("invalid url: %s", rawURL)
}

func buildWSProxyURL(targetService, targetPrefix, apiPath string) string {
	if !strings.Contains(targetService, "://") {
		targetService = "https://" + targetService
	}
	return targetService + targetPrefix + apiPath
}

func (c *Client) pingLoop(run *clientRun) {
	ctx := run.ctx
	defer func() {
		if err := recover(); err != nil {
			c.logger.Warn(ctx, c.fmtLog("ping loop panic, err: %v, stack: %s", err, string(debug.Stack()))...)
			if ctx.Err() == nil {
				run.finish(fmt.Errorf("ping loop panic: %v", err))
			}
		}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return
		}
		conn, serviceID := c.connStateForRun(run)
		if conn != nil {
			i, _ := strconv.ParseInt(serviceID, 10, 32)
			frame := NewPingFrame(int32(i))
			bs, _ := frame.Marshal()

			err := c.writeMessageForRun(run, conn, ws.BinaryMessage, bs)
			if err != nil {
				c.logger.Warn(ctx, c.fmtLog("ping failed, err: %v", err)...)
			} else {
				c.logger.Debug(ctx, c.fmtLog("ping success")...)
			}
		}
		_, _, _, pingInterval := c.configSnapshot()
		if err := sleepWithContext(ctx, pingInterval); err != nil {
			return
		}
	}
}

func (c *Client) receiveMessageLoop(run *clientRun) {
	ctx := run.ctx
	defer func() {
		if err := recover(); err != nil {
			c.logger.Error(ctx, c.fmtLog("receive message loop panic, err: %v, stack: %s", err, string(debug.Stack()))...)
			c.disconnectRun(ctx, run, nil)
			if ctx.Err() == nil {
				run.finish(fmt.Errorf("receive message loop panic: %v", err))
			}
		}
	}()
	for {
		conn, _ := c.connStateForRun(run)
		if conn == nil {
			if ctx.Err() != nil {
				return
			}
			c.logger.Error(ctx, c.fmtLog("connection is closed, receive message loop exit")...)
			run.finish(fmt.Errorf("connection is closed"))
			return
		}

		mt, msg, err := conn.ReadMessage()
		if err != nil {
			c.logger.Error(ctx, c.fmtLog("receive message failed, err: %v", err)...)
			c.disconnectRun(ctx, run, conn)
			if ctx.Err() != nil {
				return
			}
			if !c.autoReconnectEnabled() {
				run.finish(err)
				return
			}
			if err = c.reconnect(run); err != nil {
				if ctx.Err() != nil {
					return
				}
				c.logger.Error(ctx, err)
				run.finish(err)
				return
			}
			continue
		}

		if mt != ws.BinaryMessage {
			c.logger.Warn(ctx, c.fmtLog("receive unknown message, message_type: %d, message: %s", mt, msg)...)
			continue
		}

		run.wg.Add(1)
		go func() {
			defer run.wg.Done()
			c.handleMessage(run, conn, msg)
		}()
	}
}

func (c *Client) handleMessage(run *clientRun, conn *ws.Conn, msg []byte) {
	ctx := run.ctx
	defer func() {
		if err := recover(); err != nil {
			c.logger.Error(ctx, c.fmtLog("handle message panic, err: %v, stack: %s", err, string(debug.Stack()))...)
		}
	}()

	var frame Frame
	if err := frame.Unmarshal(msg); err != nil {
		c.logger.Error(ctx, c.fmtLog("unmarshal message failed, error: %v", err)...)
		return
	}

	switch FrameType(frame.Method) {
	case FrameTypeControl:
		c.handleControlFrame(run, conn, frame)
	case FrameTypeData:
		c.handleDataFrame(run, conn, frame)
	default:
	}
}

func (c *Client) handleControlFrame(run *clientRun, conn *ws.Conn, frame Frame) {
	ctx := run.ctx
	hs := Headers(frame.Headers)
	t := hs.GetString(HeaderType)

	switch MessageType(t) {
	case MessageTypePong:
		c.logger.Debug(ctx, c.fmtLog("receive pong")...)
		if len(frame.Payload) == 0 {
			return
		}
		conf := &ClientConfig{}
		if err := json.Unmarshal(frame.Payload, conf); err != nil {
			c.logger.Warn(ctx, c.fmtLog("unmarshal client config failed, err: %v", err)...)
			return
		}
		c.configureForRun(run, conn, conf)
	default:
	}
}

func (c *Client) handleDataFrame(run *clientRun, conn *ws.Conn, frame Frame) {
	ctx := run.ctx
	hs := Headers(frame.Headers)
	sum := hs.GetInt(HeaderSum)
	seq := hs.GetInt(HeaderSeq)
	msgID := hs.GetString(HeaderMessageID)
	traceID := hs.GetString(HeaderTraceID)
	type_ := hs.GetString(HeaderType)

	pl := frame.Payload
	if sum > 1 {
		// 合包
		pl = c.combine(msgID, sum, seq, pl)
		if pl == nil {
			return
		}
	}

	c.logger.Debug(ctx, c.fmtLog("receive message, message_type: %s, message_id: %s, trace_id: %s, payload: %s",
		type_, msgID, traceID, pl))

	var err error
	var rsp interface{}
	start := time.Now().UnixNano() / int64(time.Millisecond) // 兼容 go < 1.17
	switch MessageType(type_) {
	case MessageTypeEvent:
		rsp, err = c.eventHandler.Do(ctx, pl)
	case MessageTypeCard:
		return
	default:
		return
	}
	end := time.Now().UnixNano() / int64(time.Millisecond)
	hs.Add(HeaderBizRt, strconv.FormatInt(end-start, 10))

	resp := NewResponseByCode(http.StatusOK)
	if err != nil {
		c.logger.Error(ctx, c.fmtLog("handle message failed, message_type: %s, message_id: %s, trace_id: %s, err: %v",
			type_, msgID, traceID, err)...)
		resp = NewResponseByCode(http.StatusInternalServerError)
	} else {
		if rsp != nil { // for cardCallback
			resp.Data, err = json.Marshal(rsp)
			if err != nil {
				c.logger.Error(ctx, c.fmtLog("handle message failed, message_type: %s, message_id: %s, trace_id: %s, err: %v",
					type_, msgID, traceID, err)...)
				resp = NewResponseByCode(http.StatusInternalServerError)
			}
		}
	}

	p, _ := json.Marshal(resp)
	frame.Payload = p
	frame.Headers = hs
	bs, _ := frame.Marshal()

	err = c.writeMessageForRun(run, conn, ws.BinaryMessage, bs)
	if err != nil {
		c.logger.Error(ctx, c.fmtLog("response message failed, type: %s, message_id: %s, trace_id: %s, err: %v", type_, msgID, traceID, err)...)
		return
	}
}

func (c *Client) combine(msgID string, sum, seq int, bs []byte) []byte {
	val := c.cache.Get(msgID)
	if val == nil {
		buf := make([][]byte, sum)
		buf[seq] = bs
		c.cache.Set(msgID, buf, 5*time.Second)
		return nil
	}

	buf := val.([][]byte)
	buf[seq] = bs
	capacity := 0
	for _, v := range buf {
		if len(v) == 0 {
			c.cache.Set(msgID, buf, 5*time.Second)
			return nil
		}
		capacity += len(v)
	}

	pl := make([]byte, 0, capacity)
	for _, v := range buf {
		pl = append(pl, v...)
	}

	return pl
}

func (c *Client) writeMessageForRun(run *clientRun, conn *ws.Conn, messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	c.mu.Lock()
	active := c.conn == conn && c.connRun == run
	c.mu.Unlock()
	if !active {
		return fmt.Errorf("connection is closed")
	}
	return conn.WriteMessage(messageType, data)
}

func (c *Client) fmtLog(format string, i ...interface{}) []interface{} {
	log := []interface{}{fmt.Sprintf(format, i...)}
	if connID := c.connIDSnapshot(); connID != "" {
		log = append(log, fmt.Sprintf("[conn_id=%s]", connID))
	}

	return log
}

func (c *Client) callbacksSnapshot() clientCallbacks {
	c.mu.Lock()
	defer c.mu.Unlock()
	return clientCallbacks{
		onReady:        c.onReady,
		onError:        c.onError,
		onReconnecting: c.onReconnecting,
		onReconnected:  c.onReconnected,
	}
}

func (c *Client) callOnError(err error) {
	if onError := c.callbacksSnapshot().onError; onError != nil {
		onError(err)
	}
}

func (c *Client) callConnectionCallback(run *clientRun, callback func()) bool {
	run.eventsMu.Lock()
	c.mu.Lock()
	active := run.ctx.Err() == nil && c.conn != nil && c.connRun == run
	if active {
		run.connectionCallbackActive = true
	}
	c.mu.Unlock()
	run.eventsMu.Unlock()
	if !active {
		return false
	}

	func() {
		defer c.finishConnectionCallback(run)
		if callback != nil {
			callback()
		}
	}()
	return true
}

func (c *Client) finishConnectionCallback(run *clientRun) {
	run.eventsMu.Lock()
	run.connectionCallbackActive = false
	pendingDisconnected := run.pendingDisconnected
	run.pendingDisconnected = nil
	run.eventsMu.Unlock()
	if pendingDisconnected != nil {
		pendingDisconnected()
	}
}

func (c *Client) autoReconnectEnabled() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.autoReconnect
}

func (c *Client) connIDSnapshot() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connID
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// configure applies a server-pushed connection config.
func (c *Client) configure(conf *ClientConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configureLocked(conf)
}

func (c *Client) configureForRun(run *clientRun, conn *ws.Conn, conf *ClientConfig) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connRun != run || c.conn != conn {
		return false
	}
	c.configureLocked(conf)
	return true
}

// configureLocked writes the connection-control fields. The caller must hold c.mu.
func (c *Client) configureLocked(conf *ClientConfig) {
	c.reconnectCount = conf.ReconnectCount
	c.reconnectInterval = time.Duration(conf.ReconnectInterval) * time.Second
	c.reconnectNonce = conf.ReconnectNonce
	c.pingInterval = time.Duration(conf.PingInterval) * time.Second
}

// configSnapshot reads the connection-control fields under c.mu and returns a
// consistent copy, so readers (pingLoop, reconnect) never observe a torn write.
func (c *Client) configSnapshot() (reconnectCount int, reconnectInterval time.Duration, reconnectNonce int, pingInterval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.reconnectCount, c.reconnectInterval, c.reconnectNonce, c.pingInterval
}

// connState reads conn/serviceID under c.mu. Callers must release the lock (i.e.
// use the returned snapshot) before any blocking call such as writeMessage or
// ReadMessage, since writeMessage takes c.mu again and the mutex is not reentrant.
func (c *Client) connState() (conn *ws.Conn, serviceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn, c.serviceID
}

func (c *Client) connStateForRun(run *clientRun) (conn *ws.Conn, serviceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.connRun != run {
		return nil, ""
	}
	return c.conn, c.serviceID
}

func parseErr(resp *http.Response) error {
	code, _ := strconv.Atoi(resp.Header.Get(HeaderHandshakeStatus))
	msg := resp.Header.Get(HeaderHandshakeMsg)
	switch code {
	case AuthFailed:
		// Auth失败
		authCode, _ := strconv.Atoi(resp.Header.Get(HeaderHandshakeAuthErrCode))
		if authCode == ExceedConnLimit {
			return NewClientError(code, msg)
		} else {
			return NewServerError(code, msg)
		}
	case Forbidden:
		// 被封禁
		return NewClientError(code, msg)
	default:
		return NewServerError(code, msg)
	}
}

// EventHandler returns the configured event dispatcher.
func (c *Client) EventHandler() *dispatcher.EventDispatcher {
	return c.eventHandler
}
