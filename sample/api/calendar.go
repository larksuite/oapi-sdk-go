package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	calendar "github.com/larksuite/oapi-sdk-go/service/calendar/v4"
)

var calendarService = calendar.NewService(conf)

func main() {
	testCalendarList()
}

func testCalendarList() {
	ctx := context.Background()
	coreCtx := core.WrapContext(ctx)
	pageToken := ""
	syncToken := ""
	hasMore := true
	count := 0
	for hasMore {
		reqCall := calendarService.Calendars.List(coreCtx,
			request.SetUserAccessToken("u-xxxxxxxxx"))
		reqCall.SetPageSize(50)
		reqCall.SetPageToken("")
		reqCall.SetSyncToken("")
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
		pageToken = result.PageToken
		syncToken = result.SyncToken
		hasMore = result.HasMore
		fmt.Printf("calendar list finish, count = %d, calendars len = %d, pageToken = %s, syncToken = %s \n",
			count, len(result.CalendarList), pageToken, syncToken)
	}
	fmt.Printf("calendar list finish\n")
}
