package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	eventhttpserver "github.com/larksuite/oapi-sdk-go/event/http/native"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
	"net/http"
	"path"
)

func main() {

	var conf = test.GetISVConf("online")

	contact.SetUserCreateEventHandler(conf, func(coreCtx *core.Context, event *contact.UserCreateEvent) error {
		fmt.Println(coreCtx.GetRequestID())
		fmt.Println(event)
		fmt.Println(tools.Prettify(event))
		return nil
	})

	eventhttpserver.Register(path.Join("/", conf.GetAppSettings().AppID, "webhook/event"), conf)
	err := http.ListenAndServe(":8089", nil)
	if err != nil {
		panic(err)
	}

}
