package normalize

import (
	"encoding/json"
	"testing"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func TestParseMessage(t *testing.T) {
	msgType := "text"
	content := `{"text":"hello world @_all"}`
	msgId := "om_123"
	chatId := "oc_456"
	userId := "ou_789"
	mentionUserId := "ou_abc"
	mentionOpenId := "ou_open_abc"
	mentionName := "Alice"
	mentionKey := "@_user_1"
	allKey := "@_all"

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
						Key: &mentionKey,
						Id: &larkim.UserId{
							UserId: &mentionUserId,
							OpenId: &mentionOpenId,
						},
						Name: &mentionName,
					},
					{
						Key: &allKey,
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
	if norm.Content != "hello world @_all" {
		t.Errorf("expected Content 'hello world @_all', got %s", norm.Content)
	}
	if norm.RawContentType != "text" {
		t.Errorf("expected RawContentType 'text', got %s", norm.RawContentType)
	}
	if !norm.MentionAll {
		t.Errorf("expected MentionAll to be true")
	}
	if len(norm.Mentions) != 2 {
		t.Fatalf("expected 2 mentions, got %d", len(norm.Mentions))
	}

	// Test first mention (Alice)
	if norm.Mentions[0].UserID != mentionUserId {
		t.Errorf("expected mention UserID %s, got %s", mentionUserId, norm.Mentions[0].UserID)
	}
	if norm.Mentions[0].OpenID != mentionOpenId {
		t.Errorf("expected mention OpenID %s, got %s", mentionOpenId, norm.Mentions[0].OpenID)
	}
	if norm.Mentions[0].Key != mentionKey {
		t.Errorf("expected mention Key %s, got %s", mentionKey, norm.Mentions[0].Key)
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

func TestParseReaction(t *testing.T) {
	eventId := "event_123"
	createTime := "1610000000000"
	msgId := "om_123"
	userId := "ou_789"
	reactionType := "SMILE"

	eventPayload := &larkim.P2MessageReactionCreatedV1{
		EventV2Base: &larkevent.EventV2Base{
			Header: &larkevent.EventHeader{
				EventID:    eventId,
				CreateTime: createTime,
			},
		},
		Event: &larkim.P2MessageReactionCreatedV1Data{
			MessageId: &msgId,
			UserId: &larkim.UserId{
				UserId: &userId,
			},
			ReactionType: &larkim.Emoji{
				EmojiType: &reactionType,
			},
			ActionTime: &createTime,
		},
	}

	norm := ParseReaction(eventPayload)
	if norm == nil {
		t.Fatal("expected norm to not be nil")
	}

	if norm.EventID != eventId {
		t.Errorf("expected EventID %s, got %s", eventId, norm.EventID)
	}
	if norm.MessageID != msgId {
		t.Errorf("expected MessageID %s, got %s", msgId, norm.MessageID)
	}
	if norm.UserID != userId {
		t.Errorf("expected UserID %s, got %s", userId, norm.UserID)
	}
	if norm.ReactionType != reactionType {
		t.Errorf("expected ReactionType %s, got %s", reactionType, norm.ReactionType)
	}
	if norm.Action != "add" {
		t.Errorf("expected Action add, got %s", norm.Action)
	}
}
