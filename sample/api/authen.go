package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/sample"
	authen "github.com/larksuite/oapi-sdk-go/service/authen/v1"
)

// for redis store and logrus
// sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// sample.TestConfig("https://open.feishu.cn")
var authenService = authen.NewService(sample.TestConfig(lark.DomainFeiShu))

func main() {
	testAccessToken()
	//testFlushAccessToken()
	//testUserInfo()
}

func testAccessToken() {
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	body := &authen.AuthenAccessTokenReqBody{
		GrantType: "authorization_code",
		Code:      "476Bsaz9mCDIAOmjIOjD4a",
	}
	reqCall := authenService.Authens.AccessToken(coreCtx, body)

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result))
}

func testFlushAccessToken() {
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	body := &authen.AuthenRefreshAccessTokenReqBody{
		GrantType:    "refresh_token",
		RefreshToken: "[refresh_token]",
	}
	reqCall := authenService.Authens.RefreshAccessToken(coreCtx, body)

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result))
}

func testUserInfo() {
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	reqCall := authenService.Authens.UserInfo(coreCtx, lark.SetUserAccessToken("[user_access_token]"))

	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result))
}
