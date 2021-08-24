package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"github.com/larksuite/oapi-sdk-go/core/errors"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event/core/model"

	"net/http"
)

var defaultHandlers = &Handlers{
	unmarshal:  unmarshalFunc,
	handler:    handlerFunc,
	complement: complementFunc,
}

type handler func(*core.Context, *model.HTTPEvent)

type Handlers struct {
	unmarshal  handler
	handler    handler
	complement handler
}

func Handle(ctx *core.Context, httpEvent *model.HTTPEvent) {
	defer defaultHandlers.complement(ctx, httpEvent)
	if httpEvent.Err != nil {
		return
	}
	defaultHandlers.unmarshal(ctx, httpEvent)
	if httpEvent.Err != nil {
		return
	}
	defaultHandlers.handler(ctx, httpEvent)
}

func unmarshalFunc(ctx *core.Context, httpEvent *model.HTTPEvent) {
	request := httpEvent.Request
	ctx.Set(constants.HTTPHeader, request.Header)
	conf := config.ByCtx(ctx)
	conf.GetLogger().Debug(ctx, fmt.Sprintf("[unmarshal] event: %s", request.Body))
	body := []byte(request.Body)
	var err error
	if conf.GetAppSettings().EncryptKey != "" {
		body, err = tools.Decrypt(body, conf.GetAppSettings().EncryptKey)
		if err != nil {
			httpEvent.Err = err
			return
		}
		conf.GetLogger().Debug(ctx, fmt.Sprintf("[unmarshal] decrypt event: %s", string(body)))
	}
	httpEvent.Input = body
	fuzzy := &model.Fuzzy{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&fuzzy)
	if err != nil {
		httpEvent.Err = err
		return
	}
	schema := model.Version1
	token := fuzzy.Token
	if fuzzy.Schema != "" {
		schema = fuzzy.Schema
	}
	var eventType string
	if fuzzy.Event != nil {
		eventType = fuzzy.Event.Type
	}
	if fuzzy.Header != nil {
		token = fuzzy.Header.Token
		eventType = fuzzy.Header.EventType
	}
	httpEvent.Schema = schema
	httpEvent.EventType = eventType
	httpEvent.Type = fuzzy.Type
	httpEvent.Challenge = fuzzy.Challenge
	if token != conf.GetAppSettings().VerificationToken {
		httpEvent.Err = errors.NewTokenInvalidErr(token)
		return
	}
}

func handlerFunc(ctx *core.Context, httpEvent *model.HTTPEvent) {
	if constants.CallbackType(httpEvent.Type) == constants.CallbackTypeChallenge {
		return
	}
	conf := config.ByCtx(ctx)
	var handler Handler
	if type2EventHandler, ok := getType2EventHandler(conf); ok {
		h, ok := type2EventHandler[httpEvent.EventType]
		if ok {
			handler = h
		}
	}
	if handler == nil {
		httpEvent.Err = newNotHandlerErr(httpEvent.EventType)
		return
	}
	e := handler.GetEvent()
	err := json.NewDecoder(bytes.NewBuffer(httpEvent.Input)).Decode(e)
	if err != nil {
		httpEvent.Err = err
		return
	}
	err = handler.Handle(ctx, e)
	httpEvent.Err = err
}

func complementFunc(ctx *core.Context, httpEvent *model.HTTPEvent) {
	conf := config.ByCtx(ctx)
	err := httpEvent.Err
	if err != nil {
		if _, ok := err.(*NotFoundHandlerErr); ok {
			httpEvent.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, err.Error()))
			return
		}
		httpEvent.Response.Write(http.StatusInternalServerError, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, err.Error()))
		conf.GetLogger().Error(ctx, err.Error())
		return
	}
	if constants.CallbackType(httpEvent.Type) == constants.CallbackTypeChallenge {
		httpEvent.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ChallengeResponseFormat, httpEvent.Challenge))
		return
	}
	httpEvent.Response.Write(http.StatusOK, constants.DefaultContentType, fmt.Sprintf(core.ResponseFormat, "successed"))
}
