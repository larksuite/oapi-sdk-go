package ws

import (
	"context"
	"errors"
	"testing"
	"time"

	websocket "github.com/gorilla/websocket"
)

func newSocketClientConn(socket *websocket.Conn, endpoint string) *clientConn {
	return &clientConn{
		socket:       socket,
		readResult:   make(chan error, 1),
		safeEndpoint: safeEndpoint(endpoint),
	}
}

func activeRunAndSocket(socket *websocket.Conn, endpoint string) (*clientRun, *clientConn) {
	runCtx, runCancel := context.WithCancel(context.Background())
	run := &clientRun{ctx: runCtx, cancel: runCancel, everConnected: true}
	conn := newSocketClientConn(socket, endpoint)
	run.conn = conn
	return run, conn
}

func waitAcceptedGatewayConnection(t *testing.T, gateway *lifecycleGateway) *websocket.Conn {
	t.Helper()
	select {
	case conn := <-gateway.accepted:
		return conn
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("timed out waiting for accepted gateway connection")
		return nil
	}
}

func waitGatewayFrame(t *testing.T, gateway *lifecycleGateway, name string) lifecycleGatewayFrame {
	t.Helper()
	select {
	case frame := <-gateway.frames:
		return frame
	case <-time.After(lifecycleTestTimeout):
		t.Fatalf("timed out waiting for %s", name)
		return lifecycleGatewayFrame{}
	}
}

func assertNoGatewayFrame(t *testing.T, gateway *lifecycleGateway, target *websocket.Conn, timeout time.Duration) {
	t.Helper()
	select {
	case frame := <-gateway.frames:
		if frame.conn == target {
			t.Fatalf("connection %p unexpectedly received payload %q", target, frame.payload)
		}
		t.Fatalf("unexpected payload %q on connection %p", frame.payload, frame.conn)
	case <-time.After(timeout):
	}
}

func TestWriteConnectionRejectsReplacedSnapshot(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	firstSocket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, firstSocket)
	firstServerConn := waitAcceptedGatewayConnection(t, gateway)
	secondSocket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, secondSocket)
	secondServerConn := waitAcceptedGatewayConnection(t, gateway)

	run, firstConn := activeRunAndSocket(firstSocket, gateway.endpoint())
	defer run.cancel()
	secondConn := newSocketClientConn(secondSocket, gateway.endpoint())
	client := NewClient("app-id", "app-secret")
	client.run = run

	snapshot, ok := client.currentConnection(run)
	if !ok || snapshot != firstConn {
		t.Fatalf("currentConnection = (%p, %v), want first connection", snapshot, ok)
	}
	client.stateMu.Lock()
	run.conn = secondConn
	client.stateMu.Unlock()

	if current, currentOK := client.currentConnection(run); !currentOK || current != secondConn {
		t.Fatalf("currentConnection = (%p, %v), want replacement", current, currentOK)
	}
	if err := client.writeConnection(run, snapshot, websocket.BinaryMessage, []byte("stale payload")); !errors.Is(err, errConnectionClosed) {
		t.Fatalf("writeConnection with replaced snapshot returned %v, want errConnectionClosed", err)
	}
	assertNoGatewayFrame(t, gateway, firstServerConn, 100*time.Millisecond)
	assertNoGatewayFrame(t, gateway, secondServerConn, 100*time.Millisecond)
}

func TestStopAndDetachDoNotWaitForWriter(t *testing.T) {
	gateway, cleanupGateway := newLifecycleGateway(t)
	defer cleanupGateway()

	socket := dialLifecycleCandidate(t, gateway)
	defer closeGatewayConn(t, socket)
	waitAcceptedGatewayConnection(t, gateway)
	run, conn := activeRunAndSocket(socket, gateway.endpoint())
	defer run.cancel()
	client := NewClient("app-id", "app-secret")
	client.run = run

	conn.writeMu.Lock()
	writerEntered := make(chan struct{})
	writeResult := make(chan error, 1)
	go func() {
		close(writerEntered)
		writeResult <- client.writeConnection(run, conn, websocket.BinaryMessage, []byte("blocked payload"))
	}()
	waitLifecycleSignal(t, writerEntered, "writer entry")

	stopResult := make(chan bool, 1)
	go func() {
		stopResult <- client.stopRun(run, runStopByClose, nil)
	}()
	select {
	case accepted := <-stopResult:
		if !accepted {
			conn.writeMu.Unlock()
			t.Fatal("stopRun rejected the active run")
		}
	case <-time.After(300 * time.Millisecond):
		conn.writeMu.Unlock()
		t.Fatal("stopRun waited for the connection writer")
	}

	detachDone := make(chan struct{})
	go func() {
		client.deactivateConnection(run)
		close(detachDone)
	}()
	select {
	case <-detachDone:
	case <-time.After(300 * time.Millisecond):
		conn.writeMu.Unlock()
		t.Fatal("deactivateConnection waited for the connection writer")
	}

	conn.writeMu.Unlock()
	select {
	case err := <-writeResult:
		if !errors.Is(err, errConnectionClosed) {
			t.Fatalf("old writer returned %v after detach, want errConnectionClosed", err)
		}
	case <-time.After(lifecycleTestTimeout):
		t.Fatal("old writer remained blocked after releasing writeMu")
	}
}
