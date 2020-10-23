package gin

import (
	"github.com/gin-gonic/gin"
	. "github.com/larksuite/oapi-sdk-go/card/http"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

func Register(path string, conf *config.Config, g *gin.Engine) {
	g.POST(path, func(context *gin.Context) {
		Handle(conf, context.Request, context.Writer)
	})
}
