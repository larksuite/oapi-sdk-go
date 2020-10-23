package handlers

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const responseFormat = `{"codemsg":"%s"}`
const challengeResponseFormat = `{"challenge":"%s"}`

var defaultHandlers = &Handlers{
	init:       initFunc,
	validate:   validateFunc,
	unmarshal:  unmarshalFunc,
	handler:    handlerFunc,
	complement: complementFunc,
}

type handler func(*core.Context, *model.HTTPCard)

type Handlers struct {
	init       handler
	validate   handler
	unmarshal  handler
	handler    handler
	complement handler
}

func initFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	request := httpCard.HTTPRequest
	header := &model.Header{
		Timestamp:    request.Header.Get(model.LarkRequestTimestamp),
		Nonce:        request.Header.Get(model.LarkRequestRequestNonce),
		Signature:    request.Header.Get(model.LarkSignature),
		RefreshToken: request.Header.Get(model.LarkRefreshToken),
		RequestID:    request.Header.Get(constants.HTTPHeaderKeyRequestID),
	}
	httpCard.Header = header
	ctx.Set(model.LarkRequestTimestamp, header.Timestamp)
	ctx.Set(model.LarkRequestRequestNonce, header.Nonce)
	ctx.Set(model.LarkSignature, header.Signature)
	ctx.Set(model.LarkRefreshToken, header.RefreshToken)
	ctx.Set(constants.HTTPHeaderKeyRequestID, header.RequestID)
	body, err := ioutil.ReadAll(httpCard.HTTPRequest.Body)
	if err != nil {
		httpCard.Err = err
		return
	}
	httpCard.Input = body
	config.ByCtx(ctx).GetLogger().Debug(ctx, fmt.Sprintf("[init] card: %s", string(httpCard.Input)))
}

func validateFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	if httpCard.Header.Signature == "" {
		return
	}
	appSettings := config.ByCtx(ctx).GetAppSettings()
	err := verify(appSettings.VerificationToken, httpCard.Header, httpCard.Input)
	if err != nil {
		httpCard.Err = err
		return
	}
}

func unmarshalFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	out := &model.Challenge{}
	err := json.NewDecoder(bytes.NewBuffer(httpCard.Input)).Decode(&out)
	if err != nil {
		httpCard.Err = err
		return
	}
	httpCard.Type = constants.CallbackType(out.Type)
	httpCard.Challenge = out.Challenge
	if httpCard.Type == constants.CallbackTypeChallenge {
		appSettings := config.ByCtx(ctx).GetAppSettings()
		if appSettings.VerificationToken != out.Token {
			httpCard.Err = errors.NewTokenInvalidErr()
			return
		}
	}
}

func verify(verifyToken string, header *model.Header, body []byte) error {
	targetSig := signature(header.Nonce, header.Timestamp, string(body), verifyToken)
	if header.Signature == targetSig {
		return nil
	}
	return newSignatureErr()
}

func signature(nonce string, timestamp string, body string, token string) string {
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

func Handle(ctx *core.Context, httpCard *model.HTTPCard) {
	defer defaultHandlers.complement(ctx, httpCard)
	defaultHandlers.init(ctx, httpCard)
	if httpCard.Err != nil {
		return
	}
	defaultHandlers.validate(ctx, httpCard)
	if httpCard.Err != nil {
		return
	}
	defaultHandlers.unmarshal(ctx, httpCard)
	if httpCard.Err != nil {
		return
	}
	defaultHandlers.handler(ctx, httpCard)
}

func handlerFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	if httpCard.Type == constants.CallbackTypeChallenge {
		return
	}
	conf := config.ByCtx(ctx)
	h, ok := card.GetHandler(conf)
	if !ok {
		httpCard.Err = newNotHandlerErr()
		return
	}
	out := &model.Card{}
	err := json.NewDecoder(bytes.NewBuffer(httpCard.Input)).Decode(&out)
	if err != nil {
		httpCard.Err = err
		return
	}
	httpCard.Output, httpCard.Err = h(ctx, out)
}

func writeHTTPResponse(ctx *core.Context, httpCard *model.HTTPCard, code int, message string) {
	conf := config.ByCtx(ctx)
	httpCard.HTTPResponse.Header().Set(constants.ContentType, constants.DefaultContentType)
	httpCard.HTTPResponse.WriteHeader(code)
	_, err := httpCard.HTTPResponse.Write([]byte(message))
	if err != nil {
		conf.GetLogger().Error(ctx, err.Error())
	}
}

func complementFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	err := httpCard.Err
	conf := config.ByCtx(ctx)
	if err != nil {
		switch e := err.(type) {
		case *NotHandlerErr:
			conf.GetLogger().Info(ctx, e.Error())
			writeHTTPResponse(ctx, httpCard, http.StatusOK, fmt.Sprintf(responseFormat, err.Error()))
			return
		}
		conf.GetLogger().Error(ctx, err.Error())
		writeHTTPResponse(ctx, httpCard, http.StatusInternalServerError, fmt.Sprintf(responseFormat, err.Error()))
		return
	}
	if httpCard.Type == constants.CallbackTypeChallenge {
		writeHTTPResponse(ctx, httpCard, http.StatusOK, fmt.Sprintf(challengeResponseFormat, httpCard.Challenge))
		return
	}
	if httpCard.Output != nil {
		var bs []byte
		switch output := httpCard.Output.(type) {
		case string:
			bs = []byte(output)
		default:
			bs, err = json.Marshal(httpCard.Output)
			if err != nil {
				conf.GetLogger().Error(ctx, err.Error())
				writeHTTPResponse(ctx, httpCard, http.StatusInternalServerError, fmt.Sprintf(responseFormat, err.Error()))
				return
			}
		}
		writeHTTPResponse(ctx, httpCard, http.StatusOK, string(bs))
		return
	}
	writeHTTPResponse(ctx, httpCard, http.StatusOK, fmt.Sprintf(responseFormat, "successed"))
}
