package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	task "github.com/larksuite/oapi-sdk-go/service/task/v1"
)

var taskService = task.NewService(configs.TestConfig("https://open.feishu-boe.cn"))

func main() {
	testCreateTask()
}

func testCreateTask() {
	coreCtx := core.WrapContext(context.Background())
	// 测试创建任务
	createTaskReqCall := taskService.Tasks.Create(coreCtx, &task.Task{
		Summary: "测试新建任务",
		Due:     &task.Due{Time: 10000000, Timezone: "Asia/Shanghai", IsAllDay: false},
		Origin:  &task.Origin{PlatformI18nName: "{\"zh_cn\": \"IT 工作台\", \"en_us\": \"IT Workspace\"}", Href: &task.Href{Url: "https://www.feishu-boe.cn", Title: "test issue"}},
	})
	createTaskMessage, createTaskErr := createTaskReqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if createTaskErr != nil {
		fmt.Println(tools.Prettify(createTaskErr))
		e := createTaskErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println("createTaskMessage : " + tools.Prettify(createTaskMessage))

	// 测试创建协作者
	createCollaboratorReqCall := taskService.TaskCollaborators.Create(
		coreCtx,
		&task.Collaborator{Id: "ou_038f24348fc076ea1f5b83ac2d7a343c"},
	)
	createCollaboratorReqCall.SetTaskId(createTaskMessage.Task.Id)
	createCollaboratorMessage, createCollaboratorErr := createCollaboratorReqCall.Do()
	if createCollaboratorErr != nil {
		fmt.Println(tools.Prettify(createCollaboratorErr))
		e := createCollaboratorErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println("createCollaboratorMessage : " + tools.Prettify(createCollaboratorMessage))

	// 测试创建关注者
	createFollowerReqCall := taskService.TaskFollowers.Create(
		coreCtx,
		&task.Follower{Id: "ou_038f24348fc076ea1f5b83ac2d7a343c"},
	)
	createFollowerReqCall.SetTaskId(createTaskMessage.Task.Id)
	createFollowerMessage, createFollowerErr := createFollowerReqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if createFollowerErr != nil {
		fmt.Println(tools.Prettify(createFollowerErr))
		e := createFollowerErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println("createFollowerMessage : " + tools.Prettify(createFollowerMessage))

	// 测试创建提醒时间
	createReminderReqCall := taskService.TaskReminders.Create(
		coreCtx,
		&task.Reminder{RelativeFireMinute: 30},
	)
	createReminderReqCall.SetTaskId(createTaskMessage.Task.Id)
	createReminderMessage, createReminderErr := createReminderReqCall.Do()
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if createReminderErr != nil {
		fmt.Println(tools.Prettify(createReminderErr))
		e := createReminderErr.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println("createReminderMessage : " + tools.Prettify(createReminderMessage))
}
