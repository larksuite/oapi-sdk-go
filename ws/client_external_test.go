package ws_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

	client := ws.NewClient("app-id", "app-secret",
		ws.WithDomain(server.URL),
		ws.WithAutoReconnect(false),
	)
	err := client.Start(context.Background())
	if err == nil {
		t.Fatal("Start returned nil after bootstrap HTTP failure")
	}

	var serverErr *ws.ServerError
	if !errors.As(err, &serverErr) {
		t.Fatalf("errors.As(%T, *ws.ServerError) = false", err)
	}
	if serverErr.Code != http.StatusInternalServerError {
		t.Fatalf("ServerError.Code = %d, want %d", serverErr.Code, http.StatusInternalServerError)
	}
	if serverErr.Msg != serverMessage {
		t.Fatalf("ServerError.Msg = %q, want %q", serverErr.Msg, serverMessage)
	}
	if strings.Contains(err.Error(), serverMessage) {
		t.Fatalf("outer error exposed server response: %s", err)
	}
}

func TestStartPreservesAssertionProviderError(t *testing.T) {
	providerErr := errors.New("assertion provider network failure")
	client := ws.NewClient("app-id", "",
		ws.WithClientAssertionProvider(failingAssertionProvider{err: providerErr}),
		ws.WithAutoReconnect(false),
	)

	err := client.Start(context.Background())
	if !errors.Is(err, providerErr) {
		t.Fatalf("errors.Is(Start error, provider error) = false: %v", err)
	}
	if strings.Contains(err.Error(), providerErr.Error()) {
		t.Fatalf("outer error exposed assertion provider error: %s", err)
	}
}
