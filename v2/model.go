package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type RawRequest struct {
	Header http.Header
	body   []byte
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
	logID := resp.Header.Get(httpHeaderKeyLogID)
	if logID != "" {
		return logID
	}
	return resp.Header.Get(httpHeaderKeyRequestID)
}

func (resp RawResponse) String() string {
	return fmt.Sprintf("StatusCode: %d, Header:%v, Content-Type: %s, Body: %v", resp.StatusCode,
		resp.Header, resp.Header.Get(contentTypeHeader), string(resp.RawBody))
}

type CodeError struct {
	Code int               `json:"code"`
	Msg  string            `json:"msg"`
	Err  *InnerDetailError `json:"error"`
}

func (ce CodeError) Error() string {
	return Prettify(ce)
}

type InnerDetailError struct {
	Details              []*CodeErrorDetail              `json:"details,omitempty"`
	PermissionViolations []*CodeErrorPermissionViolation `json:"permission_violations,omitempty"`
	FieldViolations      []*CodeErrorFieldViolation      `json:"field_violations,omitempty"`
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

type CardAction struct {
	OpenID        string `json:"open_id"`
	UserID        string `json:"user_id"`
	OpenMessageID string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Timezone      string `json:"timezone"`

	Action *struct {
		Value    map[string]interface{} `json:"value"`
		Tag      string                 `json:"tag"`
		Option   string                 `json:"option"`
		Timezone string                 `json:"timezone"`
	} `json:"action"`
}

type Formdata struct {
	fields map[string]interface{}
	data   *struct {
		content     []byte
		contentType string
	}
}

func NewFormdata() *Formdata {
	return &Formdata{}
}

func (fd *Formdata) AddField(field string, val interface{}) *Formdata {
	if fd.fields == nil {
		fd.fields = map[string]interface{}{}
	}
	fd.fields[field] = val
	return fd
}

func (fd *Formdata) AddFile(field string, r io.Reader) *Formdata {
	return fd.AddField(field, r)
}

func (fd *Formdata) content() (string, []byte, error) {
	if fd.data != nil {
		return fd.data.contentType, fd.data.content, nil
	}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for key, val := range fd.fields {
		if r, ok := val.(io.Reader); ok {
			part, err := writer.CreateFormFile(key, "unknown-file")
			if err != nil {
				return "", nil, err
			}
			_, err = io.Copy(part, r)
			if err != nil {
				return "", nil, err
			}
			continue
		}
		err := writer.WriteField(key, fmt.Sprint(val))
		if err != nil {
			return "", nil, err
		}
	}
	contentType := writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		return "", nil, err
	}
	fd.data = &struct {
		content     []byte
		contentType string
	}{content: buf.Bytes(), contentType: contentType}
	return fd.data.contentType, fd.data.content, nil
}
