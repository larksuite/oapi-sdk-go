package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	authen "github.com/larksuite/oapi-sdk-go/service/authen/v1"
)

var authenService = authen.NewService(conf)

func main() {
	testAccessToken()
	testUserInfo()
}

func testAccessToken() {
	ctx := context.Background()
	coreCtx := core.WarpContext(ctx)
	body := &authen.AccessTokenReqBody{
		GrantType: "authorization_code",
		Code:      "[code]",
	}
	reqCall := authenService.Authens.AccessToken(coreCtx, body)

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		return
	}
	fmt.Println(tools.Prettify(result))
}

func testUserInfo() {
	ctx := context.Background()
	coreCtx := core.WarpContext(ctx)
	reqCall := authenService.Authens.UserInfo(coreCtx, request.SetUserAccessToken("[user_access_token]"))

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		return
	}
	fmt.Println(tools.Prettify(result))
}
