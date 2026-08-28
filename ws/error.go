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

func safeEndpoint(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "websocket endpoint"
	}
	return (&url.URL{Scheme: u.Scheme, Host: u.Host}).String()
}
