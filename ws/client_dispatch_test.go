package ws

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
)

func activeRunAndConn() (*clientRun, *clientConn) {
	runCtx, runCancel := context.WithCancel(context.Background())
	run := &clientRun{ctx: runCtx, cancel: runCancel}
	conn := &clientConn{
		connectionResult: make(chan error, 1),
	}
	run.conn = conn
	return run, conn
}

func TestConnectionAdmissionAfterStop(t *testing.T) {
	run, conn := activeRunAndConn()
	defer run.cancel()
	client := NewClient("app-id", "app-secret")
	client.run = run
	if !client.isConnectionActive(run, conn) {
		t.Fatal("active run/conn was not admitted")
	}
	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}
	if client.isConnectionActive(run, conn) {
		t.Fatal("stopped run/conn remained eligible for handler dispatch")
	}
}

func TestStoppedRunRejectsNewMessageWorker(t *testing.T) {
	run, _ := activeRunAndConn()
	defer run.cancel()
	client := NewClient("app-id", "app-secret")
	client.run = run
	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}
	if client.startMessageTask(run, nil) {
		t.Fatal("stopped run admitted a new message worker")
	}
}

func TestMessageTaskWaitsForEventHandler(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	socket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, socket)
	serverConn := waitAcceptedGatewayConnection(t, gateway)
	run, _ := activeRunAndSocket(socket, gateway.endpoint())
	defer run.cancel()

	started := make(chan struct{})
	release := make(chan struct{})
	client := NewClient("app-id", "app-secret", WithEventHandler(
		dispatcher.NewEventDispatcher("", "").OnCustomizedEvent("test.event", func(context.Context, *larkevent.EventReq) error {
			close(started)
			<-release
			return nil
		}),
	))
	client.run = run

	headers := Headers{}
	headers.Add(HeaderType, string(MessageTypeEvent))
	headers.Add(HeaderMessageID, "message-1")
	headers.Add(HeaderTraceID, "trace-1")
	message, err := (&Frame{
		Method:  int32(FrameTypeData),
		Headers: headers,
		Payload: []byte(`{"header":{"event_type":"test.event"}}`),
	}).Marshal()
	if err != nil {
		t.Fatalf("marshal event frame: %v", err)
	}
	if !client.startMessageTask(run, message) {
		t.Fatal("active run did not admit a message task")
	}
	waitLifecycleSignal(t, started, "event handler start")

	finished := make(chan struct{})
	go func() {
		run.wg.Wait()
		close(finished)
	}()
	select {
	case <-finished:
		t.Fatal("message task finished before its event handler returned")
	case <-time.After(100 * time.Millisecond):
	}
	assertNoGatewayFrame(t, gateway, serverConn, 100*time.Millisecond)

	close(release)
	waitLifecycleSignal(t, finished, "message task completion")
	frame := waitGatewayFrame(t, gateway, "event response")
	if frame.conn != serverConn {
		t.Fatalf("event response used connection %p, want %p", frame.conn, serverConn)
	}
	responseFrame := &Frame{}
	if err := responseFrame.Unmarshal(frame.payload); err != nil {
		t.Fatalf("decode event response frame: %v", err)
	}
	if responseFrame.Method != int32(FrameTypeData) {
		t.Fatalf("event response method = %d, want data frame", responseFrame.Method)
	}
}

func TestStoppedRunSkipsLifecycleCallbacks(t *testing.T) {
	run, _ := activeRunAndConn()
	defer run.cancel()
	var callbacks int32
	client := NewClient("app-id", "app-secret",
		WithOnReconnecting(func() { atomic.AddInt32(&callbacks, 1) }),
		WithOnError(func(error) { atomic.AddInt32(&callbacks, 1) }),
	)
	client.run = run
	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}

	client.invokeReconnectingCallback(run)
	client.invokeRecoverableErrorCallback(run, errors.New("retry failed"))
	if got := atomic.LoadInt32(&callbacks); got != 0 {
		t.Fatalf("stopped run invoked %d queued lifecycle callbacks, want 0", got)
	}
}
