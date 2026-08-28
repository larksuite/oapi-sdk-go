package ws

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	websocket "github.com/gorilla/websocket"
)

func dialLifecycleCandidate(t *testing.T, gateway *lifecycleGateway) *websocket.Conn {
	t.Helper()
	conn, resp, err := websocket.DefaultDialer.Dial(gateway.endpoint(), nil)
	if resp != nil && resp.Body != nil {
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				t.Errorf("close candidate response body: %v", closeErr)
			}
		}()
	}
	if err != nil {
		t.Fatalf("dial candidate: %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("candidate handshake status = %v, want %d", resp, http.StatusSwitchingProtocols)
	}
	return conn
}

func TestActivateConnectionRejectsStoppedRun(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()
	socket := dialLifecycleCandidate(t, gateway)
	runCtx, runCancel := context.WithCancel(context.Background())
	run := &clientRun{ctx: runCtx, cancel: runCancel}
	conn := &clientConn{
		socket:           socket,
		connectionResult: make(chan error, 1),
		safeEndpoint:     safeEndpoint(gateway.endpoint()),
	}
	var readyCount int32
	var disconnectedCount int32
	client := NewClient("app-id", "app-secret",
		WithOnReady(func() { atomic.AddInt32(&readyCount, 1) }),
		WithOnDisconnected(func() { atomic.AddInt32(&disconnectedCount, 1) }),
	)
	client.run = run

	if !client.stopRun(run, runStopByClose, nil) {
		t.Fatal("stopRun rejected the active run")
	}
	if client.activateConnection(run, conn, false) {
		t.Fatal("activateConnection accepted a stopped run")
	}
	if run.conn != nil {
		t.Fatal("activateConnection stored a candidate on a stopped run")
	}
	if err := socket.WriteMessage(websocket.BinaryMessage, []byte("probe")); err == nil {
		t.Fatal("rejected candidate socket remained writable")
	}
	if got := atomic.LoadInt32(&readyCount); got != 0 {
		t.Errorf("rejected activation invoked Ready %d times, want 0", got)
	}
	if got := atomic.LoadInt32(&disconnectedCount); got != 0 {
		t.Errorf("rejected activation invoked Disconnected %d times, want 0", got)
	}

	waitDone := make(chan struct{})
	go func() {
		run.wg.Wait()
		close(waitDone)
	}()
	select {
	case <-waitDone:
	case <-time.After(time.Second):
		t.Fatal("rejected activation registered workers")
	}
}
