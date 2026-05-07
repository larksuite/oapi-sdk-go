package channel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/channel/normalize"
	"github.com/larksuite/oapi-sdk-go/v3/channel/outbound"
	"github.com/larksuite/oapi-sdk-go/v3/channel/pipeline"
	"github.com/larksuite/oapi-sdk-go/v3/channel/safety"
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
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

	botIdentity *types.BotIdentity
	botMu       sync.Mutex

	// Handler registries
	onMessageHandlers    []func(ctx context.Context, msg *types.NormalizedMessage) error
	onCommentHandlers    []func(ctx context.Context, event *types.CommentEvent) error
	onReactionHandlers   []func(ctx context.Context, event *types.ReactionEvent) error
	onBotAddedHandlers   []func(ctx context.Context, event *types.BotAddedEvent) error
	onCardActionHandlers []func(ctx context.Context, msg *types.NormalizedMessage) error

	onErrorHandlers        []func(err error)
	onReconnectingHandlers []func()
	onReconnectedHandlers  []func()
	onDisconnectedHandlers []func()

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

// GetBotIdentity fetches and caches the bot's identity from the server.
func (ch *channelImpl) GetBotIdentity(ctx context.Context) *types.BotIdentity {
	if ch.botIdentity != nil {
		return ch.botIdentity
	}

	ch.botMu.Lock()
	defer ch.botMu.Unlock()
	if ch.botIdentity != nil {
		return ch.botIdentity
	}

	// Fetch bot info using the official SDK raw method since bot/v3/info is not generated
	resp, err := ch.client.Get(ctx, "/open-apis/bot/v3/info", nil, larkcore.AccessTokenTypeTenant)
	if err == nil && resp != nil && resp.StatusCode == 200 {
		var result struct {
			Code int `json:"code"`
			Data struct {
				Bot struct {
					AppId  string `json:"app_id"`
					OpenId string `json:"open_id"`
				} `json:"bot"`
			} `json:"data"`
		}
		if err := json.Unmarshal(resp.RawBody, &result); err == nil && result.Code == 0 {
			ch.botIdentity = &types.BotIdentity{
				OpenID: result.Data.Bot.OpenId,
				AppID:  result.Data.Bot.AppId,
			}
		}
	}

	return ch.botIdentity
}

// OnError registers a handler for WS error events.
func (ch *channelImpl) OnError(handler func(err error)) {
	ch.onErrorHandlers = append(ch.onErrorHandlers, handler)
}

// OnReconnecting registers a handler for WS reconnecting events.
func (ch *channelImpl) OnReconnecting(handler func()) {
	ch.onReconnectingHandlers = append(ch.onReconnectingHandlers, handler)
}

// OnReconnected registers a handler for WS reconnected events.
func (ch *channelImpl) OnReconnected(handler func()) {
	ch.onReconnectedHandlers = append(ch.onReconnectedHandlers, handler)
}

// OnDisconnected registers a handler for WS disconnected events.
func (ch *channelImpl) OnDisconnected(handler func()) {
	ch.onDisconnectedHandlers = append(ch.onDisconnectedHandlers, handler)
}

// UpdatePolicy updates the policy configuration for the channel.
func (ch *channelImpl) UpdatePolicy(cfg types.PolicyConfig) {
	ch.policyGate.UpdateConfig(cfg)
}

// GetPolicy returns the current policy configuration.
func (ch *channelImpl) GetPolicy() types.PolicyConfig {
	return ch.policyGate.GetConfig()
}

// Start starts the underlying WebSocket client and wires up lifecycle events.
func (ch *channelImpl) Start(ctx context.Context) error {
	if ch.wsClient == nil {
		larkcore.NewEventLogger().Info(ctx, "[Channel] Start called but wsClient is nil, skipping WebSocket connection.")
		return nil
	}
	ch.wsClient.SetOnError(func(err error) {
		for _, h := range ch.onErrorHandlers {
			h(err)
		}
	})
	ch.wsClient.SetOnReconnecting(func() {
		for _, h := range ch.onReconnectingHandlers {
			h()
		}
	})
	ch.wsClient.SetOnReconnected(func() {
		for _, h := range ch.onReconnectedHandlers {
			h()
		}
	})
	ch.wsClient.SetOnDisconnected(func() {
		for _, h := range ch.onDisconnectedHandlers {
			h()
		}
	})
	return ch.wsClient.Start(ctx)
}

// Stop gracefully stops the underlying WebSocket client.
func (ch *channelImpl) Stop(ctx context.Context) error {
	if ch.wsClient != nil {
		ch.wsClient.Close()
	}
	return nil
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
					botInfo := ch.GetBotIdentity(ctx)
					if botInfo != nil {
						// 1. Self-reply loop prevention
						if normMsg.UserID == botInfo.OpenID {
							return nil
						}
						// 2. MentionedBot check
						for _, m := range normMsg.Mentions {
							if m.UserID == botInfo.OpenID {
								normMsg.MentionedBot = true
								break
							}
						}
					}

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

// DownloadFile downloads media by key and type (e.g., "image", "file").
func (ch *channelImpl) DownloadFile(ctx context.Context, fileKey string, mediaType string) ([]byte, error) {
	if fileKey == "" {
		return nil, fmt.Errorf("fileKey cannot be empty")
	}

	if mediaType == "image" {
		req := larkim.NewGetImageReqBuilder().
			ImageKey(fileKey).
			Build()
		resp, err := ch.client.Im.V1.Image.Get(ctx, req)
		if err != nil {
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("download image API error: %d - %s", resp.Code, resp.Msg)
		}
		// Write the stream to byte array
		var buf bytes.Buffer
		_, err = buf.ReadFrom(resp.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read image stream: %w", err)
		}
		return buf.Bytes(), nil

	} else if mediaType == "file" || mediaType == "audio" || mediaType == "video" || mediaType == "media" {
		req := larkim.NewGetFileReqBuilder().
			FileKey(fileKey).
			Build()
		resp, err := ch.client.Im.V1.File.Get(ctx, req)
		if err != nil {
			return nil, err
		}
		if !resp.Success() {
			return nil, fmt.Errorf("download file API error: %d - %s", resp.Code, resp.Msg)
		}
		var buf bytes.Buffer
		_, err = buf.ReadFrom(resp.File)
		if err != nil {
			return nil, fmt.Errorf("failed to read file stream: %w", err)
		}
		return buf.Bytes(), nil
	}

	return nil, fmt.Errorf("unsupported mediaType: %s", mediaType)
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
