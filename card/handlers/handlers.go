package handlers

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/card/model"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/errors"
	"net/http"
	"strings"
)

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
	request := httpCard.Request
	ctx.Set(constants.HTTPHeader, request.Header)
	header := &model.Header{
		Timestamp:    request.Header.GetFirstValue(model.LarkRequestTimestamp),
		Nonce:        request.Header.GetFirstValue(model.LarkRequestRequestNonce),
		Signature:    request.Header.GetFirstValue(model.LarkSignature),
		RefreshToken: request.Header.GetFirstValue(model.LarkRefreshToken),
	}
	httpCard.Header = header
	config.ByCtx(ctx).GetLogger().Debug(ctx, fmt.Sprintf("[init] card: %s", request.Body))
	httpCard.Input = []byte(request.Body)
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
			httpCard.Err = errors.NewTokenInvalidErr(out.Token)
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
	h, ok := getHandler(conf)
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

func complementFunc(ctx *core.Context, httpCard *model.HTTPCard) {
	err := httpCard.Err
	conf := config.ByCtx(ctx)
	if err != nil {
		switch e := err.(type) {
		case *NotFoundHandlerErr:
			conf.GetLogger().Info(ctx, e.Error())
			httpCard.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, err.Error()))
			return
		}
		conf.GetLogger().Error(ctx, err.Error())
		httpCard.Response.Write(http.StatusInternalServerError, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, err.Error()))
		return
	}
	if httpCard.Type == constants.CallbackTypeChallenge {
		httpCard.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ChallengeResponseFormat, httpCard.Challenge))
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
				httpCard.Response.Write(http.StatusInternalServerError, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, err.Error()))
				return
			}
		}
		httpCard.Response.Write(http.StatusOK, constants.DefaultContentType, string(bs))
		return
	}
	httpCard.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, "successed"))
}
