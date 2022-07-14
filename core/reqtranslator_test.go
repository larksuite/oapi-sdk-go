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
	ctx := context.Background()
	_, err := reqTranslator.translate(ctx, map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"hello<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at>\"}",
	}, AccessTokenTypeUser, config, http.MethodPost, "/open-apis/im/v1/messages?receive_id_type=open_id", &RequestOption{})

	if err != nil {
		t.Errorf("TestTranslate failed ,%v", err)
	}

	//fmt.Println(req, err)
}

func TestPathUrlEncode(t *testing.T) {
	url, _ := reqTranslator.getFullReqUrl("https://open.feishu.com", "open-apis/:a/:b/:c/:d", map[string]interface{}{"a": 12, "b": "sssss", "c": "12121wwww", "d": "加多"}, map[string]interface{}{"user_type": "open_id"})
	fmt.Println(url)

	if url != "https://open.feishu.comopen-apis/12/sssss/12121wwww/%E5%8A%A0%E5%A4%9A?user_type=open_id" {
		t.Errorf("TestPathUrlEncode failed ")
	}
}
