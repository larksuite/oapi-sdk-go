package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/card"
	cardhttpserver "github.com/larksuite/oapi-sdk-go/card/http/native"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	"net/http"
)

func main() {

	// for redis store and logrus
	// var conf = configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
	// var conf = configs.TestConfig("https://open.feishu.cn")
	var conf = configs.TestConfig(core.DomainFeiShu)

	card.SetHandler(conf, func(ctx *core.Context, card *model.Card) (interface{}, error) {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(card))
		return "{\"config\":{\"wide_screen_mode\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"div\",\"text\":{\"tag\":\"lark_md\",\"content\":\"[飞书golang](https://www.feishu.cn)整合即时沟通、日历、音视频会议、云文档、云盘、工作台等功能于一体，成就组织和个人，更高效、更愉悦。\"}}]}}", nil
	})

	cardhttpserver.Register("/webhook/card", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
