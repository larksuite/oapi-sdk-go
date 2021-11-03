package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"github.com/larksuite/oapi-sdk-go/v2/service/contact/v3"
	"os"
)

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	larkApp := lark.NewApp(lark.DomainFeiShu, appID, appSecret,
		lark.WithLogger(lark.NewDefaultLogger(), lark.LogLevelDebug))
	ctx := context.Background()
	userCreateResp, err := contact.New(larkApp).Users.Create(ctx, &contact.UserCreateReq{
		UserIdType:       lark.StringPtr("user_id"),
		DepartmentIdType: lark.StringPtr("open_department_id"),
		User: &contact.User{
			Name:          lark.StringPtr("test-name"),
			EnName:        lark.StringPtr("test-en-name"),
			Email:         lark.StringPtr("test-email@126.com"),
			Mobile:        lark.StringPtr("1234567890"),
			Gender:        lark.IntPtr(1),
			DepartmentIds: []string{"0"},
			EmployeeType:  lark.IntPtr(1),
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("request id: %s \n", userCreateResp.RequestId())
	if userCreateResp.Code != 0 {
		fmt.Println(userCreateResp.CodeError)
		return
	}
	fmt.Println(lark.Prettify(userCreateResp.Data))
}
