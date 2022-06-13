package core

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"strings"
)

type RawResponse struct {
	StatusCode int         `json:"-"`
	Header     http.Header `json:"-"`
	RawBody    []byte      `json:"-"`
}

func (resp RawResponse) Write(writer http.ResponseWriter) {
	writer.WriteHeader(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}
	if _, err := writer.Write(resp.RawBody); err != nil {
		panic(err)
	}
}

func (resp RawResponse) JSONUnmarshalBody(val interface{}) error {
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
	contentType := resp.Header.Get(contentTypeHeader)
	body := fmt.Sprintf("<binary> len %d", len(resp.RawBody))
	if strings.Contains(contentType, "json") || strings.Contains(contentType, "text") {
		body = string(resp.RawBody)
	}
	return fmt.Sprintf("StatusCode: %d, Header:%v, Content-Type: %s, Body: %v", resp.StatusCode,
		resp.Header, resp.Header.Get(contentTypeHeader), body)
}

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  *struct {
		Details              []*CodeErrorDetail              `json:"details,omitempty"`
		PermissionViolations []*CodeErrorPermissionViolation `json:"permission_violations,omitempty"`
		FieldViolations      []*CodeErrorFieldViolation      `json:"field_violations,omitempty"`
	} `json:"error"`
}

func (ce CodeError) Error() string {
	return ce.String()
}

func (ce CodeError) String() string {
	return ""
}

type CodeErrorDetail struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type CodeErrorPermissionViolation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type CodeErrorFieldViolation struct {
	Field       string `json:"field,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

func FileNameByHeader(header http.Header) string {
	filename := ""
	_, media, _ := mime.ParseMediaType(header.Get("Content-Disposition"))
	if len(media) > 0 {
		filename = media["filename"]
	}
	return filename
}
