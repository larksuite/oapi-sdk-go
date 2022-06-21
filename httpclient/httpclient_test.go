package httpclient

import (
	"fmt"
	"testing"

	"github.com/larksuite/oapi-sdk-go/core"
)

func mockConfig() *core.Config {
	config := &core.Config{
		AppId:            "xxx",
		AppSecret:        "xxxx",
		Logger:           core.NewDefaultLogger(core.LogLevelDebug),
		LogLevel:         core.LogLevelDebug,
		EnableTokenCache: true,
		AppType:          core.AppTypeCustom,
		Domain:           "https://open.feishu.cn",
	}
	return config
}
func TestHttpClient(t *testing.T) {

	config := mockConfig()
	httpClient := NewHttpClient(config)
	resp, err := httpClient.Get("http://www.baidu.com")
	if err != nil {
		t.Errorf("TestHttpClient failed ,%v", err)
	}

	fmt.Println(resp.StatusCode)

}
