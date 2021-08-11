package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	contact "github.com/larksuite/oapi-sdk-go/service/contact/v3"
)

// for redis store and logrus
// sample.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// sample.TestConfig("https://open.feishu.cn")
var contactService = contact.NewService(lark.NewInternalAppConfigByEnv(lark.DomainFeiShu))

func main() {
	testUserServiceList()
	//testDepartmentServiceList()
}
func testUserServiceList() {
	coreCtx := lark.WrapContext(context.Background())
	reqCall := contactService.Users.List(coreCtx)
	reqCall.SetDepartmentIdType("open_id")
	reqCall.SetPageSize(20)
	reqCall.SetDepartmentIdType("open_department_id")
	reqCall.SetDepartmentId("0")
	reqCall.SetUserIdType("open_id")
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(lark.APIError)
		fmt.Printf(lark.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", lark.Prettify(result))
}

func testDepartmentServiceList() {
	coreCtx := lark.WrapContext(context.Background())
	reqCall := contactService.Departments.List(coreCtx)
	reqCall.SetDepartmentIdType("open_department_id")
	reqCall.SetUserIdType("open_id")
	result, err := reqCall.Do()
	fmt.Printf("request_id:%s\n", coreCtx.GetRequestID())
	fmt.Printf("http status code:%d", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(lark.APIError)
		fmt.Printf(lark.Prettify(e))
		return
	}
	fmt.Printf("reault:%s", lark.Prettify(result))
}
