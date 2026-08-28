package ws

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"
)

type eventMessage struct {
	frame       Frame
	headers     Headers
	messageType string
	messageID   string
	traceID     string
	payload     []byte
	startedAt   int64
}

// startMessageTask atomically admits one SDK-owned message task. Once a stop
// wins, the same run cannot add more tracked work.
func (c *Client) startMessageTask(run *clientRun, msg []byte) bool {
	if run.ctx.Err() != nil {
		return false
	}
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if run.stopReason != runStopNone {
		return false
	}
	run.wg.Add(1)
	go c.handleMessageTask(run, msg)
	return true
}

func (c *Client) handleMessageTask(run *clientRun, msg []byte) {
	defer run.wg.Done()
	defer func() {
		if recover() != nil {
			c.logger.Error(run.ctx, c.fmtLog("websocket message handler panicked")...)
		}
	}()

	var frame Frame
	if err := frame.Unmarshal(msg); err != nil {
		c.logger.Error(run.ctx, c.fmtLog("websocket message decode failed")...)
		return
	}

	switch FrameType(frame.Method) {
	case FrameTypeControl:
		c.handleControlFrame(run, frame)
	case FrameTypeData:
		c.handleDataFrame(run, frame)
	default:
	}
}

func (c *Client) handleDataFrame(run *clientRun, frame Frame) {
	hs := Headers(frame.Headers)
	sum := hs.GetInt(HeaderSum)
	seq := hs.GetInt(HeaderSeq)
	msgID := hs.GetString(HeaderMessageID)
	traceID := hs.GetString(HeaderTraceID)
	messageType := hs.GetString(HeaderType)

	payload := frame.Payload
	if sum > 1 {
		payload = c.combine(msgID, sum, seq, payload)
		if payload == nil {
			return
		}
	}

	c.logger.Debug(run.ctx, c.fmtLog("receive message, message_type: %s, message_id: %s, trace_id: %s, payload: %s",
		messageType, msgID, traceID, payload))
	if MessageType(messageType) != MessageTypeEvent || c.eventHandler == nil {
		return
	}

	startedAt := time.Now().UnixNano() / int64(time.Millisecond) // 兼容 go < 1.17
	handler := c.eventHandler
	event := eventMessage{
		frame:       frame,
		headers:     hs,
		messageType: messageType,
		messageID:   msgID,
		traceID:     traceID,
		payload:     payload,
		startedAt:   startedAt,
	}
	response, err := handler.Do(run.ctx, event.payload)
	c.writeEventResponse(run, event, response, err)
}

func (c *Client) writeEventResponse(
	run *clientRun,
	event eventMessage,
	response interface{},
	handlerErr error,
) {
	endedAt := time.Now().UnixNano() / int64(time.Millisecond)
	event.headers.Add(HeaderBizRt, strconv.FormatInt(endedAt-event.startedAt, 10))

	resp := NewResponseByCode(http.StatusOK)
	if handlerErr != nil {
		c.logger.Error(run.ctx, c.fmtLog("handle message failed, message_type: %s, message_id: %s, trace_id: %s, err: %v",
			event.messageType, event.messageID, event.traceID, handlerErr)...)
		resp = NewResponseByCode(http.StatusInternalServerError)
	} else if response != nil {
		data, err := json.Marshal(response)
		if err != nil {
			c.logger.Error(run.ctx, c.fmtLog("handle message failed, message_type: %s, message_id: %s, trace_id: %s, err: %v",
				event.messageType, event.messageID, event.traceID, err)...)
			resp = NewResponseByCode(http.StatusInternalServerError)
		} else {
			resp.Data = data
		}
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		c.logger.Error(run.ctx, c.fmtLog("response message encode failed")...)
		return
	}
	event.frame.Payload = payload
	event.frame.Headers = event.headers
	message, err := event.frame.Marshal()
	if err != nil {
		c.logger.Error(run.ctx, c.fmtLog("response frame encode failed")...)
		return
	}

	if err := c.writeMessage(run, ws.BinaryMessage, message); err != nil {
		c.logger.Error(run.ctx, c.fmtLog("response message failed, type: %s, message_id: %s, trace_id: %s", event.messageType, event.messageID, event.traceID)...)
	}
}

func (c *Client) combine(msgID string, sum, seq int, bs []byte) []byte {
	val := c.cache.Get(msgID)
	if val == nil {
		buf := make([][]byte, sum)
		buf[seq] = bs
		c.cache.Set(msgID, buf, 5*time.Second)
		return nil
	}

	buf := val.([][]byte)
	buf[seq] = bs
	capacity := 0
	for _, v := range buf {
		if len(v) == 0 {
			c.cache.Set(msgID, buf, 5*time.Second)
			return nil
		}
		capacity += len(v)
	}

	payload := make([]byte, 0, capacity)
	for _, v := range buf {
		payload = append(payload, v...)
	}

	return payload
}
