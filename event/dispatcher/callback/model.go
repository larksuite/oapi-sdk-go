package callback

import larkevent "github.com/larksuite/oapi-sdk-go/v3/event"

type CardActionTriggerEvent struct {
	*larkevent.EventV2Base                           // 事件基础数据
	*larkevent.EventReq                              // 请求原生数据
	Event                  *CardActionTriggerRequest `json:"event"` // 事件内容
}

func (m *CardActionTriggerEvent) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
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
	Value      map[string]interface{} `json:"value"`
	Tag        string                 `json:"tag"`
	Option     string                 `json:"option"`
	Timezone   string                 `json:"timezone"`
	Name       string                 `json:"name"`
	FormValue  map[string]interface{} `json:"form_value"`
	InputValue string                 `json:"input_value"`
	Options    []string               `json:"options"`
	Checked    bool                   `json:"checked"`
}

type Context struct {
	URL           string `json:"url,omitempty"`
	PreviewToken  string `json:"preview_token,omitempty"`
	OpenMessageID string `json:"open_message_id,omitempty"`
	OpenChatID    string `json:"open_chat_id,omitempty"`
}

type CardActionTriggerResponse struct {
	Toast *Toast `json:"toast,omitempty"`
	Card  *Card  `json:"card,omitempty"`
}

type Toast struct {
	Type        string            `json:"type,omitempty"`
	Content     string            `json:"content,omitempty"`
	I18nContent map[string]string `json:"i18n,omitempty"`
}

type Card struct {
	// 卡片类型, 可选值为 raw 和 template。
	//
	// raw: 由 JSON 构建的卡片。
	// template: 搭建工具构建的卡片，可视为一个卡片模板。
	Type string `json:"type,omitempty"`
	// 卡片的数据，不同的卡片类型所需填写的字段不同。
	//
	// 当 type 字段的值为 raw 时，data 中需传入卡片 JSON 的数据。参考卡片 JSON 结构按需传入数据。
	// 当 type 字段的值为 template 时，data 中可传入的字段:
	// template_id: 搭建工具中创建的卡片（也称卡片模板）的 ID，如 AAqigYkzabcef。可在搭建工具中通过复制卡片 ID 获取。
	// template_version_name: 搭建工具中创建的卡片的版本号，如 1.0.0，注意：若不填此字段，系统将默认使用该卡片的最新版本。即在搭建工具发布卡片新版本后，该新版本的卡片内容将立即对卡片 API 调用生效。
	// template_variable: 如果卡片模板内设置了变量，则可以在此处为变量名（key）赋值（value）。
	Data interface{} `json:"data,omitempty"`
}

type TemplateCard struct {
	TemplateID          string                 `json:"template_id,omitempty"`
	TemplateVariable    map[string]interface{} `json:"template_variable,omitempty"`
	TemplateVersionName string                 `json:"template_version_name,omitempty"`
}

type URLPreviewGetEvent struct {
	*larkevent.EventV2Base                       // 事件基础数据
	*larkevent.EventReq                          // 请求原生数据
	Event                  *URLPreviewGetRequest `json:"event"` // 事件内容
}

func (m *URLPreviewGetEvent) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
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
