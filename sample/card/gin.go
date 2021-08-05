package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go"
)

var appConf *lark.AppConfig

func main() {
	appConf = lark.NewInternalAppConfigByEnv(lark.DomainFeiShu)
	appConf.SetLogLevel(lark.LogLevelDebug)

	lark.WebHook.SetCardActionHandler(appConf, func(ctx *lark.Context, cardAction *lark.CardAction) (interface{}, error) {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(lark.Prettify(cardAction))
		return "{\"config\":{\"wide_screen_mode\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"div\",\"text\":{\"tag\":\"lark_md\",\"content\":\"[飞书golang](https://www.feishu.cn)整合即时沟通、日历、音视频会议、云文档、云盘、工作台等功能于一体，成就组织和个人，更高效、更愉悦。\"}}]}}", nil
	})

	g := gin.Default()
	g.POST("/webhook/card", func(context *gin.Context) {
		lark.WebHook.CardWebServeHandler(appConf, context.Request, context.Writer)
	})
	err := g.Run(":8089")
	if err != nil {
		panic(err)
	}
}
