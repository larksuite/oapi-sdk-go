package larkim

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

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
