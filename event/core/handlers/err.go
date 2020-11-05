package handlers

import "fmt"

type NotFoundHandlerErr struct {
	EventType string
}

func newNotHandlerErr(eventType string) *NotFoundHandlerErr {
	return &NotFoundHandlerErr{EventType: eventType}
}

func (e NotFoundHandlerErr) Error() string {
	return fmt.Sprintf("event type:%s, not found handler", e.EventType)
}
