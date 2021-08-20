package lark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type RawRequest struct {
	Header  http.Header
	RawBody []byte
}

func (req RawRequest) String() string {
	return fmt.Sprintf("Header:%v, Content-Type: %s, Body: %v",
		req.Header, req.Header.Get(contentTypeHeader), string(req.RawBody))
}

type RawResponse struct {
	StatusCode int         `json:"-"`
	Header     http.Header `json:"-"`
	RawBody    []byte      `json:"-"`
}

func (resp RawResponse) JsonUnmarshal(val interface{}) error {
	if !strings.Contains(resp.Header.Get(contentTypeHeader), contentTypeJson) {
		return fmt.Errorf("response content-type not json, response: %v", resp)
	}
	return json.Unmarshal(resp.RawBody, val)
}

func (resp RawResponse) RequestId() string {
	logID := resp.Header.Get(httpHeaderKeyLogId)
	if logID != "" {
		return logID
	}
	return resp.Header.Get(httpHeaderKeyRequestId)
}

func (resp RawResponse) String() string {
	return fmt.Sprintf("StatusCode: %d, Header:%v, Content-Type: %s, Body: %v", resp.StatusCode,
		resp.Header, resp.Header.Get(contentTypeHeader), string(resp.RawBody))
}
