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
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/event/core/model"
	"io/ioutil"
	"net/http"
)

const responseFormat = `{"codemsg":"%s"}`
const challengeResponseFormat = `{"challenge":"%s"}`

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
	requestID := httpEvent.HTTPRequest.Header.Get(constants.HTTPHeaderKeyRequestID)
	ctx.Set(constants.HTTPHeaderKeyRequestID, requestID)
	body, err := ioutil.ReadAll(httpEvent.HTTPRequest.Body)
	if err != nil {
		httpEvent.Err = err
		return
	}
	conf := config.ByCtx(ctx)
	conf.GetLogger().Debug(ctx, fmt.Sprintf("[unmarshal] event: %s", string(body)))
	if conf.GetAppSettings().EncryptKey != "" {
		body, err = tools.Decrypt(body, conf.GetAppSettings().EncryptKey)
		if err != nil {
			httpEvent.Err = err
			return
		}
		conf.GetLogger().Debug(ctx, fmt.Sprintf("[unmarshal] decrypt event: %s", string(body)))
	}
	httpEvent.Input = body
	notData := &model.NotData{}
	err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&notData)
	if err != nil {
		httpEvent.Err = err
		return
	}
	version := model.Version1
	token := notData.Token
	if notData.Version != "" {
		version = notData.Version
	}
	var eventType string
	if notData.Event != nil {
		eventType = notData.Event.Type
	}
	if notData.Header != nil {
		token = notData.Header.Token
		eventType = notData.Header.EventType
	}
	httpEvent.Version = version
	httpEvent.EventType = eventType
	httpEvent.Type = notData.Type
	httpEvent.Challenge = notData.Challenge
	if token != conf.GetAppSettings().VerificationToken {
		httpEvent.Err = errors.NewTokenInvalidErr()
		return
	}
}

func handlerFunc(ctx *core.Context, httpEvent *model.HTTPEvent) {
	if constants.CallbackType(httpEvent.Type) == constants.CallbackTypeChallenge {
		return
	}
	conf := config.ByCtx(ctx)
	var handler event.Handler
	if type2EventHandler, ok := event.GetType2EventHandler(conf); ok {
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
		switch e := err.(type) {
		case *NotHandlerErr:
			conf.GetLogger().Info(ctx, e.Error())
			writeHTTPResponse(httpEvent, http.StatusOK, fmt.Sprintf(responseFormat, err.Error()))
			return
		}
		writeHTTPResponse(httpEvent, http.StatusInternalServerError, fmt.Sprintf(responseFormat, err.Error()))
		conf.GetLogger().Error(ctx, err.Error())
		return
	}
	if constants.CallbackType(httpEvent.Type) == constants.CallbackTypeChallenge {
		writeHTTPResponse(httpEvent, http.StatusOK, fmt.Sprintf(challengeResponseFormat, httpEvent.Challenge))
		return
	}
	writeHTTPResponse(httpEvent, http.StatusOK, fmt.Sprintf(responseFormat, "successed"))
}

func writeHTTPResponse(httpEvent *model.HTTPEvent, code int, message string) {
	httpEvent.HTTPResponse.Header().Set(constants.ContentType, constants.DefaultContentType)
	httpEvent.HTTPResponse.WriteHeader(code)
	_, httpEvent.Err = fmt.Fprint(httpEvent.HTTPResponse, message)
}
