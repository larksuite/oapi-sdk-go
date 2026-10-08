package larkcore

import (
	"context"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

// SSEOptions controls network waits and the bounded event parser. Zero uses defaults.
type SSEOptions struct {
	OpenTimeout     time.Duration
	IdleReadTimeout time.Duration
	TotalTimeout    time.Duration
	MaxEventBytes   int
	LastEventID     string
}

// WithSSEOptions configures a streaming call without affecting ordinary requests.
func WithSSEOptions(options SSEOptions) RequestOptionFunc {
	return func(option *RequestOption) { option.SSE = &options }
}

// SSEStream owns one connection. Next has one reader; Close may run concurrently.
type SSEStream struct {
	ctx     context.Context
	mu      sync.Mutex
	body    io.ReadCloser
	cancel  context.CancelFunc
	parser  *sseParser
	headers http.Header
	event   *SSEEvent
	err     error
	closed  bool
	reading bool
	cursor  string
	retry   *int64
}

type sseReader struct {
	ctx    context.Context
	body   io.Reader
	cancel context.CancelFunc
	idle   time.Duration
}

func (r *sseReader) Read(buffer []byte) (int, error) {
	var mu sync.Mutex
	finished := false
	timedOut := false
	timer := time.AfterFunc(r.idle, func() {
		mu.Lock()
		defer mu.Unlock()
		if !finished {
			timedOut = true
			r.cancel()
		}
	})
	n, err := r.body.Read(buffer)
	mu.Lock()
	finished = true
	idleExpired := timedOut
	mu.Unlock()
	timer.Stop()
	if idleExpired {
		return n, context.DeadlineExceeded
	}
	if err != nil && r.ctx.Err() != nil {
		err = r.ctx.Err()
	}
	return n, err
}

// Next waits for one complete event, without waiting for the HTTP response to end.
func (s *SSEStream) Next() bool {
	s.mu.Lock()
	if s.reading {
		s.mu.Unlock()
		return false
	}
	if s.closed {
		s.mu.Unlock()
		return false
	}
	s.reading = true
	s.mu.Unlock()
	var event *SSEEvent
	err := s.ctx.Err()
	if err == nil {
		event, err = s.parser.next()
	}
	if contextErr := s.ctx.Err(); contextErr != nil && err == nil {
		event, err = nil, contextErr
	}
	s.mu.Lock()
	s.reading = false
	s.event = event
	s.cursor = s.parser.cursor
	s.retry = s.parser.retry
	if err != nil && !errors.Is(err, io.EOF) && !s.closed {
		s.err = err
	}
	done := err != nil || s.closed
	s.mu.Unlock()
	if done {
		if closeErr := s.Close(); closeErr != nil {
			s.mu.Lock()
			if s.err == nil {
				s.err = closeErr
			}
			s.mu.Unlock()
		}
	}
	return !done
}

// Event returns the event delivered by the last successful Next call.
func (s *SSEStream) Event() *SSEEvent { s.mu.Lock(); defer s.mu.Unlock(); return s.event }

// Err reports the terminal read error; a normal EOF is not an error.
func (s *SSEStream) Err() error { s.mu.Lock(); defer s.mu.Unlock(); return s.err }

// Headers returns a copy of the response headers.
func (s *SSEStream) Headers() http.Header { return s.headers.Clone() }

// LastEventID returns the protocol cursor, which is not a business acknowledgement.
func (s *SSEStream) LastEventID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cursor
}

// RetryDelay returns the server's suggested delay; the SDK never reconnects automatically.
func (s *SSEStream) RetryDelay() *int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.retry == nil {
		return nil
	}
	value := *s.retry
	return &value
}

// Close cancels pending reads and releases the response body once.
func (s *SSEStream) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()
	s.cancel()
	return s.body.Close()
}
