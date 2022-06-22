package card

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
)

func TestVerifyUrlOk(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("12", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, nil
	})

	//plainEventJsonStr := "{\"open_id\":\"ou_sdfimx9948345\",\"user_id\":\"eu_sd923r0sdf5\",\"open_message_id\":\"om_abcdefg1234567890\",\"tenant_key\":\"d32004232\",\"token\":\"12\",\"timezone\":\"\",\"action\":{\"value\":{\"tag\":\"button\",\"value\":\"sdfsfd\"},\"tag\":\"button\",\"option\":\"\",\"timezone\":\"\"},\"challenge\":\"121212\",\"type\":\"url_verification\"}"
	cardAction := mockCardAction()

	resp, err := cardHandler.AuthByChallenge(context.Background(), cardAction)
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
	}

	if resp.Body == nil {
		t.Errorf("verfiy url failed")
	}

}

func TestVerifyUrlFailed(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("121", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, nil
	})

	cardAction := mockCardAction()
	_, err := cardHandler.AuthByChallenge(context.Background(), cardAction)
	if err == nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}

}

func mockEventReq(token string) *event.EventReq {
	value := map[string]interface{}{}
	value["value"] = "sdfsfd"
	value["tag"] = "button"

	cardAction := &CardAction{
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Token:         token,
		Action: &struct {
			Value    map[string]interface{} `json:"value"`
			Tag      string                 `json:"tag"`
			Option   string                 `json:"option"`
			Timezone string                 `json:"timezone"`
		}{
			Value: value,
			Tag:   "button",
		},
	}

	cardActionBody := &CardActionBody{
		CardAction: cardAction,
		Challenge:  "121212",
		Type:       "",
	}

	body, _ := json.Marshal(cardActionBody)
	var timestamp = "timestamp"
	var nonce = "nonce"

	sign := Signature(timestamp, nonce, token, string(body))

	header := http.Header{}
	header.Set(event.EventRequestTimestamp, timestamp)
	header.Set(event.EventRequestNonce, nonce)
	header.Set(event.EventSignature, sign)
	req := &event.EventReq{
		Header: header,
		Body:   body,
	}

	return req
}

func TestVerifySignOk(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("121", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, nil
	})

	config := &core.Config{}
	core.NewLogger(config)
	cardHandler.Config = config

	req := mockEventReq("121")
	err := cardHandler.VerifySign(context.Background(), req)
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}
}

func TestVerifySignFailed(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("121", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, nil
	})

	config := &core.Config{}
	core.NewLogger(config)
	cardHandler.Config = config

	req := mockEventReq("12")
	err := cardHandler.VerifySign(context.Background(), req)
	if err == nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}
}

func TestDoHandleResultNilOk(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("12", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, nil
	})

	cardAction := mockCardAction()
	resp, err := cardHandler.DoHandle(context.Background(), cardAction)
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(core.Prettify(resp.Header))
	fmt.Println(string(resp.Body))
}

func TestDoHandleResultError(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("121", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		return nil, errors.New("im an error ")
	})

	cardAction := &CardAction{}
	_, err := cardHandler.DoHandle(context.Background(), cardAction)
	if err == nil {
		t.Errorf("handler error  ,%v", err)
		return
	}

}

func TestDoHandleResultCustomRespOk(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("12", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		body := make(map[string]interface{})
		body["content"] = "hello"

		i18n := make(map[string]string)
		i18n["zh_cn"] = "你好"
		i18n["en_us"] = "hello"
		i18n["ja_jp"] = "こんにちは"
		body["i18n"] = i18n

		resp := CustomResp{
			StatusCode: 400,
			Body:       body,
		}

		return &resp, nil
	})

	cardAction := mockCardAction()
	resp, err := cardHandler.DoHandle(context.Background(), cardAction)
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(core.Prettify(resp.Header))
	fmt.Println(string(resp.Body))
}

func mockCardAction() *CardAction {
	// 构建card，并返回
	value := map[string]interface{}{}
	value["value"] = "1111sdfsfd"
	value["tag"] = "b11111utton"
	cardAction := &CardAction{
		Type:          string(event.ReqTypeChallenge),
		Token:         "12",
		OpenID:        "ou_sdfimx9948345",
		UserID:        "eu_sd923r0sdf5",
		OpenMessageID: "om_abcdefg1234567890",
		TenantKey:     "d32004232",
		Action: &struct {
			Value    map[string]interface{} `json:"value"`
			Tag      string                 `json:"tag"`
			Option   string                 `json:"option"`
			Timezone string                 `json:"timezone"`
		}{
			Value: value,
			Tag:   "button",
		},
	}

	return cardAction
}
func TestDoHandleResultCardOk(t *testing.T) {
	// 创建card处理器
	cardHandler := NewCardActionHandler("12", "", func(ctx context.Context, cardAction *CardAction) (interface{}, error) {
		// 构建card，并返回
		value := map[string]interface{}{}
		value["value"] = "1111sdfsfd"
		value["tag"] = "b11111utton"

		cardActionResult := &CardAction{
			OpenID:        "ou_sdfimx9948345",
			UserID:        "eu_sd923r0sdf5",
			OpenMessageID: "om_abcdefg1234567890",
			TenantKey:     "d32004232",
			Action: &struct {
				Value    map[string]interface{} `json:"value"`
				Tag      string                 `json:"tag"`
				Option   string                 `json:"option"`
				Timezone string                 `json:"timezone"`
			}{
				Value: value,
				Tag:   "button",
			},
		}
		return cardActionResult, nil
	})

	resp, err := cardHandler.DoHandle(context.Background(), mockCardAction())
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
		return
	}

	fmt.Println(resp.StatusCode)
	fmt.Println(core.Prettify(resp.Header))
	fmt.Println(string(resp.Body))
}
