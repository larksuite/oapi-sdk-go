package ws

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	websocket "github.com/gorilla/websocket"
)

func activeRunAndConn() (*clientRun, *clientConn) {
	runCtx, runCancel := context.WithCancel(context.Background())
	run := &clientRun{ctx: runCtx, cancel: runCancel}
	conn := &clientConn{
		readResult: make(chan error, 1),
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

func TestHandlerAdmissionIsAtomicWithStop(t *testing.T) {
	run, _ := activeRunAndConn()
	defer run.cancel()
	client := NewClient("app-id", "app-secret")
	client.run = run

	started := make(chan struct{})
	release := make(chan struct{})
	finished := make(chan struct{})
	admitted := client.admitEventTask(run, func() {
		defer close(finished)
		close(started)
		<-release
	})
	if !admitted {
		t.Fatal("active run did not admit handler")
	}
	waitLifecycleSignal(t, started, "admitted handler start")
	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}
	close(release)
	waitLifecycleSignal(t, finished, "admitted handler completion")

	var calls int32
	if client.admitEventTask(run, func() {
		atomic.AddInt32(&calls, 1)
	}) {
		t.Fatal("stopped run admitted a new handler")
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("stopped run invoked handler %d times, want 0", got)
	}
}

func TestAdmittedHandlerUsesCurrentConnectionAfterReconnect(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	firstSocket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, firstSocket)
	firstServerConn := waitAcceptedGatewayConnection(t, gateway)
	secondSocket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, secondSocket)
	secondServerConn := waitAcceptedGatewayConnection(t, gateway)

	run, _ := activeRunAndSocket(firstSocket, gateway.endpoint())
	defer run.cancel()
	secondConn := newSocketClientConn(secondSocket, gateway.endpoint())
	client := NewClient("app-id", "app-secret")
	client.run = run

	started := make(chan struct{})
	release := make(chan struct{})
	writeResult := make(chan error, 1)
	if !client.admitEventTask(run, func() {
		close(started)
		<-release
		writeResult <- client.writeMessage(run, websocket.BinaryMessage, []byte("late response"))
	}) {
		t.Fatal("active run did not admit handler")
	}
	waitLifecycleSignal(t, started, "handler start on first connection")

	client.stateMu.Lock()
	run.conn = secondConn
	client.stateMu.Unlock()
	close(release)

	select {
	case err := <-writeResult:
		if err != nil {
			t.Fatalf("admitted handler write returned %v", err)
		}
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("admitted handler did not finish its write")
	}
	frame := waitGatewayFrame(t, gateway, "handler response on replacement connection")
	if frame.conn != secondServerConn {
		t.Fatalf("handler response used connection %p, want replacement %p", frame.conn, secondServerConn)
	}
	if string(frame.payload) != "late response" {
		t.Fatalf("handler response payload = %q, want %q", frame.payload, "late response")
	}
	assertNoGatewayFrame(t, gateway, firstServerConn, 100*time.Millisecond)
}

func TestAdmittedHandlerCannotWriteAfterRunStops(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	socket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, socket)
	serverConn := waitAcceptedGatewayConnection(t, gateway)
	run, _ := activeRunAndSocket(socket, gateway.endpoint())
	defer run.cancel()
	client := NewClient("app-id", "app-secret")
	client.run = run

	started := make(chan struct{})
	release := make(chan struct{})
	writeResult := make(chan error, 1)
	if !client.admitEventTask(run, func() {
		close(started)
		<-release
		writeResult <- client.writeMessage(run, websocket.BinaryMessage, []byte("stale response"))
	}) {
		t.Fatal("active run did not admit handler")
	}
	waitLifecycleSignal(t, started, "handler start before stop")
	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}
	close(release)

	select {
	case err := <-writeResult:
		if !errors.Is(err, errConnectionClosed) {
			t.Fatalf("handler write returned %v, want errConnectionClosed", err)
		}
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("admitted handler did not finish after stop")
	}
	assertNoGatewayFrame(t, gateway, serverConn, 100*time.Millisecond)
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
