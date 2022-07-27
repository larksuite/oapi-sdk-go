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
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetAppAccessTokenByMarketplaceApp(context.Background(), &larkcore.MarketplaceAppAccessTokenReq{
		AppID:     appID,
		AppSecret: appSecret,
		AppTicket: "ssss",
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
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret, lark.WithLogLevel(larkcore.LogLevelDebug))
	var resp, err = client.GetTenantAccessTokenByMarketplaceApp(context.Background(), &larkcore.MarketplaceTenantAccessTokenReq{
		AppAccessToken: "ssss",
		TenantKey:      appSecret,
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

func main() {
	GetTenantAccessTokenBySelfBuiltApp()
}
