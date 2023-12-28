package dispatcher

import (
	"context"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
)

type CardActionTriggerEvent struct {
	*larkevent.EventV2Base                           // 事件基础数据
	*larkevent.EventReq                              // 请求原生数据
	Event                  *CardActionTriggerRequest `json:"event"` // 事件内容
}

type CardActionTriggerRequest struct {
	Operator     *Operator       `json:"operator,omitempty"`
	Token        string          `json:"token,omitempty"` // 更新卡片用的token(凭证)
	Action       *CallBackAction `json:"action,omitempty"`
	Host         string          `json:"host,omitempty"`          // 宿主: im_message/im_top_notice
	DeliveryType string          `json:"delivery_type,omitempty"` // 卡片发送渠道: url_preview/
	Context      *Context        `json:"context,omitempty"`
}

type Operator struct {
	TenantKey *string `json:"tenant_key,omitempty"`
	UserID    *string `json:"user_id,omitempty"`
	OpenID    string  `json:"open_id,omitempty"`
}

type CallBackAction struct {
	Value    map[string]interface{} `json:"value,omitempty" validate:"omitempty"`
	Tag      *string                `json:"tag,omitempty" validate:"omitempty"`
	Option   *string                `json:"option,omitempty" validate:"omitempty"`
	Timezone *string                `json:"timezone,omitempty" validate:"omitempty"`
}

type Context struct {
	URL           string `json:"url,omitempty"`
	PreviewToken  string `json:"preview_token,omitempty"`
	OpenMessageID string `json:"open_message_id,omitempty"`
	OpenChatID    string `json:"open_chat_id,omitempty"`
}

type CardActionTriggerReponse struct {
	Toast *Toast `json:"toast,omitempty"`
	Card  *Card  `json:"card,omitempty"`
}

type Toast struct {
	Type        string            `json:"type,omitempty"`
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i_18_n_content,omitempty"`
}

type Card struct {
	// template/raw
	Type string `json:"type,omitempty"`
	// type为raw时：data为larkcard.MessageCard; type为raw时，data为TemplateCard
	Data interface{} `json:"data,omitempty"`
}

type TemplateCard struct {
	TemplateID          string                 `json:"template_id,omitempty"`
	TemplateVariable    map[string]interface{} `json:"template_variable,omitempty"`
	TemplateVersionName string                 `json:"template_version_name,omitempty"`
}

func (dispatcher *EventDispatcher) OnP2CardNewProtocalCardActionTrigger(handler func(ctx context.Context, event *CardActionTriggerEvent) (*CardActionTriggerReponse, error)) *EventDispatcher {
	_, existed := dispatcher.callbackType2CallbackHandler["card.action.trigger"]
	if existed {
		panic("event: multiple handler registrations for " + "card.action.trigger")
	}
	dispatcher.callbackType2CallbackHandler["card.action.trigger"] = NewCardActionTriggerEventHandler(handler)
	return dispatcher
}

// 消息处理器定义
type CardActionTriggerEventHandler struct {
	handler func(context.Context, *CardActionTriggerEvent) (*CardActionTriggerReponse, error)
}

func NewCardActionTriggerEventHandler(handler func(context.Context, *CardActionTriggerEvent) (*CardActionTriggerReponse, error)) *CardActionTriggerEventHandler {
	h := &CardActionTriggerEventHandler{handler: handler}
	return h
}

// 返回事件的消息体的实例，用于反序列化用
func (h *CardActionTriggerEventHandler) Event() interface{} {
	return &CardActionTriggerEvent{}
}

// 回调开发者注册的handle
func (h *CardActionTriggerEventHandler) Handle(ctx context.Context, event interface{}) (interface{}, error) {
	return h.handler(ctx, event.(*CardActionTriggerEvent))
}

func (m *CardActionTriggerEvent) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}
