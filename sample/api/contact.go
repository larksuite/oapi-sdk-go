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
	testUserServiceList()
	testDepartmentServiceUpdate()
}
func testUserServiceList() {

	coreCtx := core.WarpContext(context.Background())
	reqCall := contactService.Users.List(coreCtx)
	reqCall.SetDepartmentIdType("open_id")
	reqCall.SetPageSize(20)
	reqCall.SetDepartmentIdType("open_department_id")
	reqCall.SetDepartmentId("od_XXXXXXXXX")
	reqCall.SetUserIdType("open_id")
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Printf(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))
}

func testDepartmentServiceUpdate() {
	coreCtx := core.WarpContext(context.Background())
	updateBody := &contact.Department{
		Name:               "xxxxx",
		ParentDepartmentId: "od_xxxxxx",
		LeaderUserId:       "ou-xxxxxxx",
		Order:              666,
		CreateGroupChat:    false,
	}
	reqCall := contactService.Departments.Create(coreCtx, updateBody)
	reqCall.SetDepartmentIdType("open_department_id")
	reqCall.SetUserIdType("open_id")
	reqCall.SetClientToken("xxxxxxx")
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Printf(tools.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", tools.Prettify(result))
}
