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

package tokenlog

import (
	"strings"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/core/accesstoken"
)

func TestSummaryOmitsTokenValues(t *testing.T) {
	summary := Summary(&accesstoken.AccessTokenRespData{
		AccessToken:  larkcore.StringPtr("access-token-secret"),
		RefreshToken: larkcore.StringPtr("refresh-token-secret"),
		ExpiresIn:    larkcore.IntPtr(7200),
	})

	if strings.Contains(summary, "access-token-secret") || strings.Contains(summary, "refresh-token-secret") {
		t.Fatalf("summary leaked token value: %s", summary)
	}
	for _, want := range []string{
		"access_token returned:true",
		"refresh_token returned:true",
		"expires_in:7200",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("summary missing %q: %s", want, summary)
		}
	}
}
