package main

import (
	"context"
	"fmt"
	"os"

	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func rawApiTenantCall1() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret)

	//发起请求
	//1.支持httpPath传递完整请求url
	//2.body 传递http的body内容
	//3.传递租户权限
	resp, err := cli.Post(context.Background(), "https://www.feishu.cn/approval/openapi/v2/approval/get", map[string]interface{}{
		"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
	}, larkcore.AccessTokenTypeTenant)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取请求ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode) // http status code
	fmt.Println(resp.Header)     // http header
	fmt.Println(resp.RawBody)    // http body
}

func rawApiTenantCall2() {
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
	//1.如有path参数，query参数，则需要拼接到url上
	//2.body 传递http的body内容
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

	// 获取请求ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}

func rawApiUserCall1() {
	// 创建 API Client
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var cli = lark.NewClient(appID, appSecret)

	// 发起请求
	//1.httpPath:如有path参数，query参数，则需要拼接到url上
	//2.body:传递http的body内容
	//3.accessTokenType:传递User权限
	//4.options:传递User token
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

	// 获取请求ID
	fmt.Println(resp.RequestId())

	// 处理请求结果
	fmt.Println(resp.StatusCode)      // http status code
	fmt.Println(resp.Header)          // http header
	fmt.Println(string(resp.RawBody)) // http body
}

func main() {
	rawApiUserCall1()
}
