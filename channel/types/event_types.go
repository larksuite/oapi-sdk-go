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

// CommentEvent represents a comment/reply on a message or a thread.
type CommentEvent struct {
	EventID      string      `json:"event_id"`
	MessageID    string      `json:"message_id"`
	ParentID     string      `json:"parent_id"`
	ChatID       string      `json:"chat_id"`
	UserID       string      `json:"user_id"`
	Content      string      `json:"content"`
	CreateTimeMs int64       `json:"create_time_ms"`
	RawEvent     interface{} `json:"raw_event"`
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
