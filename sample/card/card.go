package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/httpserverext"
)

func getCard() *larkcard.MessageCard {
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := larkcard.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template("turquoise").
		Title(larkcard.NewMessageCardPlainText().
			Content("[已处理] 1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content("**🕐 时间：**\\n2021-02-23 20:17:51").
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 谁处理了问题
	content := "✅ " + "name" + "已处理了此告警"
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	messageCard := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
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

	resp := larkcard.CustomResp{
		StatusCode: 400,
		Body:       body,
	}
	return &resp
}
func main() {
	// 创建card处理器
	cardHandler := larkcard.NewCardActionHandler("v", "", func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))
		fmt.Println(cardAction.RequestId())

		// 返回卡片消息
		return getCard(), nil

		//custom resp
		//return getCustomResp(), nil

		// 无返回值
		return nil, nil
	})

	// 注册处理器
	http.HandleFunc("/webhook/card", httpserverext.NewCardActionHandlerFunc(cardHandler, larkevent.WithLogLevel(larkcore.LogLevelDebug)))

	// 启动http服务
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
