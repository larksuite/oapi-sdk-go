package native

import (
	"github.com/larksuite/oapi-sdk-go/core/config"
	app "github.com/larksuite/oapi-sdk-go/event/app/v1"
	. "github.com/larksuite/oapi-sdk-go/event/http"
	"net/http"
)

func Register(path string, conf *config.Config) {
	app.SetAppTicketEventHandler(conf)
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		Handle(conf, request, writer)
	})
}
