package lark

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type cardActionHandler func(context.Context, *RawRequest, *CardAction) (interface{}, error)

func (wh *webhook) CardActionHandle(ctx context.Context, req *RawRequest) *RawResponse {
	card := httpCard{
		request:  req,
		response: &RawResponse{},
	}
	card.do(ctx, wh)
	return card.response
}

func (wh *webhook) CardActionHandler(handler cardActionHandler) {
	wh.actionHandler = handler
}

type httpCard struct {
	request  *RawRequest
	response *RawResponse
}

func (c httpCard) do(ctx context.Context, wh *webhook) {
	var err error
	var type_ webhookType
	var challenge string
	var resultBS []byte
	var result interface{}
	defer func() {
		c.response.StatusCode = http.StatusOK
		c.response.Header = map[string][]string{}
		c.response.Header.Set(contentTypeHeader, defaultContentType)
		if err != nil {
			if err == notFoundCardHandlerErr {
				c.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, err.Error()))
				return
			}
			wh.app.logger.Error(ctx, fmt.Sprintf("card action handle err: %v", err))
			c.response.StatusCode = http.StatusInternalServerError
			c.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, err.Error()))
			return
		}
		if type_ == webhookTypeChallenge {
			c.response.RawBody = []byte(fmt.Sprintf(challengeResponseFormat, challenge))
			return
		}
		if len(resultBS) > 0 {
			c.response.RawBody = resultBS
			return
		}
		c.response.RawBody = []byte(fmt.Sprintf(webhookResponseFormat, "success"))
		return
	}()
	wh.app.logger.Debug(ctx, fmt.Sprintf("card action: %v", c.request))
	out := &cardChallenge{}
	err = json.Unmarshal(c.request.RawBody, out)
	if err != nil {
		return
	}
	type_ = webhookType(out.Type)
	challenge = out.Challenge
	if type_ == webhookTypeChallenge {
		if wh.app.settings.verificationToken != out.Token {
			err = errors.New("card challenge token not equal app settings token")
			return
		}
		return
	}
	err = c.verify(wh.app)
	if err != nil {
		return
	}
	h := wh.actionHandler
	if h == nil {
		err = notFoundCardHandlerErr
		return
	}
	cardAction := &CardAction{}
	err = json.Unmarshal(c.request.RawBody, cardAction)
	if err != nil {
		return
	}
	result, err = h(ctx, c.request, cardAction)
	if err != nil {
		return
	}
	if result == nil {
		return
	}
	switch r := result.(type) {
	case string:
		resultBS = []byte(r)
	default:
		resultBS, err = json.Marshal(result)
	}
}

func (c httpCard) verify(app *App) error {
	if app.settings.verificationToken == "" {
		return nil
	}
	targetSig := c.signature(c.request.Header.Get(larkRequestNonce), c.request.Header.Get(larkRequestTimestamp),
		string(c.request.RawBody), app.settings.verificationToken)
	if c.request.Header.Get(larkSignature) == targetSig {
		return nil
	}
	return errors.New("signature error")
}

func (c httpCard) signature(nonce string, timestamp string, body string, token string) string {
	var b strings.Builder
	b.WriteString(timestamp)
	b.WriteString(nonce)
	b.WriteString(token)
	b.WriteString(body)
	bs := []byte(b.String())
	h := sha1.New()
	_, _ = h.Write(bs)
	bs = h.Sum(nil)
	return fmt.Sprintf("%x", bs)
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

type cardChallenge struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

var notFoundCardHandlerErr = errors.New("card not found handler")
