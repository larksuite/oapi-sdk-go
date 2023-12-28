package dispatcher

import (
	"context"

	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
)

type URLPreviewGetEvent struct {
	*larkevent.EventV2Base                       // 事件基础数据
	*larkevent.EventReq                          // 请求原生数据
	Event                  *URLPreviewGetRequest `json:"event"` // 事件内容
}

type URLPreviewGetRequest struct {
	Operator *Operator `json:"operator,omitempty"`
	Host     string    `json:"host,omitempty"` // 宿主: im_message/im_top_notice
	Context  *Context  `json:"context,omitempty"`
}

type URLPreviewGetResponse struct {
	Inline *Inline `json:"inline,omitempty"`
	Card   *Card   `json:"card,omitempty"`
}

type Inline struct {
	Title     string            `json:"title,omitempty"`
	I18nTitle map[string]string `json:"i18n_title,omitempty"`
	ImageKey  string            `json:"image_key,omitempty"`
	URL       *URL              `json:"url,omitempty"`
}

type URL struct {
	CopyURL string `json:"copy_url,omitempty"`
	IOS     string `json:"ios,omitempty"`
	Android string `json:"android,omitempty"`
	PC      string `json:"pc,omitempty"`
	Web     string `json:"web,omitempty"`
}

func (dispatcher *EventDispatcher) OnP2CardNewProtocalURLPreviewGet(handler func(ctx context.Context, event *URLPreviewGetEvent) (*URLPreviewGetResponse, error)) *EventDispatcher {
	_, existed := dispatcher.callbackType2CallbackHandler["url.preview.get"]
	if existed {
		panic("event: multiple handler registrations for " + "url.preview.get")
	}
	dispatcher.callbackType2CallbackHandler["url.preview.get"] = NewURLPreviewGetEventHandler(handler)
	return dispatcher
}

// 消息处理器定义
type URLPreviewGetEventHandler struct {
	handler func(context.Context, *URLPreviewGetEvent) (*URLPreviewGetResponse, error)
}

func NewURLPreviewGetEventHandler(handler func(context.Context, *URLPreviewGetEvent) (*URLPreviewGetResponse, error)) *URLPreviewGetEventHandler {
	h := &URLPreviewGetEventHandler{handler: handler}
	return h
}

// 返回事件的消息体的实例，用于反序列化用
func (h *URLPreviewGetEventHandler) Event() interface{} {
	return &URLPreviewGetEvent{}
}

// 回调开发者注册的handle
func (h *URLPreviewGetEventHandler) Handle(ctx context.Context, event interface{}) (interface{}, error) {
	return h.handler(ctx, event.(*URLPreviewGetEvent))
}

func (m *URLPreviewGetEvent) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}
