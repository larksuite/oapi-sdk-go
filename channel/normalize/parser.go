package normalize

import (
	"encoding/json"
	"strconv"

	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// ParseMessage normalizes a P2MessageReceiveV1 event.
func ParseMessage(event *larkim.P2MessageReceiveV1) *types.NormalizedMessage {
	if event == nil || event.Event == nil || event.Event.Message == nil {
		return nil
	}

	msg := event.Event.Message
	sender := event.Event.Sender

	norm := &types.NormalizedMessage{
		RawEvent: event,
	}

	if event.EventV2Base != nil && event.EventV2Base.Header != nil {
		norm.EventID = event.EventV2Base.Header.EventID
		if event.EventV2Base.Header.CreateTime != "" {
			if ms, err := strconv.ParseInt(event.EventV2Base.Header.CreateTime, 10, 64); err == nil {
				norm.CreateTimeMs = ms
			}
		}
	}

	if msg.MessageId != nil {
		norm.MessageID = *msg.MessageId
	}
	if msg.ChatId != nil {
		norm.ChatID = *msg.ChatId
	}
	if msg.ChatType != nil {
		norm.ChatType = *msg.ChatType
	}
	if sender != nil && sender.SenderId != nil {
		if sender.SenderId.UserId != nil {
			norm.UserID = *sender.SenderId.UserId
		} else if sender.SenderId.OpenId != nil {
			norm.UserID = *sender.SenderId.OpenId
		}
	}
	if msg.Mentions != nil {
		for _, m := range msg.Mentions {
			mention := types.Mention{}
			if m.Id != nil {
				if m.Id.UserId != nil {
					mention.UserID = *m.Id.UserId
				} else if m.Id.OpenId != nil {
					mention.UserID = *m.Id.OpenId
				}
			}
			if m.Name != nil {
				mention.Name = *m.Name
			}
			norm.Mentions = append(norm.Mentions, mention)

			if m.Key != nil && *m.Key == "@all" {
				norm.MentionAll = true
			}
		}
	}

	// Extract content and resources
	if msg.Content != nil && msg.MessageType != nil {
		content, resources := ParseContent(*msg.MessageType, *msg.Content)
		if content != "" {
			norm.Content = content
		}
		if len(resources) > 0 {
			norm.Resources = append(norm.Resources, resources...)
		}
	}

	return norm
}

// ParseCardAction normalizes a CardActionTriggerEvent.
func ParseCardAction(event *callback.CardActionTriggerEvent) *types.NormalizedMessage {
	if event == nil || event.Event == nil {
		return nil
	}

	req := event.Event
	norm := &types.NormalizedMessage{
		RawEvent: event,
	}

	if event.EventV2Base != nil && event.EventV2Base.Header != nil {
		norm.EventID = event.EventV2Base.Header.EventID
	}

	if req.Context != nil {
		norm.MessageID = req.Context.OpenMessageID
		norm.ChatID = req.Context.OpenChatID
	}

	if req.Operator != nil {
		if req.Operator.UserID != nil {
			norm.UserID = *req.Operator.UserID
		} else {
			norm.UserID = req.Operator.OpenID
		}
	}

	if req.Action != nil && req.Action.Value != nil {
		if valBytes, err := json.Marshal(req.Action.Value); err == nil {
			norm.Content = string(valBytes)
		}
	}

	return norm
}
