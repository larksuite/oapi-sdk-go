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

package larkcore

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestTranslate(t *testing.T) {
	config := mockConfig()
	reqTranslator := ReqTranslator{}
	_, err := reqTranslator.translate(context.Background(), &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    "https://www.feishu.cn/approval/openapi/v2/approval/get",
		Body: map[string]interface{}{
			"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		}}, AccessTokenTypeTenant, config, &RequestOption{
		TenantAccessToken: "ssss",
	})

	if err != nil {
		t.Errorf("TestTranslate failed ,%v", err)
	}

}

func TestPathUrlEncode(t *testing.T) {
	url, _ := reqTranslator.getFullReqUrl("https://open.feishu.com", "open-apis/:a/:b/:c/:d", map[string]interface{}{"a": 12, "b": "sssss", "c": "12121wwww", "d": "加多"}, map[string]interface{}{"user_type": "open_id"})
	fmt.Println(url)

	if url != "https://open.feishu.comopen-apis/12/sssss/12121wwww/%E5%8A%A0%E5%A4%9A?user_type=open_id" {
		t.Errorf("TestPathUrlEncode failed ")
	}
}
