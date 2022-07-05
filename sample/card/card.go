package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/httpserverext"
)

func getCard() *card.MessageCard {
	// config
	config := card.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := card.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := card.NewMessageCardHeader().
		Template("turquoise").
		Title(card.NewMessageCardPlainText().
			Content("[已处理] 1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement := card.NewMessageCardDiv().
		Fields([]*card.MessageCardField{card.NewMessageCardField().
			Text(card.NewMessageCardLarkMd().
				Content("**🕐 时间：**\\n2021-02-23 20:17:51").
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 谁处理了问题
	content := "✅ " + "name" + "已处理了此告警"
	processPersonElement := card.NewMessageCardDiv().
		Fields([]*card.MessageCardField{card.NewMessageCardField().
			Text(card.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	messageCard := card.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]card.MessageCardElement{divElement, processPersonElement}).
		CardLink(cardLink).
		Build()

	return messageCard
}

func getCustomResp() interface{} {
	body := make(map[string]interface{})
	body["content"] = "hello"

	i18n := make(map[string]string)
	i18n["zh_cn"] = "你好"
	i18n["en_us"] = "hello"
	i18n["ja_jp"] = "こんにちは"
	body["i18n"] = i18n

	resp := card.CustomResp{
		StatusCode: 400,
		Body:       body,
	}
	return &resp
}
func main() {

	// 创建card处理器
	cardHandler := card.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *card.CardAction) (interface{}, error) {
		fmt.Println(core.Prettify(cardAction))

		// 返回卡片消息
		return getCard(), nil

		//custom resp
		//return getCustomResp(), nil

		// 无返回值
		return nil, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, event.WithLogLevel(core.LogLevelDebug)))

	// 开发者启动服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
