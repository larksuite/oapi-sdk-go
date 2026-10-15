package ws

import (
	"net/http"
	"net/url"
	"strings"
)

// validateConnectionHeaders keeps handshake protocol fields owned by the SDK.
// Other header names and values are passed through without SDK format checks.
func validateConnectionHeaders(header http.Header) error {
	for name := range header {
		if isReservedWebSocketHeader(name) {
			return errInvalidConnectionHeaders
		}
	}
	return nil
}

func isReservedWebSocketHeader(name string) bool {
	if strings.EqualFold(name, "Host") || strings.EqualFold(name, "Connection") || strings.EqualFold(name, "Upgrade") {
		return true
	}
	const webSocketHeaderPrefix = "Sec-WebSocket-"
	return len(name) >= len(webSocketHeaderPrefix) && strings.EqualFold(name[:len(webSocketHeaderPrefix)], webSocketHeaderPrefix)
}

var (
	errInvalidConnectionHost     = NewClientError(http.StatusBadRequest, "invalid websocket connection host")
	errInvalidConnectionHeaders  = NewClientError(http.StatusBadRequest, "invalid websocket connection headers")
	errInvalidConnectionEndpoint = NewClientError(http.StatusBadRequest, "invalid websocket connection endpoint")
)

func (c *Client) validateConnectionConfig() error {
	if c.connectionHost != nil {
		if err := validateConnectionHost(*c.connectionHost); err != nil {
			return err
		}
	}
	return validateConnectionHeaders(c.connectionHeaders)
}

func validateConnectionHost(host string) error {
	if strings.TrimSpace(host) == "" || strings.ContainsAny(host, "/?") {
		return errInvalidConnectionHost
	}
	return nil
}

func buildConnectionURL(endpoint string, connectionHost *string) (*url.URL, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil || strings.Contains(endpoint, "#") {
		return nil, errInvalidConnectionEndpoint
	}
	if parsed.User != nil || parsed.Hostname() == "" {
		return nil, errInvalidConnectionEndpoint
	}
	if parsed.Scheme != "ws" && parsed.Scheme != "wss" {
		return nil, errInvalidConnectionEndpoint
	}
	if connectionHost != nil {
		parsed.Host = *connectionHost
	}
	return parsed, nil
}
