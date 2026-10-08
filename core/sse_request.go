package larkcore

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"
)

func normalizeSSEOptions(value *SSEOptions) (SSEOptions, error) {
	options := SSEOptions{}
	if value != nil {
		options = *value
	}
	if options.OpenTimeout < 0 || options.IdleReadTimeout < 0 || options.TotalTimeout < 0 || options.MaxEventBytes < 0 {
		return options, fmt.Errorf("SSE options must not be negative")
	}
	if strings.ContainsAny(options.LastEventID, "\r\n\x00") || len(options.LastEventID) > 8192 {
		return options, fmt.Errorf("invalid SSE last event ID")
	}
	if options.OpenTimeout == 0 {
		options.OpenTimeout = 10 * time.Second
	}
	if options.IdleReadTimeout == 0 {
		options.IdleReadTimeout = 60 * time.Second
	}
	if options.MaxEventBytes == 0 {
		options.MaxEventBytes = 1024 * 1024
	}
	return options, nil
}

func prepareSSERequest(ctx context.Context, req *ApiReq, config *Config, option *RequestOption) (*http.Request, error) {
	copyReq := *req
	if copyReq.Body != nil && config.Serializable == nil {
		return nil, fmt.Errorf("SSE request body requires a serializer")
	}
	if len(copyReq.SupportedAccessTokenTypes) == 0 {
		copyReq.SupportedAccessTokenTypes = []AccessTokenType{AccessTokenTypeNone}
	}
	if err := validateTokenType(copyReq.SupportedAccessTokenTypes, option); err != nil {
		return nil, fmt.Errorf("SSE token type: %w", err)
	}
	tokenType, err := determineTokenType(copyReq.SupportedAccessTokenTypes, option, config)
	if err != nil {
		return nil, fmt.Errorf("SSE token selection: %w", err)
	}
	if err = validate(config, option, tokenType); err != nil {
		return nil, fmt.Errorf("SSE request validation: %w", err)
	}
	raw, err := reqTranslator.translate(ctx, &copyReq, tokenType, config, option)
	if err != nil {
		return nil, fmt.Errorf("SSE request translation: %w", err)
	}
	raw.GetBody = nil // A stream request must not be replayed by the HTTP transport.
	return raw, nil
}

func sseHTTPClient(client HttpClient) HttpClient {
	if standard, ok := client.(*http.Client); ok {
		clone := *standard
		clone.Timeout = 0
		clone.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
		return &clone
	}
	if client == nil {
		return &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }}
	}
	return client
}

// RequestSSE uses existing authentication, then opens one non-retrying SSE response.
func RequestSSE(ctx context.Context, req *ApiReq, config *Config, funcs ...RequestOptionFunc) (*SSEStream, error) {
	if ctx == nil || req == nil || config == nil {
		return nil, fmt.Errorf("SSE context, request and config are required")
	}
	option := &RequestOption{}
	for _, fn := range funcs {
		if fn == nil {
			return nil, fmt.Errorf("nil SSE request option")
		}
		fn(option)
	}
	options, err := normalizeSSEOptions(option.SSE)
	if err != nil {
		return nil, err
	}
	streamCtx, cancel := context.WithCancel(ctx)
	if options.TotalTimeout > 0 {
		cancel()
		streamCtx, cancel = context.WithTimeout(ctx, options.TotalTimeout)
	}
	raw, err := prepareSSERequest(streamCtx, req, config, option)
	if err != nil {
		cancel()
		return nil, err
	}
	raw.Header.Set("Accept", "text/event-stream")
	if options.LastEventID != "" {
		if raw.Header.Get("Last-Event-ID") != "" {
			cancel()
			return nil, fmt.Errorf("conflicting SSE last event ID")
		}
		raw.Header.Set("Last-Event-ID", options.LastEventID)
	}
	timer := time.AfterFunc(options.OpenTimeout, cancel)
	response, err := sseHTTPClient(config.HttpClient).Do(raw)
	stopped := timer.Stop()
	if err != nil {
		cancel()
		if !stopped {
			err = context.DeadlineExceeded
		}
		return nil, closeSSEFailure(response, cancel, fmt.Errorf("open SSE response: %w", err))
	}
	if !stopped || streamCtx.Err() != nil {
		return nil, closeSSEFailure(response, cancel, fmt.Errorf("SSE open cancelled or timed out"))
	}
	return newSSEStream(streamCtx, response, cancel, options)
}

func closeSSEFailure(response *http.Response, cancel context.CancelFunc, cause error) error {
	cancel()
	if response != nil && response.Body != nil {
		if err := response.Body.Close(); err != nil {
			return fmt.Errorf("%w (close SSE response: %v)", cause, err)
		}
	}
	return cause
}

func newSSEStream(ctx context.Context, response *http.Response, cancel context.CancelFunc, options SSEOptions) (*SSEStream, error) {
	if response == nil || response.Body == nil {
		return nil, closeSSEFailure(response, cancel, fmt.Errorf("missing SSE response body"))
	}
	media, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if response.StatusCode < 200 || response.StatusCode >= 300 || err != nil || media != "text/event-stream" {
		return nil, closeSSEFailure(response, cancel, fmt.Errorf("invalid SSE response: HTTP %d, expected text/event-stream", response.StatusCode))
	}
	reader := &sseReader{ctx: ctx, body: response.Body, cancel: cancel, idle: options.IdleReadTimeout}
	stream := &SSEStream{ctx: ctx, body: response.Body, cancel: cancel, headers: response.Header.Clone(), cursor: options.LastEventID, parser: newSSEParser(reader, options.MaxEventBytes, options.LastEventID)}
	stream.parser.onControl = func(cursor string, retry *int64) {
		stream.mu.Lock()
		defer stream.mu.Unlock()
		stream.cursor = cursor
		stream.retry = retry
	}
	return stream, nil
}
