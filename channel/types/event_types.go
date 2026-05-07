package types

// ReactionEvent represents a reaction added or removed from a message.
type ReactionEvent struct {
	EventID      string      `json:"event_id"`
	MessageID    string      `json:"message_id"`
	ReactionType string      `json:"reaction_type"` // e.g. "SMILE"
	UserID       string      `json:"user_id"`
	Action       string      `json:"action"` // "add" or "remove"
	CreateTimeMs int64       `json:"create_time_ms"`
	RawEvent     interface{} `json:"raw_event"`
}

// OperatorInfo represents the user who triggered the event.
type OperatorInfo struct {
	OpenID  string `json:"open_id"`
	UserID  string `json:"user_id,omitempty"`
	UnionID string `json:"union_id,omitempty"`
}

// CommentEvent represents a comment/reply on a drive document.
type CommentEvent struct {
	EventID      string       `json:"event_id"`
	CommentID    string       `json:"comment_id"`
	FileToken    string       `json:"file_token"`
	FileType     string       `json:"file_type"`
	Operator     OperatorInfo `json:"operator"`
	ReplyID      string       `json:"reply_id"`
	MentionedBot bool         `json:"mentioned_bot"`
	Timestamp    int64        `json:"timestamp"`
	RawEvent     interface{}  `json:"raw_event"`
}

// BotAddedEvent represents an event when the bot is added to a chat.
type BotAddedEvent struct {
	EventID      string      `json:"event_id"`
	ChatID       string      `json:"chat_id"`
	ChatName     string      `json:"chat_name"`
	UserID       string      `json:"user_id"` // User who added the bot
	CreateTimeMs int64       `json:"create_time_ms"`
	RawEvent     interface{} `json:"raw_event"`
}

// RejectEvent represents a message rejected by safety policies.
type RejectEvent struct {
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
	SenderID  string `json:"sender_id"`
	Reason    string `json:"reason"`
}
