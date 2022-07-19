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

func TestSendPost(t *testing.T) {
	config := mockConfig()
	_, err := Request(context.Background(), &ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    "/",
		Body: map[string]interface{}{
			"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		},
		SupportedAccessTokenTypes: []AccessTokenType{AccessTokenTypeUser},
	}, config, WithUserAccessToken("key"))

	if err != nil {
		t.Errorf("TestSendPost failed ,%v", err)
		return
	}
	fmt.Println("ok")

}
