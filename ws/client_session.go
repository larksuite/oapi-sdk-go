package ws

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	ws "github.com/gorilla/websocket"
)

// clientConn represents one physical websocket connection within a clientRun.
// A reconnect replaces run.conn with a new clientConn instead of mutating this
// value in place.
type clientConn struct {
	socket *ws.Conn

	readResult chan error
	writeMu    sync.Mutex

	connID       string
	serviceID    string
	safeEndpoint string
}

func (c *Client) fmtConnLog(conn *clientConn, format string, args ...interface{}) []interface{} {
	log := c.fmtLog(format, args...)
	if conn != nil && conn.connID != "" {
		log = append(log, "[conn_id="+conn.connID+"]")
	}
	return log
}

// activateConnection atomically publishes a physical connection and admits
// its receive worker before a concurrent stop can win. Ping belongs to the
// run and starts only with the first published connection.
func (c *Client) activateConnection(run *clientRun, conn *clientConn, reconnected bool) (activated bool) {
	c.stateMu.Lock()
	if run.stopReason != runStopNone || run.ctx.Err() != nil {
		c.stateMu.Unlock()
		c.closeSocket(run.ctx, conn.socket, conn.safeEndpoint, "reject connection candidate")
		return false
	}

	firstConnection := !run.everConnected
	run.conn = conn
	run.everConnected = true
	run.wg.Add(1)
	if firstConnection {
		run.wg.Add(1)
	}
	callback := c.onReady
	if reconnected {
		callback = c.onReconnected
	}
	go c.receiveMessageLoop(run, conn)
	if firstConnection {
		go c.pingLoop(run)
	}
	c.stateMu.Unlock()

	c.logger.Info(run.ctx, c.fmtConnLog(conn, "connected to %s", conn.safeEndpoint)...)
	c.runReadyCallback(run, conn, callback)
	return true
}

// deactivateConnection removes and closes run's current physical connection.
// The run coordinator calls it before beginning a reconnect attempt.
func (c *Client) deactivateConnection(run *clientRun) {
	c.stateMu.Lock()
	conn := run.conn
	run.conn = nil
	callback := c.onDisconnected
	c.stateMu.Unlock()

	c.closeSocket(run.ctx, conn.socket, conn.safeEndpoint, "close connection")
	c.logger.Info(run.ctx, c.fmtConnLog(conn, "disconnected to %s", conn.safeEndpoint)...)
	go c.runDisconnectedCallback(run.ctx, callback)
}

func (c *Client) closeSocket(ctx context.Context, socket *ws.Conn, endpoint, operation string) {
	if err := socket.Close(); err != nil {
		c.logger.Warn(ctx, c.fmtLog("websocket %s failed at %s", operation, endpoint)...)
	}
}

// pingLoop sends an immediate Ping and then uses the latest server-provided
// interval. Each iteration resolves the run's current physical connection.
func (c *Client) pingLoop(run *clientRun) {
	defer run.wg.Done()

	for {
		c.sendPing(run)

		_, _, _, pingInterval := c.configSnapshot()
		if !c.waitRunDelay(run, pingInterval) {
			return
		}
	}
}

func (c *Client) sendPing(run *clientRun) {
	defer func() {
		if recover() != nil {
			c.logger.Warn(run.ctx, c.fmtLog("websocket ping failed")...)
		}
	}()

	conn, ok := c.currentConnection(run)
	if !ok {
		return
	}

	var serviceID int64
	if conn.serviceID != "" {
		parsedServiceID, err := strconv.ParseInt(conn.serviceID, 10, 32)
		if err != nil {
			c.logger.Warn(run.ctx, c.fmtConnLog(conn, "websocket ping failed")...)
			return
		}
		serviceID = parsedServiceID
	}

	frame := NewPingFrame(int32(serviceID))
	payload, err := frame.Marshal()
	if err != nil {
		c.logger.Warn(run.ctx, c.fmtConnLog(conn, "websocket ping failed")...)
		return
	}
	if err := c.writeConnection(run, conn, ws.BinaryMessage, payload); err != nil {
		c.logger.Warn(run.ctx, c.fmtConnLog(conn, "websocket ping failed")...)
		return
	}
	c.logger.Debug(run.ctx, c.fmtConnLog(conn, "ping success")...)
}

func (c *Client) receiveMessageLoop(run *clientRun, conn *clientConn) {
	defer run.wg.Done()
	defer func() {
		if recover() != nil {
			c.reportReadExit(conn, newLifecycleError("read websocket", "connection lost", 0, nil, false))
		}
	}()

	for {
		messageType, msg, err := conn.socket.ReadMessage()
		if err != nil {
			if c.isConnectionActive(run, conn) {
				c.reportReadExit(conn, newLifecycleError("read websocket", "connection lost", 0, nil, false))
			}
			return
		}
		if messageType != ws.BinaryMessage {
			c.logger.Warn(run.ctx, c.fmtConnLog(conn, "websocket received unsupported message type %d", messageType)...)
			continue
		}
		c.startMessageTask(run, msg)
	}
}

func (c *Client) reportReadExit(conn *clientConn, err error) {
	conn.readResult <- err
}

func (c *Client) handleControlFrame(run *clientRun, frame Frame) {
	hs := Headers(frame.Headers)
	t := hs.GetString(HeaderType)

	switch MessageType(t) {
	case MessageTypePong:
		conn, _ := c.currentConnection(run)
		c.logger.Debug(run.ctx, c.fmtConnLog(conn, "receive pong")...)
		if len(frame.Payload) == 0 {
			return
		}
		conf := &ClientConfig{}
		if err := json.Unmarshal(frame.Payload, conf); err != nil {
			c.logger.Warn(run.ctx, c.fmtConnLog(conn, "websocket client config decode failed")...)
			return
		}
		c.applyConfig(run, conf)
	default:
	}
}

func (c *Client) runReadyCallback(run *clientRun, conn *clientConn, callback func()) {
	if !c.isConnectionActive(run, conn) {
		return
	}
	defer c.recoverCallback(run.ctx, "ready")
	if callback != nil {
		callback()
	}
}

func (c *Client) runDisconnectedCallback(ctx context.Context, callback func()) {
	defer c.recoverCallback(ctx, "disconnected")
	if callback != nil {
		callback()
	}
}

func (c *Client) isConnectionActive(run *clientRun, conn *clientConn) bool {
	if run.ctx.Err() != nil {
		return false
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return run.stopReason == runStopNone && run.conn == conn
}

func (c *Client) currentConnection(run *clientRun) (*clientConn, bool) {
	if run.ctx.Err() != nil {
		return nil, false
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if run.stopReason != runStopNone || run.conn == nil {
		return nil, false
	}
	return run.conn, true
}

func (c *Client) writeMessage(run *clientRun, messageType int, data []byte) error {
	conn, ok := c.currentConnection(run)
	if !ok {
		return errConnectionClosed
	}
	return c.writeConnection(run, conn, messageType, data)
}

func (c *Client) writeConnection(run *clientRun, conn *clientConn, messageType int, data []byte) error {
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()

	if !c.isConnectionActive(run, conn) {
		return errConnectionClosed
	}
	if err := conn.socket.WriteMessage(messageType, data); err != nil {
		return newLifecycleError("write websocket", "connection lost", 0, nil, false)
	}
	return nil
}
