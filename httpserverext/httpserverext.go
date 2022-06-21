package httpserverext

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"github.com/larksuite/oapi-sdk-go/event/dispatcher"
)

func doProcess(writer http.ResponseWriter, req *http.Request, reqHandler *event.ReqHandler) {
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
		reqHandler.Logger.Error(ctx, fmt.Sprintf("write resp result error:%s", err.Error()))
	}
}

func NewCardActionHandlerFunc(cardActionHandler *card.CardActionHandler, options ...event.OptionFunc) func(writer http.ResponseWriter, req *http.Request) {
	reqHandler := card.NewReqHandlerTemplate(cardActionHandler, options...)
	return func(writer http.ResponseWriter, req *http.Request) {
		doProcess(writer, req, reqHandler)
	}
}

func NewEventHandlerFunc(eventReqDispatcher *dispatcher.EventReqDispatcher, options ...event.OptionFunc) func(writer http.ResponseWriter, req *http.Request) {
	reqHandler := dispatcher.NewReqHandlerTemplate(eventReqDispatcher, options...)
	return func(writer http.ResponseWriter, req *http.Request) {
		doProcess(writer, req, reqHandler)
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
