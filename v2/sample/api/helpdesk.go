package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"os"
)

func main() {
	appID, appSecret, helpDeskID, helpDeskToken := os.Getenv("HD_APP_ID"), os.Getenv("HD_APP_SECRET"),
		os.Getenv("HD_HELP_DESK_ID"), os.Getenv("HD_HELP_DESK_TOKEN")
	larkApp := lark.NewApp("https://open.feishu-boe.cn",
		lark.WithAppCredential(appID, appSecret),
		lark.WithAppHelpdeskCredential(helpDeskID, helpDeskToken),
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	resp, err := larkApp.SendRequest(context.TODO(), "GET", "/open-apis/helpdesk/v1/tickets/6971250929135779860",
		nil, lark.AccessTokenTypeTenant, lark.WithNeedHelpDeskAuth())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
	fmt.Println()
	fmt.Println()
}
