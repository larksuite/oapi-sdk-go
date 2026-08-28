package channel

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	websocket "github.com/gorilla/websocket"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

const channelLifecycleTimeout = 2 * time.Second

type channelLifecycleFixture struct {
	bootstrap *httptest.Server
	gateway   *httptest.Server
	accepted  chan *websocket.Conn

	mu    sync.Mutex
	conns []*websocket.Conn
}

func newChannelLifecycleFixture(t *testing.T) (*channelLifecycleFixture, *larkws.Client, func()) {
	t.Helper()

	fixture := &channelLifecycleFixture{accepted: make(chan *websocket.Conn, 4)}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	fixture.gateway = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		fixture.mu.Lock()
		fixture.conns = append(fixture.conns, conn)
		fixture.mu.Unlock()
		fixture.accepted <- conn
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	endpoint := "ws" + strings.TrimPrefix(fixture.gateway.URL, "http") + "/channel"
	fixture.bootstrap = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&larkws.EndpointResp{
			Code: larkws.OK,
			Data: &larkws.Endpoint{
				Url: endpoint,
				ClientConfig: &larkws.ClientConfig{
					ReconnectCount:    0,
					ReconnectInterval: 1,
					PingInterval:      3600,
				},
			},
		}); err != nil {
			return
		}
	}))
	cleanup := func() {
		fixture.mu.Lock()
		conns := append([]*websocket.Conn(nil), fixture.conns...)
		fixture.mu.Unlock()
		for _, conn := range conns {
			if err := conn.Close(); err != nil && !isClosedChannelLifecycleConnectionError(err) {
				t.Errorf("close channel test connection: %v", err)
			}
		}
		fixture.bootstrap.Close()
		fixture.gateway.Close()
	}

	wsClient := larkws.NewClient("app-id", "app-secret",
		larkws.WithDomain(fixture.bootstrap.URL),
		larkws.WithAutoReconnect(false),
	)
	return fixture, wsClient, cleanup
}

func isClosedChannelLifecycleConnectionError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "use of closed network connection") ||
		strings.Contains(message, "closed network connection") ||
		strings.Contains(message, "connection already closed")
}

func newChannelLifecycleAPIClient() *lark.Client {
	mockHTTP := &MockHttpClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			header := make(http.Header)
			header.Set("Content-Type", "application/json")
			body := `{"code":0,"msg":"success"}`
			switch req.URL.Path {
			case "/open-apis/auth/v3/tenant_access_token/internal":
				body = `{"code":0,"msg":"success","tenant_access_token":"test-token","expire":7200}`
			case "/open-apis/bot/v3/info":
				body = `{"code":0,"msg":"success","bot":{"open_id":"test-bot","app_name":"Test Bot","activate_status":2}}`
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     header,
			}, nil
		},
	}
	return lark.NewClient("app-id", "app-secret", lark.WithHttpClient(mockHTTP))
}

func startChannelLifecycle(ch interface {
	Start(context.Context) error
}, ctx context.Context) <-chan error {
	result := make(chan error, 1)
	go func() {
		result <- ch.Start(ctx)
	}()
	return result
}

func waitChannelLifecycleSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(channelLifecycleTimeout):
		t.Fatalf("timed out waiting for %s", name)
	}
}

func waitChannelLifecycleResult(result <-chan error) (error, bool) {
	select {
	case err := <-result:
		return err, true
	case <-time.After(channelLifecycleTimeout):
		return nil, false
	}
}

func TestChannelStartReturnsContextCancellationAndBridgesCallbacksOnce(t *testing.T) {
	_, wsClient, cleanupFixture := newChannelLifecycleFixture(t)
	defer cleanupFixture()
	ch := NewChannel(newChannelLifecycleAPIClient(), wsClient)

	ready := make(chan struct{}, 2)
	disconnected := make(chan struct{}, 2)
	var errorsSeen int32
	ch.OnReady(func() { ready <- struct{}{} })
	ch.OnDisconnected(func() { disconnected <- struct{}{} })
	ch.OnError(func(error) { atomic.AddInt32(&errorsSeen, 1) })

	ctx, cancel := context.WithCancel(context.Background())
	result := startChannelLifecycle(ch, ctx)
	waitChannelLifecycleSignal(t, ready, "Channel Ready")
	cancel()

	err, ok := waitChannelLifecycleResult(result)
	if !ok {
		t.Error("Channel.Start did not return after context cancellation")
		if stopErr := ch.Stop(context.Background()); stopErr != nil {
			t.Errorf("Channel.Stop cleanup returned %v", stopErr)
		}
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("Channel.Start returned %v, want context.Canceled", err)
	}
	waitChannelLifecycleSignal(t, disconnected, "Channel Disconnected")
	select {
	case <-ready:
		t.Error("Channel bridged Ready more than once")
	default:
	}
	select {
	case <-disconnected:
		t.Error("Channel bridged Disconnected more than once")
	default:
	}
	if got := atomic.LoadInt32(&errorsSeen); got != 0 {
		t.Errorf("normal context cancellation bridged OnError %d times, want 0", got)
	}
}

func TestChannelStopEndsBlockingStartWithNil(t *testing.T) {
	_, wsClient, cleanupFixture := newChannelLifecycleFixture(t)
	defer cleanupFixture()
	ch := NewChannel(newChannelLifecycleAPIClient(), wsClient)

	ready := make(chan struct{}, 1)
	disconnected := make(chan struct{}, 1)
	ch.OnReady(func() { ready <- struct{}{} })
	ch.OnDisconnected(func() { disconnected <- struct{}{} })

	result := startChannelLifecycle(ch, context.Background())
	waitChannelLifecycleSignal(t, ready, "Channel Ready")
	if err := ch.Stop(context.Background()); err != nil {
		t.Fatalf("Channel.Stop returned %v", err)
	}
	if err, ok := waitChannelLifecycleResult(result); !ok {
		t.Error("Channel.Start did not return after Channel.Stop")
	} else if err != nil {
		t.Errorf("Channel.Start returned %v after Channel.Stop, want nil", err)
	}
	waitChannelLifecycleSignal(t, disconnected, "Channel Disconnected")
}

func TestChannelStartAndStopKeepNilWebSocketClientNoOp(t *testing.T) {
	ch := NewChannel(newChannelLifecycleAPIClient(), nil)
	if err := ch.Start(context.Background()); err != nil {
		t.Errorf("Channel.Start with nil WebSocket Client returned %v", err)
	}
	if err := ch.Stop(context.Background()); err != nil {
		t.Errorf("Channel.Stop with nil WebSocket Client returned %v", err)
	}
}
