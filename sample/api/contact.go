package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

var contactService = contact.NewService(conf)

func main() {
	testUserGet()
	testUserInfo()
}

func testUserGet() {
	ctx := context.Background()
	coreCtx := core.WarpContext(ctx)
	reqCall := contactService.Users.Get(coreCtx)
	reqCall.SetUserId("ou_xxxxxxxxxxxx")
	reqCall.SetUserIdType("open_id")
	reqCall.SetDepartmentIdType("open_id")
	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(result))
}
