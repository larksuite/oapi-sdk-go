package ws

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

type runStopReason uint8

const (
	runStopNone runStopReason = iota
	runStopByContext
	runStopByClose
	runStopByFailure
)

type clientRun struct {
	ctx    context.Context
	cancel context.CancelFunc

	stopReason    runStopReason
	stopErr       error
	conn          *clientConn
	everConnected bool

	wg sync.WaitGroup
}

func (c *Client) Start(callerCtx context.Context) error {
	if callerCtx == nil {
		return errNilContext
	}

	run, err := c.beginRun(callerCtx)
	if err != nil {
		return err
	}
	if run.ctx.Err() == nil {
		c.runCoordinator(run)
	}
	return c.finishRun(run)
}

func (c *Client) Close() {
	c.stateMu.Lock()
	run := c.run
	c.stateMu.Unlock()
	if run != nil {
		c.stopRun(run, runStopByClose, nil)
	}
}

func (c *Client) beginRun(callerCtx context.Context) (*clientRun, error) {
	c.stateMu.Lock()
	if c.run != nil {
		c.stateMu.Unlock()
		return nil, errClientRunning
	}
	if c.terminal {
		c.stateMu.Unlock()
		return nil, errClientTerminal
	}

	runCtx, cancel := context.WithCancel(callerCtx)
	run := &clientRun{ctx: runCtx, cancel: cancel}
	c.run = run
	c.stateMu.Unlock()
	return run, nil
}

func (c *Client) stopRun(run *clientRun, reason runStopReason, err error) (accepted bool) {
	if reason != runStopByContext && run.ctx.Err() != nil {
		return false
	}

	c.stateMu.Lock()
	if c.run != run || run.stopReason != runStopNone {
		c.stateMu.Unlock()
		return false
	}

	run.stopReason = reason
	run.stopErr = err
	if reason == runStopByContext || reason == runStopByClose ||
		(reason == runStopByFailure && run.everConnected) {
		c.terminal = true
	}
	c.stateMu.Unlock()

	run.cancel()
	return true
}

func (c *Client) finishRun(run *clientRun) error {
	if contextErr := run.ctx.Err(); contextErr != nil {
		c.stopRun(run, runStopByContext, contextErr)
	}
	run.cancel()
	run.wg.Wait()

	c.stateMu.Lock()
	err := run.stopErr
	c.run = nil
	c.stateMu.Unlock()
	return err
}

// runCoordinator exposes the lifecycle in three steps: establish a physical
// connection, activate its Ping/Receive/Pong session, then stop or reconnect.
func (c *Client) runCoordinator(run *clientRun) {
	conn, err := c.establishConnectionOnce(run)
	reconnected := false
	if err != nil {
		if !c.isRunActive(run) {
			return
		}
		if c.autoReconnect && isRetryableConnectionError(err) {
			c.invokeRecoverableErrorCallback(run, err)
		}
		conn = c.reconnectAfterFailure(run, err)
		reconnected = true
	}

	for conn != nil {
		activated := c.activateConnection(run, conn, reconnected)
		if !activated {
			return
		}
		exitErr := c.waitForConnectionExit(run, conn)
		if exitErr == nil {
			return
		}

		conn = c.reconnectAfterFailure(run, exitErr)
		reconnected = true
	}
}

// reconnectAfterFailure owns one reconnect sequence. The failed read that starts
// a reconnect is not reported as a recoverable dial error. An initial failed
// connection is reported by runCoordinator before this sequence
// begins; subsequent failed connection attempts are reported below.
func (c *Client) reconnectAfterFailure(run *clientRun, failure error) *clientConn {
	if !c.autoReconnect || !isRetryableConnectionError(failure) {
		c.stopRunWithError(run, failure)
		return nil
	}
	c.invokeReconnectingCallback(run)

	attempts := 0

	for {
		reconnectCount, reconnectInterval, reconnectNonce, _ := c.configSnapshot()
		if reconnectCount >= 0 && attempts >= reconnectCount {
			err := newLifecycleError("reconnect", "attempts exhausted", 0, nil, false)
			c.stopRunWithError(run, err)
			return nil
		}

		delay := reconnectInterval
		if attempts == 0 {
			delay = randomReconnectDelay(reconnectNonce)
		}
		if !c.waitRunDelay(run, delay) {
			return nil
		}
		attempts++

		conn, err := c.establishConnectionOnce(run)
		if err == nil {
			return conn
		}
		if !c.isRunActive(run) {
			return nil
		}
		failure = err
		if !c.autoReconnect || !isRetryableConnectionError(failure) {
			c.stopRunWithError(run, failure)
			return nil
		}
		c.invokeRecoverableErrorCallback(run, failure)
	}
}

// waitForConnectionExit returns nil when the run has stopped. Otherwise, it
// returns the read error that starts a reconnect sequence.
func (c *Client) waitForConnectionExit(run *clientRun, conn *clientConn) error {
	var exitErr error
	select {
	case <-run.ctx.Done():
	case exitErr = <-conn.readResult:
	}

	c.deactivateConnection(run)
	if !c.isRunActive(run) {
		return nil
	}
	return exitErr
}

func (c *Client) stopRunWithError(run *clientRun, err error) {
	if !c.stopRun(run, runStopByFailure, err) {
		return
	}
	c.invokeTerminalErrorCallback(run, err)
}

func (c *Client) isRunActive(run *clientRun) bool {
	if run.ctx.Err() != nil {
		return false
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return run.stopReason == runStopNone
}

func (c *Client) waitRunDelay(run *clientRun, delay time.Duration) bool {
	if delay <= 0 {
		return c.isRunActive(run)
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return c.isRunActive(run)
	case <-run.ctx.Done():
		return false
	}
}

func (c *Client) invokeReconnectingCallback(run *clientRun) {
	if run.ctx.Err() != nil {
		return
	}
	c.stateMu.Lock()
	if run.stopReason != runStopNone {
		c.stateMu.Unlock()
		return
	}
	callback := c.onReconnecting
	c.stateMu.Unlock()
	if callback != nil {
		c.runLifecycleCallback(run, callback)
	}
}

func (c *Client) invokeRecoverableErrorCallback(run *clientRun, err error) {
	if run.ctx.Err() != nil {
		return
	}
	c.stateMu.Lock()
	if run.stopReason != runStopNone {
		c.stateMu.Unlock()
		return
	}
	callback := c.onError
	c.stateMu.Unlock()
	if callback != nil {
		c.runRecoverableErrorCallback(run, callback, err)
	}
}

func (c *Client) invokeTerminalErrorCallback(run *clientRun, err error) {
	c.stateMu.Lock()
	callback := c.onError
	c.stateMu.Unlock()
	if callback != nil {
		c.runErrorCallback(run.ctx, callback, err)
	}
}

func (c *Client) runCallback(ctx context.Context, callback func()) {
	defer c.recoverCallback(ctx, "lifecycle")
	callback()
}

func (c *Client) runLifecycleCallback(run *clientRun, callback func()) {
	if !c.isRunActive(run) {
		return
	}
	c.runCallback(run.ctx, callback)
}

func (c *Client) runErrorCallback(ctx context.Context, callback func(error), err error) {
	defer c.recoverCallback(ctx, "error")
	callback(err)
}

func (c *Client) runRecoverableErrorCallback(run *clientRun, callback func(error), err error) {
	if !c.isRunActive(run) {
		return
	}
	c.runErrorCallback(run.ctx, callback, err)
}

func (c *Client) recoverCallback(ctx context.Context, category string) {
	if recover() != nil {
		c.logger.Error(ctx, c.fmtLog("websocket %s callback panicked", category)...)
	}
}

func randomReconnectDelay(reconnectNonce int) time.Duration {
	if reconnectNonce <= 0 {
		return 0
	}
	source := rand.New(rand.NewSource(time.Now().UnixNano()))
	return time.Duration(source.Intn(reconnectNonce*1000)) * time.Millisecond
}

func isRetryableConnectionError(err error) bool {
	var clientErr *ClientError
	if errors.As(err, &clientErr) {
		return false
	}
	var lifecycleErr *lifecycleError
	if errors.As(err, &lifecycleErr) {
		switch {
		case lifecycleErr.operation == "bootstrap" && lifecycleErr.category == "request build":
			return false
		case lifecycleErr.operation == "bootstrap credential" && lifecycleErr.category != "token retrieval":
			return false
		case lifecycleErr.operation == "parse endpoint":
			return false
		}
	}
	return true
}
