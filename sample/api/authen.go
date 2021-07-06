package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	authen "github.com/larksuite/oapi-sdk-go/service/authen/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var authenService = authen.NewService(configs.TestConfig(core.DomainFeiShu))

func main() {
	testAccessToken()
	//testFlushAccessToken()
	//testUserInfo()
}

func testAccessToken() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	body := &authen.AuthenAccessTokenReqBody{
		GrantType: "authorization_code",
		Code:      "476Bsaz9mCDIAOmjIOjD4a",
	}
	reqCall := authenService.Authens.AccessToken(coreCtx, body)

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(result))
}

func testFlushAccessToken() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	body := &authen.AuthenRefreshAccessTokenReqBody{
		GrantType:    "refresh_token",
		RefreshToken: "[refresh_token]",
	}
	reqCall := authenService.Authens.RefreshAccessToken(coreCtx, body)

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(result))
}

func testUserInfo() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	reqCall := authenService.Authens.UserInfo(coreCtx, request.SetUserAccessToken("[user_access_token]"))

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(result))
}
