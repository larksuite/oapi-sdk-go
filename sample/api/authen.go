package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/authen/v1"
	"os"
)

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))

	ctx := context.Background()
	accessToken(ctx, larkApp)
}

func accessToken(ctx context.Context, larkApp *lark.App) {
	authenAccessTokenResp, err := authen.New(larkApp).Authen.AccessToken(ctx, &authen.AuthenAccessTokenReq{
		Body: &authen.AuthenAccessTokenReqBody{
			GrantType: lark.StringPtr("authorization_code"),
			Code:      lark.StringPtr("1234"),
		}})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", authenAccessTokenResp.RequestId())
	if authenAccessTokenResp.Code != 0 {
		fmt.Println(authenAccessTokenResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(authenAccessTokenResp.Data))
	fmt.Println()
	fmt.Println()
}
