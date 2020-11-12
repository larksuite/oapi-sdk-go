package main

import (
	"fmt"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	eventginserver "github.com/larksuite/oapi-sdk-go/event/http/gin"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

func main() {

	conf := test.GetInternalConf("staging")

	contact.SetDepartmentCreatedEventHandler(conf, func(ctx *core.Context, event *contact.DepartmentCreatedEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(event))
		return nil
	})

	contact.SetUserCreateEventHandler(conf, func(ctx *core.Context, event *contact.UserCreateEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(event))
		return nil
	})
	contact.SetDepartmentDeletedEventHandler(conf, func(ctx *core.Context, event *contact.DepartmentDeletedEvent) error {
		fmt.Println(ctx.GetRequestID())
		fmt.Println(tools.Prettify(event))
		return nil
	})

	g := gin.Default()
	eventginserver.Register(path.Join("/", conf.GetAppSettings().AppID, "webhook/event"), conf, g)
	err := g.Run(":8089")
	if err != nil {
		fmt.Println(err)
	}

}
