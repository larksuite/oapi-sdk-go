package dispatcher

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
)

func makeTokenPayload(t *testing.T, schema string, token interface{}, include bool) []byte {
	t.Helper()
	payload := map[string]interface{}{"schema": "2.0", "header": map[string]interface{}{"event_type": "audit.test"}, "event": map[string]interface{}{}}
	target := payload["header"].(map[string]interface{})
	if schema == "v1" {
		payload = map[string]interface{}{"uuid": "audit-id", "type": "event_callback", "event": map[string]interface{}{"type": "audit.test"}}
		target = payload
	} else if schema == "challenge" {
		payload = map[string]interface{}{"type": "url_verification", "challenge": "fixture"}
		target = payload
	}
	if include {
		target["token"] = token
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal token payload: %v", err)
	}
	return body
}

func TestWebhookVerificationToken(t *testing.T) {
	cases := []struct {
		name    string
		token   interface{}
		include bool
		allowed bool
	}{
		{name: "correct", token: "audit-fixture", include: true, allowed: true},
		{name: "wrong", token: "wrong", include: true},
		{name: "missing"},
		{name: "empty", token: "", include: true},
		{name: "null", include: true},
		{name: "number", token: 123, include: true},
		{name: "object", token: map[string]string{}, include: true},
		{name: "array", token: []string{}, include: true},
	}
	for _, schema := range []string{"v1", "v2", "challenge"} {
		for _, tc := range cases {
			t.Run(schema+"/"+tc.name, func(t *testing.T) {
				calls := 0
				dispatcher := NewEventDispatcher("audit-fixture", "").OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { calls++; return nil })
				dispatcher.InitConfig(larkevent.WithSkipSignVerify(true))
				resp := dispatcher.Handle(context.Background(), &larkevent.EventReq{Header: map[string][]string{}, Body: makeTokenPayload(t, schema, tc.token, tc.include)})
				if (resp.StatusCode == http.StatusOK) != tc.allowed {
					t.Errorf("status=%d allowed=%v", resp.StatusCode, tc.allowed)
				}
				expected := 0
				if tc.allowed && schema != "challenge" {
					expected = 1
				}
				if calls != expected {
					t.Errorf("handler calls=%d want=%d", calls, expected)
				}
			})
		}
	}
}

func TestWebhookPreservesEmptyVerificationTokenConfiguration(t *testing.T) {
	for _, schema := range []string{"v1", "v2", "challenge"} {
		for _, token := range []string{"", "wrong"} {
			t.Run(schema+"/"+token, func(t *testing.T) {
				calls := 0
				dispatcher := NewEventDispatcher("", "").OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { calls++; return nil })
				resp := dispatcher.Handle(context.Background(), &larkevent.EventReq{Header: map[string][]string{}, Body: makeTokenPayload(t, schema, token, true)})
				allowed := schema != "challenge" || token == ""
				expectedCalls := 0
				if schema != "challenge" {
					expectedCalls = 1
				}
				if (resp.StatusCode == http.StatusOK) != allowed || calls != expectedCalls {
					t.Fatalf("legacy behavior changed: status=%d calls=%d", resp.StatusCode, calls)
				}
			})
		}
	}
}

func TestTrustedDispatchDoesNotRequireWebhookToken(t *testing.T) {
	calls := 0
	dispatcher := NewEventDispatcher("", "").OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { calls++; return nil })
	if _, err := dispatcher.Do(context.Background(), makeTokenPayload(t, "v2", nil, false)); err != nil {
		t.Fatalf("trusted dispatch: %v", err)
	}
	if calls != 1 {
		t.Fatalf("handler calls=%d want=1", calls)
	}
}

func TestEncryptedWebhookPreservesSignatureValidation(t *testing.T) {
	payload := map[string]interface{}{}
	if err := json.Unmarshal(makeTokenPayload(t, "v2", nil, false), &payload); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	encrypted, err := larkcore.EncryptedEventMsg(context.Background(), payload, "encryption-fixture")
	if err != nil {
		t.Fatalf("encrypt fixture: %v", err)
	}
	body, err := json.Marshal(larkevent.EventEncryptMsg{Encrypt: encrypted})
	if err != nil {
		t.Fatalf("marshal encrypted fixture: %v", err)
	}
	for _, valid := range []bool{true, false} {
		calls := 0
		dispatcher := NewEventDispatcher("", "encryption-fixture").OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { calls++; return nil })
		signature := "wrong"
		if valid {
			signature = larkevent.Signature("timestamp", "nonce", "encryption-fixture", string(body))
		}
		req := &larkevent.EventReq{Body: body, Header: map[string][]string{larkevent.EventRequestTimestamp: {"timestamp"}, larkevent.EventRequestNonce: {"nonce"}, larkevent.EventSignature: {signature}}}
		resp := dispatcher.Handle(context.Background(), req)
		if (resp.StatusCode == http.StatusOK) != valid || (calls == 1) != valid {
			t.Errorf("encrypted signature valid=%v status=%d calls=%d", valid, resp.StatusCode, calls)
		}
	}
}

func TestWebhookV2DoesNotFallBackToTopLevelToken(t *testing.T) {
	calls := 0
	dispatcher := NewEventDispatcher("audit-fixture", "").OnCustomizedEvent("audit.test", func(ctx context.Context, req *larkevent.EventReq) error { calls++; return nil })
	for _, body := range []string{
		`{"schema":"2.0","token":"audit-fixture","header":{"event_type":"audit.test"},"event":{}}`,
		`{"schema":"2.0","token":"audit-fixture","event":{"type":"audit.test"}}`,
	} {
		resp := dispatcher.Handle(context.Background(), &larkevent.EventReq{Header: map[string][]string{}, Body: []byte(body)})
		if resp.StatusCode == http.StatusOK || calls != 0 {
			t.Fatalf("v2 top-level token accepted: status=%d calls=%d", resp.StatusCode, calls)
		}
	}
}

func TestWebhookCallbackVerifiesTokenBeforeHandler(t *testing.T) {
	for _, valid := range []bool{true, false} {
		calls := 0
		dispatcher := NewEventDispatcher("audit-fixture", "").OnP2CardActionTrigger(func(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
			calls++
			return &callback.CardActionTriggerResponse{}, nil
		})
		payload := map[string]interface{}{}
		token := "wrong"
		if valid {
			token = "audit-fixture"
		}
		if err := json.Unmarshal(makeTokenPayload(t, "v2", token, true), &payload); err != nil {
			t.Fatalf("decode callback fixture: %v", err)
		}
		payload["header"].(map[string]interface{})["event_type"] = "card.action.trigger"
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal callback fixture: %v", err)
		}
		resp := dispatcher.Handle(context.Background(), &larkevent.EventReq{Header: map[string][]string{}, Body: body})
		if (resp.StatusCode == http.StatusOK) != valid || (calls == 1) != valid {
			t.Errorf("callback valid=%v status=%d calls=%d", valid, resp.StatusCode, calls)
		}
	}
}
