package main

import (
	"context"
	"fmt"
	coremodel "github.com/larksuite/oapi-sdk-go/core/model"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/event"
)

func main() {
	var conf = test.GetISVConf("staging")
	header := make(map[string][]string)
	// from http request header
	header["X-Request-Id"] = []string{"63278309j-yuewuyeu-7828389"}
	req := &coremodel.OapiRequest{
		Ctx:    context.Background(),
		Header: coremodel.NewOapiHeader(header),
		Body:   "{json}",
	}
	resp := event.Handle(conf, req)
	fmt.Println(tools.Prettify(resp))
}
