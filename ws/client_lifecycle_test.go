package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gws "github.com/gorilla/websocket"
)

type lifecycleTestServer struct {
	t               *testing.T
	server          *httptest.Server
	bootstrap       func(*lifecycleTestServer, int, http.ResponseWriter, *http.Request)
	connections     chan *gws.Conn
	bootstrapCalls  int32
	connectionsMu   sync.Mutex
	openConnections []*gws.Conn
}

type blockingConnectedLogger struct {
	connected chan struct{}
	release   chan struct{}
	once      sync.Once
}

func (l *blockingConnectedLogger) Debug(context.Context, ...interface{}) {}
func (l *blockingConnectedLogger) Warn(context.Context, ...interface{})  {}
func (l *blockingConnectedLogger) Error(context.Context, ...interface{}) {}

func (l *blockingConnectedLogger) Info(_ context.Context, args ...interface{}) {
	if len(args) == 0 {
		return
	}
	message, ok := args[0].(string)
	if !ok || !strings.HasPrefix(message, "connected to ") {
		return
	}
	l.once.Do(func() { close(l.connected) })
	<-l.release
}

func newLifecycleTestServer(
	t *testing.T,
	bootstrap func(*lifecycleTestServer, int, http.ResponseWriter, *http.Request),
) *lifecycleTestServer {
	t.Helper()
	s := &lifecycleTestServer{
		t:           t,
		bootstrap:   bootstrap,
		connections: make(chan *gws.Conn, 8),
	}
	s.server = httptest.NewServer(http.HandlerFunc(s.handle))
	return s
}

func (s *lifecycleTestServer) handle(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case GenEndpointUri:
		call := int(atomic.AddInt32(&s.bootstrapCalls, 1))
		if s.bootstrap != nil {
			s.bootstrap(s, call, w, r)
			return
		}
		s.writeEndpoint(w)
	case "/ws":
		conn, err := (&gws.Upgrader{}).Upgrade(w, r, nil)
		if err != nil {
			s.t.Errorf("upgrade websocket: %v", err)
			return
		}
		s.connectionsMu.Lock()
		s.openConnections = append(s.openConnections, conn)
		s.connectionsMu.Unlock()
		s.connections <- conn
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	default:
		http.NotFound(w, r)
	}
}

func (s *lifecycleTestServer) writeEndpoint(w http.ResponseWriter) {
	endpointURL := "ws" + strings.TrimPrefix(s.server.URL, "http") + "/ws?" + DeviceID + "=device&" + ServiceID + "=1"
	_ = json.NewEncoder(w).Encode(&EndpointResp{
		Code: OK,
		Data: &Endpoint{
			Url: endpointURL,
			ClientConfig: &ClientConfig{
				ReconnectCount:    -1,
				ReconnectInterval: 0,
				ReconnectNonce:    0,
				PingInterval:      60,
			},
		},
	})
}

func (s *lifecycleTestServer) close() {
	s.connectionsMu.Lock()
	connections := append([]*gws.Conn(nil), s.openConnections...)
	s.connectionsMu.Unlock()
	for _, conn := range connections {
		_ = conn.Close()
	}
	s.server.Close()
}

func useLifecycleBootstrapClient(server *httptest.Server) func() {
	originalClient := bootstrapHTTPClient
	bootstrapHTTPClient = server.Client()
	return func() { bootstrapHTTPClient = originalClient }
}

func startClient(client *Client, ctx context.Context) <-chan error {
	done := make(chan error, 1)
	go func() {
		done <- client.Start(ctx)
	}()
	return done
}

func waitForSignal(t *testing.T, signal <-chan struct{}, name string) {
	t.Helper()
	select {
	case <-signal:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for %s", name)
	}
}

func waitForConnection(t *testing.T, connections <-chan *gws.Conn) *gws.Conn {
	t.Helper()
	select {
	case conn := <-connections:
		return conn
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for websocket connection")
		return nil
	}
}

func waitForStartResult(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for Client.Start to return")
		return nil
	}
}

func TestStartReturnsOnContextCancelAfterReady(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	var readyOnce sync.Once
	var disconnected int32
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
		WithOnDisconnected(func() { atomic.AddInt32(&disconnected, 1) }),
	)
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := startClient(client, ctx)
	waitForSignal(t, ready, "ready callback")
	waitForConnection(t, server.connections)
	cancel()

	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&disconnected); got != 1 {
		t.Fatalf("OnDisconnected calls = %d, want 1", got)
	}
	if got := atomic.LoadInt32(&server.bootstrapCalls); got != 1 {
		t.Fatalf("bootstrap calls = %d, want 1", got)
	}
}

func TestConcurrentCloseIsIdempotentAndReentrant(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	var readyOnce sync.Once
	var disconnected int32
	var client *Client
	client = NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
		WithOnDisconnected(func() {
			atomic.AddInt32(&disconnected, 1)
			client.Close()
		}),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, ready, "ready callback")
	waitForConnection(t, server.connections)

	var closeWG sync.WaitGroup
	for i := 0; i < 32; i++ {
		closeWG.Add(1)
		go func() {
			defer closeWG.Done()
			client.Close()
		}()
	}
	closed := make(chan struct{})
	go func() {
		closeWG.Wait()
		close(closed)
	}()
	waitForSignal(t, closed, "concurrent Close calls")

	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&disconnected); got != 1 {
		t.Fatalf("OnDisconnected calls = %d, want 1", got)
	}
}

func TestCancelInterruptsInitialReconnectBackoff(t *testing.T) {
	secondAttempt := make(chan struct{})
	var secondAttemptOnce sync.Once
	server := newLifecycleTestServer(t, func(_ *lifecycleTestServer, call int, w http.ResponseWriter, _ *http.Request) {
		if call >= 2 {
			secondAttemptOnce.Do(func() { close(secondAttempt) })
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: SystemBusy, Msg: "system busy"})
	})
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	client := NewClient("app-id", "app-secret", WithDomain(server.server.URL))
	client.configure(&ClientConfig{
		ReconnectCount:    -1,
		ReconnectInterval: 3600,
		ReconnectNonce:    0,
		PingInterval:      60,
	})
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := startClient(client, ctx)
	waitForSignal(t, secondAttempt, "second bootstrap attempt")
	cancel()

	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
	calls := atomic.LoadInt32(&server.bootstrapCalls)
	time.Sleep(50 * time.Millisecond)
	if got := atomic.LoadInt32(&server.bootstrapCalls); got != calls {
		t.Fatalf("bootstrap calls continued after cancellation: before=%d after=%d", calls, got)
	}
}

func TestCloseCancelsBlockedBootstrap(t *testing.T) {
	requestStarted := make(chan struct{})
	releaseRequest := make(chan struct{})
	var requestStartedOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStartedOnce.Do(func() { close(requestStarted) })
		select {
		case <-r.Context().Done():
		case <-releaseRequest:
		}
	}))
	defer server.Close()
	defer close(releaseRequest)
	defer useLifecycleBootstrapClient(server)()

	client := NewClient("app-id", "app-secret", WithDomain(server.URL))
	defer client.Close()
	done := startClient(client, context.Background())
	waitForSignal(t, requestStarted, "bootstrap request")
	client.Close()

	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
}

func TestRuntimeTerminalReconnectErrorReturnsFromStart(t *testing.T) {
	server := newLifecycleTestServer(t, func(s *lifecycleTestServer, call int, w http.ResponseWriter, _ *http.Request) {
		if call == 1 {
			s.writeEndpoint(w)
			return
		}
		_ = json.NewEncoder(w).Encode(&EndpointResp{Code: Forbidden, Msg: "app is forbidden"})
	})
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	var readyOnce sync.Once
	onError := make(chan error, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
		WithOnError(func(err error) { onError <- err }),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, ready, "ready callback")
	firstConn := waitForConnection(t, server.connections)
	_ = firstConn.Close()

	startErr := waitForStartResult(t, done)
	var clientErr *ClientError
	if !errors.As(startErr, &clientErr) {
		t.Fatalf("Start error type = %T, want *ClientError", startErr)
	}
	if clientErr.Code != Forbidden {
		t.Fatalf("ClientError code = %d, want %d", clientErr.Code, Forbidden)
	}
	select {
	case callbackErr := <-onError:
		if callbackErr != startErr {
			t.Fatalf("OnError error = %v, want the terminal Start error %v", callbackErr, startErr)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for OnError callback")
	}
	if got := atomic.LoadInt32(&server.bootstrapCalls); got != 2 {
		t.Fatalf("bootstrap calls = %d, want 2", got)
	}
}

func TestRuntimeReconnectSuccessKeepsStartRunning(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	reconnected := make(chan struct{})
	var readyOnce sync.Once
	var reconnectedOnce sync.Once
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
		WithOnReconnected(func() { reconnectedOnce.Do(func() { close(reconnected) }) }),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, ready, "ready callback")
	firstConn := waitForConnection(t, server.connections)
	_ = firstConn.Close()
	waitForSignal(t, reconnected, "reconnected callback")
	waitForConnection(t, server.connections)

	select {
	case err := <-done:
		t.Fatalf("Start returned after a successful reconnect: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	client.Close()
	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
}

func TestConcurrentStartRejectedAndSequentialRestartAllowed(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{}, 2)
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { ready <- struct{}{} }),
	)
	defer client.Close()

	firstDone := startClient(client, context.Background())
	waitForSignal(t, ready, "first ready callback")
	waitForConnection(t, server.connections)

	if err := client.Start(context.Background()); !errors.Is(err, ErrClientAlreadyStarted) {
		t.Fatalf("concurrent Start error = %v, want ErrClientAlreadyStarted", err)
	}
	if got := atomic.LoadInt32(&server.bootstrapCalls); got != 1 {
		t.Fatalf("bootstrap calls after rejected Start = %d, want 1", got)
	}

	client.Close()
	if err := waitForStartResult(t, firstDone); !errors.Is(err, context.Canceled) {
		t.Fatalf("first Start error = %v, want context.Canceled", err)
	}

	secondDone := startClient(client, context.Background())
	waitForSignal(t, ready, "second ready callback")
	waitForConnection(t, server.connections)
	select {
	case err := <-secondDone:
		t.Fatalf("second Start returned before Close: %v", err)
	case <-time.After(50 * time.Millisecond):
	}

	client.Close()
	if err := waitForStartResult(t, secondDone); !errors.Is(err, context.Canceled) {
		t.Fatalf("second Start error = %v, want context.Canceled", err)
	}
}

func TestCloseBeforeReadyCallbackSuppressesReady(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	logger := &blockingConnectedLogger{
		connected: make(chan struct{}),
		release:   make(chan struct{}),
	}
	var ready int32
	var disconnected int32
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithLogger(logger),
		WithOnReady(func() { atomic.AddInt32(&ready, 1) }),
		WithOnDisconnected(func() { atomic.AddInt32(&disconnected, 1) }),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, logger.connected, "connected log")
	client.Close()
	close(logger.release)

	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
	if got := atomic.LoadInt32(&ready); got != 0 {
		t.Fatalf("OnReady calls = %d, want 0", got)
	}
	if got := atomic.LoadInt32(&disconnected); got != 1 {
		t.Fatalf("OnDisconnected calls = %d, want 1", got)
	}
}

func TestOldConnectionCannotWriteAfterReconnect(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	reconnected := make(chan struct{})
	var readyOnce sync.Once
	var reconnectedOnce sync.Once
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
		WithOnReconnected(func() { reconnectedOnce.Do(func() { close(reconnected) }) }),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, ready, "ready callback")
	firstServerConn := waitForConnection(t, server.connections)
	run := client.currentRun()
	firstClientConn, _ := client.connStateForRun(run)
	_ = firstServerConn.Close()
	waitForSignal(t, reconnected, "reconnected callback")
	waitForConnection(t, server.connections)

	if err := client.writeMessageForRun(run, firstClientConn, gws.BinaryMessage, []byte("stale")); err == nil {
		t.Fatal("old connection write unexpectedly succeeded after reconnect")
	}
	if client.configureForRun(run, firstClientConn, &ClientConfig{ReconnectCount: 99}) {
		t.Fatal("old connection config unexpectedly applied after reconnect")
	}
	if reconnectCount, _, _, _ := client.configSnapshot(); reconnectCount != -1 {
		t.Fatalf("reconnect count = %d after stale config, want -1", reconnectCount)
	}

	client.Close()
	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
}

func TestCloseDoesNotWaitForWriterLock(t *testing.T) {
	server := newLifecycleTestServer(t, nil)
	defer server.close()
	defer useLifecycleBootstrapClient(server.server)()

	ready := make(chan struct{})
	var readyOnce sync.Once
	client := NewClient("app-id", "app-secret",
		WithDomain(server.server.URL),
		WithOnReady(func() { readyOnce.Do(func() { close(ready) }) }),
	)
	defer client.Close()

	done := startClient(client, context.Background())
	waitForSignal(t, ready, "ready callback")
	waitForConnection(t, server.connections)
	run := client.currentRun()
	conn, _ := client.connStateForRun(run)

	client.writeMu.Lock()
	writeDone := make(chan error, 1)
	go func() {
		writeDone <- client.writeMessageForRun(run, conn, gws.BinaryMessage, []byte("blocked"))
	}()
	closeDone := make(chan struct{})
	go func() {
		client.Close()
		close(closeDone)
	}()
	waitForSignal(t, closeDone, "Close while writer is waiting")
	client.writeMu.Unlock()

	select {
	case err := <-writeDone:
		if err == nil {
			t.Fatal("write unexpectedly succeeded after Close")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("writer did not unblock after Close")
	}
	if err := waitForStartResult(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Start error = %v, want context.Canceled", err)
	}
}
