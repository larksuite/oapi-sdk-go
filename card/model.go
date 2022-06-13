package card

import "encoding/json"

type MessageCard struct {
	Config       *MessageCardConfig       `json:"config,omitempty"`
	Header       *MessageCardHeader       `json:"header,omitempty"`
	Elements     []MessageCardElement     `json:"elements,omitempty"`
	I18nElements *MessageCardI18nElements `json:"i18n_elements,omitempty"`
	CardLink     *MessageCardURL          `json:"card_link,omitempty"`
}

type MessageCardI18nElements struct {
	ZhCN []MessageCardElement `json:"zh_cn,omitempty"`
	EnUS []MessageCardElement `json:"en_us,omitempty"`
	JaJP []MessageCardElement `json:"ja_jp,omitempty"`
}

type MessageCardElement interface {
	Tag() string
	MarshalJSON() ([]byte, error)
}

type MessageCardURL struct {
	URL        string `json:"url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
	IOSURL     string `json:"ios_url,omitempty"`
	PCURL      string `json:"pc_url,omitempty"`
}

type MessageCardConfig struct {
	EnableForward  *bool `json:"enable_forward,omitempty"`
	UpdateMulti    *bool `json:"update_multi,omitempty"`
	WideScreenMode *bool `json:"wide_screen_mode,omitempty"`
}

type MessageCardHeader struct {
	Template *string               `json:"template,omitempty"`
	Title    *MessageCardPlainText `json:"title,omitempty"`
}

func (m *MessageCard) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageCardPlainText struct {
	Content string                    `json:"content,omitempty"`
	Lines   *int                      `json:"lines,omitempty"`
	I18n    *MessageCardPlainTextI18n `json:"i18n,omitempty"`
}

type MessageCardPlainTextI18n struct {
	ZhCN string `json:"zh_cn,omitempty"`
	EnUS string `json:"en_us,omitempty"`
	JaJP string `json:"ja_jp,omitempty"`
}

type CardAction struct {
	OpenID        string `json:"open_id"`
	UserID        string `json:"user_id"`
	OpenMessageID string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Timezone      string `json:"timezone"`

	Action *struct {
		Value    map[string]interface{} `json:"value"`
		Tag      string                 `json:"tag"`
		Option   string                 `json:"option"`
		Timezone string                 `json:"timezone"`
	} `json:"action"`
}

type cardChallenge struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}
type CustomResp struct {
	StatusCode int
	Body       []byte
}

type CustomToastBody struct {
	Content string `json:"content"`
	I18n    *I18n  `json:"i18n"`
}

type I18n struct {
	ZhCn string `json:"zh_cn"`
	EnCn string `json:"en_us"`
	JaJp string `json:"ja_jp"`
}

type CardActionBody struct {
	*CardAction
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}
