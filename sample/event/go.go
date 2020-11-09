package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event"
	application "github.com/larksuite/oapi-sdk-go/service/application/v1"
)

func main() {
	var conf = test.GetISVConf("staging")

	application.SetAppOpenEventHandler(conf, func(coreCtx *core.Context, appOpenEvent *application.AppOpenEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appOpenEvent)
		fmt.Println(tools.Prettify(appOpenEvent))
		return nil
	})

	application.SetAppStatusChangeEventHandler(conf, func(coreCtx *core.Context, appStatusChangeEvent *application.AppStatusChangeEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(appStatusChangeEvent.Event.AppId)
		fmt.Println(appStatusChangeEvent.Event.Status)
		fmt.Println(tools.Prettify(appStatusChangeEvent))
		return nil
	})

	header := make(map[string][]string)
	// from http request header
	header["X-Request-Id"] = []string{"63278309j-yuewuyeu-7828389"}
	req := &coremodel.OapiRequest{
		Ctx:    context.Background(),
		Header: coremodel.NewOapiHeader(header),
		Body:   "{json}", // from http request body
	}
	resp := event.Handle(conf, req)
	fmt.Println(tools.Prettify(resp))
}
