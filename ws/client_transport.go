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
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("websocket endpoint is invalid")
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
		return nil, dialErr
	}
	if socket == nil || resp == nil || resp.StatusCode != http.StatusSwitchingProtocols {
		if socket != nil {
			c.closeSocket(run.ctx, socket, safeEndpoint(rawEndpoint), "discard handshake candidate")
		}
		if resp != nil {
			return nil, parseHandshakeError(resp)
		}
		return nil, errors.New("websocket handshake response is invalid")
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
		return "", nil, NewClientError(larkcore.ErrCodeAppSecretAndClientAssertionEmpty, "appSecret and clientAssertionProvider cannot be nil")
	}

	if c.clientAssertionProvider != nil {
		aud, extractErr := extractAudFromWSURL(c.domain)
		if extractErr != nil {
			return "", nil, extractErr
		}
		clientAssertionToken, retrieveErr := c.clientAssertionProvider.RetrieveToken(ctx, aud)
		if retrieveErr != nil {
			return "", nil, retrieveErr
		}
		if clientAssertionToken == nil || clientAssertionToken.Value == "" {
			return "", nil, NewClientError(larkcore.ErrCodeClientAssertionTokenEmpty, "client assertion token is empty")
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
		return "", nil, marshalErr
	}

	req, requestErr := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(bs))
	if requestErr != nil {
		return "", nil, requestErr
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
		return "", nil, requestErr
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn(ctx, "websocket bootstrap response close failed")
		}
	}()
	respBody, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return "", nil, readErr
	}
	if resp.StatusCode != http.StatusOK {
		message := "system busy"
		bootstrapResponse := &EndpointResp{}
		if err := json.Unmarshal(respBody, bootstrapResponse); err == nil && bootstrapResponse.Msg != "" {
			message = bootstrapResponse.Msg
		}
		return "", nil, NewServerError(resp.StatusCode, message)
	}

	endpointResp := &EndpointResp{}
	if unmarshalErr := json.Unmarshal(respBody, endpointResp); unmarshalErr != nil {
		return "", nil, unmarshalErr
	}

	switch endpointResp.Code {
	case OK:
	case SystemBusy:
		return "", nil, NewServerError(endpointResp.Code, "system busy")
	case InternalError:
		return "", nil, NewServerError(endpointResp.Code, endpointResp.Msg)
	default:
		return "", nil, NewClientError(endpointResp.Code, endpointResp.Msg)
	}

	endpoint := endpointResp.Data
	if endpoint == nil || endpoint.Url == "" {
		return "", nil, NewServerError(http.StatusInternalServerError, "endpoint is null")
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
		return errors.New("websocket handshake response is invalid")
	}
	code, parseErr := strconv.Atoi(resp.Header.Get(HeaderHandshakeStatus))
	if parseErr != nil || code == 0 {
		code = resp.StatusCode
	}
	message := resp.Header.Get(HeaderHandshakeMsg)
	switch code {
	case AuthFailed:
		authCode, authErr := strconv.Atoi(resp.Header.Get(HeaderHandshakeAuthErrCode))
		if authErr != nil {
			authCode = 0
		}
		if authCode == ExceedConnLimit {
			return NewClientError(code, message)
		}
		return NewServerError(code, message)
	case Forbidden:
		return NewClientError(code, message)
	default:
		return NewServerError(code, message)
	}
}
