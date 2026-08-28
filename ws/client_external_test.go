package ws_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	ws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

type failingAssertionProvider struct {
	err error
}

func (p failingAssertionProvider) RetrieveToken(context.Context, string) (*larkcore.Token, error) {
	return nil, p.err
}

func TestStartHTTPFailurePreservesServerError(t *testing.T) {
	const serverMessage = "target service unavailable"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 20050,
			"msg":  serverMessage,
		}); err != nil {
			t.Errorf("encode bootstrap error: %v", err)
		}
	}))
	defer server.Close()

	callbackErr := make(chan error, 1)
	client := ws.NewClient("app-id", "app-secret",
		ws.WithDomain(server.URL),
		ws.WithAutoReconnect(false),
		ws.WithOnError(func(err error) { callbackErr <- err }),
	)
	err := client.Start(context.Background())
	if err == nil {
		t.Fatal("Start returned nil after bootstrap HTTP failure")
	}

	serverErr, ok := err.(*ws.ServerError)
	if !ok {
		t.Fatalf("Start returned %T, want *ws.ServerError", err)
	}
	if serverErr.Code != http.StatusInternalServerError {
		t.Fatalf("ServerError.Code = %d, want %d", serverErr.Code, http.StatusInternalServerError)
	}
	if serverErr.Msg != serverMessage {
		t.Fatalf("ServerError.Msg = %q, want %q", serverErr.Msg, serverMessage)
	}
	select {
	case callbackErr := <-callbackErr:
		if _, ok := callbackErr.(*ws.ServerError); !ok {
			t.Errorf("OnError received %T, want *ws.ServerError", callbackErr)
		}
	default:
		t.Fatal("OnError was not invoked")
	}
}

func TestStartClientFailurePreservesClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"code": 19001,
			"msg":  "invalid credential",
		}); err != nil {
			t.Errorf("encode bootstrap error: %v", err)
		}
	}))
	defer server.Close()

	callbackErr := make(chan error, 1)
	client := ws.NewClient("app-id", "app-secret",
		ws.WithDomain(server.URL),
		ws.WithOnError(func(err error) { callbackErr <- err }),
	)
	err := client.Start(context.Background())
	if err == nil {
		t.Fatal("Start returned nil after bootstrap client failure")
	}
	clientErr, ok := err.(*ws.ClientError)
	if !ok {
		t.Fatalf("Start returned %T, want *ws.ClientError", err)
	}
	if clientErr.Code != 19001 || clientErr.Msg != "invalid credential" {
		t.Fatalf("ClientError = %#v, want code 19001 and original message", clientErr)
	}
	select {
	case callbackErr := <-callbackErr:
		if _, ok := callbackErr.(*ws.ClientError); !ok {
			t.Errorf("OnError received %T, want *ws.ClientError", callbackErr)
		}
	default:
		t.Fatal("OnError was not invoked")
	}
}

func TestStartPreservesAssertionProviderError(t *testing.T) {
	providerErr := errors.New("assertion provider network failure")
	client := ws.NewClient("app-id", "",
		ws.WithClientAssertionProvider(failingAssertionProvider{err: providerErr}),
		ws.WithAutoReconnect(false),
	)

	err := client.Start(context.Background())
	if err != providerErr {
		t.Fatalf("Start error = %v, want original provider error %v", err, providerErr)
	}
}
