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
