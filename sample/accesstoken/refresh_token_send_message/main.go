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
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/refreshtoken"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type envClientAssertionProvider struct{}

func (p *envClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	return &larkcore.Token{Value: os.Getenv("CLIENT_ASSERTION")}, nil
}

func main() {
	ctx := context.Background()
	client := newClient()

	tokenResp, err := client.AccessToken.Refresh(ctx, refreshtoken.NewTokenRequestBuilder().
		RefreshToken(os.Getenv("REFRESH_TOKEN")).
		Build())
	if err != nil {
		fmt.Println(err)
		return
	}
	if tokenResp.Data == nil || tokenResp.Data.AccessToken == nil {
		fmt.Println("empty user access token response")
		return
	}

	msgResp, err := sendTextMessage(ctx, client, larkcore.StringValue(tokenResp.Data.AccessToken))
	if err != nil {
		fmt.Println(err)
		return
	}
	if !msgResp.Success() {
		fmt.Println(msgResp.Code, msgResp.Msg, msgResp.RequestId())
		return
	}
	fmt.Println(larkcore.Prettify(msgResp))
}

func sendTextMessage(ctx context.Context, client *lark.Client, userAccessToken string) (*larkim.CreateMessageResp, error) {
	content := larkim.NewTextMsgBuilder().
		Text(envOrDefault("MESSAGE_TEXT", "hello from refreshed user access token")).
		Build()

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(os.Getenv("RECEIVE_ID")).
			Content(content).
			Build()).
		Build()

	return client.Im.Message.Create(ctx, req, larkcore.WithUserAccessToken(userAccessToken))
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func newClient() *lark.Client {
	options := []lark.ClientOptionFunc{}
	if os.Getenv("CLIENT_ASSERTION") != "" {
		options = append(options, lark.WithClientAssertionProvider(&envClientAssertionProvider{}))
	}
	return lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), options...)
}
