package dispatcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
)

func TestWebhookProviderRequestTypes(t *testing.T) {
	for _, schema := range []string{"v1", "v2", "challenge", "callback"} {
		t.Run(schema, func(t *testing.T) {
			providerCalls, handlerCalls := 0, 0
			dispatcher := NewEventDispatcher("static-token", "static-key").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
				providerCalls++
				return WebhookCredentials{VerificationToken: "dynamic-token"}, nil
			}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error {
				handlerCalls++
				return nil
			}).OnP2CardActionTrigger(func(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
				handlerCalls++
				return &callback.CardActionTriggerResponse{}, nil
			})
			if providerCalls != 0 {
				t.Fatal("initialization invoked provider")
			}
			body := makeTokenPayload(t, schema, "dynamic-token", true)
			if schema == "callback" {
				body = []byte(strings.Replace(string(body), "audit.test", "card.action.trigger", 1))
			}
			resp := dispatcher.Handle(context.Background(), &larkevent.EventReq{Body: body})
			wantHandlers := 1
			if schema == "challenge" {
				wantHandlers = 0
			}
			if resp.StatusCode != http.StatusOK || providerCalls != 1 || handlerCalls != wantHandlers {
				t.Fatalf("status=%d provider=%d handler=%d", resp.StatusCode, providerCalls, handlerCalls)
			}
			if dispatcher.verificationToken != "static-token" || dispatcher.eventEncryptKey != "static-key" {
				t.Fatal("shared constructor credentials changed")
			}
		})
	}
}

func makeProviderEncryptedRequest(t *testing.T, token, key string) *larkevent.EventReq {
	t.Helper()
	payload := map[string]interface{}{}
	if err := json.Unmarshal(makeTokenPayload(t, "v2", token, true), &payload); err != nil {
		t.Fatalf("decode provider fixture: %v", err)
	}
	encrypted, err := larkcore.EncryptedEventMsg(context.Background(), payload, key)
	if err != nil {
		t.Fatalf("encrypt provider fixture: %v", err)
	}
	body, err := json.Marshal(larkevent.EventEncryptMsg{Encrypt: encrypted})
	if err != nil {
		t.Fatalf("encode provider fixture: %v", err)
	}
	return &larkevent.EventReq{Body: body, Header: map[string][]string{
		larkevent.EventRequestTimestamp: {"timestamp"},
		larkevent.EventRequestNonce:     {"nonce"},
		larkevent.EventSignature:        {larkevent.Signature("timestamp", "nonce", key, string(body))},
	}}
}

func TestWebhookProviderRotation(t *testing.T) {
	var current atomic.Value
	current.Store(WebhookCredentials{VerificationToken: "first"})
	providerCalls, handlerCalls := 0, 0
	dispatcher := NewEventDispatcher("", "").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
		providerCalls++
		credentials := current.Load().(WebhookCredentials)
		// Rotate before processing begins: this request must retain the old snapshot.
		current.Store(WebhookCredentials{VerificationToken: "second", EncryptKey: "second-key"})
		return credentials, nil
	}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error {
		handlerCalls++
		return nil
	})
	first := dispatcher.Handle(context.Background(), &larkevent.EventReq{Body: makeTokenPayload(t, "v2", "first", true)})
	second := dispatcher.Handle(context.Background(), makeProviderEncryptedRequest(t, "second", "second-key"))
	if first.StatusCode != http.StatusOK || second.StatusCode != http.StatusOK || providerCalls != 2 || handlerCalls != 2 {
		t.Fatalf("rotation: statuses=%d/%d provider=%d handler=%d", first.StatusCode, second.StatusCode, providerCalls, handlerCalls)
	}
}

func TestWebhookProviderSingleEncryptedSnapshot(t *testing.T) {
	providerCalls, handlerCalls := 0, 0
	dispatcher := NewEventDispatcher("", "static-key").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
		providerCalls++
		if providerCalls == 1 {
			return WebhookCredentials{VerificationToken: "first", EncryptKey: "first-key"}, nil
		}
		return WebhookCredentials{VerificationToken: "second", EncryptKey: "second-key"}, nil
	}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error {
		handlerCalls++
		return nil
	})
	resp := dispatcher.Handle(context.Background(), makeProviderEncryptedRequest(t, "first", "first-key"))
	if resp.StatusCode != http.StatusOK || providerCalls != 1 || handlerCalls != 1 {
		t.Fatalf("snapshot: status=%d provider=%d handler=%d", resp.StatusCode, providerCalls, handlerCalls)
	}
}

func TestWebhookProviderEmptyPreservesStaticBehavior(t *testing.T) {
	for _, schema := range []string{"v1", "v2", "challenge"} {
		for _, token := range []string{"", "wrong"} {
			t.Run(schema+"/"+token, func(t *testing.T) {
				calls := 0
				dynamic := NewEventDispatcher("static-token", "static-key").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
					calls++
					return WebhookCredentials{}, nil
				})
				static := NewEventDispatcher("", "")
				for _, dispatcher := range []*EventDispatcher{dynamic, static} {
					dispatcher.OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { return nil })
				}
				req := &larkevent.EventReq{Body: makeTokenPayload(t, schema, token, true)}
				got := dynamic.Handle(context.Background(), req)
				want := static.Handle(context.Background(), req)
				if got.StatusCode != want.StatusCode || string(got.Body) != string(want.Body) || calls != 1 {
					t.Fatalf("empty provider differs from static: got=%v want=%v calls=%d", got.StatusCode, want.StatusCode, calls)
				}
			})
		}
	}
}

type providerTestLogger struct{ entries []string }

func (logger *providerTestLogger) Debug(ctx context.Context, args ...interface{}) {
	logger.entries = append(logger.entries, fmt.Sprint(args...))
}
func (logger *providerTestLogger) Info(ctx context.Context, args ...interface{}) {
	logger.entries = append(logger.entries, fmt.Sprint(args...))
}
func (logger *providerTestLogger) Warn(ctx context.Context, args ...interface{}) {
	logger.entries = append(logger.entries, fmt.Sprint(args...))
}
func (logger *providerTestLogger) Error(ctx context.Context, args ...interface{}) {
	logger.entries = append(logger.entries, fmt.Sprint(args...))
}

func TestWebhookProviderFailureSanitized(t *testing.T) {
	for _, mode := range []string{"error", "panic", "panic-nil", "canceled-before", "canceled-after"} {
		for _, schema := range []string{"v2", "challenge"} {
			t.Run(mode+"/"+schema, func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				providerCalls, handlerCalls := 0, 0
				logger := &providerTestLogger{}
				dispatcher := NewEventDispatcher("static-token", "").WithWebhookCredentialsProvider(func(providerCtx context.Context) (WebhookCredentials, error) {
					providerCalls++
					if providerCtx != ctx {
						t.Error("request context not propagated")
					}
					switch mode {
					case "error":
						return WebhookCredentials{}, errors.New("secret-fixture")
					case "panic":
						panic("secret-fixture")
					case "panic-nil":
						panic(nil)
					case "canceled-after":
						cancel()
					}
					return WebhookCredentials{VerificationToken: "static-token"}, nil
				}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { handlerCalls++; return nil })
				dispatcher.Config.Logger = logger
				if mode == "canceled-before" {
					cancel()
				}
				resp := dispatcher.Handle(ctx, &larkevent.EventReq{Body: makeTokenPayload(t, schema, "static-token", true)})
				wantCalls := 1
				if mode == "canceled-before" {
					wantCalls = 0
				}
				if resp.StatusCode == http.StatusOK || handlerCalls != 0 || providerCalls != wantCalls {
					t.Fatalf("failure bypassed: status=%d handler=%d provider=%d", resp.StatusCode, handlerCalls, providerCalls)
				}
				if strings.Contains(string(resp.Body)+strings.Join(logger.entries, "\n"), "secret-fixture") {
					t.Fatal("provider secret leaked")
				}
			})
		}
	}
}

func TestWebhookProviderTrustedBypass(t *testing.T) {
	providerCalls, handlerCalls := 0, 0
	dispatcher := NewEventDispatcher("", "").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
		providerCalls++
		return WebhookCredentials{}, errors.New("must not be invoked")
	}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { handlerCalls++; return nil })
	if _, err := dispatcher.Do(context.Background(), makeTokenPayload(t, "v2", nil, false)); err != nil {
		t.Fatalf("trusted dispatch failed: %v", err)
	}
	if providerCalls != 0 || handlerCalls != 1 {
		t.Fatalf("provider=%d handler=%d", providerCalls, handlerCalls)
	}
}

func TestWebhookProviderConcurrentSnapshots(t *testing.T) {
	var providerCalls, handlerCalls int32
	// Each request reads distinct credentials while sharing one dispatcher. The barrier ensures
	// they finish reading before processing, exposing shared-field replacement.
	const requestCount = 16
	var arrived sync.WaitGroup
	arrived.Add(requestCount)
	var next int32
	entered := make(chan struct{})
	dispatcher := NewEventDispatcher("", "").WithWebhookCredentialsProvider(func(ctx context.Context) (WebhookCredentials, error) {
		atomic.AddInt32(&providerCalls, 1)
		index := atomic.AddInt32(&next, 1)
		credentials := WebhookCredentials{EncryptKey: fmt.Sprintf("key-%d", index)}
		entered <- struct{}{}
		arrived.Done()
		arrived.Wait()
		return credentials, nil
	}).OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error {
		atomic.AddInt32(&handlerCalls, 1)
		return nil
	})
	// Launch sequentially up to the provider barrier so payloads align with the
	// provider's deterministic sequence without using context to pass credentials.
	var completed sync.WaitGroup
	completed.Add(requestCount)
	for index := 1; index <= requestCount; index++ {
		req := makeProviderEncryptedRequest(t, "", fmt.Sprintf("key-%d", index))
		go runConcurrentWebhook(t, dispatcher, req, &completed)
		<-entered
	}
	completed.Wait()
	if providerCalls != requestCount || handlerCalls != requestCount {
		t.Fatalf("provider=%d handler=%d", providerCalls, handlerCalls)
	}
}

func runConcurrentWebhook(t *testing.T, dispatcher *EventDispatcher, req *larkevent.EventReq, completed *sync.WaitGroup) {
	defer completed.Done()
	resp := dispatcher.Handle(context.Background(), req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("concurrent status=%d", resp.StatusCode)
	}
}
