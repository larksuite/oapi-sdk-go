package ws

import (
	"context"
	"sync"
	"testing"
	"time"
)

// These tests exercise the concurrent read/write paths on the connection-control
// fields and the active run/connection identity. They are meant to be run under
// `go test -race`; the data
// race only surfaces under the race detector, so a plain run passing is not
// sufficient evidence the bug is fixed.
//
// The reader side goes through configSnapshot, currentConnection and exact
// connection checks. The writer side uses configure/applyConfig and connection
// replacement, which share stateMu with lifecycle identity.

// newRaceClient builds an in-memory Client without establishing a connection.
// conn stays nil on purpose: these tests cover the field-synchronization paths,
// not real dialing. The initial config matches validConfigPayloads[0] so a reader
// that runs before the first configure() still observes a valid whole-config view.
func newRaceClient() *Client {
	p := validConfigPayloads[0]
	return &Client{
		reconnectCount:    p.ReconnectCount,
		reconnectInterval: time.Duration(p.ReconnectInterval) * time.Second,
		reconnectNonce:    p.ReconnectNonce,
		pingInterval:      time.Duration(p.PingInterval) * time.Second,
	}
}

// validConfigPayloads is the closed set of configs the writer goroutines push.
// configSnapshot must always read back one of these combinations as a whole —
// never a torn mix of fields from different writes.
var validConfigPayloads = []*ClientConfig{
	{ReconnectCount: -1, ReconnectInterval: 1, ReconnectNonce: 5, PingInterval: 10},
	{ReconnectCount: 3, ReconnectInterval: 2, ReconnectNonce: 30, PingInterval: 120},
	{ReconnectCount: 0, ReconnectInterval: 5, ReconnectNonce: 0, PingInterval: 60},
}

func isValidSnapshot(count int, interval time.Duration, nonce int, ping time.Duration) bool {
	for _, p := range validConfigPayloads {
		if count == p.ReconnectCount &&
			interval == time.Duration(p.ReconnectInterval)*time.Second &&
			nonce == p.ReconnectNonce &&
			ping == time.Duration(p.PingInterval)*time.Second {
			return true
		}
	}
	return false
}

// configure (writer, models pong config push) racing configSnapshot (reader,
// models pingLoop/reconnect reading the config fields) must not data-race, and
// every snapshot must be a consistent whole-config view, never a torn read.
func TestConfigureAndConfigSnapshotNoRace(t *testing.T) {
	c := newRaceClient()

	const goroutines = 8
	const iterations = 200
	var wg sync.WaitGroup

	// writers: push whole ClientConfig values
	for w := 0; w < goroutines; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				c.configure(validConfigPayloads[(w+i)%len(validConfigPayloads)])
			}
		}(w)
	}

	// readers: take config snapshots and assert each is a consistent whole
	for r := 0; r < goroutines; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				count, interval, nonce, ping := c.configSnapshot()
				if !isValidSnapshot(count, interval, nonce, ping) {
					t.Errorf("torn config snapshot: count=%d interval=%v nonce=%d ping=%v",
						count, interval, nonce, ping)
					return
				}
			}
		}()
	}

	wg.Wait()
}

// Run-level config updates and current-connection snapshots racing connection
// replacement must not data-race.
func TestConfigureAndConnStateNoRace(t *testing.T) {
	c := newRaceClient()
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()
	run := &clientRun{ctx: runCtx, cancel: runCancel}
	conn := &clientConn{
		connectionResult: make(chan error, 1),
		serviceID:        "42",
	}
	replacement := &clientConn{
		connectionResult: make(chan error, 1),
		serviceID:        "84",
	}
	c.run = run
	run.conn = conn

	const goroutines = 8
	const iterations = 200
	var wg sync.WaitGroup

	// config writers
	for w := 0; w < goroutines; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				c.configure(validConfigPayloads[(w+i)%len(validConfigPayloads)])
			}
		}(w)
	}

	// Active-run config writers are independent of which physical connection
	// supplied a Pong frame.
	for w := 0; w < goroutines; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				if !c.applyConfig(run, validConfigPayloads[(w+i)%len(validConfigPayloads)]) {
					t.Errorf("active run rejected config update")
					return
				}
			}
		}(w)
	}

	// Replace the current connection while readers snapshot it.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			c.stateMu.Lock()
			if run.conn == conn {
				run.conn = replacement
			} else {
				run.conn = conn
			}
			c.stateMu.Unlock()
		}
	}()

	// Snapshot readers accept either fully published connection and also exercise
	// the exact physical-connection check used by receive cleanup.
	for r := 0; r < goroutines; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				current, ok := c.currentConnection(run)
				if !ok || (current != conn && current != replacement) {
					t.Errorf("currentConnection returned (%p, %v)", current, ok)
					return
				}
				c.isConnectionActive(run, current)
			}
		}()
	}

	wg.Wait()
}

// Single-threaded behavior guard: locking the write path must not change
// configure's semantics — a snapshot taken right after a configure reads back
// exactly what was written, with the second-based unit conversion intact.
func TestConfigureThenSnapshotReadsBackWrittenValues(t *testing.T) {
	c := newRaceClient()

	c.configure(&ClientConfig{
		ReconnectCount:    7,
		ReconnectInterval: 3,
		ReconnectNonce:    15,
		PingInterval:      90,
	})

	count, interval, nonce, ping := c.configSnapshot()
	if count != 7 {
		t.Errorf("reconnectCount = %d, want 7", count)
	}
	if interval != 3*time.Second {
		t.Errorf("reconnectInterval = %v, want 3s", interval)
	}
	if nonce != 15 {
		t.Errorf("reconnectNonce = %d, want 15", nonce)
	}
	if ping != 90*time.Second {
		t.Errorf("pingInterval = %v, want 90s", ping)
	}
}
