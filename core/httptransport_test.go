package larkcore

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestSendPost(t *testing.T) {
	config := mockConfig()
	_, err := Request(context.Background(), &HttpReq{
		HttpMethod: http.MethodPost,
		ApiPath:    "https://www.feishu.cn/approval/openapi/v2/approval/get",
		Body: map[string]interface{}{
			"approval_code": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		},
		SupportedAccessTokenTypes: []AccessTokenType{AccessTokenTypeTenant},
	}, config)

	if err != nil {
		t.Errorf("TestSendPost failed ,%v", err)
		return
	}
	fmt.Println("ok")

}
