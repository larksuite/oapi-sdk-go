// Package approval code generated by oapi sdk gen
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

	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
)

// POST /open-apis/approval/v4/approvals
func main() {
	// 创建 Client
	client := lark.NewClient("appID", "appSecret")
	// 创建请求对象
	req := larkapproval.NewCreateApprovalReqBuilder().
		DepartmentIdType("open_department_id").
		UserIdType("user_id").

		ApprovalCreate(larkapproval.NewApprovalCreateBuilder().
			ApprovalName("@i18n@approval_name").
			ApprovalCode("7C468A54-8745-2245-9675-08B7C63E7A85").
			Description("@i18n@description").
			Viewers([]*larkapproval.ApprovalCreateViewers{larkapproval.NewApprovalCreateViewersBuilder().Build()}).
			Form(larkapproval.NewApprovalFormBuilder().Build()).
			NodeList([]*larkapproval.ApprovalNode{larkapproval.NewApprovalNodeBuilder().Build()}).
			Settings(larkapproval.NewApprovalSettingBuilder().Build()).
			Config(larkapproval.NewApprovalConfigBuilder().Build()).
			Icon(0).
			I18nResources([]*larkapproval.I18nResource{larkapproval.NewI18nResourceBuilder().Build()}).
			ProcessManagerIds([]string{}).
			Build()).
		Build()
	// 发起请求
	resp, err := client.Approval.Approval.Create(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
}
