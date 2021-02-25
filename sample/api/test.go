package main

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/log"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	"net/http"
)

func main() {

	// 企业自建应用的配置
	// AppID、AppSecret: 开发者后台的应用凭证（App ID、App Secret）
	// VerificationToken、EncryptKey：开发者后台的事件订阅（Verification Token、Encrypt Key），可以为空字符串。
	appSetting := config.NewISVAppSettings("AppID", "AppSecret", "VerificationToken", "EncryptKey")

	// 当前访问的是飞书，使用默认存储、默认日志（Debug级别），更多可选配置：config.NewConfig()。
	conf := config.NewConfigWithDefaultStore(constants.DomainFeiShu, appSetting, log.NewDefaultLogger(), log.LevelInfo)

	// 启动httpServer，开发者后台的事件订阅，请求网址：https://domain/webhook/event
	eventhttpserver.Register("/webhook/event", conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}
}
