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

	// Handler registries
	onMessageHandlers []func(ctx context.Context, msg *types.NormalizedMessage) error
	onCommentHandlers []func(ctx context.Context, event *types.CommentEvent) error
	onReactionHandlers []func(ctx context.Context, event *types.ReactionEvent) error
	onBotAddedHandlers []func(ctx context.Context, event *types.BotAddedEvent) error

	messageHandlerReg  bool
	reactionHandlerReg bool
	botAddedHandlerReg bool
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
	ch.onMessageHandlers = append(ch.onMessageHandlers, handler)
	ch.ensureMessageHandler()
}

// OnComment registers a handler for CommentEvent.
func (ch *channelImpl) OnComment(handler func(ctx context.Context, event *types.CommentEvent) error) {
	ch.onCommentHandlers = append(ch.onCommentHandlers, handler)
	ch.ensureMessageHandler()
}

func (ch *channelImpl) ensureMessageHandler() {
	if ch.messageHandlerReg || ch.wsClient == nil {
		return
	}
	ch.messageHandlerReg = true
	dispatcher := ch.wsClient.EventHandler()
	if dispatcher != nil {
		dispatcher.OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
			// Handle Message
			if len(ch.onMessageHandlers) > 0 {
				normMsg := normalize.ParseMessage(event)
				if normMsg != nil {
					if safety.IsStale(normMsg.CreateTimeMs, ch.staleWindow) {
						// do nothing
					} else if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(normMsg.EventID) {
						// do nothing
					} else {
						decision := ch.policyGate.Evaluate(normMsg)
						if decision.Allowed && ch.processLock.Acquire(normMsg.MessageID) {
							dispatchHandler := func(ctx context.Context, batch *types.BatchedDispatch) error {
								defer func() {
									for _, id := range batch.SourceIDs {
										ch.processLock.Release(id)
									}
								}()
								for _, h := range ch.onMessageHandlers {
									h(ctx, batch.Message)
								}
								return nil
							}
							ch.pipelineManager.Push(ctx, normMsg.ChatID, normMsg, dispatchHandler)
						}
					}
				}
			}

			// Handle Comment
			if len(ch.onCommentHandlers) > 0 {
				commentEvent := normalize.ParseComment(event)
				// If it's a valid comment (e.g. has ParentID)
				if commentEvent != nil && commentEvent.ParentID != "" {
					if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(commentEvent.EventID) {
						// do nothing
					} else {
						if ch.processLock.Acquire(commentEvent.EventID) {
							defer ch.processLock.Release(commentEvent.EventID)

							// Serialize per chat
							err := ch.pipelineManager.Run(ctx, commentEvent.ChatID, func() error {
								for _, h := range ch.onCommentHandlers {
									if err := h(ctx, commentEvent); err != nil {
										return err
									}
								}
								return nil
							})
							if err != nil {
								// handle error if needed
							}
						}
					}
				}
			}

			return nil
		})
	}
}

// OnBotAdded registers a handler for BotAddedEvent.
func (ch *channelImpl) OnBotAdded(handler func(ctx context.Context, event *types.BotAddedEvent) error) {
	ch.onBotAddedHandlers = append(ch.onBotAddedHandlers, handler)
	if ch.botAddedHandlerReg || ch.wsClient == nil {
		return
	}
	ch.botAddedHandlerReg = true
	dispatcher := ch.wsClient.EventHandler()
	if dispatcher != nil {
		dispatcher.OnP2ChatMemberBotAddedV1(func(ctx context.Context, event *larkim.P2ChatMemberBotAddedV1) error {
			if len(ch.onBotAddedHandlers) > 0 {
				botAddedEvent := normalize.ParseBotAdded(event)
				if botAddedEvent != nil {
					if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(botAddedEvent.EventID) {
						return nil
					}
					if !ch.processLock.Acquire(botAddedEvent.EventID) {
						return nil
					}
					defer ch.processLock.Release(botAddedEvent.EventID)

					// Serialize per chat
					err := ch.pipelineManager.Run(ctx, botAddedEvent.ChatID, func() error {
						for _, h := range ch.onBotAddedHandlers {
							if err := h(ctx, botAddedEvent); err != nil {
								return err
							}
						}
						return nil
					})
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	}
}

// OnReaction registers a handler for ReactionEvent.
func (ch *channelImpl) OnReaction(handler func(ctx context.Context, event *types.ReactionEvent) error) {
	ch.onReactionHandlers = append(ch.onReactionHandlers, handler)
	if ch.reactionHandlerReg || ch.wsClient == nil {
		return
	}
	ch.reactionHandlerReg = true
	dispatcher := ch.wsClient.EventHandler()
	if dispatcher != nil {
		handleReaction := func(ctx context.Context, reactionEvent *types.ReactionEvent) error {
			if reactionEvent != nil {
				if ch.dedupCache != nil && ch.dedupCache.IsDuplicate(reactionEvent.EventID) {
					return nil
				}
				if !ch.processLock.Acquire(reactionEvent.EventID) {
					return nil
				}
				defer ch.processLock.Release(reactionEvent.EventID)

				// Serialize per message
				err := ch.pipelineManager.Run(ctx, reactionEvent.MessageID, func() error {
					for _, h := range ch.onReactionHandlers {
						if err := h(ctx, reactionEvent); err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					return err
				}
			}
			return nil
		}

		dispatcher.OnP2MessageReactionCreatedV1(func(ctx context.Context, event *larkim.P2MessageReactionCreatedV1) error {
			if len(ch.onReactionHandlers) > 0 {
				return handleReaction(ctx, normalize.ParseReaction(event))
			}
			return nil
		})
		dispatcher.OnP2MessageReactionDeletedV1(func(ctx context.Context, event *larkim.P2MessageReactionDeletedV1) error {
			if len(ch.onReactionHandlers) > 0 {
				return handleReaction(ctx, normalize.ParseReaction(event))
			}
			return nil
		})
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

	if input.Card != "" {
		res, err := ch.Send(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to send initial card message: %w", err)
		}
		return NewCardStreamController(ch.client, res.MessageID), nil
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
