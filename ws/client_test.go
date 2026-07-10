package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gorillaws "github.com/gorilla/websocket"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type wsMockClientAssertionProvider struct {
	mu     sync.Mutex
	tokens []*larkcore.Token
	index  int
	calls  int
	auds   []string
}

func (p *wsMockClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	p.auds = append(p.auds, aud)
	token := p.tokens[p.index]
	if p.index < len(p.tokens)-1 {
		p.index++
	}
	return token, nil
}

func TestGetConnURLWithAppSecret(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != GenEndpointUri {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.AppSecret == "" || req.ClientAssertion != "" {
			t.Fatalf("unexpected bootstrap request: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret", WithDomain(server.URL))
	connURL, err := client.getConnURL(context.Background())
	if err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
	if connURL != "wss://example.com/ws" {
		t.Fatalf("unexpected conn url: %s", connURL)
	}
}

func TestGetConnURLWithCustomHeadersAndUserAgent(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Tt-Env") != "ppe_test" {
			t.Fatalf("missing X-Tt-Env header: %s", r.Header.Get("X-Tt-Env"))
		}
		if r.Header.Get("X-Use-Ppe") != "1" {
			t.Fatalf("missing X-Use-Ppe header: %s", r.Header.Get("X-Use-Ppe"))
		}
		userAgent := r.Header.Get("User-Agent")
		if !strings.HasPrefix(userAgent, "oapi-sdk-go/") {
			t.Fatalf("unexpected user agent: %s", userAgent)
		}
		if !strings.Contains(userAgent, " source/ws-test") {
			t.Fatalf("missing source in user agent: %s", userAgent)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	headers := make(http.Header)
	headers.Set("X-Tt-Env", "ppe_test")
	headers.Set("X-Use-Ppe", "1")

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHeaders(headers),
		WithSource("ws-test"),
	)
	if _, err := client.getConnURL(context.Background()); err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
}

func TestGetConnURLWithChannelTag(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.ChannelTag != "web_boe_channel" {
			t.Fatalf("unexpected channel tag: %s", req.ChannelTag)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret", WithDomain(server.URL), WithChannelTag("web_boe_channel"))
	if _, err := client.getConnURL(context.Background()); err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
}

func TestPingFrameCarriesChannelTagOnWire(t *testing.T) {
	frame := newPingFrameWithChannelTag(42, "web_boe_channel")
	encoded, err := frame.Marshal()
	if err != nil {
		t.Fatalf("marshal ping frame: %v", err)
	}

	decoded := &Frame{}
	if err := decoded.Unmarshal(encoded); err != nil {
		t.Fatalf("unmarshal ping frame: %v", err)
	}
	headers := Headers(decoded.Headers)
	if got := headers.GetString(HeaderType); got != string(MessageTypePing) {
		t.Fatalf("ping type = %q, want %q", got, MessageTypePing)
	}
	if got := headers.GetString(HeaderChannelTag); got != "web_boe_channel" {
		t.Fatalf("channel_tag = %q, want web_boe_channel", got)
	}
}

func TestAppMainPingFrameOmitsChannelTagOnWire(t *testing.T) {
	frame := NewPingFrame(42)
	encoded, err := frame.Marshal()
	if err != nil {
		t.Fatalf("marshal ping frame: %v", err)
	}

	decoded := &Frame{}
	if err := decoded.Unmarshal(encoded); err != nil {
		t.Fatalf("unmarshal ping frame: %v", err)
	}
	for _, header := range decoded.Headers {
		if header.Key == HeaderChannelTag {
			t.Fatalf("app main ping unexpectedly contains channel_tag header: %#v", header)
		}
	}
}

func TestPingLoopCarriesChannelTagAcrossIntervals(t *testing.T) {
	upgrader := gorillaws.Upgrader{}
	frames := make(chan *Frame, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for i := 0; i < 2; i++ {
			_, data, err := conn.ReadMessage()
			if err != nil {
				return
			}
			frame := &Frame{}
			if err := frame.Unmarshal(data); err != nil {
				return
			}
			frames <- frame
		}
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}
	client := NewClient("app-id", "app-secret", WithChannelTag("web_boe_channel"))
	client.mu.Lock()
	client.conn = conn
	client.serviceID = "42"
	client.pingInterval = 10 * time.Millisecond
	client.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go client.pingLoop(ctx)
	for i := 0; i < 2; i++ {
		select {
		case frame := <-frames:
			headers := Headers(frame.Headers)
			if got := headers.GetString(HeaderChannelTag); got != "web_boe_channel" {
				t.Fatalf("ping %d channel_tag = %q", i+1, got)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for ping %d", i+1)
		}
	}
	cancel()
	client.disconnect(context.Background())
}

func TestReadyAndReconnectedCallbacks(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	upgrader := gorillaws.Upgrader{}
	closeFirst := make(chan struct{})
	var connectionCount atomic.Int32
	var server *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc(GenEndpointUri, func(w http.ResponseWriter, r *http.Request) {
		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?device_id=12345&service_id=42"
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: wsURL}})
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		if connectionCount.Add(1) == 1 {
			<-closeFirst
			return
		}
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	})
	server = httptest.NewServer(mux)
	defer server.Close()
	bootstrapHTTPClient = server.Client()

	ready := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithOnReady(func() { ready <- struct{}{} }),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
	)
	client.reconnectNonce = 0
	client.reconnectInterval = time.Millisecond
	client.reconnectCount = 2

	ctx, cancel := context.WithCancel(context.Background())
	startDone := make(chan error, 1)
	go func() { startDone <- client.Start(ctx) }()

	select {
	case <-ready:
	case <-time.After(time.Second):
		t.Fatal("OnReady was not called")
	}
	close(closeFirst)
	select {
	case <-reconnected:
	case <-time.After(time.Second):
		t.Fatal("OnReconnected was not called")
	}
	if got := connectionCount.Load(); got != 2 {
		t.Fatalf("connection count = %d, want 2", got)
	}
	select {
	case <-ready:
		t.Fatal("OnReady called more than once")
	default:
	}

	cancel()
	select {
	case err := <-startDone:
		if err != context.Canceled {
			t.Fatalf("Start returned %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Start did not return after context cancellation")
	}
}

func TestGetConnURLWithClientAssertionProxy(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/proxy"+GenEndpointUri {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get(larkcore.HeaderXTargetService) != "open.feishu.cn" {
			t.Fatalf("unexpected target service: %s", r.Header.Get(larkcore.HeaderXTargetService))
		}
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.ClientAssertion != "assertion" || req.AppSecret != "" {
			t.Fatalf("unexpected bootstrap request: %#v", req)
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer proxyServer.Close()

	bootstrapHTTPClient = proxyServer.Client()
	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion", TargetInfo: &larkcore.TargetInfo{TargetService: proxyServer.URL, TargetPrefix: "/proxy"}}}}
	client := NewClient("app-id", "", WithDomain("https://open.feishu.cn"), WithClientAssertionProvider(provider))
	if _, err := client.getConnURL(context.Background()); err != nil {
		t.Fatalf("get conn url failed: %v", err)
	}
}

func TestGetConnURLWithClientAssertionProxyHTTPErrorMsg(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{
			Code: 20050,
			Msg:  "target service unavailable",
		})
	}))
	defer proxyServer.Close()

	bootstrapHTTPClient = proxyServer.Client()
	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion", TargetInfo: &larkcore.TargetInfo{TargetService: proxyServer.URL, TargetPrefix: "/proxy"}}}}
	client := NewClient("app-id", "", WithDomain("https://open.feishu.cn"), WithClientAssertionProvider(provider))
	_, err := client.getConnURL(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	serverErr, ok := err.(*ServerError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if serverErr.Code != http.StatusInternalServerError {
		t.Fatalf("unexpected server error code: %d", serverErr.Code)
	}
	if serverErr.Msg != "target service unavailable" {
		t.Fatalf("unexpected server error msg: %s", serverErr.Msg)
	}
}

func TestGetConnURLRetrieveTokenEachTime(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	bodyAssertions := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req BootstrapRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		bodyAssertions = append(bodyAssertions, req.ClientAssertion)
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "wss://example.com/ws"}})
	}))
	defer server.Close()

	bootstrapHTTPClient = server.Client()
	provider := &wsMockClientAssertionProvider{tokens: []*larkcore.Token{{Value: "assertion-1"}, {Value: "assertion-2"}}}
	client := NewClient("app-id", "", WithDomain(server.URL), WithClientAssertionProvider(provider))
	if _, err := client.getConnURL(context.Background()); err != nil {
		t.Fatalf("first get conn url failed: %v", err)
	}
	if _, err := client.getConnURL(context.Background()); err != nil {
		t.Fatalf("second get conn url failed: %v", err)
	}
	if provider.calls != 2 {
		t.Fatalf("expected provider called twice, got %d", provider.calls)
	}
	if len(bodyAssertions) != 2 || bodyAssertions[0] != "assertion-1" || bodyAssertions[1] != "assertion-2" {
		t.Fatalf("unexpected assertions: %#v", bodyAssertions)
	}
}

func TestAttachUser(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-apis/event/v1/connections/12345/bind_user" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer user-token" {
			t.Fatalf("unexpected authorization header: %s", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-Tt-Env") != "boe_sup_user_channel" {
			t.Fatalf("missing route header: %s", r.Header.Get("X-Tt-Env"))
		}
		userAgent := r.Header.Get("User-Agent")
		if !strings.HasPrefix(userAgent, "oapi-sdk-go/") || !strings.Contains(userAgent, " source/ws-test") {
			t.Fatalf("unexpected user agent: %s", userAgent)
		}
		var req AttachUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.ChannelTag != "web_boe_channel" {
			t.Fatalf("unexpected channel tag: %s", req.ChannelTag)
		}
		_ = json.NewEncoder(w).Encode(&AttachUserResp{Code: OK})
	}))
	defer server.Close()

	headers := make(http.Header)
	headers.Set("X-Tt-Env", "boe_sup_user_channel")

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHeaders(headers),
		WithSource("ws-test"),
		WithChannelTag("web_boe_channel"),
	)
	client.mu.Lock()
	client.connID = "12345"
	client.mu.Unlock()

	if err := client.AttachUser(context.Background(), "user-token"); err != nil {
		t.Fatalf("attach user failed: %v", err)
	}
}

func TestDetachUser(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open-apis/event/v1/connections/12345/unbind_user" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer user-token" {
			t.Fatalf("missing user token")
		}
		if r.Header.Get("X-Tt-Env") != "boe_sup_user_channel" {
			t.Fatalf("missing route header")
		}
		userAgent := r.Header.Get("User-Agent")
		if !strings.HasPrefix(userAgent, "oapi-sdk-go/") || !strings.Contains(userAgent, " source/ws-test") {
			t.Fatalf("unexpected user agent: %s", userAgent)
		}
		var req AttachUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request failed: %v", err)
		}
		if req.ChannelTag != "web_boe_channel" {
			t.Fatalf("unexpected channel tag: %s", req.ChannelTag)
		}
		_ = json.NewEncoder(w).Encode(&AttachUserResp{Code: OK})
	}))
	defer server.Close()

	headers := make(http.Header)
	headers.Set("X-Tt-Env", "boe_sup_user_channel")

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHeaders(headers),
		WithSource("ws-test"),
		WithChannelTag("web_boe_channel"),
	)
	client.mu.Lock()
	client.connID = "12345"
	client.mu.Unlock()

	if err := client.DetachUser(context.Background(), "user-token"); err != nil {
		t.Fatalf("detach user failed: %v", err)
	}
}

func TestAttachDetachUseSnapshotWhileClosing(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	requestStarted := make(chan struct{}, 2)
	releaseRequests := make(chan struct{})
	bindingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req AttachUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		if req.ChannelTag != "web_boe_channel" {
			t.Errorf("channel_tag = %q", req.ChannelTag)
			return
		}
		requestStarted <- struct{}{}
		<-releaseRequests
		_ = json.NewEncoder(w).Encode(&AttachUserResp{Code: OK})
	}))
	defer bindingServer.Close()
	bootstrapHTTPClient = bindingServer.Client()

	upgrader := gorillaws.Upgrader{}
	webSocketServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer webSocketServer.Close()
	wsURL := "ws" + strings.TrimPrefix(webSocketServer.URL, "http")
	conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial websocket: %v", err)
	}

	client := NewClient("app-id", "app-secret",
		WithDomain(bindingServer.URL),
		WithChannelTag("web_boe_channel"),
	)
	client.mu.Lock()
	client.conn = conn
	client.connID = "12345"
	client.mu.Unlock()

	errs := make(chan error, 2)
	go func() { errs <- client.AttachUser(context.Background(), "user-token-a") }()
	go func() { errs <- client.DetachUser(context.Background(), "user-token-b") }()
	for i := 0; i < 2; i++ {
		select {
		case <-requestStarted:
		case <-time.After(time.Second):
			t.Fatalf("request %d did not reach server", i+1)
		}
	}

	closeDone := make(chan struct{})
	go func() {
		client.Close()
		close(closeDone)
	}()
	select {
	case <-closeDone:
	case <-time.After(time.Second):
		t.Fatal("Close blocked behind AttachUser/DetachUser HTTP requests")
	}
	close(releaseRequests)
	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			t.Fatalf("binding request %d failed: %v", i+1, err)
		}
	}
	if got := client.ConnectionID(); got != "" {
		t.Fatalf("connection id after close = %q", got)
	}
}

func TestAttachUserRequiresConnection(t *testing.T) {
	client := NewClient("app-id", "app-secret")
	err := client.AttachUser(context.Background(), "user-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if clientErr.Code != http.StatusBadRequest || clientErr.Msg != "connection is not ready" {
		t.Fatalf("unexpected error: %#v", clientErr)
	}
}

func TestAttachUserRequiresToken(t *testing.T) {
	client := NewClient("app-id", "app-secret")
	client.mu.Lock()
	client.connID = "12345"
	client.mu.Unlock()

	err := client.AttachUser(context.Background(), "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if clientErr.Code != http.StatusBadRequest || clientErr.Msg != "userAccessToken is required" {
		t.Fatalf("unexpected error: %#v", clientErr)
	}
}

func TestAttachUserOpenAPIError(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(&AttachUserResp{Code: 999, Msg: "user is already bound to another connection"})
	}))
	defer server.Close()

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret", WithDomain(server.URL))
	client.mu.Lock()
	client.connID = "12345"
	client.mu.Unlock()

	err := client.AttachUser(context.Background(), "user-token")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("unexpected error type: %#v", err)
	}
	if clientErr.Code != 999 || clientErr.Msg != "user is already bound to another connection" {
		t.Fatalf("unexpected error: %#v", clientErr)
	}
}

func TestAppendWebSocketQueryParams(t *testing.T) {
	params := url.Values{}
	params.Set("x-tt-env", "boe_sup_user_channel")

	connURL, err := appendWebSocketQueryParams("wss://example.com/ws?device_id=12345&service_id=67890", params)
	if err != nil {
		t.Fatalf("append websocket query params failed: %v", err)
	}
	query := connURL.Query()
	if query.Get(DeviceID) != "12345" || query.Get(ServiceID) != "67890" {
		t.Fatalf("lost connection query: %s", connURL.String())
	}
	if query.Get("x-tt-env") != "boe_sup_user_channel" {
		t.Fatalf("missing route query: %s", connURL.String())
	}
}

func TestSanitizeConnURLForLog(t *testing.T) {
	connURL, err := url.Parse("wss://example.com/ws?access_key=secret&ticket=ticket-value&device_id=12345&x-tt-env=boe_sup_user_channel")
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}

	sanitized := sanitizeConnURLForLog(connURL)
	if strings.Contains(sanitized, "secret") || strings.Contains(sanitized, "ticket-value") {
		t.Fatalf("sensitive query values leaked: %s", sanitized)
	}
	if !strings.Contains(sanitized, "access_key=%3Credacted%3E") || !strings.Contains(sanitized, "ticket=%3Credacted%3E") {
		t.Fatalf("sensitive query keys were not redacted: %s", sanitized)
	}
	if !strings.Contains(sanitized, "device_id=12345") || !strings.Contains(sanitized, "x-tt-env=boe_sup_user_channel") {
		t.Fatalf("non-sensitive query values were not preserved: %s", sanitized)
	}
}

func TestConnectionID(t *testing.T) {
	client := NewClient("app-id", "app-secret")
	client.mu.Lock()
	client.connID = "12345"
	client.mu.Unlock()

	if client.ConnectionID() != "12345" {
		t.Fatalf("unexpected connection id: %s", client.ConnectionID())
	}
}

func TestBuildWSProxyURL(t *testing.T) {
	testCases := map[string]string{
		"proxy.example.com":         "https://proxy.example.com/v1" + GenEndpointUri,
		"https://proxy.example.com": "https://proxy.example.com/v1" + GenEndpointUri,
		"http://proxy.example.com":  "http://proxy.example.com/v1" + GenEndpointUri,
	}

	for targetService, expected := range testCases {
		proxyURL := buildWSProxyURL(targetService, "/v1", GenEndpointUri)
		if proxyURL != expected {
			t.Fatalf("unexpected proxy url for %s: %s", targetService, proxyURL)
		}
	}
}
