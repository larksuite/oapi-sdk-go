package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	websocket "github.com/gorilla/websocket"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

const lifecycleTestTimeout = 2 * time.Second
const lifecycleNoProgressTimeout = 100 * time.Millisecond

type lifecycleGateway struct {
	server   *httptest.Server
	accepted chan *websocket.Conn
	requests chan lifecycleGatewayRequest
	messages chan []byte
	frames   chan lifecycleGatewayFrame

	mu    sync.Mutex
	conns []*websocket.Conn
}

type lifecycleGatewayFrame struct {
	conn    *websocket.Conn
	payload []byte
}

type lifecycleGatewayRequest struct {
	path      string
	deviceID  string
	serviceID string
	username  string
	password  string
}

type lifecycleAssertionProvider struct {
	mu       sync.Mutex
	failures []error
	token    *larkcore.Token
	calls    int
}

func (p *lifecycleAssertionProvider) RetrieveToken(context.Context, string) (*larkcore.Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	if index := p.calls - 1; index < len(p.failures) && p.failures[index] != nil {
		return nil, p.failures[index]
	}
	return p.token, nil
}

func (p *lifecycleAssertionProvider) callCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls
}

func newLifecycleGateway(t *testing.T) (*lifecycleGateway, func()) {
	t.Helper()

	gateway := &lifecycleGateway{
		accepted: make(chan *websocket.Conn, 8),
		requests: make(chan lifecycleGatewayRequest, 8),
		messages: make(chan []byte, 8),
		frames:   make(chan lifecycleGatewayFrame, 16),
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	gateway.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		username, password, _ := r.BasicAuth()
		gateway.mu.Lock()
		gateway.conns = append(gateway.conns, conn)
		gateway.mu.Unlock()

		gateway.requests <- lifecycleGatewayRequest{
			path:      r.URL.Path,
			deviceID:  r.URL.Query().Get(DeviceID),
			serviceID: r.URL.Query().Get(ServiceID),
			username:  username,
			password:  password,
		}
		gateway.accepted <- conn

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			if messageType == websocket.BinaryMessage {
				payload := append([]byte(nil), message...)
				select {
				case gateway.messages <- payload:
				default:
				}
				select {
				case gateway.frames <- lifecycleGatewayFrame{conn: conn, payload: payload}:
				default:
				}
			}
		}
	}))
	cleanup := func() {
		gateway.mu.Lock()
		conns := append([]*websocket.Conn(nil), gateway.conns...)
		gateway.mu.Unlock()
		for _, conn := range conns {
			if err := conn.Close(); err != nil && !isClosedLifecycleConnectionError(err) {
				t.Errorf("close test gateway connection: %v", err)
			}
		}
		gateway.server.Close()
	}
	return gateway, cleanup
}

func (g *lifecycleGateway) endpoint() string {
	return "ws" + strings.TrimPrefix(g.server.URL, "http") + "/callback"
}

func installLifecycleBootstrap(t *testing.T, endpoint string, requests *int32) (string, func()) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(requests, 1)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&EndpointResp{
			Code: OK,
			Data: &Endpoint{
				Url: endpoint,
				ClientConfig: &ClientConfig{
					ReconnectCount:    0,
					ReconnectInterval: 1,
					PingInterval:      3600,
				},
			},
		}); err != nil {
			return
		}
	}))
	originalClient := bootstrapHTTPClient
	bootstrapHTTPClient = server.Client()
	cleanup := func() {
		bootstrapHTTPClient = originalClient
		server.Close()
	}
	return server.URL, cleanup
}

func installLifecycleBootstrapHandler(t *testing.T, handler http.Handler) (string, func()) {
	t.Helper()

	server := httptest.NewServer(handler)
	originalClient := bootstrapHTTPClient
	bootstrapHTTPClient = server.Client()
	cleanup := func() {
		bootstrapHTTPClient = originalClient
		server.Close()
	}
	return server.URL, cleanup
}

func startLifecycleClient(client *Client, ctx context.Context) <-chan error {
	result := make(chan error, 1)
	go func() {
		result <- client.Start(ctx)
	}()
	return result
}

func waitLifecycleSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(lifecycleTestTimeout):
		t.Fatalf("timed out waiting for %s", name)
	}
}

func waitLifecycleResult(result <-chan error) (error, bool) {
	select {
	case err := <-result:
		return err, true
	case <-time.After(lifecycleTestTimeout):
		return nil, false
	}
}

func assertNoLifecycleSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), lifecycleNoProgressTimeout)
	defer cancel()
	select {
	case <-signal:
		t.Fatalf("%s happened while it should have been blocked", name)
	case <-ctx.Done():
	}
}

func assertNoLifecycleResult(t *testing.T, result <-chan error, name string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), lifecycleNoProgressTimeout)
	defer cancel()
	select {
	case err := <-result:
		t.Fatalf("%s returned %v while it should have been blocked", name, err)
	case <-ctx.Done():
	}
}

func closeGatewayConn(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	if err := conn.Close(); err != nil && !isClosedLifecycleConnectionError(err) {
		t.Errorf("close gateway connection: %v", err)
	}
}

func isClosedLifecycleConnectionError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "use of closed network connection") ||
		strings.Contains(message, "closed network connection") ||
		strings.Contains(message, "connection already closed")
}

func TestStartNilDoesNotConsumeLaterValidStart(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	var callbackCount int32
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() {
			atomic.AddInt32(&callbackCount, 1)
			ready <- struct{}{}
		}),
		WithOnError(func(error) {
			atomic.AddInt32(&callbackCount, 1)
		}),
	)

	if err := client.Start(nil); err == nil {
		t.Fatal("Start(nil) returned nil, want a parameter error")
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 0 {
		t.Fatalf("Start(nil) made %d bootstrap requests, want 0", got)
	}
	if got := atomic.LoadInt32(&callbackCount); got != 0 {
		t.Fatalf("Start(nil) invoked %d lifecycle callbacks, want 0", got)
	}

	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready after a valid Start")
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Error("valid Start did not return after Close")
	} else if err != nil {
		t.Errorf("valid Start returned %v after Close, want nil", err)
	}
}

func TestInitialFailureThenReconnectUsesReconnectedCallback(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		endpoint := gateway.endpoint()
		if requestNumber == 1 {
			endpoint = "ws://127.0.0.1:1/unreachable"
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&EndpointResp{
			Code: OK,
			Data: &Endpoint{
				Url: endpoint,
				ClientConfig: &ClientConfig{
					ReconnectCount:    1,
					ReconnectInterval: 0,
					ReconnectNonce:    0,
					PingInterval:      3600,
				},
			},
		}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	callbackDone := make(chan struct{}, 1)
	var callbackMu sync.Mutex
	callbackOrder := make([]string, 0, 3)
	recordCallback := func(name string) {
		callbackMu.Lock()
		callbackOrder = append(callbackOrder, name)
		callbackMu.Unlock()
	}
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnError(func(error) { recordCallback("error") }),
		WithOnReconnecting(func() { recordCallback("reconnecting") }),
		WithOnReconnected(func() {
			recordCallback("reconnected")
			callbackDone <- struct{}{}
		}),
	)
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, callbackDone, "ordered recovery callbacks after initial failure")
	callbackMu.Lock()
	callbackSequence := strings.Join(callbackOrder, ",")
	callbackMu.Unlock()
	if callbackSequence != "error,reconnecting,reconnected" {
		t.Fatalf("recovery callback order = %q, want error,reconnecting,reconnected", callbackSequence)
	}
	select {
	case <-ready:
		t.Fatal("initial reconnect invoked Ready, want Reconnected only")
	default:
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 2 {
		t.Fatalf("bootstrap requests = %d, want 2", got)
	}

	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestRecoverableErrorBlocksReconnectUntilReturned(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	secondBootstrap := make(chan struct{}, 1)
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		if requestNumber == 2 {
			secondBootstrap <- struct{}{}
		}
		endpoint := gateway.endpoint()
		if requestNumber == 1 {
			endpoint = "ws://127.0.0.1:1/unreachable"
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
			Url: endpoint,
			ClientConfig: &ClientConfig{
				ReconnectCount:    1,
				ReconnectInterval: 0,
				ReconnectNonce:    0,
				PingInterval:      3600,
			},
		}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	errorEntered := make(chan struct{}, 1)
	releaseError := make(chan struct{})
	reconnecting := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnError(func(error) {
			errorEntered <- struct{}{}
			<-releaseError
		}),
		WithOnReconnecting(func() { reconnecting <- struct{}{} }),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
	)
	defer client.Close()
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, errorEntered, "recoverable Error callback entry")
	assertNoLifecycleSignal(t, reconnecting, "Reconnecting callback")
	assertNoLifecycleSignal(t, secondBootstrap, "second bootstrap")
	close(releaseError)

	waitLifecycleSignal(t, reconnecting, "Reconnecting after Error returned")
	waitLifecycleSignal(t, secondBootstrap, "second bootstrap after Error returned")
	waitLifecycleSignal(t, reconnected, "Reconnected after Error returned")
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestReconnectingBlocksDialUntilReturned(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	secondBootstrap := make(chan struct{}, 1)
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		if requestNumber == 2 {
			secondBootstrap <- struct{}{}
		}
		endpoint := gateway.endpoint()
		if requestNumber == 1 {
			endpoint = "ws://127.0.0.1:1/unreachable"
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
			Url: endpoint,
			ClientConfig: &ClientConfig{
				ReconnectCount:    1,
				ReconnectInterval: 0,
				ReconnectNonce:    0,
				PingInterval:      3600,
			},
		}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	reconnectingEntered := make(chan struct{}, 1)
	releaseReconnecting := make(chan struct{})
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReconnecting(func() {
			reconnectingEntered <- struct{}{}
			<-releaseReconnecting
		}),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
	)
	defer client.Close()
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, reconnectingEntered, "Reconnecting callback entry")
	assertNoLifecycleSignal(t, secondBootstrap, "second bootstrap")
	close(releaseReconnecting)

	waitLifecycleSignal(t, secondBootstrap, "second bootstrap after Reconnecting returned")
	waitLifecycleSignal(t, reconnected, "Reconnected after Reconnecting returned")
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestReconnectSequenceNotifiesReconnectingOnce(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		endpoint := gateway.endpoint()
		if requestNumber < 3 {
			endpoint = "ws://127.0.0.1:1/unreachable"
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
			Url: endpoint,
			ClientConfig: &ClientConfig{
				ReconnectCount:    2,
				ReconnectInterval: 0,
				ReconnectNonce:    0,
				PingInterval:      3600,
			},
		}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	var callbackMu sync.Mutex
	callbackOrder := make([]string, 0, 4)
	record := func(name string) {
		callbackMu.Lock()
		callbackOrder = append(callbackOrder, name)
		callbackMu.Unlock()
	}
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnError(func(error) { record("error") }),
		WithOnReconnecting(func() { record("reconnecting") }),
		WithOnReconnected(func() {
			record("reconnected")
			reconnected <- struct{}{}
		}),
	)
	defer client.Close()
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, reconnected, "Reconnected after retry sequence")
	callbackMu.Lock()
	sequence := strings.Join(callbackOrder, ",")
	callbackMu.Unlock()
	if sequence != "error,reconnecting,error,reconnected" {
		t.Fatalf("reconnect callback sequence = %q, want error,reconnecting,error,reconnected", sequence)
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 3 {
		t.Fatalf("bootstrap requests = %d, want 3", got)
	}

	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestRetrieveTokenFailureRetriesWhenRecoverable(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	provider := &lifecycleAssertionProvider{
		failures: []error{errors.New("temporary assertion provider failure")},
		token:    &larkcore.Token{Value: "assertion"},
	}
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "", WithDomain(domain), WithClientAssertionProvider(provider), WithOnReconnected(func() {
		reconnected <- struct{}{}
	}))
	client.reconnectCount = 1
	client.reconnectInterval = 0
	client.reconnectNonce = 0
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, reconnected, "Reconnected after assertion token retry")
	if got := provider.callCount(); got != 2 {
		t.Fatalf("RetrieveToken calls = %d, want 2", got)
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 1 {
		t.Fatalf("bootstrap requests = %d, want 1", got)
	}

	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestRetrieveTokenClientErrorDoesNotRetry(t *testing.T) {
	provider := &lifecycleAssertionProvider{
		failures: []error{NewClientError(19001, "invalid assertion")},
	}
	client := NewClient("app-id", "", WithClientAssertionProvider(provider))

	err := client.Start(context.Background())
	if err == nil {
		t.Fatal("Start returned nil after ClientAssertionProvider ClientError")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("Start returned %T, want *ClientError", err)
	}
	if clientErr.Code != 19001 || clientErr.Msg != "invalid assertion" {
		t.Fatalf("ClientError = %#v, want original assertion provider error", clientErr)
	}
	if got := provider.callCount(); got != 1 {
		t.Fatalf("RetrieveToken calls = %d, want 1", got)
	}
}

func TestConnectedSessionStartsPingAndAppliesPongConfig(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	logger := &lifecycleRecordingLogger{}
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithLogger(logger),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	var serverConn *websocket.Conn
	select {
	case serverConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for gateway connection")
	}

	var pingPayload []byte
	select {
	case pingPayload = <-gateway.messages:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("connected session did not send Ping")
	}
	var ping Frame
	if err := ping.Unmarshal(pingPayload); err != nil {
		t.Fatalf("decode Ping: %v", err)
	}
	if FrameType(ping.Method) != FrameTypeControl || Headers(ping.Headers).GetString(HeaderType) != string(MessageTypePing) {
		t.Fatalf("first session frame = method %d, type %q; want control Ping", ping.Method, Headers(ping.Headers).GetString(HeaderType))
	}
	if logs := logger.String(); !strings.Contains(logs, "ping success") || strings.Contains(logs, "websocket ping succeeded") {
		t.Fatalf("Ping logs = %q, want v3_main label ping success", logs)
	}

	pongConfig, err := json.Marshal(&ClientConfig{PingInterval: 7})
	if err != nil {
		t.Fatalf("encode Pong config: %v", err)
	}
	pongHeaders := Headers{}
	pongHeaders.Add(HeaderType, string(MessageTypePong))
	pong := &Frame{Method: int32(FrameTypeControl), Headers: pongHeaders, Payload: pongConfig}
	pongPayload, err := pong.Marshal()
	if err != nil {
		t.Fatalf("encode Pong: %v", err)
	}
	if err := serverConn.WriteMessage(websocket.BinaryMessage, pongPayload); err != nil {
		t.Fatalf("send Pong: %v", err)
	}

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.NewTimer(lifecycleTestTimeout)
	defer timeout.Stop()
	for {
		_, _, _, pingInterval := client.configSnapshot()
		if pingInterval == 7*time.Second {
			break
		}
		select {
		case <-ticker.C:
		case <-timeout.C:
			t.Fatalf("Pong config was not applied; ping interval = %v", pingInterval)
		}
	}
	if logs := logger.String(); !strings.Contains(logs, "receive pong") || strings.Contains(logs, "websocket pong received") {
		t.Fatalf("Pong logs = %q, want v3_main label receive pong", logs)
	}

	client.Close()
	if startErr, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if startErr != nil {
		t.Fatalf("Start returned %v after Close, want nil", startErr)
	}
}

func TestConnectionTimesOutWithoutInboundFrame(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: gateway.endpoint()}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	client := NewClient("app-id", "app-secret", WithDomain(domain), WithAutoReconnect(false))
	client.pingInterval = 10 * time.Millisecond
	result := startLifecycleClient(client, context.Background())

	waitAcceptedGatewayConnection(t, gateway)
	pingFrame := waitGatewayFrame(t, gateway, "initial Ping")
	var ping Frame
	if err := ping.Unmarshal(pingFrame.payload); err != nil {
		t.Fatalf("decode initial Ping: %v", err)
	}
	if FrameType(ping.Method) != FrameTypeControl || Headers(ping.Headers).GetString(HeaderType) != string(MessageTypePing) {
		t.Fatalf("initial frame = method %d, type %q; want control Ping", ping.Method, Headers(ping.Headers).GetString(HeaderType))
	}

	select {
	case err := <-result:
		var timeoutErr net.Error
		if !errors.As(err, &timeoutErr) || !timeoutErr.Timeout() {
			t.Fatalf("Start returned %v, want read timeout after no inbound frame", err)
		}
	case <-time.After(pongWait(client.pingInterval) + 2*time.Second):
		t.Fatal("Start did not return after the read deadline elapsed")
	}
}

func TestSDKPongRefreshesReadDeadline(t *testing.T) {
	pongHeaders := Headers{}
	pongHeaders.Add(HeaderType, string(MessageTypePong))
	testInboundFrameRefreshesReadDeadline(t, &Frame{Method: int32(FrameTypeControl), Headers: pongHeaders}, "Pong")
}

func TestDataFrameRefreshesReadDeadline(t *testing.T) {
	testInboundFrameRefreshesReadDeadline(t, &Frame{Method: int32(FrameTypeData)}, "Data frame")
}

func testInboundFrameRefreshesReadDeadline(t *testing.T, frame *Frame, frameName string) {
	t.Helper()
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: gateway.endpoint()}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	client := NewClient("app-id", "app-secret", WithDomain(domain), WithAutoReconnect(false))
	client.pingInterval = 10 * time.Millisecond
	result := startLifecycleClient(client, context.Background())

	serverConn := waitAcceptedGatewayConnection(t, gateway)
	waitGatewayFrame(t, gateway, "initial Ping")

	client.stateMu.Lock()
	conn := client.run.conn
	client.stateMu.Unlock()
	if conn == nil {
		t.Fatal("connected client has no current connection")
	}
	if err := conn.socket.SetReadDeadline(time.Now().Add(500 * time.Millisecond)); err != nil {
		t.Fatalf("set short read deadline: %v", err)
	}

	payload, err := frame.Marshal()
	if err != nil {
		t.Fatalf("encode %s: %v", frameName, err)
	}
	if err := serverConn.WriteMessage(websocket.BinaryMessage, payload); err != nil {
		t.Fatalf("send %s: %v", frameName, err)
	}

	select {
	case err := <-result:
		t.Fatalf("Start returned %v after a %s refreshed the read deadline", err, frameName)
	case <-time.After(750 * time.Millisecond):
	}

	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestReconnectKeepsOneRunPingUsingCurrentConnection(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		serviceID := 101
		if requestNumber > 1 {
			serviceID = 202
		}
		endpoint := fmt.Sprintf("%s?%s=%d", gateway.endpoint(), ServiceID, serviceID)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&EndpointResp{
			Code: OK,
			Data: &Endpoint{
				Url: endpoint,
				ClientConfig: &ClientConfig{
					ReconnectCount:    1,
					ReconnectInterval: 1,
					ReconnectNonce:    0,
					PingInterval:      1,
				},
			},
		}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "initial Ready")
	firstServerConn := waitAcceptedGatewayConnection(t, gateway)
	firstPing := waitGatewayFrame(t, gateway, "initial Ping")
	assertGatewayPing(t, firstPing, firstServerConn, 101)

	closeGatewayConn(t, firstServerConn)
	waitLifecycleSignal(t, reconnected, "Reconnected")
	secondServerConn := waitAcceptedGatewayConnection(t, gateway)

	select {
	case frame := <-gateway.frames:
		t.Fatalf("replacement connection started a second immediate Ping worker: connection=%p payload=%q", frame.conn, frame.payload)
	case <-time.After(200 * time.Millisecond):
	}
	secondPing := waitGatewayFrame(t, gateway, "scheduled Ping on replacement connection")
	assertGatewayPing(t, secondPing, secondServerConn, 202)

	client.Close()
	if startErr, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if startErr != nil {
		t.Fatalf("Start returned %v after Close, want nil", startErr)
	}
	select {
	case frame := <-gateway.frames:
		t.Fatalf("Ping worker wrote after Close: connection=%p payload=%q", frame.conn, frame.payload)
	case <-time.After(200 * time.Millisecond):
	}
}

func assertGatewayPing(t *testing.T, frame lifecycleGatewayFrame, wantConn *websocket.Conn, wantService int32) {
	t.Helper()
	if frame.conn != wantConn {
		t.Fatalf("Ping used connection %p, want %p", frame.conn, wantConn)
	}
	var ping Frame
	if err := ping.Unmarshal(frame.payload); err != nil {
		t.Fatalf("decode Ping: %v", err)
	}
	if FrameType(ping.Method) != FrameTypeControl || Headers(ping.Headers).GetString(HeaderType) != string(MessageTypePing) {
		t.Fatalf("frame = method %d, type %q; want control Ping", ping.Method, Headers(ping.Headers).GetString(HeaderType))
	}
	if ping.Service != wantService {
		t.Fatalf("Ping service = %d, want %d", ping.Service, wantService)
	}
}

func TestPreCanceledStartIsTerminalWithoutNetworkOrErrorCallback(t *testing.T) {
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&bootstrapRequests, 1)
		http.Error(w, "unexpected bootstrap", http.StatusServiceUnavailable)
	}))
	defer cleanupBootstrap()

	var errorCallbacks int32
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnError(func(error) { atomic.AddInt32(&errorCallbacks, 1) }),
	)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := client.Start(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("pre-canceled Start returned %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 0 {
		t.Fatalf("pre-canceled Start made %d bootstrap requests, want 0", got)
	}
	if got := atomic.LoadInt32(&errorCallbacks); got != 0 {
		t.Fatalf("pre-canceled Start invoked OnError %d times, want 0", got)
	}

	if err := client.Start(context.Background()); err == nil {
		t.Fatal("second Start returned nil after a pre-canceled Start, want lifecycle error")
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 0 {
		t.Fatalf("terminal Client made %d bootstrap requests, want 0", got)
	}
}

func TestConcurrentStartRejectsSecondCallerBeforeNetworkUnblocks(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	requestStarted := make(chan struct{}, 1)
	releaseBootstrap := make(chan struct{})
	var releaseBootstrapOnce sync.Once
	releaseBootstrapNow := func() {
		releaseBootstrapOnce.Do(func() { close(releaseBootstrap) })
	}
	defer releaseBootstrapNow()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&bootstrapRequests, 1)
		requestStarted <- struct{}{}
		select {
		case <-releaseBootstrap:
		case <-r.Context().Done():
			return
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: gateway.endpoint()}}); err != nil {
			return
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 2)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	first := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, requestStarted, "first bootstrap request")
	second := startLifecycleClient(client, context.Background())

	select {
	case err := <-second:
		if err == nil {
			t.Error("concurrent Start returned nil, want lifecycle error")
		}
	case <-time.After(300 * time.Millisecond):
		t.Error("concurrent Start blocked behind the active lifecycle")
	}

	releaseBootstrapNow()
	waitLifecycleSignal(t, ready, "first Ready callback")
	client.Close()
	if err, ok := waitLifecycleResult(first); !ok {
		t.Error("first Start did not return after Close")
	} else if err != nil {
		t.Errorf("first Start returned %v after Close, want nil", err)
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 1 {
		t.Errorf("concurrent Starts made %d bootstrap requests, want 1", got)
	}
}

func TestContextCancellationWinsOverReadFailure(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	disconnected := make(chan struct{}, 1)
	var errorCallbacks int32
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
		WithOnError(func(error) { atomic.AddInt32(&errorCallbacks, 1) }),
	)
	ctx, cancel := context.WithCancel(context.Background())
	result := startLifecycleClient(client, ctx)
	waitLifecycleSignal(t, ready, "Ready")

	cancel()
	err, ok := waitLifecycleResult(result)
	if !ok {
		t.Error("Start did not return after caller context cancellation")
		client.Close()
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("Start returned %v, want context.Canceled", err)
	}
	waitLifecycleSignal(t, disconnected, "Disconnected after context cancellation")
	if got := atomic.LoadInt32(&errorCallbacks); got != 0 {
		t.Errorf("normal context cancellation invoked OnError %d times, want 0", got)
	}
}

func TestCloseWinsOverReadFailureAndIsIdempotent(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	disconnected := make(chan struct{}, 8)
	var errorCallbacks int32
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
		WithOnError(func(error) { atomic.AddInt32(&errorCallbacks, 1) }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	var closeWG sync.WaitGroup
	for i := 0; i < 8; i++ {
		closeWG.Add(1)
		go func() {
			defer closeWG.Done()
			client.Close()
		}()
	}
	closeDone := make(chan struct{})
	go func() {
		closeWG.Wait()
		close(closeDone)
	}()
	waitLifecycleSignal(t, closeDone, "concurrent Close calls")

	if err, ok := waitLifecycleResult(result); !ok {
		t.Error("Start did not return after Close")
	} else if err != nil {
		t.Errorf("Start returned %v after Close, want nil", err)
	}
	waitLifecycleSignal(t, disconnected, "Disconnected after Close")
	select {
	case <-disconnected:
		t.Error("Disconnected ran more than once for one connection")
	default:
	}
	if got := atomic.LoadInt32(&errorCallbacks); got != 0 {
		t.Errorf("normal Close invoked OnError %d times, want 0", got)
	}
}

func TestStartWaitsForDisconnectedCallback(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	disconnectedEntered := make(chan struct{}, 1)
	releaseDisconnected := make(chan struct{})
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnDisconnected(func() {
			disconnectedEntered <- struct{}{}
			<-releaseDisconnected
		}),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	client.Close()
	waitLifecycleSignal(t, disconnectedEntered, "Disconnected callback entry")
	assertNoLifecycleResult(t, result, "Start")

	close(releaseDisconnected)
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Disconnected callback returned")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestReadyCallbackPanicIsRedacted(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	const panicMarker = "ready-callback-sensitive-panic"
	readyEntered := make(chan struct{}, 1)
	logger := &lifecycleRecordingLogger{}
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithLogger(logger),
		WithOnReady(func() {
			readyEntered <- struct{}{}
			panic(panicMarker)
		}),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, readyEntered, "Ready callback entry")

	client.Close()
	if _, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	}
	assertLifecycleTextDoesNotContain(t, logger.String(), panicMarker)
}

func TestReconnectFailureReportsEachAttemptOnce(t *testing.T) {
	const serverMessage = "bootstrap-server-sensitive-marker"
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&bootstrapRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: InternalError, Msg: serverMessage}); err != nil {
			t.Errorf("encode bootstrap error: %v", err)
		}
	}))
	defer cleanupBootstrap()

	logger := &lifecycleRecordingLogger{}
	callbackErrs := make(chan error, 3)
	client := NewClient("app-id", "app-secret", WithDomain(domain), WithLogger(logger))
	client.SetOnError(func(err error) { callbackErrs <- err })
	client.reconnectCount = 1
	client.reconnectInterval = 0
	client.reconnectNonce = 0

	err := client.Start(context.Background())
	if err == nil {
		t.Fatal("Start returned nil after reconnect attempts were exhausted")
	}
	if err.Error() != "unable to connect to server after 1 retries" {
		t.Fatalf("Start error = %q, want reconnect exhaustion error", err)
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 2 {
		t.Fatalf("bootstrap requests = %d, want 2", got)
	}
	if got := len(callbackErrs); got != 2 {
		t.Fatalf("OnError calls = %d, want 2 for the initial and retried connection failures", got)
	}
	for index := 0; index < 2; index++ {
		if callbackErr := <-callbackErrs; callbackErr == nil {
			t.Errorf("OnError call %d received nil", index+1)
		} else if _, ok := callbackErr.(*ServerError); !ok {
			t.Errorf("OnError call %d received %T, want *ServerError", index+1, callbackErr)
		}
	}

	logs := logger.String()
	for _, message := range []string{
		"websocket initial connection failed",
		"websocket reconnecting: attempt 1/1",
		"websocket reconnect attempt 1 failed",
		"websocket reconnect attempts exhausted after 1 attempts",
	} {
		if !strings.Contains(logs, message) {
			t.Errorf("lifecycle logs missing %q: %s", message, logs)
		}
	}
	assertLifecycleTextDoesNotContain(t, logs, serverMessage)
}

func TestReconnectExhaustionWithNoAttemptsDoesNotRepeatOnError(t *testing.T) {
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&bootstrapRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: InternalError, Msg: "temporarily unavailable"}); err != nil {
			t.Errorf("encode bootstrap error: %v", err)
		}
	}))
	defer cleanupBootstrap()

	callbackErrs := make(chan error, 2)
	client := NewClient("app-id", "app-secret", WithDomain(domain), WithOnError(func(err error) {
		callbackErrs <- err
	}))
	client.reconnectCount = 0

	if err := client.Start(context.Background()); err == nil {
		t.Fatal("Start returned nil after reconnect attempts were exhausted")
	}
	if got := atomic.LoadInt32(&bootstrapRequests); got != 1 {
		t.Fatalf("bootstrap requests = %d, want 1", got)
	}
	if got := len(callbackErrs); got != 1 {
		t.Fatalf("OnError calls = %d, want 1 for the initial connection failure", got)
	}
}

func TestRuntimeDeadlineClosesConnectionAndReturnsDeadline(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	disconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	result := startLifecycleClient(client, ctx)
	waitLifecycleSignal(t, ready, "Ready before deadline")

	err, ok := waitLifecycleResult(result)
	if !ok {
		t.Error("Start did not return after its running context deadline")
		client.Close()
	} else if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Start returned %v, want context.DeadlineExceeded", err)
	}
	waitLifecycleSignal(t, disconnected, "Disconnected after deadline")
	if got := atomic.LoadInt32(&bootstrapRequests); got != 1 {
		t.Errorf("deadline triggered %d bootstrap requests, want 1", got)
	}
}

func TestTerminalReadFailureWinsOverReentrantClose(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	errorCallback := make(chan error, 2)
	var client *Client
	client = NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnError(func(err error) {
			errorCallback <- err
			client.Close()
		}),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	var serverConn *websocket.Conn
	select {
	case serverConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for gateway connection")
	}
	closeGatewayConn(t, serverConn)

	var callbackErr error
	select {
	case callbackErr = <-errorCallback:
	case <-time.After(lifecycleTestTimeout):
		t.Error("terminal read failure did not invoke OnError")
	}
	startErr, ok := waitLifecycleResult(result)
	if !ok {
		t.Error("terminal read failure did not end Start")
		client.Close()
	} else if startErr == nil {
		t.Error("terminal read failure followed by Close returned nil, want failure")
	}
	if callbackErr != nil && startErr != nil && callbackErr.Error() != startErr.Error() {
		t.Errorf("OnError received %q, Start returned %q", callbackErr, startErr)
	}
	select {
	case extra := <-errorCallback:
		t.Errorf("OnError ran more than once, extra error: %v", extra)
	default:
	}
}

func TestDisconnectedWaitsForReadyCallback(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	readyEntered := make(chan struct{}, 1)
	releaseReady := make(chan struct{})
	disconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() {
			readyEntered <- struct{}{}
			<-releaseReady
		}),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, readyEntered, "Ready callback entry")

	var serverConn *websocket.Conn
	select {
	case serverConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for gateway connection")
	}
	closeGatewayConn(t, serverConn)

	if startErr, startReturned := waitLifecycleResult(result); startReturned {
		t.Fatalf("Start returned %v while Ready callback was blocked", startErr)
	}
	close(releaseReady)

	waitLifecycleSignal(t, disconnected, "Disconnected after Ready callback returned")
	startErr, startReturned := waitLifecycleResult(result)
	if !startReturned {
		t.Error("Start did not return after Ready callback returned")
		client.Close()
	} else if startErr == nil {
		t.Error("server disconnect returned nil, want terminal failure")
	}
}

func TestDisconnectedWaitsForReconnectedCallback(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	thirdBootstrap := make(chan struct{}, 1)
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
		if requestNumber == 3 {
			thirdBootstrap <- struct{}{}
		}
		reconnectInterval := 0
		if requestNumber == 2 {
			reconnectInterval = 1
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
			Url: gateway.endpoint(),
			ClientConfig: &ClientConfig{
				ReconnectCount:    2,
				ReconnectInterval: reconnectInterval,
				ReconnectNonce:    0,
				PingInterval:      3600,
			},
		}}); err != nil {
			t.Errorf("encode bootstrap response: %v", err)
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	reconnecting := make(chan struct{}, 2)
	reconnectedEntered := make(chan struct{}, 1)
	releaseReconnected := make(chan struct{})
	disconnected := make(chan struct{}, 2)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnReconnecting(func() { reconnecting <- struct{}{} }),
		WithOnReconnected(func() {
			reconnectedEntered <- struct{}{}
			<-releaseReconnected
		}),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
	)
	defer client.Close()
	result := startLifecycleClient(client, context.Background())

	waitLifecycleSignal(t, ready, "initial Ready")
	firstServerConn := waitAcceptedGatewayConnection(t, gateway)
	closeGatewayConn(t, firstServerConn)
	waitLifecycleSignal(t, reconnecting, "first Reconnecting")
	secondServerConn := waitAcceptedGatewayConnection(t, gateway)
	waitLifecycleSignal(t, reconnectedEntered, "Reconnected callback entry")
	waitLifecycleSignal(t, disconnected, "first Disconnected")

	closeGatewayConn(t, secondServerConn)
	assertNoLifecycleSignal(t, disconnected, "second Disconnected")
	assertNoLifecycleResult(t, result, "Start")
	assertNoLifecycleSignal(t, thirdBootstrap, "third bootstrap")
	close(releaseReconnected)

	waitLifecycleSignal(t, disconnected, "second Disconnected after Reconnected returned")
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}

func TestLifecycleCallbacksCanCloseClient(t *testing.T) {
	for _, stage := range []string{"ready", "recoverable_error", "reconnecting", "reconnected"} {
		stage := stage
		t.Run(stage, func(t *testing.T) {
			gateway, cleanupGateway := newLifecycleGateway(t)
			defer cleanupGateway()

			var bootstrapRequests int32
			domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestNumber := atomic.AddInt32(&bootstrapRequests, 1)
				endpoint := gateway.endpoint()
				if stage != "ready" && requestNumber == 1 {
					endpoint = "ws://127.0.0.1:1/unreachable"
				}
				if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
					Url: endpoint,
					ClientConfig: &ClientConfig{
						ReconnectCount:    1,
						ReconnectInterval: 0,
						ReconnectNonce:    0,
						PingInterval:      3600,
					},
				}}); err != nil {
					t.Errorf("encode bootstrap response: %v", err)
				}
			}))
			defer cleanupBootstrap()

			callbackReturned := make(chan struct{}, 1)
			var client *Client
			closeFromCallback := func() {
				client.Close()
				callbackReturned <- struct{}{}
			}
			options := []ClientOption{WithDomain(domain)}
			switch stage {
			case "ready":
				options = append(options, WithOnReady(closeFromCallback))
			case "recoverable_error":
				options = append(options, WithOnError(func(error) { closeFromCallback() }))
			case "reconnecting":
				options = append(options, WithOnReconnecting(closeFromCallback))
			case "reconnected":
				options = append(options, WithOnReconnected(closeFromCallback))
			}
			client = NewClient("app-id", "app-secret", options...)
			defer client.Close()
			result := startLifecycleClient(client, context.Background())

			waitLifecycleSignal(t, callbackReturned, stage+" callback Close return")
			if err, ok := waitLifecycleResult(result); !ok {
				t.Fatal("Start did not return after callback Close")
			} else if err != nil {
				t.Fatalf("Start returned %v after callback Close, want nil", err)
			}
		})
	}
}

func TestNilReadyCallbackDoesNotBlockDisconnected(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var bootstrapRequests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &bootstrapRequests)
	defer cleanupBootstrap()

	disconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())

	var serverConn *websocket.Conn
	select {
	case serverConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for gateway connection")
	}
	closeGatewayConn(t, serverConn)
	waitLifecycleSignal(t, disconnected, "Disconnected with nil Ready callback")
	if err, ok := waitLifecycleResult(result); !ok {
		t.Error("nil Ready callback left Start blocked after disconnect")
		client.Close()
	} else if err == nil {
		t.Error("server disconnect returned nil, want terminal failure")
	}
}

func TestCloseCancelsBlockedBootstrapAndReturnsNil(t *testing.T) {
	requestStarted := make(chan struct{}, 1)
	requestCanceled := make(chan struct{}, 1)
	releaseHandler := make(chan struct{})
	var releaseHandlerOnce sync.Once
	releaseHandlerNow := func() {
		releaseHandlerOnce.Do(func() { close(releaseHandler) })
	}
	defer releaseHandlerNow()
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode blocked bootstrap request: %v", err)
			return
		}
		requestStarted <- struct{}{}
		select {
		case <-r.Context().Done():
			requestCanceled <- struct{}{}
		case <-releaseHandler:
			http.Error(w, "released by test", http.StatusServiceUnavailable)
		}
	}))
	defer cleanupBootstrap()

	client := NewClient("app-id", "app-secret", WithDomain(domain), WithAutoReconnect(false))
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, requestStarted, "bootstrap request")

	closeReturned := make(chan struct{})
	go func() {
		client.Close()
		close(closeReturned)
	}()
	select {
	case <-closeReturned:
	case <-time.After(300 * time.Millisecond):
		t.Error("Close blocked while bootstrap was in flight")
		releaseHandlerNow()
		waitLifecycleSignal(t, closeReturned, "Close after releasing bootstrap")
	}

	select {
	case <-requestCanceled:
	case <-time.After(300 * time.Millisecond):
		t.Error("Close did not cancel the in-flight bootstrap request")
	}
	if err, ok := waitLifecycleResult(result); !ok {
		t.Error("Start did not return after Close canceled bootstrap")
	} else if err != nil {
		t.Errorf("Start returned %v after Close, want nil", err)
	}
}

func TestStartupFailureCanRetryButSuccessfulRunIsTerminal(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var requests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&requests, 1)
		if requestNumber == 1 || requestNumber > 2 {
			http.Error(w, fmt.Sprintf("bootstrap failure %d", requestNumber), http.StatusServiceUnavailable)
			return
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: gateway.endpoint()}}); err != nil {
			return
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	if err := client.Start(context.Background()); err == nil {
		t.Fatal("bootstrap failure returned nil")
	}

	second := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready after retry")
	client.Close()
	if err, ok := waitLifecycleResult(second); !ok {
		t.Error("successful retry did not return after Close")
	} else if err != nil {
		t.Errorf("successful retry returned %v after Close, want nil", err)
	}

	third := startLifecycleClient(client, context.Background())
	if err, ok := waitLifecycleResult(third); !ok {
		t.Error("Start after a successful run did not immediately return a lifecycle error")
		client.Close()
	} else if err == nil {
		t.Error("Start after a successful run returned nil, want lifecycle error")
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Errorf("terminal Client made %d bootstrap requests, want 2 total", got)
	}
}

func TestContextCancellationStopsReconnectInterval(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	requestNumbers := make(chan int32, 4)
	var requests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestNumber := atomic.AddInt32(&requests, 1)
		requestNumbers <- requestNumber
		endpoint := gateway.endpoint()
		if requestNumber > 1 {
			endpoint = "ws://127.0.0.1:1/unreachable"
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{
			Code: OK,
			Data: &Endpoint{
				Url: endpoint,
				ClientConfig: &ClientConfig{
					ReconnectCount:    3,
					ReconnectInterval: 30,
					PingInterval:      3600,
				},
			},
		}); err != nil {
			return
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 1)
	reconnecting := make(chan struct{}, 1)
	retryError := make(chan error, 2)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnReconnecting(func() { reconnecting <- struct{}{} }),
		WithOnError(func(err error) { retryError <- err }),
	)
	ctx, cancel := context.WithCancel(context.Background())
	result := startLifecycleClient(client, ctx)
	waitLifecycleSignal(t, ready, "Ready before reconnect")

	var serverConn *websocket.Conn
	select {
	case serverConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for initial gateway connection")
	}
	closeGatewayConn(t, serverConn)
	waitLifecycleSignal(t, reconnecting, "Reconnecting")

	for {
		select {
		case requestNumber := <-requestNumbers:
			if requestNumber >= 2 {
				goto retryStarted
			}
		case <-time.After(lifecycleTestTimeout):
			t.Fatal("timed out waiting for reconnect bootstrap")
		}
	}

retryStarted:
	select {
	case <-retryError:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for failed reconnect attempt")
	}
	cancel()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after canceling reconnect wait")
	} else if !errors.Is(err, context.Canceled) {
		t.Fatalf("Start returned %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Fatalf("canceled reconnect made %d bootstrap requests, want 2", got)
	}
}

func TestReconnectPublishesOneReplacementConnection(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var requests int32
	domain, cleanupBootstrap := installLifecycleBootstrapHandler(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requests, 1)
		if err := json.NewEncoder(w).Encode(&EndpointResp{
			Code: OK,
			Data: &Endpoint{
				Url: gateway.endpoint(),
				ClientConfig: &ClientConfig{
					ReconnectCount:    1,
					ReconnectInterval: 1,
					PingInterval:      3600,
				},
			},
		}); err != nil {
			return
		}
	}))
	defer cleanupBootstrap()

	ready := make(chan struct{}, 2)
	reconnecting := make(chan struct{}, 2)
	reconnected := make(chan struct{}, 2)
	disconnected := make(chan struct{}, 3)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnReconnecting(func() { reconnecting <- struct{}{} }),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "initial Ready")

	var firstServerConn *websocket.Conn
	select {
	case firstServerConn = <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for initial gateway connection")
	}
	closeGatewayConn(t, firstServerConn)
	waitLifecycleSignal(t, reconnecting, "Reconnecting")
	waitLifecycleSignal(t, reconnected, "Reconnected")
	waitLifecycleSignal(t, disconnected, "first Disconnected")

	select {
	case <-gateway.accepted:
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for replacement gateway connection")
	}
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return after closing reconnected Client")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
	waitLifecycleSignal(t, disconnected, "replacement Disconnected")
	select {
	case <-ready:
		t.Fatal("reconnect invoked Ready again")
	default:
	}
	select {
	case <-reconnecting:
		t.Fatal("reconnect invoked Reconnecting more than once")
	default:
	}
	select {
	case <-reconnected:
		t.Fatal("reconnect invoked Reconnected more than once")
	default:
	}
	if got := atomic.LoadInt32(&requests); got != 2 {
		t.Fatalf("reconnect made %d bootstrap requests, want 2", got)
	}
}

func TestLifecycleSettersRaceWithClose(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	var requests int32
	domain, cleanupBootstrap := installLifecycleBootstrap(t, gateway.endpoint(), &requests)
	defer cleanupBootstrap()
	ready := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(domain),
		WithAutoReconnect(false),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	result := startLifecycleClient(client, context.Background())
	waitLifecycleSignal(t, ready, "Ready")

	const workers = 8
	const iterations = 100
	var setters sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		setters.Add(1)
		go func() {
			defer setters.Done()
			for i := 0; i < iterations; i++ {
				client.SetOnReady(func() {})
				client.SetOnError(func(error) {})
				client.SetOnReconnecting(func() {})
				client.SetOnReconnected(func() {})
				client.SetOnDisconnected(func() {})
			}
		}()
	}
	client.Close()
	setters.Wait()
	if err, ok := waitLifecycleResult(result); !ok {
		t.Fatal("Start did not return while callback setters raced with Close")
	} else if err != nil {
		t.Fatalf("Start returned %v after Close, want nil", err)
	}
}
