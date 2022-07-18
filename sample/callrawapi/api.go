package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 原生 API 调用推荐用法
func rawApiUserCallNew() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true), lark.WithLogLevel(larkcore.LogLevelDebug))

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
		larkcore.WithUserAccessToken("u-3Sr1oTO4V1FWxTFTFYuFCqhk2Vs4h5IbhMG00gmw0CXh"),
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
	fmt.Println(string(resp.RawBody)) // http body
}

// 原生 API 调用推荐用法
func rawApiTenantCallNew() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret, lark.WithLogReqAtDebug(true))

	//发起请求
	resp, err := cli.Do(context.Background(), &larkcore.ApiReq{
		HttpMethod: http.MethodGet,
		ApiPath:    "https://www.feishu.cn/approval/openapi/v2/approval/get",
		Body: map[string]interface{}{
			"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		},
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
	fmt.Println(resp.StatusCode) // http status code
	fmt.Println(resp.Header)     // http header
	fmt.Println(resp.RawBody)    // http body
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
	var cli = lark.NewClient(appID, appSecret)

	// 发起请求
	//1.httpPath:如有 path 参数，query 参数，则需要拼接到 url 上
	//2.body:传递 http 的 body 内容
	//3.accessTokenType:传递 user 权限
	//4.options:传递 user token
	resp, err := cli.Get(context.Background(),
		"https://open.feishu.cn/open-apis/contact/v3/users/ou_c245b0a7dff2725cfa2fb104f8b48b9d?user_id_type=open_id",
		nil,
		larkcore.AccessTokenTypeUser,
		larkcore.WithUserAccessToken("u-23P_0Vu5JdeGGShR0dw7.f1hmckRk5KzNww0g0wawHAU"))

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
func main() {
	rawApiUserCallNew()
}
