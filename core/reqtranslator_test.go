package larkcore

import (
	"context"
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
