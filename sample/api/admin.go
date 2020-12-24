package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	admin "github.com/larksuite/oapi-sdk-go/service/admin/v1"
	"github.com/sirupsen/logrus"
)

var adminService = admin.NewService(conf)

func main() {
	testAdminDeptStatsList()
}

func testAdminDeptStatsList() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	reqCall := adminService.AdminDeptStats.List(coreCtx, request.SetTenantKey("t-xxxxxxxxx"))
	reqCall.SetDepartmentIdType("department_id")
	reqCall.SetDepartmentId("od-xxxxx")
	reqCall.SetContainsChildDept(true)
	reqCall.SetStartDate("2020-12-10")
	reqCall.SetEndDate("2020-12-14")
	reqCall.SetPageSize(2)
	// reqCall.SetPageToken("14")
	result, err := reqCall.Do()
	fmt.Println("request_id: ", coreCtx.GetRequestID())
	fmt.Println("http status code: ", coreCtx.GetHTTPStatusCode())
	if err != nil {
		e := err.(*response.Error)
		fmt.Println(tools.Prettify(e))
		return
	}
	fmt.Println("result: ", tools.Prettify(result))
}
