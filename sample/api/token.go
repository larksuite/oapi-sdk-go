/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"context"
	"fmt"
	"os"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkext "github.com/larksuite/oapi-sdk-go/v3/service/ext"
)

func GetAppAccessTokenBySelfBuiltApp() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetAppAccessTokenBySelfBuiltApp(context.Background(), &larkcore.SelfBuiltAppAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func GetAppAccessTokenByMarketApp() {
	var appID, appSecret = "cli_a271fc46df38d017", "IH38Skhm4flL54qGP67BCxbSKKoqg4Y4"
	client := lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithOpenBaseUrl("https://open.larksuite-boe.com"))

	var resp, err = client.GetAppAccessTokenByMarketplaceApp(context.Background(), &larkcore.MarketplaceAppAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
		AppTicket: "g203arcEEE7CCKQHNBUNTS3RCURRR5MKFD6BJAQY",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func GetTenantAccessTokenBySelfBuiltApp() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetTenantAccessTokenBySelfBuiltApp(context.Background(), &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func GetTenantAccessTokenByMarketApp() {
	var appID, appSecret = "cli_a271fc46df38d017", "IH38Skhm4flL54qGP67BCxbSKKoqg4Y4"
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetTenantAccessTokenByMarketplaceApp(context.Background(), &larkcore.MarketplaceTenantAccessTokenReq{
		AppAccessToken: "a-g203ardk3WIVTKS3YB2TGKYGGSD3OAMSW66SSYS3",
		TenantKey:      "tenantkey",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func ResendAppTicket() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.ResendAppTicket(context.Background(), &larkcore.ResendAppTicketReq{
		AppID:     appID,
		AppSecret: appSecret,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func GetAuthenAccessToken() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithLogReqAtDebug(true))
	var resp, err = client.Ext.Authen.AuthenAccessToken(context.Background(),
		larkext.NewAuthenAccessTokenReqBuilder().
			Body(larkext.NewAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeAuthorizationCode).
				Code("b42j45f5df9d40979f19d3680e18584e").
				Build()).
			Build())
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp.Data))
}

func RefreshAuthenAccessToken() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithLogReqAtDebug(true))
	var resp, err = client.Ext.Authen.RefreshAuthenAccessToken(context.Background(),
		larkext.NewRefreshAuthenAccessTokenReqBuilder().
			Body(larkext.NewRefreshAuthenAccessTokenReqBodyBuilder().
				GrantType(larkext.GrantTypeRefreshCode).
				RefreshToken("ur-1pAKN26JRfWHMv_CnaF6Yxkh4.IM1le1io004lI00aYP").
				Build()).
			Build())
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(resp.Data.RefreshToken)

	fmt.Println(larkcore.Prettify(resp))
}

func AuthenUserInfo() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug), lark.WithLogReqAtDebug(true))
	var resp, err = client.Ext.Authen.AuthenUserInfo(context.Background(), larkcore.WithUserAccessToken("u-1i4E.f2MJ3DV3VOgZIab06014wBBg4Q3pq00k5Q001CO"))
	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
}

func main() {
	GetTenantAccessTokenBySelfBuiltApp()
}
