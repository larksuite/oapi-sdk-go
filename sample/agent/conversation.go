package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"

	core "github.com/larksuite/oapi-sdk-go/v3/core"
	agent "github.com/larksuite/oapi-sdk-go/v3/service/agent/v1"
)

func newUUID() (string, error) {
	var data [16]byte
	if _, err := rand.Read(data[:]); err != nil {
		return "", fmt.Errorf("generate UUID: %w", err)
	}
	data[6] = (data[6] & 15) | 64
	data[8] = (data[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", data[:4], data[4:6], data[6:8], data[8:10], data[10:]), nil
}

func (e *example) chat(ctx context.Context) error {
	sid, err := e.createSession(ctx)
	if err != nil {
		return err
	}
	queryID, err := e.sendMessage(ctx, sid)
	if err != nil {
		return err
	}
	log.Printf("AGENT_SESSION_ID=%s", sid)
	log.Printf("AGENT_QUERY_MESSAGE_ID=%s", queryID)
	return e.read(ctx, sid, queryID, "", false)
}

func (e *example) createSession(ctx context.Context) (string, error) {
	uuid, err := newUUID()
	if err != nil {
		return "", err
	}
	body := agent.NewCreateSessionReqBodyBuilder().
		Uuid(uuid).
		Name("SDK Agent example").
		Build()
	request := agent.NewCreateSessionReqBuilder().
		Body(body).
		Build()

	response, err := e.client.Agent.V1.Session.Create(
		ctx,
		request,
		core.WithUserAccessToken(e.token),
	)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	if !response.Success() {
		return "", fmt.Errorf(
			"create session: HTTP %d, code %d",
			response.StatusCode,
			response.Code,
		)
	}
	if response.Data == nil || response.Data.Id == nil {
		return "", fmt.Errorf("create session: missing id")
	}
	return *response.Data.Id, nil
}

func (e *example) sendMessage(ctx context.Context, sid string) (string, error) {
	text := agent.NewTextItemBuilder().
		Type("text").
		Text("请用三句话介绍 SSE。").
		Build()
	contents := []interface{}{text}
	// Lightweight union: use an ordinary map for a file branch not listed in Meta.
	if uri := os.Getenv("AGENT_FILE_URI"); uri != "" {
		file := map[string]interface{}{
			"type": "file",
			"uri":  uri,
		}
		contents = append(contents, file)
	}
	uuid, err := newUUID()
	if err != nil {
		return "", err
	}
	event := agent.NewEventBuilder().
		Type("user.message").
		Contents(contents).
		Build()
	body := agent.NewCreateSessionEventReqBodyBuilder().
		Uuid(uuid).
		Events([]*agent.Event{event}).
		Build()
	request := agent.NewCreateSessionEventReqBuilder().
		SessionId(sid).
		Body(body).
		Build()

	response, err := e.client.Agent.V1.SessionEvent.Create(
		ctx,
		request,
		core.WithUserAccessToken(e.token),
	)
	if err != nil {
		return "", fmt.Errorf("send message: %w", err)
	}
	if !response.Success() {
		return "", fmt.Errorf(
			"send message: HTTP %d, code %d",
			response.StatusCode,
			response.Code,
		)
	}
	if response.Data == nil || len(response.Data.Messages) == 0 {
		return "", fmt.Errorf("send message: missing user query id")
	}
	if response.Data.Messages[0].Id == nil {
		return "", fmt.Errorf("send message: missing user query id")
	}
	return *response.Data.Messages[0].Id, nil
}

func handleEvent(event *core.SSEEvent) error {
	if event.Event == "message.failed" {
		return fmt.Errorf("Agent reply failed")
	}
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(event.Data), &payload); err != nil {
		return fmt.Errorf("decode SSE data: %w", err)
	}
	if payload == nil {
		return fmt.Errorf("SSE data must be a JSON object")
	}
	if base, ok := payload["BaseResp"].(map[string]interface{}); ok {
		if code, ok := base["StatusCode"].(float64); ok && code != 0 {
			return fmt.Errorf("SSE business error: code %.0f", code)
		}
	}
	switch event.Event {
	case "text.delta":
		delta, _ := payload["delta"].(string)
		// Pass delta to your application's UI; do not log user content by default.
		log.Printf("text delta received: %d bytes", len(delta))
	case "text":
		text, _ := payload["text"].(string)
		// The complete item replaces previews for the same item id after reconnect.
		log.Printf("complete text received: %d bytes", len(text))
	default:
		log.Printf("SSE event=%q", event.Event)
	}
	message, ok := payload["message"].(map[string]interface{})
	if ok && message["role"] == "assistant" {
		if id, ok := message["id"].(string); ok {
			log.Printf("AGENT_ASSISTANT_MESSAGE_ID=%s", id)
		}
	}
	return nil
}
