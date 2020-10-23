package native

import (
	. "github.com/larksuite/oapi-sdk-go/card/http"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"net/http"
)

func Register(path string, conf *config.Config) {
	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		Handle(conf, request, writer)
	})
}
