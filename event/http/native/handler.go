package native

import (
	"github.com/larksuite/oapi-sdk-go/core"
	. "github.com/larksuite/oapi-sdk-go/event/http"
	"net/http"
)

func Register(path string, conf core.Config) {
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		Handle(conf, request, writer)
	})
}
