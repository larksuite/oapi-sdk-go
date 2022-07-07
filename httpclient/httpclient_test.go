package httpclient

import (
	"fmt"
	"net/http"
	"testing"
)

type CustomHttpClient struct {
}

func (client *CustomHttpClient) Do(*http.Request) (*http.Response, error) {
	return nil, nil
}

func TestHttpClient(t *testing.T) {

	httpClient := &CustomHttpClient{}

	req, _ := http.NewRequest("GET", "http://www.baidu.com", nil)
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Errorf("TestHttpClient failed ,%v", err)
	}

	fmt.Println(resp)

}
