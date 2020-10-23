package handlers

import "fmt"

type NotHandlerErr struct {
	EventType string
}

func newNotHandlerErr(eventType string) *NotHandlerErr {
	return &NotHandlerErr{EventType: eventType}
}

func (e NotHandlerErr) Error() string {
	return fmt.Sprintf("event type:%s, not find handler", e.EventType)
}
