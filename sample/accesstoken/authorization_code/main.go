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

	"github.com/larksuite/oapi-sdk-go/sample/accesstoken/tokenlog"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken/authorizationcode"
)

type envClientAssertionProvider struct{}

func (p *envClientAssertionProvider) RetrieveToken(ctx context.Context, aud string) (*larkcore.Token, error) {
	return &larkcore.Token{Value: os.Getenv("CLIENT_ASSERTION")}, nil
}

func main() {
	client := newClient()

	req := authorizationcode.NewTokenRequestBuilder().
		Code(os.Getenv("OAUTH_CODE")).
		RedirectUri(os.Getenv("REDIRECT_URI")).
		CodeVerifier(os.Getenv("CODE_VERIFIER")).
		Build()

	resp, err := client.AccessToken.RetrieveByAuthorizationCode(context.Background(), req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !resp.Success() {
		fmt.Println(resp.StatusCode, resp.RequestId())
		return
	}
	fmt.Println(tokenlog.Summary(resp.Data))
}

func newClient() *lark.Client {
	options := []lark.ClientOptionFunc{}
	if os.Getenv("CLIENT_ASSERTION") != "" {
		options = append(options, lark.WithClientAssertionProvider(&envClientAssertionProvider{}))
	}
	return lark.NewClient(os.Getenv("APP_ID"), os.Getenv("APP_SECRET"), options...)
}
