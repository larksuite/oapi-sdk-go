package ws

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gorillaws "github.com/gorilla/websocket"
)

const transportTestTimeout = 5 * time.Second

type transportRequest struct {
	host        string
	escapedPath string
	rawQuery    string
	header      http.Header
}

type recordingLogger struct {
	mu      sync.Mutex
	entries []string
}

func (l *recordingLogger) append(args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, fmt.Sprint(args...))
}

func (l *recordingLogger) Debug(_ context.Context, args ...interface{}) {
	l.append(args...)
}

func (l *recordingLogger) Info(_ context.Context, args ...interface{}) {
	l.append(args...)
}

func (l *recordingLogger) Warn(_ context.Context, args ...interface{}) {
	l.append(args...)
}

func (l *recordingLogger) Error(_ context.Context, args ...interface{}) {
	l.append(args...)
}

func (l *recordingLogger) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(append([]string(nil), l.entries...), "\n")
}

func transportTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), transportTestTimeout)
}

func awaitTransportRequest(t *testing.T, ctx context.Context, ch <-chan transportRequest) transportRequest {
	t.Helper()
	select {
	case request := <-ch:
		return request
	case <-ctx.Done():
		t.Fatalf("timed out waiting for websocket request: %v", ctx.Err())
		return transportRequest{}
	}
}

func awaitSignal(t *testing.T, ctx context.Context, ch <-chan struct{}, description string) {
	t.Helper()
	select {
	case <-ch:
		return
	case <-ctx.Done():
		t.Fatalf("timed out waiting for %s: %v", description, ctx.Err())
	}
}

func writeEndpointResponse(t *testing.T, w http.ResponseWriter, endpoint string, config *ClientConfig) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(&EndpointResp{
		Code: OK,
		Data: &Endpoint{Url: endpoint, ClientConfig: config},
	}); err != nil {
		t.Errorf("encode endpoint response: %v", err)
	}
}

func cloneRequestHeader(header http.Header) http.Header {
	cloned := make(http.Header, len(header))
	for name, values := range header {
		cloned[name] = append([]string(nil), values...)
	}
	return cloned
}

func serveWebSocketUntilClientCloses(t *testing.T, w http.ResponseWriter, r *http.Request, requests chan<- transportRequest) {
	t.Helper()
	upgrader := gorillaws.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		t.Errorf("upgrade websocket: %v", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close websocket server connection: %v", err)
		}
	}()
	requests <- transportRequest{
		host:        r.Host,
		escapedPath: r.URL.EscapedPath(),
		rawQuery:    r.URL.RawQuery,
		header:      cloneRequestHeader(r.Header),
	}
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func assertClientError(t *testing.T, err error, code int, message string) *ClientError {
	t.Helper()
	if err == nil {
		t.Fatal("expected client error, got nil")
	}
	clientErr, ok := err.(*ClientError)
	if !ok {
		t.Fatalf("expected *ClientError, got %T: %v", err, err)
	}
	if clientErr.Code != code {
		t.Fatalf("unexpected client error code: got %d, want %d", clientErr.Code, code)
	}
	if clientErr.Msg != message {
		t.Fatalf("unexpected client error message: got %q, want %q", clientErr.Msg, message)
	}
	wantError := fmt.Sprintf("%d: %s", code, message)
	if clientErr.Error() != wantError {
		t.Fatalf("unexpected Error(): got %q, want %q", clientErr.Error(), wantError)
	}
	return clientErr
}

func assertNotContainsAny(t *testing.T, text string, forbidden ...string) {
	t.Helper()
	for _, value := range forbidden {
		if value != "" && strings.Contains(text, value) {
			t.Fatalf("value %q leaked in %q", value, text)
		}
	}
}

func TestBootstrapAndConnectionHeadersUseIndependentSnapshots(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	bootstrapRequests := make(chan http.Header, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bootstrapRequests <- cloneRequestHeader(r.Header)
		writeEndpointResponse(t, w, "ws://public.invalid/socket", nil)
	}))
	defer server.Close()

	bootstrapValues := []string{"bootstrap-original", "bootstrap-second"}
	connectionValues := []string{"connection-original", "connection-second"}
	bootstrapHeaders := http.Header{"X-Bootstrap-Phase": bootstrapValues}
	connectionHeaders := http.Header{"X-Connection-Phase": connectionValues}

	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHeaders(bootstrapHeaders),
		WithConnectionHeaders(connectionHeaders),
		WithAutoReconnect(false),
	)

	bootstrapValues[0] = "bootstrap-mutated-through-slice"
	bootstrapHeaders.Set("X-Bootstrap-Phase", "bootstrap-mutated-through-map")
	connectionValues[0] = "connection-mutated-through-slice"
	connectionHeaders.Set("X-Connection-Phase", "connection-mutated-through-map")

	ctx, cancel := transportTestContext(t)
	defer cancel()
	if _, err := client.getConnURL(ctx); err != nil {
		t.Fatalf("get connection URL: %v", err)
	}

	select {
	case header := <-bootstrapRequests:
		if got, want := header.Values("X-Bootstrap-Phase"), []string{"bootstrap-original", "bootstrap-second"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected bootstrap header snapshot: got %#v, want %#v", got, want)
		}
		if got := header.Values("X-Connection-Phase"); len(got) != 0 {
			t.Fatalf("connection headers leaked into bootstrap request: %#v", got)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for bootstrap request: %v", ctx.Err())
	}
}

func TestConnectAppliesConnectionHeadersAndHostOverride(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	bootstrapRequests := make(chan http.Header, 1)
	connectionRequests := make(chan transportRequest, 1)
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == GenEndpointUri {
			bootstrapRequests <- cloneRequestHeader(r.Header)
			writeEndpointResponse(t, w, "ws://public.example:9443/socket%2Fpath?device_id=device-value&service_id=42&opaque=%2Fraw", nil)
			return
		}
		serveWebSocketUntilClientCloses(t, w, r, connectionRequests)
	}))
	defer server.Close()

	bootstrapValues := []string{"bootstrap-original", "bootstrap-second"}
	connectionValues := []string{"connection-original", "connection-second"}
	bootstrapHeaders := http.Header{"X-Bootstrap-Phase": bootstrapValues}
	connectionHeaders := http.Header{
		"X-Connection-Phase": connectionValues,
		"Authorization":      []string{"Bearer connection-secret"},
		"Cookie":             []string{"session=connection-cookie"},
		"Origin":             []string{"https://origin.example"},
	}

	bootstrapHTTPClient = server.Client()
	overrideHost := server.Listener.Addr().String()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithHeaders(bootstrapHeaders),
		WithConnectionHeaders(connectionHeaders),
		WithConnectionHost(overrideHost),
		WithAutoReconnect(false),
		WithLogger(&recordingLogger{}),
	)

	bootstrapValues[0] = "bootstrap-mutated-through-slice"
	bootstrapHeaders.Set("X-Bootstrap-Phase", "bootstrap-mutated-through-map")
	connectionValues[0] = "connection-mutated-through-slice"
	connectionHeaders.Set("X-Connection-Phase", "connection-mutated-through-map")

	ctx, cancel := transportTestContext(t)
	defer cancel()
	if err := client.connect(ctx); err != nil {
		t.Fatalf("connect websocket: %v", err)
	}
	defer client.disconnect(ctx)

	select {
	case header := <-bootstrapRequests:
		if got, want := header.Values("X-Bootstrap-Phase"), []string{"bootstrap-original", "bootstrap-second"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("unexpected bootstrap header snapshot: got %#v, want %#v", got, want)
		}
		if got := header.Values("X-Connection-Phase"); len(got) != 0 {
			t.Fatalf("connection headers leaked into bootstrap request: %#v", got)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for bootstrap request: %v", ctx.Err())
	}

	request := awaitTransportRequest(t, ctx, connectionRequests)
	if request.host != overrideHost {
		t.Fatalf("unexpected websocket Host: got %q, want %q", request.host, overrideHost)
	}
	if request.escapedPath != "/socket%2Fpath" {
		t.Fatalf("unexpected escaped websocket path: %q", request.escapedPath)
	}
	if request.rawQuery != "device_id=device-value&service_id=42&opaque=%2Fraw" {
		t.Fatalf("unexpected raw websocket query: %q", request.rawQuery)
	}
	if got, want := request.header.Values("X-Connection-Phase"), []string{"connection-original", "connection-second"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected connection header snapshot: got %#v, want %#v", got, want)
	}
	if got := request.header.Values("X-Bootstrap-Phase"); len(got) != 0 {
		t.Fatalf("bootstrap headers leaked into websocket handshake: %#v", got)
	}
	if request.header.Get("Authorization") != "Bearer connection-secret" ||
		request.header.Get("Cookie") != "session=connection-cookie" ||
		request.header.Get("Origin") != "https://origin.example" {
		t.Fatalf("business headers were not forwarded: %#v", request.header)
	}
	if client.connUrl == nil || client.connUrl.Scheme != "ws" || client.connUrl.Host != overrideHost ||
		client.connUrl.EscapedPath() != "/socket%2Fpath" || client.connUrl.RawQuery != "device_id=device-value&service_id=42&opaque=%2Fraw" {
		t.Fatalf("unexpected final connection URL: %#v", client.connUrl)
	}
}

func TestReconnectReusesConnectionConfiguration(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	connectionRequests := make(chan transportRequest, 2)
	firstHandshakeClosed := make(chan struct{})
	var bootstrapCount int32
	var handshakeCount int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == GenEndpointUri {
			atomic.AddInt32(&bootstrapCount, 1)
			writeEndpointResponse(t, w, "ws://public.example/reconnect?device_id=reconnect-device&service_id=7", &ClientConfig{
				ReconnectCount:    1,
				ReconnectInterval: 0,
				ReconnectNonce:    0,
				PingInterval:      3600,
			})
			return
		}

		upgrader := gorillaws.Upgrader{CheckOrigin: func(_ *http.Request) bool { return true }}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade websocket: %v", err)
			return
		}
		requestNumber := atomic.AddInt32(&handshakeCount, 1)
		connectionRequests <- transportRequest{
			host:        r.Host,
			escapedPath: r.URL.EscapedPath(),
			rawQuery:    r.URL.RawQuery,
			header:      cloneRequestHeader(r.Header),
		}
		if requestNumber == 1 {
			if err := conn.Close(); err != nil {
				t.Errorf("close first websocket connection: %v", err)
			}
			close(firstHandshakeClosed)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				t.Errorf("close reconnected websocket connection: %v", err)
			}
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	connectionValues := []string{"fixed-connection-value", "fixed-second-value"}
	connectionHeaders := http.Header{"X-Reconnect-Header": connectionValues}
	disconnected := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	overrideHost := server.Listener.Addr().String()
	bootstrapHTTPClient = server.Client()
	client := NewClient("app-id", "app-secret",
		WithDomain(server.URL),
		WithConnectionHeaders(connectionHeaders),
		WithConnectionHost(overrideHost),
		WithAutoReconnect(false),
		WithOnDisconnected(func() { disconnected <- struct{}{} }),
		WithOnReconnected(func() { reconnected <- struct{}{} }),
		WithLogger(&recordingLogger{}),
	)
	connectionValues[0] = "mutated-through-slice"
	connectionHeaders.Set("X-Reconnect-Header", "mutated-through-map")

	ctx, cancel := transportTestContext(t)
	defer cancel()
	if err := client.connect(ctx); err != nil {
		t.Fatalf("initial websocket connect: %v", err)
	}
	awaitSignal(t, ctx, firstHandshakeClosed, "server to close initial connection")
	awaitSignal(t, ctx, disconnected, "initial disconnect callback")
	if err := client.reconnect(ctx); err != nil {
		t.Fatalf("reconnect websocket: %v", err)
	}
	defer client.disconnect(ctx)
	awaitSignal(t, ctx, reconnected, "reconnected callback")

	first := awaitTransportRequest(t, ctx, connectionRequests)
	second := awaitTransportRequest(t, ctx, connectionRequests)
	for index, request := range []transportRequest{first, second} {
		if request.host != overrideHost {
			t.Fatalf("handshake %d used unexpected Host: got %q, want %q", index+1, request.host, overrideHost)
		}
		if got, want := request.header.Values("X-Reconnect-Header"), []string{"fixed-connection-value", "fixed-second-value"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("handshake %d used unexpected header snapshot: got %#v, want %#v", index+1, got, want)
		}
		if request.rawQuery != "device_id=reconnect-device&service_id=7" {
			t.Fatalf("handshake %d changed raw query: %q", index+1, request.rawQuery)
		}
	}
	if got := atomic.LoadInt32(&bootstrapCount); got != 2 {
		t.Fatalf("unexpected bootstrap count: got %d, want 2", got)
	}
	if got := atomic.LoadInt32(&handshakeCount); got != 2 {
		t.Fatalf("unexpected handshake count: got %d, want 2", got)
	}
}

func TestReservedConnectionHeadersFailClosed(t *testing.T) {
	invalidHeaders := []struct {
		name     string
		header   http.Header
		sentinel string
	}{
		{name: "host", header: http.Header{"hOsT": []string{"blocked-host-value"}}, sentinel: "blocked-host-value"},
		{name: "connection", header: http.Header{"cOnNeCtIoN": []string{"blocked-connection-value"}}, sentinel: "blocked-connection-value"},
		{name: "upgrade", header: http.Header{"uPgRaDe": []string{"blocked-upgrade-value"}}, sentinel: "blocked-upgrade-value"},
		{name: "websocket key", header: http.Header{"sEc-WeBsOcKeT-kEy": []string{"blocked-key-value"}}, sentinel: "blocked-key-value"},
		{name: "websocket protocol", header: http.Header{"sec-websocket-protocol": []string{"blocked-protocol-value"}}, sentinel: "blocked-protocol-value"},
		{name: "custom websocket protocol header", header: http.Header{"SEC-WEBSOCKET-X-CUSTOM": []string{"blocked-custom-value"}}, sentinel: "blocked-custom-value"},
		{name: "invalid field name", header: http.Header{"Bad Header": []string{"blocked-invalid-name-value"}}, sentinel: "blocked-invalid-name-value"},
		{name: "carriage return", header: http.Header{"X-Bad-Value": []string{"blocked-cr-value\rnext"}}, sentinel: "blocked-cr-value"},
		{name: "line feed", header: http.Header{"X-Bad-Value": []string{"blocked-lf-value\nnext"}}, sentinel: "blocked-lf-value"},
		{name: "nul", header: http.Header{"X-Bad-Value": []string{"blocked-nul-value\x00next"}}, sentinel: "blocked-nul-value"},
		{name: "delete", header: http.Header{"X-Bad-Value": []string{"blocked-del-value\x7fnext"}}, sentinel: "blocked-del-value"},
		{name: "c0 control", header: http.Header{"X-Bad-Value": []string{"blocked-control-value\x01next"}}, sentinel: "blocked-control-value"},
	}

	for _, testCase := range invalidHeaders {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			originalClient := bootstrapHTTPClient
			defer func() { bootstrapHTTPClient = originalClient }()

			var bootstrapCount int32
			var handshakeCount int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == GenEndpointUri {
					atomic.AddInt32(&bootstrapCount, 1)
					writeEndpointResponse(t, w, "ws://public.example/blocked", nil)
					return
				}
				atomic.AddInt32(&handshakeCount, 1)
				w.Header().Set(HeaderHandshakeStatus, "503")
				w.Header().Set(HeaderHandshakeMsg, "rejected by test server")
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer server.Close()

			logger := &recordingLogger{}
			onError := make(chan error, 1)
			bootstrapHTTPClient = server.Client()
			client := NewClient("app-id", "app-secret",
				WithDomain(server.URL),
				WithConnectionHeaders(testCase.header),
				WithConnectionHost(server.Listener.Addr().String()),
				WithAutoReconnect(false),
				WithLogger(logger),
				WithOnError(func(err error) { onError <- err }),
			)

			ctx, cancel := transportTestContext(t)
			defer cancel()
			err := client.Start(ctx)
			assertClientError(t, err, 400, "invalid websocket connection headers")
			select {
			case callbackErr := <-onError:
				assertClientError(t, callbackErr, 400, "invalid websocket connection headers")
			case <-ctx.Done():
				t.Fatalf("timed out waiting for OnError: %v", ctx.Err())
			}
			if got := atomic.LoadInt32(&bootstrapCount); got != 0 {
				t.Fatalf("invalid header reached bootstrap endpoint %d times", got)
			}
			if got := atomic.LoadInt32(&handshakeCount); got != 0 {
				t.Fatalf("invalid header reached websocket endpoint %d times", got)
			}
			assertNotContainsAny(t, err.Error(), testCase.sentinel)
			assertNotContainsAny(t, logger.String(), testCase.sentinel)
		})
	}

	t.Run("last connection header option replaces invalid configuration", func(t *testing.T) {
		originalClient := bootstrapHTTPClient
		defer func() { bootstrapHTTPClient = originalClient }()

		requests := make(chan transportRequest, 2)
		var server *httptest.Server
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == GenEndpointUri {
				writeEndpointResponse(t, w, "ws://public.example/last-wins", nil)
				return
			}
			serveWebSocketUntilClientCloses(t, w, r, requests)
		}))
		defer server.Close()
		bootstrapHTTPClient = server.Client()

		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithConnectionHeaders(http.Header{"Host": []string{"must-be-replaced"}}),
			WithConnectionHeaders(http.Header{
				"Authorization": []string{"Bearer allowed"},
				"Cookie":        []string{"session=allowed"},
				"Origin":        []string{"https://origin.example"},
			}),
			WithConnectionHost(server.Listener.Addr().String()),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		if err := client.connect(ctx); err != nil {
			t.Fatalf("connect with last valid headers: %v", err)
		}
		request := awaitTransportRequest(t, ctx, requests)
		if request.header.Get("Authorization") != "Bearer allowed" || request.header.Get("Cookie") != "session=allowed" || request.header.Get("Origin") != "https://origin.example" {
			t.Fatalf("last valid connection headers were not applied: %#v", request.header)
		}
		client.disconnect(ctx)

		cleared := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithConnectionHeaders(http.Header{"Upgrade": []string{"must-be-cleared"}}),
			WithConnectionHeaders(nil),
			WithConnectionHost(server.Listener.Addr().String()),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		if err := cleared.connect(ctx); err != nil {
			t.Fatalf("connect after clearing headers with nil: %v", err)
		}
		clearedRequest := awaitTransportRequest(t, ctx, requests)
		if got := clearedRequest.header.Get("Upgrade"); !strings.EqualFold(got, "websocket") {
			t.Fatalf("websocket protocol Upgrade header was not generated: %q", got)
		}
		cleared.disconnect(ctx)
	})

	t.Run("last invalid option fails and header error has priority", func(t *testing.T) {
		logger := &recordingLogger{}
		client := NewClient("app-id", "app-secret",
			WithConnectionHeaders(http.Header{"Authorization": []string{"Bearer allowed"}}),
			WithConnectionHeaders(http.Header{"Upgrade": []string{"priority-header-secret"}}),
			WithConnectionHost(""),
			WithAutoReconnect(false),
			WithLogger(logger),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		err := client.connect(ctx)
		assertClientError(t, err, 400, "invalid websocket connection headers")
		assertNotContainsAny(t, err.Error(), "priority-header-secret")
		assertNotContainsAny(t, logger.String(), "priority-header-secret")
	})
}

func TestConnectionHostValidationAndRewrite(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	requests := make(chan transportRequest, 1)
	var bootstrapCount int32
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == GenEndpointUri {
			atomic.AddInt32(&bootstrapCount, 1)
			writeEndpointResponse(t, w, "ws://public.example:9443/fixed%2Fpath?device_id=host-device&service_id=11&opaque=%2Fraw", nil)
			return
		}
		serveWebSocketUntilClientCloses(t, w, r, requests)
	}))
	defer server.Close()
	bootstrapHTTPClient = server.Client()

	label63 := strings.Repeat("a", 63)
	host253 := strings.Join([]string{strings.Repeat("a", 63), strings.Repeat("b", 63), strings.Repeat("c", 63), strings.Repeat("d", 61)}, ".")
	validHosts := []string{
		"service",
		"123",
		"ws.internal.example",
		"WS.Internal",
		"xn--bcher-kva.internal",
		label63,
		host253,
		"0.0.0.0",
		"127.0.0.1",
		"255.255.255.255:65535",
		"service:1",
		"service:65535",
		"[2001:db8::1]",
		"[::ffff:192.0.2.1]:8443",
	}

	for _, configuredHost := range validHosts {
		configuredHost := configuredHost
		t.Run("valid_"+configuredHost, func(t *testing.T) {
			dialTargets := make(chan string, 1)
			client := NewClient("app-id", "app-secret",
				WithDomain(server.URL),
				WithConnectionHost(configuredHost),
				WithAutoReconnect(false),
				WithLogger(&recordingLogger{}),
			)
			client.dialer.Proxy = nil
			client.dialer.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
				dialTargets <- address
				dialer := &net.Dialer{}
				return dialer.DialContext(ctx, network, server.Listener.Addr().String())
			}

			ctx, cancel := transportTestContext(t)
			defer cancel()
			if err := client.connect(ctx); err != nil {
				t.Fatalf("connect using valid host %q: %v", configuredHost, err)
			}
			request := awaitTransportRequest(t, ctx, requests)
			if request.host != configuredHost {
				t.Fatalf("HTTP Host changed configured spelling: got %q, want %q", request.host, configuredHost)
			}
			if request.escapedPath != "/fixed%2Fpath" || request.rawQuery != "device_id=host-device&service_id=11&opaque=%2Fraw" {
				t.Fatalf("host rewrite changed path/query: path=%q query=%q", request.escapedPath, request.rawQuery)
			}
			if client.connUrl == nil || client.connUrl.Host != configuredHost {
				t.Fatalf("final URL did not preserve configured host: %#v", client.connUrl)
			}
			select {
			case target := <-dialTargets:
				wantTarget := configuredHost
				if _, _, err := net.SplitHostPort(configuredHost); err != nil {
					wantTarget = configuredHost + ":80"
				}
				if target != wantTarget {
					t.Fatalf("unexpected network target: got %q, want %q", target, wantTarget)
				}
			case <-ctx.Done():
				t.Fatalf("timed out waiting for network target: %v", ctx.Err())
			}
			client.disconnect(ctx)
		})
	}

	host254 := strings.Join([]string{strings.Repeat("a", 63), strings.Repeat("b", 63), strings.Repeat("c", 63), strings.Repeat("d", 62)}, ".")
	invalidHosts := []string{
		"",
		"ws://internal",
		"user@internal",
		"internal/path",
		"internal?x",
		"internal#x",
		" internal",
		"internal ",
		"internal\rnext",
		"internal\nnext",
		"internal%2f",
		".internal",
		"internal.",
		"a..b",
		"-a",
		"a-",
		"a_b",
		strings.Repeat("a", 64),
		host254,
		"bücher.internal",
		"1.2.3",
		"256.0.0.1",
		"01.2.3.4",
		"1.2.3.999",
		"2001:db8::1",
		"[2001:db8::1",
		"[[2001:db8::1]]",
		"[fe80::1%eth0]",
		"service:",
		"service:+1",
		"service:-1",
		"service:abc",
		"service:0",
		"service:01",
		"service:065535",
		"service:65536",
	}

	for _, configuredHost := range invalidHosts {
		configuredHost := configuredHost
		t.Run("invalid_"+fmt.Sprintf("%q", configuredHost), func(t *testing.T) {
			beforeBootstrap := atomic.LoadInt32(&bootstrapCount)
			var dialCount int32
			logger := &recordingLogger{}
			client := NewClient("app-id", "app-secret",
				WithDomain(server.URL),
				WithConnectionHost(configuredHost),
				WithAutoReconnect(false),
				WithLogger(logger),
			)
			client.dialer.Proxy = nil
			client.dialer.NetDialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
				atomic.AddInt32(&dialCount, 1)
				return nil, errors.New("unexpected dial")
			}
			ctx, cancel := transportTestContext(t)
			defer cancel()
			err := client.connect(ctx)
			assertClientError(t, err, 400, "invalid websocket connection host")
			assertNotContainsAny(t, err.Error(), configuredHost)
			assertNotContainsAny(t, logger.String(), configuredHost)
			if after := atomic.LoadInt32(&bootstrapCount); after != beforeBootstrap {
				t.Fatalf("invalid host reached bootstrap endpoint: before=%d after=%d", beforeBootstrap, after)
			}
			if got := atomic.LoadInt32(&dialCount); got != 0 {
				t.Fatalf("invalid host attempted %d websocket dials", got)
			}
		})
	}

	t.Run("last host option replaces prior invalid host", func(t *testing.T) {
		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithConnectionHost(""),
			WithConnectionHost(server.Listener.Addr().String()),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		if err := client.connect(ctx); err != nil {
			t.Fatalf("connect with last valid host: %v", err)
		}
		request := awaitTransportRequest(t, ctx, requests)
		if request.host != server.Listener.Addr().String() {
			t.Fatalf("last host option did not win: %q", request.host)
		}
		client.disconnect(ctx)

		invalidLast := NewClient("app-id", "app-secret",
			WithConnectionHost(server.Listener.Addr().String()),
			WithConnectionHost(""),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		assertClientError(t, invalidLast.connect(ctx), 400, "invalid websocket connection host")
	})
}

func TestConnectionEndpointValidation(t *testing.T) {
	t.Run("default ws endpoint is preserved byte for byte", func(t *testing.T) {
		originalClient := bootstrapHTTPClient
		defer func() { bootstrapHTTPClient = originalClient }()

		requests := make(chan transportRequest, 1)
		var endpoint string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == GenEndpointUri {
				writeEndpointResponse(t, w, endpoint, nil)
				return
			}
			serveWebSocketUntilClientCloses(t, w, r, requests)
		}))
		defer server.Close()
		endpoint = "ws" + strings.TrimPrefix(server.URL, "http") + "/escaped%2Fpath?device_id=default-device&service_id=21&opaque=%2Fraw"
		bootstrapHTTPClient = server.Client()
		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		if err := client.connect(ctx); err != nil {
			t.Fatalf("connect default ws endpoint: %v", err)
		}
		request := awaitTransportRequest(t, ctx, requests)
		if request.host != server.Listener.Addr().String() || request.escapedPath != "/escaped%2Fpath" || request.rawQuery != "device_id=default-device&service_id=21&opaque=%2Fraw" {
			t.Fatalf("default ws endpoint changed: %#v", request)
		}
		if client.connUrl == nil || client.connUrl.String() != endpoint {
			t.Fatalf("default ws endpoint serialization changed: got %v, want %q", client.connUrl, endpoint)
		}
		client.disconnect(ctx)
	})

	t.Run("default wss endpoint is preserved byte for byte", func(t *testing.T) {
		originalClient := bootstrapHTTPClient
		defer func() { bootstrapHTTPClient = originalClient }()

		requests := make(chan transportRequest, 1)
		var endpoint string
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == GenEndpointUri {
				writeEndpointResponse(t, w, endpoint, nil)
				return
			}
			serveWebSocketUntilClientCloses(t, w, r, requests)
		}))
		defer server.Close()
		endpoint = "wss" + strings.TrimPrefix(server.URL, "https") + "/secure%2Fpath?device_id=secure-device&service_id=22&opaque=%2Fraw"
		bootstrapHTTPClient = server.Client()
		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		transport, ok := server.Client().Transport.(*http.Transport)
		if !ok || transport.TLSClientConfig == nil {
			t.Fatalf("httptest TLS transport is unavailable: %T", server.Client().Transport)
		}
		client.dialer.Proxy = nil
		client.dialer.TLSClientConfig = transport.TLSClientConfig.Clone()

		ctx, cancel := transportTestContext(t)
		defer cancel()
		if err := client.connect(ctx); err != nil {
			t.Fatalf("connect default wss endpoint: %v", err)
		}
		request := awaitTransportRequest(t, ctx, requests)
		if request.host != server.Listener.Addr().String() || request.escapedPath != "/secure%2Fpath" || request.rawQuery != "device_id=secure-device&service_id=22&opaque=%2Fraw" {
			t.Fatalf("default wss endpoint changed: %#v", request)
		}
		if client.connUrl == nil || client.connUrl.String() != endpoint {
			t.Fatalf("default wss endpoint serialization changed: got %v, want %q", client.connUrl, endpoint)
		}
		client.disconnect(ctx)
	})

	invalidEndpoints := []struct {
		name     string
		endpoint string
	}{
		{name: "unsupported scheme", endpoint: "http://public.example/socket?raw=unsupported-scheme-secret"},
		{name: "empty hostname", endpoint: "ws:///socket?raw=empty-host-secret"},
		{name: "userinfo", endpoint: "wss://endpoint-user:endpoint-password@public.example/socket?raw=userinfo-secret"},
		{name: "fragment", endpoint: "ws://public.example/socket?raw=fragment-query-secret#fragment-secret"},
		{name: "empty port", endpoint: "ws://public.example:/socket?raw=empty-port-secret"},
		{name: "non-decimal port", endpoint: "ws://public.example:not-a-port/socket?raw=non-decimal-port-secret"},
		{name: "zero port", endpoint: "ws://public.example:0/socket?raw=zero-port-secret"},
		{name: "out of range port", endpoint: "ws://public.example:65536/socket?raw=range-port-secret"},
		{name: "unparseable", endpoint: "ws://%zz/socket?raw=parse-secret"},
	}

	for _, testCase := range invalidEndpoints {
		testCase := testCase
		for _, withOverride := range []bool{false, true} {
			withOverride := withOverride
			t.Run(fmt.Sprintf("%s_override_%t", testCase.name, withOverride), func(t *testing.T) {
				originalClient := bootstrapHTTPClient
				defer func() { bootstrapHTTPClient = originalClient }()

				var bootstrapCount int32
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					atomic.AddInt32(&bootstrapCount, 1)
					writeEndpointResponse(t, w, testCase.endpoint, nil)
				}))
				defer server.Close()
				bootstrapHTTPClient = server.Client()

				logger := &recordingLogger{}
				onError := make(chan error, 1)
				options := []ClientOption{
					WithDomain(server.URL),
					WithHeaders(http.Header{"X-Bootstrap-Secret": []string{"bootstrap-endpoint-secret"}}),
					WithConnectionHeaders(http.Header{"Authorization": []string{"Bearer connection-endpoint-secret"}}),
					WithAutoReconnect(false),
					WithLogger(logger),
					WithOnError(func(err error) { onError <- err }),
				}
				if withOverride {
					options = append(options, WithConnectionHost("127.0.0.1:65535"))
				}
				client := NewClient("app-id", "app-secret", options...)
				var dialCount int32
				client.dialer.Proxy = nil
				client.dialer.NetDialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
					atomic.AddInt32(&dialCount, 1)
					return nil, errors.New("unexpected endpoint dial")
				}

				ctx, cancel := transportTestContext(t)
				defer cancel()
				err := client.Start(ctx)
				assertClientError(t, err, 400, "invalid websocket endpoint")
				select {
				case callbackErr := <-onError:
					assertClientError(t, callbackErr, 400, "invalid websocket endpoint")
					assertNotContainsAny(t, callbackErr.Error(), testCase.endpoint, "bootstrap-endpoint-secret", "connection-endpoint-secret")
				case <-ctx.Done():
					t.Fatalf("timed out waiting for OnError: %v", ctx.Err())
				}
				if got := atomic.LoadInt32(&bootstrapCount); got != 1 {
					t.Fatalf("unexpected bootstrap count: got %d, want 1", got)
				}
				if got := atomic.LoadInt32(&dialCount); got != 0 {
					t.Fatalf("invalid endpoint attempted %d websocket dials", got)
				}
				assertNotContainsAny(t, err.Error(), testCase.endpoint, "bootstrap-endpoint-secret", "connection-endpoint-secret")
				assertNotContainsAny(t, logger.String(), testCase.endpoint, "endpoint-user", "endpoint-password", "unsupported-scheme-secret", "empty-host-secret", "userinfo-secret", "fragment-query-secret", "fragment-secret", "empty-port-secret", "non-decimal-port-secret", "zero-port-secret", "range-port-secret", "parse-secret", "bootstrap-endpoint-secret", "connection-endpoint-secret")
			})
		}
	}
}

func TestConnectWithoutOverrides(t *testing.T) {
	t.Run("uses the server endpoint without new options", func(t *testing.T) {
		originalClient := bootstrapHTTPClient
		defer func() { bootstrapHTTPClient = originalClient }()

		requests := make(chan transportRequest, 1)
		var endpoint string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == GenEndpointUri {
				writeEndpointResponse(t, w, endpoint, nil)
				return
			}
			serveWebSocketUntilClientCloses(t, w, r, requests)
		}))
		defer server.Close()
		endpoint = "ws" + strings.TrimPrefix(server.URL, "http") + "/legacy/path?device_id=legacy-device&service_id=31&opaque=legacy-query"
		bootstrapHTTPClient = server.Client()
		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		if err := client.connect(ctx); err != nil {
			t.Fatalf("connect without new options: %v", err)
		}
		request := awaitTransportRequest(t, ctx, requests)
		if request.host != server.Listener.Addr().String() || request.escapedPath != "/legacy/path" || request.rawQuery != "device_id=legacy-device&service_id=31&opaque=legacy-query" {
			t.Fatalf("default websocket request changed: %#v", request)
		}
		client.disconnect(ctx)
	})

	t.Run("preserves legacy handshake rejection without new options", func(t *testing.T) {
		originalClient := bootstrapHTTPClient
		defer func() { bootstrapHTTPClient = originalClient }()

		var endpoint string
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == GenEndpointUri {
				writeEndpointResponse(t, w, endpoint, nil)
				return
			}
			w.Header().Set(HeaderHandshakeStatus, "403")
			w.Header().Set(HeaderHandshakeMsg, "legacy handshake message")
			w.WriteHeader(http.StatusForbidden)
		}))
		defer server.Close()
		endpoint = "ws" + strings.TrimPrefix(server.URL, "http") + "/legacy/rejected"
		bootstrapHTTPClient = server.Client()
		client := NewClient("app-id", "app-secret",
			WithDomain(server.URL),
			WithAutoReconnect(false),
			WithLogger(&recordingLogger{}),
		)
		ctx, cancel := transportTestContext(t)
		defer cancel()
		err := client.connect(ctx)
		assertClientError(t, err, http.StatusForbidden, "legacy handshake message")
	})
}

func createTransportTestCertificate(t *testing.T, dnsName string) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate test CA key: %v", err)
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "websocket transport test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create test CA certificate: %v", err)
	}
	caCertificate, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse test CA certificate: %v", err)
	}

	serverKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate websocket server key: %v", err)
	}
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: dnsName},
		DNSNames:     []string{dnsName},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	serverDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCertificate, &serverKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create websocket server certificate: %v", err)
	}
	serverKeyDER, err := x509.MarshalPKCS8PrivateKey(serverKey)
	if err != nil {
		t.Fatalf("marshal websocket server key: %v", err)
	}
	certificate, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverDER}),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverKeyDER}),
	)
	if err != nil {
		t.Fatalf("load websocket server certificate: %v", err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(caCertificate)
	return certificate, roots
}

func TestWSSOverrideUsesFinalHostForTLSAndSNI(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()
	originalDefaultDialer := gorillaws.DefaultDialer
	defer func() { gorillaws.DefaultDialer = originalDefaultDialer }()

	const certificateHost = "ws.internal.example"
	certificate, roots := createTransportTestCertificate(t, certificateHost)
	sniNames := make(chan string, 4)
	requests := make(chan transportRequest, 2)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWebSocketUntilClientCloses(t, w, r, requests)
	}))
	server.Config.ErrorLog = log.New(ioutil.Discard, "", 0)
	server.TLS = &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
		GetConfigForClient: func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
			if hello.ServerName != "" {
				sniNames <- hello.ServerName
			}
			return nil, nil
		},
	}
	server.StartTLS()
	defer server.Close()
	endpoint := "wss://public.example/secure/socket?device_id=tls-device&service_id=41&opaque=%2Fraw"
	bootstrapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		writeEndpointResponse(t, w, endpoint, nil)
	}))
	defer bootstrapServer.Close()
	bootstrapHTTPClient = bootstrapServer.Client()

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatalf("create poisoned cookie jar: %v", err)
	}
	poisonedPool := &sync.Pool{}
	gorillaws.DefaultDialer = &gorillaws.Dialer{
		HandshakeTimeout: time.Nanosecond,
		TLSClientConfig:  &tls.Config{ServerName: "wrong-before-new-client.example"},
		NetDialTLSContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return nil, errors.New("poisoned global TLS dialer was used")
		},
		Subprotocols:      []string{"poisoned-subprotocol"},
		Jar:               jar,
		WriteBufferPool:   poisonedPool,
		EnableCompression: true,
	}

	serverPort := fmt.Sprintf("%d", server.Listener.Addr().(*net.TCPAddr).Port)
	overrideHost := net.JoinHostPort(certificateHost, serverPort)
	client := NewClient("app-id", "app-secret",
		WithDomain(bootstrapServer.URL),
		WithConnectionHost(overrideHost),
		WithAutoReconnect(false),
		WithLogger(&recordingLogger{}),
	)
	// Mutating Gorilla's global after construction must not change a client's transport policy.
	gorillaws.DefaultDialer.TLSClientConfig.ServerName = "wrong-after-new-client.example"
	gorillaws.DefaultDialer.Subprotocols[0] = "mutated-subprotocol"
	gorillaws.DefaultDialer.EnableCompression = false

	if client.dialer == gorillaws.DefaultDialer {
		t.Fatal("client shares Gorilla's mutable global dialer")
	}
	if client.dialer.TLSClientConfig != nil || client.dialer.NetDialTLSContext != nil || len(client.dialer.Subprotocols) != 0 || client.dialer.Jar != nil || client.dialer.WriteBufferPool != nil || client.dialer.EnableCompression {
		t.Fatalf("client inherited mutable global dialer fields: %#v", client.dialer)
	}
	if client.dialer.HandshakeTimeout != 45*time.Second {
		t.Fatalf("unexpected fresh dialer handshake timeout: %v", client.dialer.HandshakeTimeout)
	}

	client.dialer.Proxy = nil
	client.dialer.TLSClientConfig = &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	client.dialer.NetDialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		dialer := &net.Dialer{}
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	if client.dialer.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("TLS certificate verification was disabled")
	}
	if client.dialer.TLSClientConfig.ServerName != "" || client.dialer.NetDialTLSContext != nil {
		t.Fatal("test seam would bypass URL-derived TLS ServerName")
	}

	ctx, cancel := transportTestContext(t)
	defer cancel()
	if err := client.connect(ctx); err != nil {
		t.Fatalf("connect matching wss override: %v", err)
	}
	request := awaitTransportRequest(t, ctx, requests)
	if request.host != overrideHost {
		t.Fatalf("unexpected TLS websocket Host: got %q, want %q", request.host, overrideHost)
	}
	select {
	case sni := <-sniNames:
		if sni != certificateHost {
			t.Fatalf("unexpected SNI: got %q, want %q", sni, certificateHost)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for matching SNI: %v", ctx.Err())
	}
	client.disconnect(ctx)

	mismatchHost := net.JoinHostPort("mismatch.internal.example", serverPort)
	mismatchClient := NewClient("app-id", "app-secret",
		WithDomain(bootstrapServer.URL),
		WithConnectionHost(mismatchHost),
		WithAutoReconnect(false),
		WithLogger(&recordingLogger{}),
	)
	mismatchClient.dialer.Proxy = nil
	mismatchClient.dialer.TLSClientConfig = &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	mismatchClient.dialer.NetDialContext = func(ctx context.Context, network, _ string) (net.Conn, error) {
		dialer := &net.Dialer{}
		return dialer.DialContext(ctx, network, server.Listener.Addr().String())
	}
	err = mismatchClient.connect(ctx)
	if err == nil {
		mismatchClient.disconnect(ctx)
		t.Fatal("certificate mismatch unexpectedly connected")
	}
	if mismatchClient.conn != nil {
		t.Fatal("certificate mismatch retained a websocket connection")
	}
	select {
	case sni := <-sniNames:
		if sni != "mismatch.internal.example" {
			t.Fatalf("unexpected mismatch SNI: %q", sni)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for mismatch SNI: %v", ctx.Err())
	}
	select {
	case unexpected := <-requests:
		t.Fatalf("certificate mismatch reached HTTP websocket handshake: %#v", unexpected)
	default:
	}
}

func TestConnectionLogsRedactSensitiveValues(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	const (
		bootstrapSecret  = "bootstrap-log-secret"
		connectionSecret = "connection-log-secret"
		cookieSecret     = "cookie-log-secret"
		deviceSecret     = "device-log-secret"
		querySecret      = "raw-query-log-secret"
		handshakeSecret  = "server-handshake-secret"
	)

	successRequests := make(chan transportRequest, 1)
	var successEndpoint string
	successServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == GenEndpointUri {
			writeEndpointResponse(t, w, successEndpoint, nil)
			return
		}
		serveWebSocketUntilClientCloses(t, w, r, successRequests)
	}))
	defer successServer.Close()
	successEndpoint = "ws://public.example/sensitive/path?device_id=" + deviceSecret + "&service_id=51&opaque=" + querySecret
	bootstrapHTTPClient = successServer.Client()
	successLogger := &recordingLogger{}
	successClient := NewClient("app-id", "app-secret",
		WithDomain(successServer.URL),
		WithHeaders(http.Header{"X-Bootstrap-Secret": []string{bootstrapSecret}}),
		WithConnectionHeaders(http.Header{
			"Authorization": []string{"Bearer " + connectionSecret},
			"Cookie":        []string{"session=" + cookieSecret},
		}),
		WithConnectionHost(successServer.Listener.Addr().String()),
		WithAutoReconnect(false),
		WithLogger(successLogger),
	)

	ctx, cancel := transportTestContext(t)
	defer cancel()
	if err := successClient.connect(ctx); err != nil {
		t.Fatalf("connect sensitive logging fixture: %v", err)
	}
	request := awaitTransportRequest(t, ctx, successRequests)
	if request.header.Get("Authorization") != "Bearer "+connectionSecret || request.header.Get("Cookie") != "session="+cookieSecret {
		t.Fatalf("sensitive logging fixture did not receive headers: %#v", request.header)
	}
	successClient.disconnect(ctx)

	var rejectedEndpoint string
	rejectedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == GenEndpointUri {
			writeEndpointResponse(t, w, rejectedEndpoint, nil)
			return
		}
		maliciousMessage := strings.Join([]string{
			handshakeSecret,
			r.Header.Get("Authorization"),
			r.Header.Get("Cookie"),
			r.URL.RawQuery,
		}, "|")
		w.Header().Set(HeaderHandshakeStatus, "403")
		w.Header().Set(HeaderHandshakeMsg, maliciousMessage)
		w.WriteHeader(http.StatusForbidden)
	}))
	defer rejectedServer.Close()
	rejectedEndpoint = "ws://public.example/rejected/path?device_id=" + deviceSecret + "&service_id=52&opaque=" + querySecret
	bootstrapHTTPClient = rejectedServer.Client()
	rejectedLogger := &recordingLogger{}
	onError := make(chan error, 1)
	rejectedClient := NewClient("app-id", "app-secret",
		WithDomain(rejectedServer.URL),
		WithHeaders(http.Header{"X-Bootstrap-Secret": []string{bootstrapSecret}}),
		WithConnectionHeaders(http.Header{
			"Authorization": []string{"Bearer " + connectionSecret},
			"Cookie":        []string{"session=" + cookieSecret},
		}),
		WithConnectionHost(rejectedServer.Listener.Addr().String()),
		WithAutoReconnect(false),
		WithLogger(rejectedLogger),
		WithOnError(func(err error) { onError <- err }),
	)
	err := rejectedClient.Start(ctx)
	assertClientError(t, err, http.StatusForbidden, "websocket handshake failed")
	var callbackErr error
	select {
	case callbackErr = <-onError:
		assertClientError(t, callbackErr, http.StatusForbidden, "websocket handshake failed")
	case <-ctx.Done():
		t.Fatalf("timed out waiting for rejected OnError: %v", ctx.Err())
	}

	for name, text := range map[string]string{
		"success logs": successLogger.String(),
		"failure logs": rejectedLogger.String(),
		"return error": err.Error(),
		"OnError":      callbackErr.Error(),
	} {
		t.Run(name, func(t *testing.T) {
			assertNotContainsAny(t, text,
				bootstrapSecret,
				connectionSecret,
				cookieSecret,
				deviceSecret,
				querySecret,
				handshakeSecret,
				successEndpoint,
				rejectedEndpoint,
			)
		})
	}
}

func TestConnectionLogsRedactMalformedHandshakeResponse(t *testing.T) {
	originalClient := bootstrapHTTPClient
	defer func() { bootstrapHTTPClient = originalClient }()

	const (
		headerSecret = "malformed-header-secret"
		querySecret  = "malformed-query-secret"
		deviceSecret = "malformed-device-secret"
	)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for malformed handshake fixture: %v", err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			t.Errorf("close malformed handshake listener: %v", err)
		}
	}()

	requestSeen := make(chan transportRequest, 1)
	serverDone := make(chan error, 1)
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			serverDone <- fmt.Errorf("accept malformed handshake connection: %w", err)
			return
		}
		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			if closeErr := conn.Close(); closeErr != nil {
				serverDone <- fmt.Errorf("read malformed handshake request: %v; close connection: %w", err, closeErr)
				return
			}
			serverDone <- fmt.Errorf("read malformed handshake request: %w", err)
			return
		}
		if err := request.Body.Close(); err != nil {
			if closeErr := conn.Close(); closeErr != nil {
				serverDone <- fmt.Errorf("close malformed handshake request body: %v; close connection: %w", err, closeErr)
				return
			}
			serverDone <- fmt.Errorf("close malformed handshake request body: %w", err)
			return
		}
		requestSeen <- transportRequest{
			host:        request.Host,
			escapedPath: request.URL.EscapedPath(),
			rawQuery:    request.URL.RawQuery,
			header:      cloneRequestHeader(request.Header),
		}
		reflectedValue := request.Header.Get("Authorization") + "-" + request.URL.Query().Get("opaque")
		_, writeErr := fmt.Fprintf(conn, "HTTP/1.1 %s\r\n\r\n", reflectedValue)
		closeErr := conn.Close()
		switch {
		case writeErr != nil:
			serverDone <- fmt.Errorf("write malformed handshake response: %w", writeErr)
		case closeErr != nil:
			serverDone <- fmt.Errorf("close malformed handshake connection: %w", closeErr)
		default:
			serverDone <- nil
		}
	}()

	var endpoint string
	bootstrapServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeEndpointResponse(t, w, endpoint, nil)
	}))
	defer bootstrapServer.Close()
	endpoint = "ws://public.example/malformed?device_id=" + deviceSecret + "&service_id=53&opaque=" + querySecret
	bootstrapHTTPClient = bootstrapServer.Client()

	logger := &recordingLogger{}
	onError := make(chan error, 1)
	client := NewClient("app-id", "app-secret",
		WithDomain(bootstrapServer.URL),
		WithConnectionHeaders(http.Header{"Authorization": []string{headerSecret}}),
		WithConnectionHost(listener.Addr().String()),
		WithAutoReconnect(false),
		WithLogger(logger),
		WithOnError(func(err error) { onError <- err }),
	)

	ctx, cancel := transportTestContext(t)
	defer cancel()
	err = client.Start(ctx)
	if err == nil || err.Error() != errWebSocketHandshakeFailed.Error() {
		t.Fatalf("unexpected malformed handshake error: %v", err)
	}
	var callbackErr error
	select {
	case callbackErr = <-onError:
	case <-ctx.Done():
		t.Fatalf("timed out waiting for malformed handshake OnError: %v", ctx.Err())
	}
	if callbackErr == nil || callbackErr.Error() != errWebSocketHandshakeFailed.Error() {
		t.Fatalf("unexpected malformed handshake callback error: %v", callbackErr)
	}
	request := awaitTransportRequest(t, ctx, requestSeen)
	if request.header.Get("Authorization") != headerSecret || request.rawQuery != "device_id="+deviceSecret+"&service_id=53&opaque="+querySecret {
		t.Fatalf("malformed handshake fixture did not receive expected data: %#v", request)
	}
	select {
	case serverErr := <-serverDone:
		if serverErr != nil {
			t.Fatal(serverErr)
		}
	case <-ctx.Done():
		t.Fatalf("timed out waiting for malformed handshake server: %v", ctx.Err())
	}

	for name, text := range map[string]string{
		"failure logs": logger.String(),
		"return error": err.Error(),
		"OnError":      callbackErr.Error(),
	} {
		t.Run(name, func(t *testing.T) {
			assertNotContainsAny(t, text, headerSecret, querySecret, deviceSecret, endpoint)
		})
	}
}
