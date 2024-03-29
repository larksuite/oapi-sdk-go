// Code generated by Lark OpenAPI.

package wiki

import (
	"github.com/larksuite/oapi-sdk-go/v3/core"
	v2 "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
)

type Service struct {
	*v2.V2
}

func NewService(config *larkcore.Config) *Service {
	return &Service{
		V2: v2.New(config),
	}
}
