package normalize

import (
	"encoding/json"
	"testing"

	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func TestParseMessage(t *testing.T) {
	msgType := "text"
	content := `{"text":"hello world"}`
	msgId := "om_123"
	chatId := "oc_456"
	userId := "ou_789"
	mentionUserId := "ou_abc"
	mentionName := "Alice"

	event := &larkim.P2MessageReceiveV1{
		Event: &larkim.P2MessageReceiveV1Data{
			Sender: &larkim.EventSender{
				SenderId: &larkim.UserId{
					UserId: &userId,
				},
			},
			Message: &larkim.EventMessage{
				MessageId:   &msgId,
				ChatId:      &chatId,
				MessageType: &msgType,
				Content:     &content,
				Mentions: []*larkim.MentionEvent{
					{
						Id: &larkim.UserId{
							UserId: &mentionUserId,
						},
						Name: &mentionName,
					},
				},
			},
		},
	}

	norm := ParseMessage(event)
	if norm == nil {
		t.Fatal("expected norm to not be nil")
	}

	if norm.MessageID != msgId {
		t.Errorf("expected MessageID %s, got %s", msgId, norm.MessageID)
	}
	if norm.ChatID != chatId {
		t.Errorf("expected ChatID %s, got %s", chatId, norm.ChatID)
	}
	if norm.UserID != userId {
		t.Errorf("expected UserID %s, got %s", userId, norm.UserID)
	}
	if norm.Content != "hello world" {
		t.Errorf("expected Content 'hello world', got %s", norm.Content)
	}
	if len(norm.Mentions) != 1 {
		t.Fatalf("expected 1 mention, got %d", len(norm.Mentions))
	}
	if norm.Mentions[0].UserID != mentionUserId {
		t.Errorf("expected mention UserID %s, got %s", mentionUserId, norm.Mentions[0].UserID)
	}
}

func TestParseCardAction(t *testing.T) {
	val := map[string]interface{}{"key": "value"}
	event := &callback.CardActionTriggerEvent{
		Event: &callback.CardActionTriggerRequest{
			Context: &callback.Context{
				OpenMessageID: "om_123",
				OpenChatID:    "oc_456",
			},
			Operator: &callback.Operator{
				OpenID: "ou_789",
			},
			Action: &callback.CallBackAction{
				Value: val,
			},
		},
	}

	norm := ParseCardAction(event)
	if norm == nil {
		t.Fatal("expected norm to not be nil")
	}

	if norm.MessageID != "om_123" {
		t.Errorf("expected MessageID om_123, got %s", norm.MessageID)
	}
	if norm.ChatID != "oc_456" {
		t.Errorf("expected ChatID oc_456, got %s", norm.ChatID)
	}
	if norm.UserID != "ou_789" {
		t.Errorf("expected UserID ou_789, got %s", norm.UserID)
	}

	expectedContent, _ := json.Marshal(val)
	if norm.Content != string(expectedContent) {
		t.Errorf("expected Content %s, got %s", string(expectedContent), norm.Content)
	}
}
