package ws

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestConnectionHostValidationAndRawURL(t *testing.T) {
	for _, host := range []string{"example.test", "example.test:443", "127.0.0.1:80", "[::1]", "[::1]:443"} {
		if err := validateConnectionHost(host); err != nil {
			t.Errorf("valid host %q: %v", host, err)
		}
		u, err := buildConnectionURL("wss://original.test/a%2Fb?device_id=d&x=%2F&x=2", &host)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := u.String(), "wss://"+host+"/a%2Fb?device_id=d&x=%2F&x=2"; got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
	for _, host := range []string{"", " ", "https://example.test", "x/path", "x?q=1"} {
		if validateConnectionHost(host) == nil {
			t.Errorf("accepted invalid host %q", host)
		}
	}
}

func TestConnectionConfigurationFailsBeforeBootstrap(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { atomic.AddInt32(&calls, 1) }))
	defer server.Close()
	options := []ClientOption{WithConnectionHost("")}
	for _, name := range []string{"Host", "cOnNeCtIoN", "Upgrade", "sEc-WeBsOcKeT-Key"} {
		options = append(options, WithConnectionHeaders(http.Header{name: []string{"x"}}))
	}
	for _, option := range options {
		client := NewClient("app", "placeholder", WithDomain(server.URL), option)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := client.Start(ctx)
		cancel()
		if _, ok := err.(*ClientError); !ok {
			t.Errorf("got %T, want ClientError", err)
		}
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("bootstrap calls = %d", got)
	}
}

func TestConnectionHostAndHeadersAreIsolatedAndReused(t *testing.T) {
	type request struct{ host, uri, route, bootstrap string }
	handshakes := make(chan request, 2)
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handshakes <- request{r.Host, r.RequestURI, r.Header.Get("X-Route"), r.Header.Get("X-Bootstrap")}
		w.Header().Set(HeaderHandshakeStatus, "403")
		w.WriteHeader(http.StatusForbidden)
	}))
	defer gateway.Close()
	bootstrap := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Route") != "" || r.Header.Get("X-Bootstrap") != "http-only" {
			t.Error("bootstrap headers leaked or missing")
		}
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{Url: "ws://original.invalid/a%2Fb?device_id=d&service_id=42&x=%2F&x=2"}}); err != nil {
			t.Error(err)
		}
	}))
	defer bootstrap.Close()
	headers := http.Header{"X-Route": []string{"ws-only"}}
	dialer := &websocket.Dialer{HandshakeTimeout: time.Second, NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr != "logical.test:80" {
			t.Errorf("logical dial address = %s", addr)
		}
		return (&net.Dialer{Timeout: time.Second}).DialContext(ctx, network, strings.TrimPrefix(gateway.URL, "http://"))
	}}
	client := NewClient("app", "placeholder", WithDomain(bootstrap.URL), WithWebSocketDialer(dialer),
		WithHeaders(http.Header{"X-Bootstrap": []string{"http-only"}}), WithConnectionHost("logical.test"), WithConnectionHeaders(headers))
	headers["X-Route"][0] = "mutated"
	for i := 0; i < 2; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.establishConnectionOnce(&clientRun{ctx: ctx})
		cancel()
		if err == nil {
			t.Fatal("expected rejected handshake")
		}
		select {
		case got := <-handshakes:
			if got.host != "logical.test" || got.uri != "/a%2Fb?device_id=d&service_id=42&x=%2F&x=2" || got.route != "ws-only" || got.bootstrap != "" {
				t.Errorf("unexpected handshake: %+v", got)
			}
		case <-time.After(time.Second):
			t.Fatal("handshake not received")
		}
	}
}

func TestConnectionHostControlsTLSVerification(t *testing.T) {
	gateway := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusForbidden) }))
	defer gateway.Close()
	roots := x509.NewCertPool()
	roots.AddCert(gateway.Certificate())
	var bootstrapRequests int32
	domain, cleanup := installLifecycleBootstrap(t, "wss://original.invalid/ws?device_id=d&service_id=42", &bootstrapRequests)
	defer cleanup()
	dialer := &websocket.Dialer{HandshakeTimeout: time.Second, TLSClientConfig: &tls.Config{RootCAs: roots},
		NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{Timeout: time.Second}).DialContext(ctx, network, strings.TrimPrefix(gateway.URL, "https://"))
		}}
	for _, host := range []string{"127.0.0.1", "wrong.invalid"} {
		client := NewClient("app", "placeholder", WithDomain(domain), WithWebSocketDialer(dialer), WithConnectionHost(host))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := client.establishConnectionOnce(&clientRun{ctx: ctx})
		cancel()
		if host == "127.0.0.1" {
			if _, ok := err.(*ClientError); !ok {
				t.Errorf("TLS should reach HTTP handshake, got %T: %v", err, err)
			}
		} else if err == nil || !strings.Contains(err.Error(), "certificate") {
			t.Errorf("expected certificate rejection, got %v", err)
		}
	}
}

func TestLifecycleReconnectUsesConnectionHostAndHeaders(t *testing.T) {
	requests := make(chan *http.Request, 4)
	sockets := make(chan *websocket.Conn, 4)
	upgrader := websocket.Upgrader{}
	gateway := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Error(err)
			return
		}
		requests <- r
		sockets <- socket
	}))
	defer gateway.Close()
	bootstrap := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&EndpointResp{Code: OK, Data: &Endpoint{
			Url:          "ws://original.invalid/ws?device_id=d&service_id=42",
			ClientConfig: &ClientConfig{ReconnectCount: 2, ReconnectNonce: 0, ReconnectInterval: 0, PingInterval: 3600},
		}}); err != nil {
			t.Error(err)
		}
	}))
	defer bootstrap.Close()
	dialer := &websocket.Dialer{HandshakeTimeout: time.Second, NetDialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		return (&net.Dialer{Timeout: time.Second}).DialContext(ctx, network, strings.TrimPrefix(gateway.URL, "http://"))
	}}
	ready := make(chan struct{}, 1)
	reconnected := make(chan struct{}, 1)
	client := NewClient("app", "placeholder", WithDomain(bootstrap.URL), WithWebSocketDialer(dialer),
		WithConnectionHost("logical.test"), WithConnectionHeaders(http.Header{"X-Route": []string{"ws-only"}}),
		WithOnReady(func() { ready <- struct{}{} }), WithOnReconnected(func() { reconnected <- struct{}{} }))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result := startLifecycleClient(client, ctx)
	defer client.Close()
	waitLifecycleSignal(t, ready, "Ready")
	var first *websocket.Conn
	select {
	case first = <-sockets:
	case <-ctx.Done():
		t.Fatal("first socket not received")
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	waitLifecycleSignal(t, reconnected, "Reconnected")
	var second *websocket.Conn
	select {
	case second = <-sockets:
	case <-ctx.Done():
		t.Fatal("second socket not received")
	}
	defer func() {
		if err := second.Close(); err != nil {
			t.Error(err)
		}
	}()
	for i := 0; i < 2; i++ {
		select {
		case r := <-requests:
			if r.Host != "logical.test" || r.Header.Get("X-Route") != "ws-only" {
				t.Error("reconnect lost configuration")
			}
		case <-ctx.Done():
			t.Fatal("handshake not received")
		}
	}
	client.Close()
	if err, ok := waitLifecycleResult(result); !ok || err != nil {
		t.Errorf("Close result: %v, completed: %v", err, ok)
	}
}
