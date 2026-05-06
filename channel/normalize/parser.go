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

// ParseReaction normalizes a message reaction event (Created or Deleted).
func ParseReaction(event interface{}) *types.ReactionEvent {
	norm := &types.ReactionEvent{
		RawEvent: event,
	}

	switch ev := event.(type) {
	case *larkim.P2MessageReactionCreatedV1:
		if ev.EventV2Base != nil && ev.EventV2Base.Header != nil {
			norm.EventID = ev.EventV2Base.Header.EventID
			if ev.EventV2Base.Header.CreateTime != "" {
				if ms, err := strconv.ParseInt(ev.EventV2Base.Header.CreateTime, 10, 64); err == nil {
					norm.CreateTimeMs = ms
				}
			}
		}
		if ev.Event != nil {
			if ev.Event.MessageId != nil {
				norm.MessageID = *ev.Event.MessageId
			}
			if ev.Event.ReactionType != nil && ev.Event.ReactionType.EmojiType != nil {
				norm.ReactionType = *ev.Event.ReactionType.EmojiType
			}
			if ev.Event.UserId != nil {
				if ev.Event.UserId.UserId != nil {
					norm.UserID = *ev.Event.UserId.UserId
				} else if ev.Event.UserId.OpenId != nil {
					norm.UserID = *ev.Event.UserId.OpenId
				}
			}
			norm.Action = "add"
		}

	case *larkim.P2MessageReactionDeletedV1:
		if ev.EventV2Base != nil && ev.EventV2Base.Header != nil {
			norm.EventID = ev.EventV2Base.Header.EventID
			if ev.EventV2Base.Header.CreateTime != "" {
				if ms, err := strconv.ParseInt(ev.EventV2Base.Header.CreateTime, 10, 64); err == nil {
					norm.CreateTimeMs = ms
				}
			}
		}
		if ev.Event != nil {
			if ev.Event.MessageId != nil {
				norm.MessageID = *ev.Event.MessageId
			}
			if ev.Event.ReactionType != nil && ev.Event.ReactionType.EmojiType != nil {
				norm.ReactionType = *ev.Event.ReactionType.EmojiType
			}
			if ev.Event.UserId != nil {
				if ev.Event.UserId.UserId != nil {
					norm.UserID = *ev.Event.UserId.UserId
				} else if ev.Event.UserId.OpenId != nil {
					norm.UserID = *ev.Event.UserId.OpenId
				}
			}
			norm.Action = "remove"
		}
	default:
		return nil
	}

	return norm
}

// ParseComment normalizes a P2MessageReceiveV1 event as a comment/reply.
func ParseComment(event *larkim.P2MessageReceiveV1) *types.CommentEvent {
	if event == nil || event.Event == nil || event.Event.Message == nil {
		return nil
	}

	msg := event.Event.Message
	sender := event.Event.Sender

	norm := &types.CommentEvent{
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
	if msg.ParentId != nil {
		norm.ParentID = *msg.ParentId
	}
	if msg.ChatId != nil {
		norm.ChatID = *msg.ChatId
	}
	if sender != nil && sender.SenderId != nil {
		if sender.SenderId.UserId != nil {
			norm.UserID = *sender.SenderId.UserId
		} else if sender.SenderId.OpenId != nil {
			norm.UserID = *sender.SenderId.OpenId
		}
	}
	if msg.Content != nil && msg.MessageType != nil {
		content, _ := ParseContent(*msg.MessageType, *msg.Content)
		norm.Content = content
	}

	return norm
}

// ParseBotAdded normalizes a P2ChatMemberBotAddedV1 event.
func ParseBotAdded(event *larkim.P2ChatMemberBotAddedV1) *types.BotAddedEvent {
	if event == nil || event.Event == nil {
		return nil
	}

	norm := &types.BotAddedEvent{
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

	if event.Event.ChatId != nil {
		norm.ChatID = *event.Event.ChatId
	}
	if event.Event.Name != nil {
		norm.ChatName = *event.Event.Name
	}
	if event.Event.OperatorId != nil {
		if event.Event.OperatorId.UserId != nil {
			norm.UserID = *event.Event.OperatorId.UserId
		} else if event.Event.OperatorId.OpenId != nil {
			norm.UserID = *event.Event.OperatorId.OpenId
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
