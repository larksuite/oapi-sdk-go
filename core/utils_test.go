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
	"testing"
)

func TestEncryptedEventMsg(t *testing.T) {
	en, err := EncryptedEventMsg(context.Background(), map[string]string{"key1": "value1", "key2": "value2"}, "encrypteKey")
	if err != nil {
		t.Errorf("TestEncryptedEventMsg failed ,%v", err)
	}

	if en == "" {
		t.Errorf("TestEncryptedEventMsg failed ,%v", err)
	}
}

func TestExtractAudFromURL(t *testing.T) {
	testCases := map[string]string{
		"https://open.feishu.cn/open-apis": "open.feishu.cn",
		"open.larksuite.com":               "open.larksuite.com",
		"https://fsopen.bytedance.net":     "fsopen.bytedance.net",
		"fsopen.bytedance.net/path/to/api": "fsopen.bytedance.net",
	}

	for rawURL, expected := range testCases {
		aud, err := extractAudFromURL(rawURL)
		if err != nil {
			t.Fatalf("extract aud failed for %s: %v", rawURL, err)
		}
		if aud != expected {
			t.Fatalf("unexpected aud for %s: %s", rawURL, aud)
		}
	}
}

func TestBuildProxyURL(t *testing.T) {
	proxyURL := buildProxyURL("proxy.example.com", "/v1", OAuthTokenUrlPath)
	if proxyURL != "https://proxy.example.com/v1/open-apis/authen/v2/oauth/token" {
		t.Fatalf("unexpected proxy url: %s", proxyURL)
	}
}
