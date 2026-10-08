package larkcore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestSSEParserFraming(t *testing.T) {
	p := newSSEParser(strings.NewReader("\ufeff: hello\r\nid: cursor\r\n\r\nevent: text.delta\ndata: 中文\ndata: second\nid:\nretry: 1000\n\ndata: incomplete"), 1024, "old")
	e, err := p.next()
	if err != nil || e.Event != "text.delta" || e.Data != "中文\nsecond" || e.ID == nil || *e.ID != "" || p.cursor != "" || p.retry == nil || *p.retry != 1000 {
		t.Fatalf("unexpected framing: %+v, %v", e, err)
	}
	if _, err := p.next(); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF: %v", err)
	}
}

func TestSSEParserLimit(t *testing.T) {
	p := newSSEParser(strings.NewReader("data: "+strings.Repeat("x", 20)+"\n\n"), 16, "")
	if _, err := p.next(); err == nil {
		t.Fatal("expected bounded parser error")
	}
}

func openTestSSE(t *testing.T, ctx context.Context, handler http.HandlerFunc, options SSEOptions) (*SSEStream, func()) {
	t.Helper()
	server := httptest.NewServer(handler)
	config := &Config{Serializable: &DefaultSerialization{}, AppId: "test-app", AppSecret: "test-secret", BaseUrl: server.URL, HttpClient: server.Client()}
	stream, err := RequestSSE(ctx, &ApiReq{ApiPath: "/events", HttpMethod: http.MethodGet}, config, WithSSEOptions(options))
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	return stream, server.Close
}

func TestSSEImmediateDeliveryAndCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, stop := openTestSSE(t, ctx, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		if _, err := io.WriteString(w, "id: first\ndata: first\n\n"); err != nil {
			return
		}
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}, SSEOptions{})
	defer stop()
	defer func() {
		if err := stream.Close(); err != nil {
			t.Error(err)
		}
	}()
	result := make(chan bool, 1)
	go func() { result <- stream.Next() }()
	select {
	case ok := <-result:
		if !ok || stream.Event().Data != "first" {
			t.Fatal("first event not delivered")
		}
	case <-time.After(time.Second):
		t.Fatal("waited for response completion")
	}
	go func() { result <- stream.Next() }()
	if stream.LastEventID() != "first" {
		t.Fatal("cursor vanished during blocked read")
	}
	cancel()
	select {
	case ok := <-result:
		if ok || !errors.Is(stream.Err(), context.Canceled) {
			t.Fatalf("cancel: %v", stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("cancel did not interrupt read")
	}
}

func TestSSEIdleTimeout(t *testing.T) {
	stream, stop := openTestSSE(t, context.Background(), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}, SSEOptions{IdleReadTimeout: 30 * time.Millisecond})
	defer stop()
	if stream.Next() || !errors.Is(stream.Err(), context.DeadlineExceeded) {
		t.Fatalf("idle timeout: %v", stream.Err())
	}
}

func TestSSERejectOrdinaryResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, "{}"); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	_, err := RequestSSE(context.Background(), &ApiReq{ApiPath: "/", HttpMethod: http.MethodGet}, &Config{Serializable: &DefaultSerialization{}, AppId: "test-app", AppSecret: "test-secret", BaseUrl: server.URL, HttpClient: server.Client()})
	if err == nil {
		t.Fatal("accepted ordinary JSON response")
	}
}

func TestSSEPostHeadersAndNoRetry(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.Method != http.MethodPost || r.Header.Get("Accept") != "text/event-stream" || r.Header.Get("Last-Event-ID") != "resume" || r.Header.Get("Authorization") != "Bearer test-user-token" {
			t.Error("request translation failed")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if !strings.Contains(string(body), `"text":"hello"`) {
			t.Error("missing JSON request body")
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()
	req := &ApiReq{ApiPath: "/events", HttpMethod: http.MethodPost, Body: map[string]string{"text": "hello"}, SupportedAccessTokenTypes: []AccessTokenType{AccessTokenTypeUser}}
	_, err := RequestSSE(context.Background(), req, &Config{Serializable: &DefaultSerialization{}, AppId: "test-app", AppSecret: "test-secret", BaseUrl: server.URL, HttpClient: server.Client()}, WithUserAccessToken("test-user-token"), WithSSEOptions(SSEOptions{LastEventID: "resume"}))
	if err == nil || calls != 1 {
		t.Fatalf("calls=%d, error=%v", calls, err)
	}
}

func TestSSEOpenTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { <-r.Context().Done() }))
	defer server.Close()
	_, err := RequestSSE(context.Background(), &ApiReq{ApiPath: "/", HttpMethod: http.MethodGet}, &Config{Serializable: &DefaultSerialization{}, AppId: "test-app", AppSecret: "test-secret", BaseUrl: server.URL, HttpClient: server.Client()}, WithSSEOptions(SSEOptions{OpenTimeout: 30 * time.Millisecond}))
	if err == nil {
		t.Fatal("expected opening timeout")
	}
}
