//go:build go1.23

package larkim

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestListByIterYieldsNextPageRequestError(t *testing.T) {
	nextPageErr := errors.New("next page failed")
	calls := 0
	pageToken := "next-page"
	chatID := "chat-id"

	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			calls++
			if calls > 1 {
				return nil, nextPageErr
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"code":0,"data":{"items":[{"chat_id":"` + chatID + `"}],"page_token":"` + pageToken + `","has_more":true}}`))),
				Request:    req,
			}, nil
		}),
	}

	c := &chat{config: &larkcore.Config{
		AppId:        "cli_test",
		BaseUrl:      "https://open.feishu.test",
		HttpClient:   client,
		Logger:       larkcore.NewDefaultLogger(larkcore.LogLevelError),
		Serializable: &larkcore.DefaultSerialization{},
	}}
	seq, err := c.ListByIter(context.Background(), NewListChatReqBuilder().Build(), larkcore.WithTenantAccessToken("tenant-token"))
	if err != nil {
		t.Fatalf("ListByIter returned error: %v", err)
	}

	var got []*ListChat
	var gotErr error
	for item, err := range seq {
		if err != nil {
			gotErr = err
			break
		}
		got = append(got, item)
	}

	if len(got) != 1 || got[0].ChatId == nil || *got[0].ChatId != chatID {
		t.Fatalf("unexpected yielded items: %#v", got)
	}
	if !errors.Is(gotErr, nextPageErr) {
		t.Fatalf("expected next page error %v, got %v", nextPageErr, gotErr)
	}
	if calls != 2 {
		t.Fatalf("expected initial request plus next page request, got %d", calls)
	}
}
