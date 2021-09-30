package lark

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type CustomerBot struct {
	webhook string
	secret  string
}

type CustomerBotSendMessageResp struct {
	*RawResponse `json:"-"`
	CodeError
}

type customerBotSendMessageReq struct {
	Timestamp string      `json:"timestamp,omitempty"`
	Sign      string      `json:"sign,omitempty"`
	MsgType   string      `json:"msg_type,omitempty"`
	Content   interface{} `json:"content"`
}

func (c *CustomerBot) sign(timestamp string) (string, error) {
	sign := fmt.Sprintf("%s\n%s", timestamp, c.secret)
	var data []byte
	h := hmac.New(sha256.New, []byte(sign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (c *CustomerBot) newSendMessageReq(msgType string, content interface{}) (*customerBotSendMessageReq, error) {
	timestamp, sign := "", ""
	if c.secret != "" {
		var err error
		timestamp = strconv.FormatInt(time.Now().Unix(), 10)
		sign, err = c.sign(timestamp)
		if err != nil {
			return nil, err
		}
	}
	return &customerBotSendMessageReq{
		Timestamp: timestamp,
		Sign:      sign,
		MsgType:   msgType,
		Content:   content,
	}, nil
}

func (c *CustomerBot) SendMessage(ctx context.Context, msgType string, content interface{}) (*CustomerBotSendMessageResp, error) {
	req, err := c.newSendMessageReq(msgType, content)
	if err != nil {
		return nil, err
	}
	reqBs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhook, bytes.NewBuffer(reqBs))
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set(userAgentHeader, userAgent())
	httpRequest.Header.Set(contentTypeHeader, defaultContentType)
	rawResp, err := sendHTTPRequest(httpRequest)
	if err != nil {
		return nil, err
	}
	codeError := CodeError{}
	err = json.Unmarshal(rawResp.RawBody, &codeError)
	if err != nil {
		return nil, err
	}
	return &CustomerBotSendMessageResp{
		RawResponse: rawResp,
		CodeError:   codeError,
	}, nil
}

func NewCustomerBot(webhook string, secret string) *CustomerBot {
	return &CustomerBot{
		webhook: webhook,
		secret:  secret,
	}
}
