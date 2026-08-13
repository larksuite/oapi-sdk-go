package ws

import (
	"errors"
	"fmt"
)

// ErrClientAlreadyStarted is returned when Start is called while the client is running.
var ErrClientAlreadyStarted = errors.New("ws client is already started")

type ClientError struct {
	Code int
	Msg  string
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func NewClientError(code int, msg string) *ClientError {
	return &ClientError{
		code,
		msg,
	}
}

type ServerError struct {
	Code int
	Msg  string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func NewServerError(code int, msg string) *ServerError {
	return &ServerError{
		code,
		msg,
	}
}
