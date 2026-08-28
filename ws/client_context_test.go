package ws

import (
	"context"
	"errors"
	"testing"
	"time"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type lifecycleContextKey struct{}

type contextInspectingAssertionProvider struct {
	observed chan contextObservation
	waitDone bool
}

type contextObservation struct {
	value       interface{}
	deadline    time.Time
	hasDeadline bool
	err         error
}

func (p *contextInspectingAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	if p.waitDone {
		<-ctx.Done()
	}
	deadline, hasDeadline := ctx.Deadline()
	p.observed <- contextObservation{
		value:       ctx.Value(lifecycleContextKey{}),
		deadline:    deadline,
		hasDeadline: hasDeadline,
		err:         ctx.Err(),
	}
	return nil, errors.New("context inspection complete")
}

func TestStartPropagatesContextValueAndDeadlineToAssertionProvider(t *testing.T) {
	deadline := time.Now().Add(time.Minute)
	callerCtx := context.WithValue(context.Background(), lifecycleContextKey{}, "caller-value")
	callerCtx, cancel := context.WithDeadline(callerCtx, deadline)
	defer cancel()

	provider := &contextInspectingAssertionProvider{observed: make(chan contextObservation, 1)}
	client := NewClient("app-id", "", WithClientAssertionProvider(provider), WithAutoReconnect(false))
	if err := client.Start(callerCtx); err == nil {
		t.Fatal("Start returned nil after assertion provider failure")
	}
	observation := <-provider.observed
	if observation.value != "caller-value" {
		t.Fatalf("provider context value = %#v, want caller-value", observation.value)
	}
	if !observation.hasDeadline || !observation.deadline.Equal(deadline) {
		t.Fatalf("provider deadline = %v, %v; want %v, true", observation.deadline, observation.hasDeadline, deadline)
	}
}

func TestStartPropagatesDeadlineErrorToAssertionProvider(t *testing.T) {
	callerCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	provider := &contextInspectingAssertionProvider{
		observed: make(chan contextObservation, 1),
		waitDone: true,
	}
	client := NewClient("app-id", "", WithClientAssertionProvider(provider), WithAutoReconnect(false))
	if err := client.Start(callerCtx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Start returned %v, want context.DeadlineExceeded", err)
	}
	observation := <-provider.observed
	if !errors.Is(observation.err, context.DeadlineExceeded) {
		t.Fatalf("provider context Err = %v, want context.DeadlineExceeded", observation.err)
	}
}

func TestRunContextFollowsCallerCancellation(t *testing.T) {
	callerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewClient("app-id", "app-secret")
	run, err := client.beginRun(callerCtx)
	if err != nil {
		t.Fatalf("beginRun returned error: %v", err)
	}

	cancel()
	select {
	case <-run.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("run context did not follow caller cancellation")
	}
	if !errors.Is(run.ctx.Err(), context.Canceled) {
		t.Fatalf("run context Err = %v, want context.Canceled", run.ctx.Err())
	}
	if finishErr := client.finishRun(run); !errors.Is(finishErr, context.Canceled) {
		t.Fatalf("finishRun returned %v, want context.Canceled", finishErr)
	}
}

func TestConnectionStateRejectsCallerCancellation(t *testing.T) {
	callerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := NewClient("app-id", "app-secret")
	run, err := client.beginRun(callerCtx)
	if err != nil {
		t.Fatalf("beginRun returned error: %v", err)
	}
	conn := &clientConn{
		readResult: make(chan error, 1),
	}
	client.stateMu.Lock()
	run.conn = conn
	client.stateMu.Unlock()

	cancel()
	if client.isConnectionActive(run, conn) {
		t.Fatal("canceled caller context left connection active")
	}
	if client.stopRun(run, runStopByFailure, errors.New("late connection failure")) {
		t.Fatal("canceled caller context accepted a later connection failure")
	}
	if current, ok := client.currentConnection(run); ok || current != nil {
		t.Fatalf("currentConnection = (%p, %v), want (nil, false)", current, ok)
	}
	if client.startMessageTask(run, nil) {
		t.Fatal("canceled caller context admitted a message task")
	}
	if finishErr := client.finishRun(run); !errors.Is(finishErr, context.Canceled) {
		t.Fatalf("finishRun returned %v, want context.Canceled", finishErr)
	}
}
