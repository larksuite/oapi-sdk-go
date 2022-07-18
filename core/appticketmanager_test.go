package larkcore

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type MockHttpClient struct {
}

func (client *MockHttpClient) Do(*http.Request) (*http.Response, error) {
	return nil, nil
}
func mockConfig() *Config {

	config := &Config{
		AppId:            "xxx",
		AppSecret:        "xxx",
		Logger:           newLoggerProxy(LogLevelInfo, NewEventLogger()),
		LogLevel:         LogLevelInfo,
		EnableTokenCache: true,
		HttpClient:       &MockHttpClient{},
		AppType:          AppTypeSelfBuilt,
		BaseUrl:          "https://www.baidu.com",
	}
	return config
}

func TestAppTicketManagerSetAndGet(t *testing.T) {
	config := mockConfig()
	cache := &localCache{}
	appTicketManager := AppTicketManager{cache: cache}

	err := appTicketManager.Set(context.Background(), config.AppId, "appTicketValue", time.Minute)
	if err != nil {
		t.Errorf("set key failed ,%v", err)
	}

	appTicket, err := appTicketManager.Get(context.Background(), config)
	if err != nil {
		t.Errorf("get key failed ,%v", err)
	}

	fmt.Println(appTicket)
}

//
//func TestAppTicketTimeOutAPiGet(t *testing.T) {
//	config := mockConfig()
//	cache := &localCache{}
//	appTicketManager := AppTicketManager{cache: cache}
//
//	err := appTicketManager.Set(context.Background(), config.AppId, "appTicketValue", time.Second)
//	if err != nil {
//		t.Errorf("set key failed ,%v", err)
//	}
//
//	time.Sleep(time.Second * 2)
//
//	appTicket, err := appTicketManager.Get(context.Background(), config)
//	if err != nil {
//		t.Errorf("get key failed ,%v", err)
//	}
//
//	fmt.Println(appTicket)
//}
