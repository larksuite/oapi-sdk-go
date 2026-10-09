package main

import (
	"context"
	"fmt"
	"log"
	"time"

	core "github.com/larksuite/oapi-sdk-go/v3/core"
	agent "github.com/larksuite/oapi-sdk-go/v3/service/agent/v1"
)

func (e *example) read(
	ctx context.Context,
	sid, queryID, cursor string,
	closeEarly bool,
) (result error) {
	request := agent.NewStreamSessionEventReqBuilder().
		SessionId(sid).
		MessageId(queryID).
		Build()
	bounds := core.SSEOptions{
		OpenTimeout:     10 * time.Second,
		IdleReadTimeout: 60 * time.Second,
		TotalTimeout:    180 * time.Second,
		MaxEventBytes:   1024 * 1024,
		LastEventID:     cursor,
	}
	stream, err := e.client.Agent.V1.SessionEvent.Stream(
		ctx,
		request,
		core.WithUserAccessToken(e.token),
		core.WithSSEOptions(bounds),
	)
	if err != nil {
		return fmt.Errorf("open stream: %w", err)
	}
	defer func() {
		if err := stream.Close(); err != nil {
			if result == nil {
				result = fmt.Errorf("close stream: %w", err)
			} else {
				log.Printf("stream cleanup failed: %T", err)
			}
		}
	}()
	log.Printf("stream content type=%q", stream.Headers().Get("Content-Type"))
	for stream.Next() {
		event := stream.Event()
		if err := handleEvent(event); err != nil {
			return err
		}
		// Checkpoint only after your application has processed this event.
		if event.ID != nil {
			log.Printf("AGENT_LAST_EVENT_ID=%q", *event.ID)
		}
		if closeEarly {
			log.Print("closing receiver after one event; Agent may continue generating")
			break
		}
	}
	if err := stream.Err(); err != nil {
		return fmt.Errorf("read stream: %w", err)
	}
	return nil
}

func (e *example) interrupt(ctx context.Context) error {
	sid, err := required("AGENT_SESSION_ID")
	if err != nil {
		return err
	}
	assistantID, err := required("AGENT_ASSISTANT_MESSAGE_ID")
	if err != nil {
		return err
	}
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	event := agent.NewEventBuilder().
		Type("user.interrupt").
		MessageId(assistantID).
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
		return fmt.Errorf("interrupt reply: %w", err)
	}
	if !response.Success() {
		return fmt.Errorf("interrupt reply: HTTP %d, code %d", response.StatusCode, response.Code)
	}
	log.Print("interrupt accepted; immediate snapshot may still be in_progress")
	return e.queryReplyStatus(ctx, sid, assistantID)
}

func (e *example) queryReplyStatus(ctx context.Context, sid, assistantID string) error {
	request := agent.NewGetSessionMessageReqBuilder().
		SessionId(sid).
		MessageId(assistantID).
		Build()
	state, err := e.client.Agent.V1.SessionMessage.Get(
		ctx,
		request,
		core.WithUserAccessToken(e.token),
	)
	if err != nil {
		return fmt.Errorf("query reply: %w", err)
	}
	if !state.Success() {
		return fmt.Errorf("query reply: code %d", state.Code)
	}
	if state.Data == nil || state.Data.Message == nil || state.Data.Message.Status == nil {
		return fmt.Errorf("query reply: missing status")
	}
	log.Printf(
		"reply status=%q; completed means it finished before cancellation",
		*state.Data.Message.Status,
	)
	return nil
}
