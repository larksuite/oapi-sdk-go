package core

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestSendPost(t *testing.T) {
	config := mockConfig()
	resp, err := SendPost(context.Background(), config, "/open-apis/im/v1/messages?receive_id_type=open_id", AccessTokenTypeTenant, map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"hello<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at>\"}",
	})

	if err != nil {
		t.Errorf("TestSendPost failed ,%v", err)
		return
	}

	fmt.Println(string(resp.RawBody))
}

func TestSendRequest(t *testing.T) {
	config := mockConfig()
	resp, err := SendRequest(context.Background(), config, http.MethodPost, "/open-apis/im/v1/messages?receive_id_type=open_id", []AccessTokenType{AccessTokenTypeTenant}, map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"hello<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at>\"}",
	})

	if err != nil {
		t.Errorf("TestSendRequest failed ,%v", err)
		return
	}
	fmt.Println(string(resp.RawBody))
}
