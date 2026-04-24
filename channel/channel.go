package channel

import (
	"context"
	"fmt"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/channel/normalize"
	"github.com/larksuite/oapi-sdk-go/v3/channel/outbound"
	"github.com/larksuite/oapi-sdk-go/v3/channel/pipeline"
	"github.com/larksuite/oapi-sdk-go/v3/channel/safety"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

// channelImpl is the default implementation of the Channel interface.
type channelImpl struct {
	client          *lark.Client
	wsClient        *larkws.Client
	uploader        outbound.Uploader
	dedupCache      *safety.DedupCache
	pipelineManager *pipeline.ChatPipelineManager
	policyGate      *safety.PolicyGate
	processLock     *safety.ProcessingLock
	staleWindow     time.Duration
}

// NewChannel creates a new Channel instance with the provided Lark API client and WebSocket client.
func NewChannel(client *lark.Client, wsClient *larkws.Client) types.Channel {
	return &channelImpl{
		client:          client,
		wsClient:        wsClient,
		uploader:        outbound.NewUploader(client),
		dedupCache:      safety.NewDedupCache(10000, 1*time.Hour),
		pipelineManager: pipeline.NewChatPipelineManager(types.DefaultBatchConfig()),
		policyGate:      safety.NewPolicyGate(&types.PolicyConfig{}, nil),
		processLock:     safety.NewProcessingLock(types.DefaultLockTTL, 1*time.Minute),
		staleWindow:     types.DefaultStaleWindow,
	}
}

// UpdatePolicy updates the policy configuration for the channel.
func (ch *channelImpl) UpdatePolicy(cfg types.PolicyConfig) {
	ch.policyGate.UpdateConfig(cfg)
}

// OnMessage registers a handler for NormalizedMessage events.
func (ch *channelImpl) OnMessage(handler func(ctx context.Context, msg *types.NormalizedMessage) error) {
	if ch.wsClient != nil {
		dispatcher := ch.wsClient.EventHandler()
		if dispatcher != nil {
			dispatcher.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
				normMsg := normalize.ParseMessage(event)
				if normMsg != nil {
					// 1. Stale Check
					if safety.IsStale(normMsg.CreateTimeMs, ch.staleWindow) {
						return nil // drop stale message
					}

					// 2. Dedup Cache
					if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(normMsg.EventID) {
						return nil
					}

					// 3. Policy Gate
					decision := ch.policyGate.Evaluate(normMsg)
					if !decision.Allowed {
						// For now we just drop it if policy fails.
						// Optionally we could trigger a callback.
						return nil
					}

					// 4. Processing Lock
					if !ch.processLock.Acquire(normMsg.MessageID) {
						return nil // drop in-flight message
					}

					// 5. Batch Pipeline
					dispatchHandler := func(ctx context.Context, batch *types.BatchedDispatch) error {
						defer func() {
							for _, id := range batch.SourceIDs {
								ch.processLock.Release(id)
							}
						}()
						return handler(ctx, batch.Message)
					}

					ch.pipelineManager.Push(ctx, normMsg.ChatID, normMsg, dispatchHandler)
				}
				return nil
			})
		}
	}
}

// OnCardAction registers a handler for CardActionTriggerEvent events.
func (ch *channelImpl) OnCardAction(handler func(ctx context.Context, msg *types.NormalizedMessage) error) {
	if ch.wsClient != nil {
		dispatcher := ch.wsClient.EventHandler()
		if dispatcher != nil {
			dispatcher.OnP2CardActionTrigger(func(ctx context.Context, event *callback.CardActionTriggerEvent) (*callback.CardActionTriggerResponse, error) {
				normMsg := normalize.ParseCardAction(event)
				if normMsg != nil {
					// Card actions don't use batching but we can use queueing and locks.
					if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(normMsg.EventID) {
						return nil, nil
					}

					if !ch.processLock.Acquire(normMsg.EventID) {
						return nil, nil
					}
					defer ch.processLock.Release(normMsg.EventID)

					// Queue to serialize per chat
					err := ch.pipelineManager.Run(ctx, normMsg.ChatID, func() error {
						return handler(ctx, normMsg)
					})

					if err != nil {
						return nil, err
					}
				}
				return nil, nil
			})
		}
	}
}

// Stream initiates a streaming message session. It returns a StreamController to append and flush content.
func (ch *channelImpl) Stream(ctx context.Context, input *types.SendInput) (types.StreamController, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	// Ensure we send an initial message to get a MessageID
	// For streaming, we typically want to start with empty or initial markdown text
	if input.Markdown == "" && input.Text == "" {
		input.Markdown = "..."
	}

	res, err := ch.Send(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to send initial message for streaming: %w", err)
	}

	return NewMarkdownStreamController(ch.client, res.MessageID, input.Markdown, input.Title), nil
}
