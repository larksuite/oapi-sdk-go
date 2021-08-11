package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"net/http"
)

func main() {

	// for redis store and logrus
	// var conf = sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
	// var conf = sample.TestConfig("https://open.feishu.cn")
	var conf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)

	lark.WebHook.SetCardActionHandler(conf, func(ctx *core.Context, cardAction *lark.CardAction) (interface{}, error) {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(cardAction))
		return "{\"config\":{\"wide_screen_mode\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"div\",\"text\":{\"tag\":\"lark_md\",\"content\":\"[飞书golang](https://www.feishu.cn)整合即时沟通、日历、音视频会议、云文档、云盘、工作台等功能于一体，成就组织和个人，更高效、更愉悦。\"}}]}}", nil
	})

	lark.WebHook.CardWebServeRouter("/webhook/card", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
