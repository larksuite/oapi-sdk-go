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
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

func getUserInfo() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var feishu_client = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqAtDebug(false))

	resp, err := feishu_client.Contact.User.Get(context.Background(), larkcontact.NewGetUserReqBuilder().
		UserIdType("open_id").
		UserId("ou_xxx").
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}

}

func patchUser() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var feishu_client = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqAtDebug(true))

	user := larkcontact.NewUserBuilder().Build()
	resp, err := feishu_client.Contact.User.Patch(context.Background(),
		larkcontact.NewPatchUserReqBuilder().
			UserId("ou_155184d1e73cbfb8973e5a9e698e74f2").
			UserIdType(larkcontact.UserIdTypeOpenId).
			User(user).
			Build(), larkcore.WithUserAccessToken("ssss"))

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}
}
func createUser() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var client = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqAtDebug(true))

	resp, err := client.Contact.User.Create(context.Background(),
		larkcontact.NewCreateUserReqBuilder().UserIdType(larkcontact.UserIdTypeOpenId).User(larkcontact.NewUserBuilder().Build()).Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(resp.Data.User)
	} else {
		fmt.Println(resp.Msg, resp.Code, resp.RequestId())
	}
}
func main() {
	createUser()

}
