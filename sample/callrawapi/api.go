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
	"net/http"
	"os"

	larkapproval "github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 原生 API 调用推荐用法
func rawApiUserCallNew() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret)

	// 发起请求
	resp, err := cli.Do(context.Background(),
		&larkcore.ApiReq{
			HttpMethod:                http.MethodGet,
			ApiPath:                   "https://open.feishu.cn/open-apis/contact/v3/users/:user_id",
			Body:                      nil,
			QueryParams:               larkcore.QueryParams{"user_id_type": []string{"open_id"}},
			PathParams:                larkcore.PathParams{"user_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d"},
			SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeUser},
		},
	)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body,二进制数据
}

// 原生 API 调用推荐用法
func rawApiTenantCallNew() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true))

	//发起请求
	body := map[string]interface{}{}
	body["approval_code"] = "ou_c245b0a7dff2725cfa2fb104f8b48b9d"

	resp, err := cli.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodGet,
		ApiPath:                   "https://www.feishu.cn/approval/openapi/v2/approval/get",
		Body:                      body,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
	})

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body,二进制数据
}

// 老的原生调用方法，仅做兼容使用
func rawApiTenantCallOld() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret)

	// 构建消息体
	content := larkim.NewTextMsgBuilder().
		Text("加多").
		Line().
		TextLine("hello").
		TextLine("world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "陆续").
		Build()

	// 发起请求
	//1.如有 path 参数，query 参数，则需要拼接到 url 上
	//2.body 传递 http 的 body 内容
	//3.传递租户权限
	resp, err := cli.Post(context.Background(), "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id", map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    content,
	}, larkcore.AccessTokenTypeTenant)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}

// 老的原生调用方法，仅做兼容使用
func rawApiUserCallOld() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(false))

	// 发起请求
	//1.httpPath:如有 path 参数，query 参数，则需要拼接到 url 上
	//2.body:传递 http 的 body 内容
	//3.accessTokenType:传递 user 权限
	//4.options:传递 user token
	resp, err := cli.Get(context.Background(),
		"https://open.feishu.cn/open-apis/contact/v3/users/ou_c245b0a7dff2725cfa2fb104f8b48b9d?user_id_type=open_id",
		nil,
		larkcore.AccessTokenTypeTenant,
		larkcore.WithTenantAccessToken("t-cbaf253c7d6f99e1c579950d67e86e56da64a395"))

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}

// 老的原生调用方法，仅做兼容使用
func rawApiGetTokenCallOld() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithEnableTokenCache(false))

	// 发起请求
	body := map[string]interface{}{
		"app_id":     "a",
		"app_secret": "b",
	}
	resp, err := cli.Post(context.Background(),
		"/open-apis/auth/v3/tenant_access_token/internal",
		&body,
		larkcore.AccessTokenTypeNone)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}

// 原生 API 调用推荐用法
func CreateDoc() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true))

	//发起请求
	body := map[string]interface{}{}
	body["folder_token"] = "fldcnqquW1svRIYVT2Np6IuLCKd"
	body["title"] = "哈喽哈喽"
	body["ss"] = []larkapproval.NodeApprover{*larkapproval.NewNodeApproverBuilder().Key("").Value([]string{"aa"}).Build()}

	resp, err := cli.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodPost,
		ApiPath:                   "https://open.feishu.cn/open-apis/docx/v1/documents",
		Body:                      body,
		PathParams:                nil,
		QueryParams:               nil,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeApp},
	}, larkcore.WithTenantKey("sss"))

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body,二进制数据
}

// 原生 API 调用推荐用法
func GetRootFolderMeta() {
	// 创建 API Client
	var appID, appSecret = "cli_a0b9187584b8d01c", "8d0JWJrjIppJRPXDFmSmsb1IQDlm5gkv"
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true), lark.WithOpenBaseUrl("https://open.feishu-boe.cn"))

	//发起请求
	resp, err := cli.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod:                http.MethodGet,
		ApiPath:                   "/open-apis/drive/explorer/v2/root_folder/meta",
		Body:                      nil,
		PathParams:                nil,
		QueryParams:               nil,
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
	})

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求 ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body,二进制数据
}

func main() {
	GetRootFolderMeta()
}
