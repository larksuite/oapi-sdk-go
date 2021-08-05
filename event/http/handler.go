package http

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/event"
	"net/http"
)

func Handle(conf core.Config, request *http.Request, response http.ResponseWriter) {
	req, err := core.ToOapiRequest(request)
	if err != nil {
		err = core.NewOapiResponseOfErr(err).WriteTo(response)
		if err != nil {
			conf.GetLogger().Error(context.TODO(), err)
		}
		return
	}
	err = event.Handle(conf, req).WriteTo(response)
	if err != nil {
		conf.GetLogger().Error(context.TODO(), err)
	}
}
