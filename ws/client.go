package ws

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	autoReconnect           bool          // 是否自动重连，默认开启
	reconnectNonce          int           // 首次重连抖动，单位秒
	reconnectCount          int           // 重连次数，负数无限次
	reconnectInterval       time.Duration // 重连间隔
	pingInterval            time.Duration // Ping间隔
	writeTimeout            time.Duration // 单次写入超时
	httpClient              larkcore.HttpClient
	websocketDialer         *websocket.Dialer
	cache                   *larkcache.Cache
	onReady                 func()
	onError                 func(err error)
	onReconnecting          func()
	onReconnected           func()
	onDisconnected          func()

	stateMu  sync.Mutex
	run      *clientRun
	terminal bool
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

// WithWriteTimeout sets the maximum duration of one websocket write. A
// non-positive timeout keeps the default value.
func WithWriteTimeout(timeout time.Duration) ClientOption {
	return func(cli *Client) {
		if timeout > 0 {
			cli.writeTimeout = timeout
		}
	}
}

// WithHttpClient sets the HTTP client used to fetch the websocket endpoint.
// When httpClient is an *http.Client, its Timeout also bounds that request.
func WithHttpClient(httpClient larkcore.HttpClient) ClientOption {
	return func(cli *Client) {
		if httpClient != nil {
			cli.httpClient = httpClient
		}
	}
}

// WithWebSocketDialer sets the dialer used to establish websocket connections.
// Its Proxy, TLSClientConfig and HandshakeTimeout settings apply only to this
// Client instance.
func WithWebSocketDialer(dialer *websocket.Dialer) ClientOption {
	return func(cli *Client) {
		if dialer != nil {
			cli.websocketDialer = dialer
		}
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
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onReady = f
}

func (c *Client) SetOnReconnecting(f func()) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onReconnecting = f
}

func (c *Client) SetOnReconnected(f func()) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onReconnected = f
}

func (c *Client) SetOnError(f func(err error)) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onError = f
}

func (c *Client) SetOnDisconnected(f func()) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.onDisconnected = f
}

func NewClient(appId, appSecret string, opts ...ClientOption) *Client {
	dialer := *websocket.DefaultDialer
	cli := &Client{
		appID:             appId,
		appSecret:         appSecret,
		autoReconnect:     true,
		reconnectNonce:    30,
		reconnectCount:    -1,
		reconnectInterval: 2 * time.Minute,
		pingInterval:      2 * time.Minute,
		writeTimeout:      10 * time.Second,
		httpClient:        &http.Client{Timeout: 10 * time.Second},
		websocketDialer:   &dialer,
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

func (c *Client) fmtLog(format string, i ...interface{}) []interface{} {
	return []interface{}{fmt.Sprintf(format, i...)}
}

func (c *Client) configure(conf *ClientConfig) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.applyConfigLocked(conf)
}

func (c *Client) applyConfig(run *clientRun, conf *ClientConfig) bool {
	if run.ctx.Err() != nil {
		return false
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if run.stopReason != runStopNone {
		return false
	}
	if conf != nil {
		c.applyConfigLocked(conf)
	}
	return true
}

func (c *Client) applyConfigLocked(conf *ClientConfig) {
	c.reconnectCount = conf.ReconnectCount
	c.reconnectInterval = time.Duration(conf.ReconnectInterval) * time.Second
	c.reconnectNonce = conf.ReconnectNonce
	c.pingInterval = time.Duration(conf.PingInterval) * time.Second
}

func (c *Client) configSnapshot() (reconnectCount int, reconnectInterval time.Duration, reconnectNonce int, pingInterval time.Duration) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.reconnectCount, c.reconnectInterval, c.reconnectNonce, c.pingInterval
}

// EventHandler returns the configured event dispatcher.
func (c *Client) EventHandler() *dispatcher.EventDispatcher {
	return c.eventHandler
}
