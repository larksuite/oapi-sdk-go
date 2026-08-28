package ws

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	errNilContext       = errors.New("websocket start context is nil")
	errClientRunning    = errors.New("websocket client is already running")
	errClientTerminal   = errors.New("websocket client cannot be restarted")
	errConnectionClosed = errors.New("websocket connection is closed")
)

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

type lifecycleError struct {
	operation  string
	category   string
	statusCode int
	cause      error
	unwrap     bool
}

func (e *lifecycleError) Error() string {
	if e.statusCode != 0 {
		return fmt.Sprintf("websocket %s failed: %s (status=%d)", e.operation, e.category, e.statusCode)
	}
	return fmt.Sprintf("websocket %s failed: %s", e.operation, e.category)
}

func (e *lifecycleError) Unwrap() error {
	if !e.unwrap {
		return nil
	}
	return e.cause
}

func newLifecycleError(operation, category string, statusCode int, cause error, unwrap bool) error {
	return &lifecycleError{
		operation:  operation,
		category:   category,
		statusCode: statusCode,
		cause:      cause,
		unwrap:     unwrap,
	}
}

func safeEndpoint(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "websocket endpoint"
	}
	return (&url.URL{Scheme: u.Scheme, Host: u.Host}).String()
}
