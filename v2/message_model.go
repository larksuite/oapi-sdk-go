package lark

import "encoding/json"

type LanguageType string

const (
	LanguageTypeZhCN LanguageType = "zh_cn"
	LanguageTypeEnUS LanguageType = "en_us"
	LanguageTypeJaJP LanguageType = "ja_jp"
)

type MessageText struct {
	Text string `json:"text"`
}

func (m *MessageText) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessagePost map[LanguageType]*MessagePostContent

type MessagePostContent struct {
	Title   string                 `json:"title,omitempty"`
	Content [][]MessagePostElement `json:"content,omitempty"`
}

func (m *MessagePost) JSON() (string, error) {
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
	UserId string `json:"user_id,omitempty"`
}

func (m *MessagePostAt) Tag() string {
	return "at"
}

func (m *MessagePostAt) IsPost() {
}

func (m *MessagePostAt) MarshalJSON() ([]byte, error) {
	return messagePostElementJson(m)
}

type MessagePostImg struct {
	ImageKey string `json:"image_key,omitempty"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

func (m *MessagePostImg) Tag() string {
	return "img"
}

func (m *MessagePostImg) IsPost() {
}

func (m *MessagePostImg) MarshalJSON() ([]byte, error) {
	return messagePostElementJson(m)
}

type MessageShareChat struct {
	ChatId string `json:"chat_id,omitempty"`
}

func (m *MessageShareChat) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageShareUser struct {
	UserId string `json:"user_id,omitempty"`
}

func (m *MessageShareUser) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageImage struct {
	ImageKey string `json:"image_key,omitempty"`
}

func (m *MessageImage) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageAudio struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageAudio) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageVideo struct {
	FileKey  string `json:"file_key,omitempty"`
	ImageKey string `json:"image_key,omitempty"`
}

func (m *MessageVideo) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageFile struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageFile) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageSticker struct {
	FileKey string `json:"file_key,omitempty"`
}

func (m *MessageSticker) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type MessageCard struct {
	Config       *MessageCardConfig                  `json:"config,omitempty"`
	Header       *MessageCardHeader                  `json:"header,omitempty"`
	Elements     []MessageCardElement                `json:"elements,omitempty"`
	I18nElements map[LanguageType]MessageCardElement `json:"i18n_elements,omitempty"`
	CardLink     *MessageCardURL                     `json:"card_link,omitempty"`
}

func (m *MessageCard) JSON() (string, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func messageCardElementJson(e MessageCardElement) ([]byte, error) {
	data, err := structToMap(e)
	if err != nil {
		return nil, err
	}
	data["tag"] = e.Tag()
	return json.Marshal(data)
}

type MessageCardElement interface {
	Tag() string
	MarshalJSON() ([]byte, error)
}

type MessageCardActionElement interface {
	MessageCardElement
	IsAction()
}

type MessageCardExtraElement interface {
	MessageCardElement
	IsExtra()
}

type MessageCardNoteElement interface {
	MessageCardElement
	IsNote()
}

type MessageCardHr struct {
}

func (m *MessageCardHr) Tag() string {
	return "hr"
}

func (m *MessageCardHr) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardMarkdown struct {
	Content string                     `json:"content,omitempty"`
	Href    map[string]*MessageCardURL `json:"href,omitempty"`
}

func (m *MessageCardMarkdown) Tag() string {
	return "markdown"
}

func (m *MessageCardMarkdown) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardDiv struct {
	Text   MessageCardText         `json:"text,omitempty"`
	Fields []*MessageCardField     `json:"fields,omitempty"`
	Extra  MessageCardExtraElement `json:"extra,omitempty"`
}

func (m *MessageCardDiv) Tag() string {
	return "div"
}

func (m *MessageCardDiv) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardText interface {
	MessageCardElement
	Content_() string
}

type MessageCardPlainText struct {
	Content string `json:"content,omitempty"`
	Lines   *int   `json:"lines,omitempty"`
}

func (m *MessageCardPlainText) Tag() string {
	return "plain_text"
}

func (m *MessageCardPlainText) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func (m *MessageCardPlainText) Content_() string {
	return m.Content
}

func (m *MessageCardPlainText) IsExtra() {
}

func (m *MessageCardPlainText) IsNote() {
}

type MessageCardLarkMd struct {
	Content string `json:"content,omitempty"`
}

func (m *MessageCardLarkMd) Tag() string {
	return "lark_md"
}

func (m *MessageCardLarkMd) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func (m *MessageCardLarkMd) Content_() string {
	return m.Content
}

func (m *MessageCardLarkMd) IsExtra() {
}

func (m *MessageCardLarkMd) IsNote() {
}

type MessageCardImageModel string

const (
	MessageCardImageModelFitHorizontal MessageCardImageModel = "fit_horizontal"
	MessageCardImageModelCropCenter    MessageCardImageModel = "crop_center"
)

func (m MessageCardImageModel) Ptr() *MessageCardImageModel {
	return &m
}

type MessageCardImage struct {
	Alt          MessageCardText        `json:"alt,omitempty"`
	Title        MessageCardText        `json:"title,omitempty"`
	ImgKey       string                 `json:"img_key,omitempty"`
	CustomWidth  *int                   `json:"custom_width,omitempty"`
	CompactWidth *bool                  `json:"compact_width,omitempty"`
	Mode         *MessageCardImageModel `json:"mode,omitempty"`
	Preview      *bool                  `json:"preview,omitempty"`
}

func (m *MessageCardImage) Tag() string {
	return "img"
}

func (m *MessageCardImage) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func (m *MessageCardImage) IsExtra() {
}

func (m *MessageCardImage) IsNote() {
}

type MessageCardNote struct {
	Elements []MessageCardNoteElement `json:"elements,omitempty"`
}

func (m *MessageCardNote) Tag() string {
	return "note"
}

func (m *MessageCardNote) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardButtonType string

const (
	MessageCardButtonTypeDefault MessageCardButtonType = "default"
	MessageCardButtonTypePrimary MessageCardButtonType = "primary"
	MessageCardButtonTypeDanger  MessageCardButtonType = "danger"
)

func (bt MessageCardButtonType) Ptr() *MessageCardButtonType {
	return &bt
}

type MessageCardActionConfirm struct {
	Title MessageCardText `json:"title,omitempty"`
	Text  MessageCardText `json:"text,omitempty"`
}

type MessageCardButton struct {
	Text     MessageCardText           `json:"text,omitempty"`
	URL      *string                   `json:"url,omitempty"`
	MultiURL *MessageCardURL           `json:"multi_url,omitempty"`
	Type     *MessageCardButtonType    `json:"type,omitempty"`
	Value    map[string]interface{}    `json:"value,omitempty"`
	Confirm  *MessageCardActionConfirm `json:"confirm,omitempty"`
}

func (m *MessageCardButton) Tag() string {
	return "button"
}

func (m *MessageCardButton) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func (m *MessageCardButton) IsAction() {
}

func (m *MessageCardButton) IsExtra() {
}

type MessageCardDatePickerBase struct {
	InitialDate     *string                   `json:"initial_date,omitempty"`
	InitialTime     *string                   `json:"initial_time,omitempty"`
	InitialDatetime *string                   `json:"initial_datetime,omitempty"`
	Placeholder     MessageCardText           `json:"placeholder,omitempty"`
	Value           map[string]interface{}    `json:"value,omitempty"`
	Confirm         *MessageCardActionConfirm `json:"confirm,omitempty"`
}

func (m *MessageCardDatePickerBase) IsAction() {
}

func (m *MessageCardDatePickerBase) IsExtra() {
}

type MessageCardDatePicker struct {
	*MessageCardDatePickerBase
}

func (m *MessageCardDatePicker) Tag() string {
	return "date_picker"
}

func (m *MessageCardDatePicker) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardPickerTime struct {
	*MessageCardDatePickerBase
}

func (m *MessageCardPickerTime) Tag() string {
	return "picker_time"
}

func (m *MessageCardPickerTime) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardPickerDatetime struct {
	*MessageCardDatePickerBase
}

func (m *MessageCardPickerDatetime) Tag() string {
	return "picker_datetime"
}

func (m *MessageCardPickerDatetime) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardSelectOption struct {
	Text     MessageCardText        `json:"text,omitempty"`
	Value    string                 `json:"value,omitempty"`
	URL      *string                `json:"url,omitempty"`
	MultiURL *MessageCardURL        `json:"multi_url,omitempty"`
	Type     *MessageCardButtonType `json:"type,omitempty"`
}

type MessageCardOverflow struct {
	Options []*MessageCardSelectOption `json:"options,omitempty"`
	Value   map[string]interface {
	} `json:"value,omitempty"`
	Confirm *MessageCardActionConfirm `json:"confirm,omitempty"`
}

func (m *MessageCardOverflow) Tag() string {
	return "overflow"
}

func (m *MessageCardOverflow) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

func (m *MessageCardOverflow) IsAction() {
}

func (m *MessageCardOverflow) IsExtra() {
}

type MessageCardSelectMenuBase struct {
	Placeholder   MessageCardText            `json:"placeholder,omitempty"`
	InitialOption string                     `json:"initial_option,omitempty"`
	Options       []*MessageCardSelectOption `json:"options,omitempty"`
	Value         map[string]interface {
	} `json:"value"`
	Confirm *MessageCardActionConfirm `json:"confirm,omitempty"`
}

func (m *MessageCardSelectMenuBase) IsAction() {
}

func (m *MessageCardSelectMenuBase) IsExtra() {
}

type MessageCardSelectMenuStatic struct {
	*MessageCardSelectMenuBase
}

func (m *MessageCardSelectMenuStatic) Tag() string {
	return "select_static"
}

func (m *MessageCardSelectMenuStatic) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardSelectMenuPerson struct {
	*MessageCardSelectMenuBase
}

func (m *MessageCardSelectMenuPerson) Tag() string {
	return "select_person"
}

func (m *MessageCardSelectMenuPerson) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardActionLayout string

const (
	MessageCardActionLayoutDisected   MessageCardActionLayout = "bisected"
	MessageCardActionLayoutTrisection MessageCardActionLayout = "trisection"
	MessageCardActionLayoutFlow       MessageCardActionLayout = "flow"
)

func (al MessageCardActionLayout) Ptr() *MessageCardActionLayout {
	return &al
}

type MessageCardAction struct {
	Actions []MessageCardActionElement `json:"actions,omitempty"`
	Layout  *MessageCardActionLayout   `json:"layout,omitempty"`
}

func (m *MessageCardAction) Tag() string {
	return "action"
}

func (m *MessageCardAction) MarshalJSON() ([]byte, error) {
	return messageCardElementJson(m)
}

type MessageCardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode,omitempty"`
}

type MessageCardHeader struct {
	Template *string               `json:"template,omitempty"`
	Title    *MessageCardPlainText `json:"title,omitempty"`
}

type MessageCardURL struct {
	URL        string `json:"url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
	IOSURL     string `json:"ios_url,omitempty"`
	PCURL      string `json:"pc_url,omitempty"`
}

type MessageCardField struct {
	IsShort bool            `json:"is_short,omitempty"`
	Text    MessageCardText `json:"text,omitempty"`
}
