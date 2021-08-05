package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	vc "github.com/larksuite/oapi-sdk-go/service/vc/v1"
)

// for redis store and logrus
// configs.TestConfigWithLogrusAndRedisStore(lark.DomainFeiShu)
// configs.TestConfig("https://open.feishu.cn")
var VCService = vc.NewService(configs.TestConfig(lark.DomainFeiShu))

func main() {
	testReserveApply()
}

func testReserveApply() {
	ctx := context.Background()
	coreCtx := lark.WrapContext(ctx)
	body := &vc.ReserveApplyReqBody{
		EndTime: 1617161325,
		MeetingSettings: &vc.ReserveMeetingSetting{
			Topic: "Test VC",
			ActionPermissions: []*vc.ReserveActionPermission{{
				Permission: 1,
				PermissionCheckers: []*vc.ReservePermissionChecker{{
					CheckField: 1,
					CheckMode:  1,
					CheckList:  []string{"77bbc392"},
				},
				},
			},
			},
		},
	}
	reqCall := VCService.Reserves.Apply(coreCtx, body, request.SetUserAccessToken("User access token"))
	reqCall.SetUserIdType("user_id")
	result, err := reqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(lark.Prettify(err))
		e := err.(*lark.APIError)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(lark.Prettify(result))
}
