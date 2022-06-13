package httpserverext

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/feishu/oapi-sdk-go/card"
	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/dispatcher"
	"github.com/feishu/oapi-sdk-go/event"
)

func newReqHandler(handler event.IReqHandler, options ...event.OptionFunc) *event.ReqHandler {
	reqHandler := event.ReqHandler{IReqHandler: handler}
	switch h := handler.(type) {
	case *dispatcher.EventReqDispatcher:
		for _, option := range options {
			option(h.Config)
		}
		reqHandler.Config = h.Config
	case *card.CardActionHandler:
		for _, option := range options {
			option(h.Config)
		}
		reqHandler.Config = h.Config
	}
	return &reqHandler
}
func doProcess(writer http.ResponseWriter, req *http.Request, handler event.IReqHandler, options ...event.OptionFunc) {
	// 构建模板类
	reqHandler := newReqHandler(handler, options...)

	// 转换http请求对象为标准请求对象
	ctx := context.Background()
	eventReq, err := translate(ctx, req)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	//处理请求
	eventResp, err := reqHandler.Handle(ctx, eventReq)
	if err != nil {
		eventResp = processError(ctx, reqHandler.Config, req.RequestURI, err)
	}

	// 回写结果
	err = write(ctx, writer, eventResp)
	if err != nil {
		panic(err)
	}
}

func NewCardActionHandlerFunc(cardActionHandler event.IReqHandler, options ...event.OptionFunc) func(writer http.ResponseWriter, req *http.Request) {
	// 类型判断
	if _, ok := cardActionHandler.(*card.CardActionHandler); !ok {
		err := errors.New("cardActionHandler type not match,please pass a card.CardActionHandler instance")
		panic(err)
	}

	// 逻辑处理
	return func(writer http.ResponseWriter, req *http.Request) {
		doProcess(writer, req, cardActionHandler, options...)
	}
}

func NewEventReqHandlerFunc(eventReqDispatcher event.IReqHandler, options ...event.OptionFunc) func(writer http.ResponseWriter, req *http.Request) {
	// 类型判断
	if _, ok := eventReqDispatcher.(*dispatcher.EventReqDispatcher); !ok {
		err := errors.New("eventReqDispatcher type not match,please pass a dispatcher.eventReqDispatcher instance")
		panic(err)
	}

	// 逻辑处理
	return func(writer http.ResponseWriter, req *http.Request) {
		doProcess(writer, req, eventReqDispatcher, options...)
	}
}

func processError(ctx context.Context, config *core.Config, path string, err error) *event.EventResp {
	header := map[string][]string{}
	header[event.ContentTypeHeader] = []string{event.DefaultContentType}
	eventResp := &event.EventResp{
		Header:     header,
		Body:       []byte(fmt.Sprintf(event.WebhookResponseFormat, err.Error())),
		StatusCode: http.StatusInternalServerError,
	}
	config.Logger.Error(ctx, fmt.Sprintf("event handle err:%s, %v", path, err))
	return eventResp
}

func write(ctx context.Context, writer http.ResponseWriter, eventResp *event.EventResp) error {
	writer.WriteHeader(eventResp.StatusCode)
	for k, vs := range eventResp.Header {
		for _, v := range vs {
			writer.Header().Add(k, v)
		}
	}

	if len(eventResp.Body) > 0 {
		_, err := writer.Write(eventResp.Body)
		return err
	}
	return nil
}
func translate(ctx context.Context, req *http.Request) (*event.EventReq, error) {
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	eventReq := &event.EventReq{
		Header: req.Header,
		Body:   rawBody,
	}

	return eventReq, nil
}
