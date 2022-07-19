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
	"fmt"
	"net/http"
	"time"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func createDefaultClient() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithMarketplaceApp())
	fmt.Println(feishu_client)
}

func createClientWithLogLevel() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithLogLevel(larkcore.LogLevelDebug))
	fmt.Println(feishu_client)
}

func createClientWithAllOptions() {
	var feishu_client = lark.NewClient("appID", "appSecret",
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithOpenBaseUrl(lark.FeishuBaseUrl),
		lark.WithMarketplaceApp(),
		lark.WithReqTimeout(3*time.Second),
		lark.WithEnableTokenCache(false),
		lark.WithHelpdeskCredential("id", "token"),
		lark.WithLogger(larkcore.NewEventLogger()),
		lark.WithHttpClient(http.DefaultClient),
		lark.WithLogReqAtDebug(true))
	fmt.Println(feishu_client)

}

func main() {

}
