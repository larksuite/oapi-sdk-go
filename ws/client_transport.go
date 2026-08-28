package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	ws "github.com/gorilla/websocket"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

var bootstrapHTTPClient = http.DefaultClient

// establishConnectionOnce performs exactly one bootstrap and websocket dial.
// Retry policy belongs to the lifecycle coordinator, not the transport layer.
func (c *Client) establishConnectionOnce(run *clientRun) (*clientConn, error) {
	rawEndpoint, conf, err := c.fetchEndpoint(run.ctx)
	if err != nil {
		return nil, err
	}
	if !c.applyConfig(run, conf) {
		return nil, errConnectionClosed
	}

	u, err := url.Parse(rawEndpoint)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, newLifecycleError("parse endpoint", "invalid endpoint", 0, nil, false)
	}

	socket, resp, dialErr := ws.DefaultDialer.DialContext(run.ctx, rawEndpoint, nil)
	if resp != nil && resp.Body != nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				c.logger.Warn(run.ctx, "websocket handshake response close failed")
			}
		}()
	}
	if dialErr != nil {
		if socket != nil {
			c.closeSocket(run.ctx, socket, safeEndpoint(rawEndpoint), "discard dial candidate")
		}
		if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
			return nil, parseHandshakeError(resp)
		}
		return nil, newLifecycleError("dial websocket", "transport", 0, nil, false)
	}
	if socket == nil || resp == nil || resp.StatusCode != http.StatusSwitchingProtocols {
		if socket != nil {
			c.closeSocket(run.ctx, socket, safeEndpoint(rawEndpoint), "discard handshake candidate")
		}
		if resp != nil {
			return nil, parseHandshakeError(resp)
		}
		return nil, newLifecycleError("handshake websocket", "invalid response", 0, nil, false)
	}

	return &clientConn{
		socket:       socket,
		readResult:   make(chan error, 1),
		connID:       u.Query().Get(DeviceID),
		serviceID:    u.Query().Get(ServiceID),
		safeEndpoint: safeEndpoint(rawEndpoint),
	}, nil
}

func (c *Client) fetchEndpoint(ctx context.Context) (endpointURL string, conf *ClientConfig, err error) {
	requestURL := strings.TrimRight(c.domain, "/") + GenEndpointUri
	body := &BootstrapRequest{AppID: c.appID}
	headers := make(http.Header)

	if c.clientAssertionProvider == nil && c.appSecret == "" {
		cause := NewClientError(larkcore.ErrCodeAppSecretAndClientAssertionEmpty, "credentials are required")
		return "", nil, newLifecycleError("bootstrap credential", "credential unavailable", 0, cause, true)
	}

	if c.clientAssertionProvider != nil {
		aud, extractErr := extractAudFromWSURL(c.domain)
		if extractErr != nil {
			return "", nil, newLifecycleError("bootstrap credential", "invalid audience", 0, nil, false)
		}
		clientAssertionToken, retrieveErr := c.clientAssertionProvider.RetrieveToken(ctx, aud)
		if retrieveErr != nil {
			var clientErr *ClientError
			if errors.As(retrieveErr, &clientErr) {
				cause := NewClientError(clientErr.Code, "client rejected")
				return "", nil, newLifecycleError("bootstrap credential", "credential unavailable", 0, cause, true)
			}
			return "", nil, newLifecycleError("bootstrap credential", "token retrieval", 0, nil, false)
		}
		if clientAssertionToken == nil || clientAssertionToken.Value == "" {
			cause := NewClientError(larkcore.ErrCodeClientAssertionTokenEmpty, "client assertion token is empty")
			return "", nil, newLifecycleError("bootstrap credential", "credential unavailable", 0, cause, true)
		}
		body.ClientAssertion = clientAssertionToken.Value
		body.AppSecret = ""
		if clientAssertionToken.TargetInfo != nil {
			requestURL = buildWSProxyURL(clientAssertionToken.TargetInfo.TargetService, clientAssertionToken.TargetInfo.TargetPrefix, GenEndpointUri)
			headers.Set(larkcore.HeaderXTargetService, aud)
		}
	} else {
		body.AppSecret = c.appSecret
	}

	bs, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return "", nil, newLifecycleError("bootstrap", "request build", 0, nil, false)
	}

	req, requestErr := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(bs))
	if requestErr != nil {
		return "", nil, newLifecycleError("bootstrap", "request build", 0, nil, false)
	}

	req.Header.Add("locale", "zh")
	req.Header.Add("Content-Type", "application/json")
	for k, values := range c.headers {
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}
	for k, values := range headers {
		req.Header.Del(k)
		for _, value := range values {
			req.Header.Add(k, value)
		}
	}
	req.Header.Set("User-Agent", larkcore.UserAgent(c.source))
	resp, requestErr := bootstrapHTTPClient.Do(req)
	if requestErr != nil {
		return "", nil, newLifecycleError("bootstrap", "transport", 0, nil, false)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn(ctx, "websocket bootstrap response close failed")
		}
	}()
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", nil, newLifecycleError("bootstrap", "invalid response", resp.StatusCode, nil, false)
	}
	if resp.StatusCode != http.StatusOK {
		return "", nil, newLifecycleError("bootstrap", "http status", resp.StatusCode, nil, false)
	}

	endpointResp := &EndpointResp{}
	if unmarshalErr := json.Unmarshal(respBody, endpointResp); unmarshalErr != nil {
		return "", nil, newLifecycleError("bootstrap", "invalid response", resp.StatusCode, nil, false)
	}

	switch endpointResp.Code {
	case OK:
	case SystemBusy:
		cause := NewServerError(endpointResp.Code, "server unavailable")
		return "", nil, newLifecycleError("bootstrap", "server unavailable", endpointResp.Code, cause, true)
	case InternalError:
		cause := NewServerError(endpointResp.Code, "server unavailable")
		return "", nil, newLifecycleError("bootstrap", "server unavailable", endpointResp.Code, cause, true)
	default:
		cause := NewClientError(endpointResp.Code, "client rejected")
		return "", nil, newLifecycleError("bootstrap", "client rejected", endpointResp.Code, cause, true)
	}

	endpoint := endpointResp.Data
	if endpoint == nil || endpoint.Url == "" {
		cause := NewServerError(http.StatusInternalServerError, "endpoint unavailable")
		return "", nil, newLifecycleError("bootstrap", "invalid response", http.StatusInternalServerError, cause, true)
	}

	return endpoint.Url, endpoint.ClientConfig, nil
}

func extractAudFromWSURL(rawURL string) (string, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Host != "" {
		return parsedURL.Host, nil
	}
	return "", fmt.Errorf("invalid url: %s", rawURL)
}

func buildWSProxyURL(targetService, targetPrefix, apiPath string) string {
	if !strings.Contains(targetService, "://") {
		targetService = "https://" + targetService
	}
	return targetService + targetPrefix + apiPath
}

func parseHandshakeError(resp *http.Response) error {
	if resp == nil {
		return newLifecycleError("handshake websocket", "invalid response", 0, nil, false)
	}
	code, parseErr := strconv.Atoi(resp.Header.Get(HeaderHandshakeStatus))
	if parseErr != nil || code == 0 {
		code = resp.StatusCode
	}
	switch code {
	case AuthFailed:
		authCode, authErr := strconv.Atoi(resp.Header.Get(HeaderHandshakeAuthErrCode))
		if authErr != nil {
			authCode = 0
		}
		if authCode == ExceedConnLimit {
			cause := NewClientError(code, "client rejected")
			return newLifecycleError("handshake websocket", "client rejected", resp.StatusCode, cause, true)
		}
		cause := NewServerError(code, "server rejected")
		return newLifecycleError("handshake websocket", "server rejected", resp.StatusCode, cause, true)
	case Forbidden:
		cause := NewClientError(code, "client rejected")
		return newLifecycleError("handshake websocket", "client rejected", resp.StatusCode, cause, true)
	default:
		cause := NewServerError(code, "server rejected")
		return newLifecycleError("handshake websocket", "server rejected", resp.StatusCode, cause, true)
	}
}
