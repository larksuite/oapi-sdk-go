package larkcore

import (
	"context"
	"net/http"
	"testing"
)

func TestSendPost(t *testing.T) {
	config := mockConfig()
	_, err := SendRequest(context.Background(), config, "POST", "/", []AccessTokenType{AccessTokenTypeUser}, map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"hello<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at>\"}",
	}, WithUserAccessToken("121"))

	if err != nil {
		t.Errorf("TestSendPost failed ,%v", err)
		return
	}

	//fmt.Println(string(resp.RawBody))
}

func TestSendRequest(t *testing.T) {
	config := mockConfig()
	_, err := SendRequest(context.Background(), config, http.MethodPost, "/", []AccessTokenType{AccessTokenTypeUser}, map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"hello<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at>\"}",
	}, WithUserAccessToken("121"))

	if err != nil {
		t.Errorf("TestSendRequest failed ,%v", err)
		return
	}
	//fmt.Println(string(resp.RawBody))
}
