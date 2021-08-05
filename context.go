package lark

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/core"
)

type Context = core.Context

func WrapContext(ctx context.Context) *Context {
	return core.WrapContext(ctx)
}
