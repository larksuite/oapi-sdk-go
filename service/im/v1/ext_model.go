package larkim

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	larkevent "github.com/larksuite/oapi-sdk-go/event"
)

type ChatI18nNames struct {
	EnUs string `json:"en_us,omitempty"`
	ZhCn string `json:"zh_cn,omitempty"`
}

type P1AddBotV1Data struct {
	Type                string         `json:"type,omitempty"`
	AppID               string         `json:"app_id,omitempty"`
	ChatI18nNames       *ChatI18nNames `json:"chat_i18n_names,omitempty"`
	ChatName            string         `json:"chat_name,omitempty"`
	ChatOwnerEmployeeID string         `json:"chat_owner_employee_id,omitempty"`
	ChatOwnerName       string         `json:"chat_owner_name,omitempty"`
	ChatOwnerOpenID     string         `json:"chat_owner_open_id,omitempty"`
	OpenChatID          string         `json:"open_chat_id,omitempty"`
	OperatorEmployeeID  string         `json:"operator_employee_id,omitempty"`
	OperatorName        string         `json:"operator_name,omitempty"`
	OperatorOpenID      string         `json:"operator_open_id,omitempty"`
	OwnerIsBot          bool           `json:"owner_is_bot,omitempty"`
	TenantKey           string         `json:"tenant_key,omitempty"`
}

type P1AddBotV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1AddBotV1Data `json:"event"`
}

func (m *P1AddBotV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1RemoveBotV1Data struct {
	Type                string         `json:"type,omitempty"`
	AppID               string         `json:"app_id,omitempty"`
	ChatI18nNames       *ChatI18nNames `json:"chat_i18n_names,omitempty"`
	ChatName            string         `json:"chat_name,omitempty"`
	ChatOwnerEmployeeID string         `json:"chat_owner_employee_id,omitempty"`
	ChatOwnerName       string         `json:"chat_owner_name,omitempty"`
	ChatOwnerOpenID     string         `json:"chat_owner_open_id,omitempty"`
	OpenChatID          string         `json:"open_chat_id,omitempty"`
	OperatorEmployeeID  string         `json:"operator_employee_id,omitempty"`
	OperatorName        string         `json:"operator_name,omitempty"`
	OperatorOpenID      string         `json:"operator_open_id,omitempty"`
	OwnerIsBot          bool           `json:"owner_is_bot,omitempty"`
	TenantKey           string         `json:"tenant_key,omitempty"`
}

type P1RemoveBotV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1RemoveBotV1Data `json:"event"`
}

func (m *P1RemoveBotV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1P2PChatCreateV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1P2PChatCreateV1Data `json:"event"`
}

type P1P2PChatCreateV1Data struct {
	AppID     string        `json:"app_id,omitempty"`
	ChatID    string        `json:"chat_id,omitempty"`
	Operator  *P1OperatorV1 `json:"operator,omitempty"`
	TenantKey string        `json:"tenant_key,omitempty"`
	Type      string        `json:"type,omitempty"`
	User      *P1UserV1     `json:"user,omitempty"`
}

func (m *P1P2PChatCreateV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1OperatorV1 struct {
	OpenId string `json:"open_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
}

type P1UserV1 struct {
	OpenId string `json:"open_id,omitempty"`
	UserId string `json:"user_id,omitempty"`
	Name   string `json:"name,omitempty"`
}

type P1UserInOutChatV1Data struct {
	Type      string        `json:"type,omitempty"`
	AppID     string        `json:"app_id,omitempty"`
	ChatId    string        `json:"chat_id,omitempty"`
	Operator  *P1OperatorV1 `json:"operator,omitempty"`
	TenantKey string        `json:"tenant_key,omitempty"`
	Users     []*P1UserV1   `json:"users,omitempty"`
}

type P1UserInOutChatV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1UserInOutChatV1Data `json:"event"`
}

func (m *P1UserInOutChatV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1ChatDisbandV1Data struct {
	Type      string        `json:"type,omitempty"`
	AppID     string        `json:"app_id,omitempty"`
	ChatId    string        `json:"chat_id,omitempty"`
	Operator  *P1OperatorV1 `json:"operator,omitempty"`
	TenantKey string        `json:"tenant_key,omitempty"`
}

type P1ChatDisbandV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1ChatDisbandV1Data `json:"event"`
}

func (m *P1ChatDisbandV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1GroupSettingChangeV1 struct {
	OwnerOpenId         string `json:"owner_open_id,omitempty"`
	OwnerUserId         string `json:"owner_user_id,omitempty"`
	AddMemberPermission string `json:"add_member_permission,omitempty"`
	MessageNotification bool   `json:"message_notification,omitempty"`
}
type P1GroupSettingUpdatedV1Data struct {
	Type         string                  `json:"type,omitempty"`
	AppID        string                  `json:"app_id,omitempty"`
	ChatId       string                  `json:"chat_id,omitempty"`
	Operator     *P1OperatorV1           `json:"operator,omitempty"`
	TenantKey    string                  `json:"tenant_key,omitempty"`
	BeforeChange *P1GroupSettingChangeV1 `json:"before_change,omitempty"`
	AfterChange  *P1GroupSettingChangeV1 `json:"after_change,omitempty"`
}

type P1GroupSettingUpdatedV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1GroupSettingUpdatedV1Data `json:"event"`
}

func (m *P1GroupSettingUpdatedV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

// protocol v1的 Message_Read
type P1MessageReadV1Data struct {
	MessageIdList []string `json:"message_id_list,omitempty"`
	AppID         string   `json:"app_id"`
	OpenAppID     string   `json:"open_chat_id"`
	OpenID        string   `json:"open_id"`
	TenantKey     string   `json:"tenant_key"`
	Type          string   `json:"type"`
}

type P1MessageReadV1 struct {
	*larkevent.EventReq
	*larkevent.EventBase
	Event *P1MessageReadV1Data `json:"event"`
}

func (m *P1MessageReadV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

// protocol v1的 message
type P1MessageReceiveV1 struct {
	*larkevent.EventBase
	*larkevent.EventReq
	Event *P1MessageReceiveV1Data `json:"event"`
}

func (m *P1MessageReceiveV1) RawReq(req *larkevent.EventReq) {
	m.EventReq = req
}

type P1MessageReceiveV1Data struct {
	Type             string   `json:"type,omitempty"`
	AppID            string   `json:"app_id,omitempty"`
	TenantKey        string   `json:"tenant_key,omitempty"`
	RootID           string   `json:"root_id,omitempty"`
	ParentID         string   `json:"parent_id,omitempty"`
	OpenChatID       string   `json:"open_chat_id,omitempty"`
	ChatType         string   `json:"chat_type,omitempty"`
	MsgType          string   `json:"msg_type,omitempty"`
	OpenID           string   `json:"open_id,omitempty"`
	EmployeeID       string   `json:"employee_id,omitempty"`
	UnionID          string   `json:"union_id,omitempty"`
	OpenMessageID    string   `json:"open_message_id,omitempty"`
	IsMention        bool     `json:"is_mention,omitempty"`
	Text             string   `json:"text,omitempty"`
	TextWithoutAtBot string   `json:"text_without_at_bot,omitempty"`
	Title            string   `json:"title,omitempty"`
	ImageKeys        []string `json:"image_keys,omitempty"`
	ImageKey         string   `json:"image_key,omitempty"`
	FileKey          string   `json:"file_key,omitempty"`
}

/**
text类型消息结构化
*/
type MessageText struct {
	builder strings.Builder
}

func NewTextMsgBuilder() *MessageText {
	m := &MessageText{}
	m.builder.WriteString("{\"text\":\"")
	return m
}

func (t *MessageText) Text(text string) *MessageText {
	t.builder.WriteString(text)
	return t
}

func (t *MessageText) TextLine(text string) *MessageText {
	t.builder.WriteString(text)
	t.builder.WriteString("\\n")
	return t
}

func (t *MessageText) Line() *MessageText {
	t.builder.WriteString("\\n")
	return t
}

func (t *MessageText) AtUser(userId, name string) *MessageText {
	t.builder.WriteString("<at user_id=\\\"")
	t.builder.WriteString(userId)
	t.builder.WriteString("\\\">")
	t.builder.WriteString(name)
	t.builder.WriteString("</at>")
	return t
	return t
}

func (t *MessageText) AtAll() *MessageText {
	t.builder.WriteString("<at user_id=\\\"all\\\">")
	t.builder.WriteString("</at>")
	return t
}

func (t *MessageText) Build() string {
	t.builder.WriteString("\"}")
	return t.builder.String()
}

/**
 post类型消息结构化
**/

func NewMessagePost() *MessagePost {
	msg := MessagePost{}
	return &msg
}

func (m *MessagePost) ZhCn(zhCn *MessagePostContent) *MessagePost {
	m.ZhCN = zhCn
	return m
}

func (m *MessagePost) EnUs(enUs *MessagePostContent) *MessagePost {
	m.EnUS = enUs
	return m
}

func (m *MessagePost) JaJs(jaJp *MessagePostContent) *MessagePost {
	m.JaJP = jaJp
	return m
}
func (m *MessagePost) Build() (string, error) {
	return m.String()
}

type MessagePost struct {
	ZhCN *MessagePostContent `json:"zh_cn,omitempty"`
	EnUS *MessagePostContent `json:"en_us,omitempty"`
	JaJP *MessagePostContent `json:"ja_jp,omitempty"`
}

func NewMessagePostContent() *MessagePostContent {
	m := MessagePostContent{}
	return &m
}

func (m *MessagePostContent) ContentTitle(title string) *MessagePostContent {
	m.Title = title
	return m
}

func (m *MessagePostContent) AppendContent(postElements []MessagePostElement) *MessagePostContent {
	m.Content = append(m.Content, postElements)
	return m
}
func (m *MessagePostContent) Build() *MessagePostContent {
	return m
}

type MessagePostContent struct {
	Title   string                 `json:"title,omitempty"`
	Content [][]MessagePostElement `json:"content,omitempty"`
}

func (m *MessagePost) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessagePostElement interface {
	Tag() string
	IsPost()
	MarshalJSON() ([]byte, error)
}

func messagePostElementJson(e MessagePostElement) ([]byte, error) {
	data, err := structToMap(e)
	if err != nil {
		return nil, err
	}
	data["tag"] = e.Tag()
	return json.Marshal(data)
}

type MessagePostText struct {
	Text     string `json:"text,omitempty"`
	UnEscape bool   `json:"un_escape,omitempty"`
}

func (m *MessagePostText) Tag() string {
	return "text"
}

func (m *MessagePostText) IsPost() {
}

func (m *MessagePostText) MarshalJSON() ([]byte, error) {
	return messagePostElementJson(m)
}

type MessagePostA struct {
	Text     string `json:"text,omitempty"`
	Href     string `json:"href,omitempty"`
	UnEscape bool   `json:"un_escape,omitempty"`
}

func (m *MessagePostA) Tag() string {
	return "a"
}

func (m *MessagePostA) IsPost() {
}

func (m *MessagePostA) MarshalJSON() ([]byte, error) {
	return messagePostElementJson(m)
}

type MessagePostAt struct {
	UserId   string `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
}

func (m *MessagePostAt) Tag() string {
	return "at"
}

func (m *MessagePostAt) IsPost() {
}

func (m *MessagePostAt) MarshalJSON() ([]byte, error) {
	return messagePostElementJson(m)
}

/**
 image类型消息结构化
**/
type MessageImage struct {
	ImageKey string `json:"image_key,omitempty"`
}

func (m *MessageImage) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
文件类型消息结构化
**/
type MessageFile struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageFile) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
audio类型消息结构化
**/
type MessageAudio struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageAudio) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
media类型消息结构化
**/
type MessageMedia struct {
	FileKey  string `json:"file_key,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
}

func (m *MessageMedia) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
sticker 类型消息结构化
**/
type MessageSticker struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageSticker) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
share_chat 类型消息结构化
**/

type MessageShareChat struct {
	ChatId string `json:"chat_id,omitempty"`
}

func (m *MessageShareChat) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/**
share_user 类型消息结构化
**/

type MessageShareUser struct {
	UserId string `json:"user_id,omitempty"`
}

func (m *MessageShareUser) String() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type jsonTag struct {
	name         string
	stringFormat bool
	ignore       bool
}

func parseJSONTag(val string) (jsonTag, error) {
	if val == "-" {
		return jsonTag{ignore: true}, nil
	}
	var tag jsonTag
	i := strings.Index(val, ",")
	if i == -1 || val[:i] == "" {
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	tag = jsonTag{
		name: val[:i],
	}
	switch val[i+1:] {
	case "omitempty":
	case "omitempty,string":
		tag.stringFormat = true
	default:
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	return tag, nil
}

func structToMap(val interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	s := reflect.Indirect(reflect.ValueOf(val))
	st := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldDesc := st.Field(i)
		fieldVal := s.Field(i)
		if fieldDesc.Anonymous {
			embeddedMap, err := structToMap(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range embeddedMap {
				m[k] = v
			}
			continue
		}
		jsonTag := fieldDesc.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		tag, err := parseJSONTag(jsonTag)
		if err != nil {
			return nil, err
		}
		if tag.ignore {
			continue
		}
		if fieldDesc.Type.Kind() == reflect.Ptr && fieldVal.IsNil() {
			continue
		}
		// nil maps are treated as empty maps.
		if fieldDesc.Type.Kind() == reflect.Map && fieldVal.IsNil() {
			continue
		}
		if fieldDesc.Type.Kind() == reflect.Slice && fieldVal.IsNil() {
			continue
		}
		if tag.stringFormat {
			m[tag.name] = formatAsString(fieldVal, fieldDesc.Type.Kind())
		} else {
			m[tag.name] = fieldVal.Interface()
		}
	}
	return m, nil
}

func formatAsString(v reflect.Value, kind reflect.Kind) string {
	if kind == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return fmt.Sprintf("%v", v.Interface())
}

/**
消息类型，枚举值
*/
const (
	MsgTypeText        string = "text"
	MsgTypePost        string = "post"
	MsgTypeImage       string = "image"
	MsgTypeFile        string = "file"
	MsgTypeAudio       string = "audio"
	MsgTypeMedia       string = "media"
	MsgTypeSticker     string = "sticker"
	MsgTypeInteractive string = "interactive"
	MsgTypeShareChat   string = "share_chat"
	MsgTypeShareUser   string = "share_user"
)

const (
	ChatTypePrivate string = "private"
	ChatTypePublic  string = "public"
)
