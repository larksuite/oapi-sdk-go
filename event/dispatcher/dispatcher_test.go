package dispatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func mockEncryptedBody(encrypteKey string) []byte {

	eventBody := ""
	en, _ := core.EncryptedEventMsg(context.Background(), eventBody, encrypteKey)
	fmt.Println(encrypteKey)

	encrypt := event.EventEncryptMsg{Encrypt: en}
	body1, _ := json.Marshal(encrypt)

	return body1
}

func mockEvent() []byte {

	eventBody := ""

	body1, _ := json.Marshal(eventBody)

	return body1
}

func TestVerifyUrlOk(t *testing.T) {
	handler := NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	//plainEventJsonStr := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"f7984f25108f8137722bb63cee927e66\",\"event_type\":\"contact.user.created_v3\",\"app_id\":\"cli_xxxxxxxx\",\"tenant_key\":\"xxxxxxx\",\"create_time\":\"1603977298000000\",\"token\":\"v\"},\"event\":{\"object\":{\"open_id\":\"ou_7dab8a3d3cdcc9da365777c7ad535d62\",\"union_id\":\"on_576833b917gda3d939b9a3c2d53e72c8\",\"user_id\":\"e33ggbyz\",\"name\":\"张三\",\"employee_no\":\"employee_no\"}},\"challenge\":\"1212\",\"type\":\"url_verification\"}"
	_, err := handler.AuthByChallenge(context.Background(), event.ReqTypeEventCallBack, "", "")
	if err != nil {
		t.Errorf("verfiy url failed ,%v", err)
	}

}

func TestVerifyUrlFailed(t *testing.T) {
	// 创建card处理器
	handler := NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})
	//plainEventJsonStr := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"f7984f25108f8137722bb63cee927e66\",\"event_type\":\"contact.user.created_v3\",\"app_id\":\"cli_xxxxxxxx\",\"tenant_key\":\"xxxxxxx\",\"create_time\":\"1603977298000000\",\"token\":\"1v\"},\"event\":{\"object\":{\"open_id\":\"ou_7dab8a3d3cdcc9da365777c7ad535d62\",\"union_id\":\"on_576833b917gda3d939b9a3c2d53e72c8\",\"user_id\":\"e33ggbyz\",\"name\":\"张三\",\"employee_no\":\"employee_no\"}},\"challenge\":\"1212\",\"type\":\"url_verification\"}"
	_, err := handler.AuthByChallenge(context.Background(), event.ReqTypeEventCallBack, "", "")
	if err == nil {
		fmt.Println(err)
		return
	}
}

func mockEventReq(token string) *event.EventReq {

	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:9999/webhook/event", nil)

	body := "{\"schema\":\"2.0\",\"header\":{\"event_id\":\"f7984f25108f8137722bb63cee927e66\",\"event_type\":\"contact.user.created_v3\",\"app_id\":\"cli_xxxxxxxx\",\"tenant_key\":\"xxxxxxx\",\"create_time\":\"1603977298000000\",\"token\":\"v\"},\"event\":{\"object\":{\"open_id\":\"ou_7dab8a3d3cdcc9da365777c7ad535d62\",\"union_id\":\"on_576833b917gda3d939b9a3c2d53e72c8\",\"user_id\":\"e33ggbyz\",\"name\":\"张三\",\"employee_no\":\"employee_no\"}},\"challenge\":\"1212\",\"type\":\"url_verification\"}"

	var timestamp = "timestamp"
	var nonce = "nonce"
	sourceSign := event.Signature(timestamp, nonce, token, string(body))

	// 添加header
	req.Header.Set(event.EventRequestTimestamp, timestamp)
	req.Header.Set(event.EventRequestNonce, nonce)
	req.Header.Set(event.EventSignature, sourceSign)

	eventReq := event.EventReq{
		Header: req.Header,
		Body:   []byte(body),
	}

	return &eventReq
}

func TestParseReq(t *testing.T) {
	// 创建card处理器
	handler := NewEventDispatcher("", "").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	config := &core.Config{}
	core.NewLogger(config)
	handler.Config = config

	// mock请求
	req := mockEventReq("121")
	resp, err := handler.ParseReq(context.Background(), req)
	if err != nil {
		t.Errorf("TestParseReq failed ,%v", err)
		return
	}

	fmt.Println(resp)
}

func TestDecryptEvent(t *testing.T) {
	// 创建card处理器
	handler := NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	config := &core.Config{}
	core.NewLogger(config)
	handler.Config = config

	resp, err := handler.DecryptEvent(context.Background(), "bZ3L7yh6m3Fkuffl4g+3uGjIHnhOm5fVZbKVuyT8t7tcd5ABMYm8l28/X900ZL3knZ7n+sCREu/H2WnIzft0amC7+xWqNH8o25IU63N4BnZWfHh+4hyG76QPd19vkw2bPJCx9aqxK8Nz+xqFNbk0RdgyWhmgd30jSxHtcQXAllkI7FMpGpOCteJad3bLXPDBQIV/xkCtKICCS7Z63gakpxZCLaRZ3qCXP1fapHh+LBIupxenrU6ysc7I3nHmjmKie41IiWwS5puG4zQHhVbq6KWLcgWm/3NBZOPQy53ucMu75SXA55I7jarVLZXWUcqBGrcgE3vouWbtwgZuzmoTQl0GSh5VYSVvpW992BuGxUWj0XjPYdICJm6Cr7xouNXwMcdb7N8caVdkdSZeEnswG19qSyDoQhklwzNGW0yiaayulBqJNjfge/G5V3401c2XaIuAeEIo+QQ4RSNpRGfnHkbu/j55FGQAGWjpuBNaIwZbaUoVP3NkGP+vM5rpEDe3sL2GN+Xsd+g9yBs7FqdMV8mXTGgLjCqjrPrke5/km76Q3Pe6KPs2YexMRG4MkSx3xUTzZnNn7zIzShPcjeSwBd2pxk6ht5N+fzueZdxl6Oo=")
	if err != nil {
		t.Errorf("TestDecryptEvent failed ,%v", err)
		return
	}

	fmt.Println(resp)
}

func TestVerifySignOk(t *testing.T) {
	// 创建card处理器
	handler := NewEventDispatcher("v", "1212121212").OnMessageReceiveV1(func(ctx context.Context, event *larkim.MessageReceiveEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	}).OnUserCreatedV3(func(ctx context.Context, event *larkcontact.UserCreatedEvent) error {
		fmt.Println(core.Prettify(event))
		return nil
	})

	config := &core.Config{}
	core.NewLogger(config)
	handler.Config = config

	req := mockEventReq("1212121212")
	err := handler.VerifySign(context.Background(), req)
	if err != nil {
		t.Errorf("TestVerifySignOk failed ,%v", err)
		return
	}
}
func TestAppTicket(t *testing.T) {

	event := appTicketEvent{
		EventBase: &event.EventBase{
			Ts:    "",
			UUID:  "",
			Token: "1212121212",
			Type:  "",
		},
		Event: &appTicketEventData{
			AppId:     "jiaduoappId",
			Type:      "app_ticket",
			AppTicket: "AppTicketvalue",
		},
	}

	body, _ := json.Marshal(event)
	fmt.Println(string(body))

}
