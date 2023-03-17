// Package hire code generated by oapi sdk gen
/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkhire

import (
	"context"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func NewService(config *larkcore.Config) *HireService {
	h := &HireService{config: config}
	h.Application = &application{service: h}
	h.ApplicationInterview = &applicationInterview{service: h}
	h.Attachment = &attachment{service: h}
	h.EhrImportTask = &ehrImportTask{service: h}
	h.Employee = &employee{service: h}
	h.Job = &job{service: h}
	h.JobManager = &jobManager{service: h}
	h.JobProcess = &jobProcess{service: h}
	h.Note = &note{service: h}
	h.OfferSchema = &offerSchema{service: h}
	h.Referral = &referral{service: h}
	h.ResumeSource = &resumeSource{service: h}
	h.Talent = &talent{service: h}
	return h
}

type HireService struct {
	config               *larkcore.Config
	Application          *application          // 入职
	ApplicationInterview *applicationInterview // application.interview
	Attachment           *attachment           // 附件
	EhrImportTask        *ehrImportTask        // 导入 e-HR
	Employee             *employee             // 入职
	Job                  *job                  // 职位
	JobManager           *jobManager           // job.manager
	JobProcess           *jobProcess           // 流程
	Note                 *note                 // 备注
	OfferSchema          *offerSchema          // offer_schema
	Referral             *referral             // 内推
	ResumeSource         *resumeSource         // 简历来源
	Talent               *talent               // 人才
}

type application struct {
	service *HireService
}
type applicationInterview struct {
	service *HireService
}
type attachment struct {
	service *HireService
}
type ehrImportTask struct {
	service *HireService
}
type employee struct {
	service *HireService
}
type job struct {
	service *HireService
}
type jobManager struct {
	service *HireService
}
type jobProcess struct {
	service *HireService
}
type note struct {
	service *HireService
}
type offerSchema struct {
	service *HireService
}
type referral struct {
	service *HireService
}
type resumeSource struct {
	service *HireService
}
type talent struct {
	service *HireService
}

// 创建投递
//
// - 根据人才 ID 和职位 ID 创建投递
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/create
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/create_application.go
func (a *application) Create(ctx context.Context, req *CreateApplicationReq, options ...larkcore.RequestOptionFunc) (*CreateApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取投递信息
//
// - 根据投递 ID 获取单个投递信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_application.go
func (a *application) Get(ctx context.Context, req *GetApplicationReq, options ...larkcore.RequestOptionFunc) (*GetApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications/:application_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取投递列表
//
// - 根据限定条件获取投递列表信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/list
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/list_application.go
func (a *application) List(ctx context.Context, req *ListApplicationReq, options ...larkcore.RequestOptionFunc) (*ListApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取 Offer 信息
//
// - 根据投递 ID 获取 Offer 信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/offer
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/offer_application.go
func (a *application) Offer(ctx context.Context, req *OfferApplicationReq, options ...larkcore.RequestOptionFunc) (*OfferApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications/:application_id/offer"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &OfferApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 终止投递
//
// - 根据投递 ID 修改投递状态为「已终止」
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/terminate
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/terminate_application.go
func (a *application) Terminate(ctx context.Context, req *TerminateApplicationReq, options ...larkcore.RequestOptionFunc) (*TerminateApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications/:application_id/terminate"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &TerminateApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 操作候选人入职
//
// - 根据投递 ID 操作候选人入职并创建员工。投递须处于「待入职」阶段，可通过「转移阶段」接口变更投递状态
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/application/transfer_onboard
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/transferOnboard_application.go
func (a *application) TransferOnboard(ctx context.Context, req *TransferOnboardApplicationReq, options ...larkcore.RequestOptionFunc) (*TransferOnboardApplicationResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications/:application_id/transfer_onboard"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &TransferOnboardApplicationResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

//
//
// -
//
// - 官网 API 文档链接:https://open.feishu.cn/api-explorer?from=op_doc_tab&apiName=list&project=hire&resource=application.interview&version=v1
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/list_applicationInterview.go
func (a *applicationInterview) List(ctx context.Context, req *ListApplicationInterviewReq, options ...larkcore.RequestOptionFunc) (*ListApplicationInterviewResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/applications/:application_id/interviews"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListApplicationInterviewResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取附件信息
//
// - 获取招聘系统中附件的元信息，比如文件名、创建时间、文件 url 等
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/attachment/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_attachment.go
func (a *attachment) Get(ctx context.Context, req *GetAttachmentReq, options ...larkcore.RequestOptionFunc) (*GetAttachmentResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/attachments/:attachment_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetAttachmentResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取附件预览信息
//
// - 根据附件 ID 获取附件预览信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/attachment/preview
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/preview_attachment.go
func (a *attachment) Preview(ctx context.Context, req *PreviewAttachmentReq, options ...larkcore.RequestOptionFunc) (*PreviewAttachmentResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/attachments/:attachment_id/preview"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &PreviewAttachmentResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新 e-HR 导入任务结果
//
// - 在处理完导入 e-HR 事件后，可调用该接口，更新  e-HR 导入任务结果
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/ehr_import_task/patch
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/patch_ehrImportTask.go
func (e *ehrImportTask) Patch(ctx context.Context, req *PatchEhrImportTaskReq, options ...larkcore.RequestOptionFunc) (*PatchEhrImportTaskResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/ehr_import_tasks/:ehr_import_task_id"
	apiReq.HttpMethod = http.MethodPatch
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, e.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &PatchEhrImportTaskResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, e.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 通过员工 ID 获取入职信息
//
// - 通过员工 ID 获取入职信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/employee/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_employee.go
func (e *employee) Get(ctx context.Context, req *GetEmployeeReq, options ...larkcore.RequestOptionFunc) (*GetEmployeeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/employees/:employee_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, e.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetEmployeeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, e.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 通过投递 ID 获取入职信息
//
// - 通过投递 ID 获取入职信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/employee/get_by_application
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/getByApplication_employee.go
func (e *employee) GetByApplication(ctx context.Context, req *GetByApplicationEmployeeReq, options ...larkcore.RequestOptionFunc) (*GetByApplicationEmployeeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/employees/get_by_application"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, e.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetByApplicationEmployeeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, e.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新入职状态
//
// - 根据员工 ID 更新员工转正、离职状态
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/employee/patch
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/patch_employee.go
func (e *employee) Patch(ctx context.Context, req *PatchEmployeeReq, options ...larkcore.RequestOptionFunc) (*PatchEmployeeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/employees/:employee_id"
	apiReq.HttpMethod = http.MethodPatch
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, e.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &PatchEmployeeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, e.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 新建职位
//
// - 新建职位，字段的是否必填，以系统中的「职位字段管理」中的设置为准。
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job/combined_create
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/combinedCreate_job.go
func (j *job) CombinedCreate(ctx context.Context, req *CombinedCreateJobReq, options ...larkcore.RequestOptionFunc) (*CombinedCreateJobResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/combined_create"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CombinedCreateJobResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新职位
//
// - 更新职位信息，该接口为全量更新，若字段没有返回值，则原有值将会被清空。字段的是否必填，将以系统中的「职位字段管理」中的设置为准。
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job/combined_update
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/combinedUpdate_job.go
func (j *job) CombinedUpdate(ctx context.Context, req *CombinedUpdateJobReq, options ...larkcore.RequestOptionFunc) (*CombinedUpdateJobResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/:job_id/combined_update"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CombinedUpdateJobResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取职位设置
//
// - 获取职位设置
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job/config
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/config_job.go
func (j *job) Config(ctx context.Context, req *ConfigJobReq, options ...larkcore.RequestOptionFunc) (*ConfigJobResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/:job_id/config"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ConfigJobResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取职位信息
//
// - 根据职位 ID 获取职位信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_job.go
func (j *job) Get(ctx context.Context, req *GetJobReq, options ...larkcore.RequestOptionFunc) (*GetJobResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/:job_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetJobResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新职位设置
//
// - 更新职位设置，包括面试评价表、Offer 申请表等。接口将按照所选择的「更新选项」进行设置参数校验和更新。
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job/update_config
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/updateConfig_job.go
func (j *job) UpdateConfig(ctx context.Context, req *UpdateConfigJobReq, options ...larkcore.RequestOptionFunc) (*UpdateConfigJobResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/:job_id/update_config"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &UpdateConfigJobResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取职位上的招聘人员信息
//
// - 根据职位 ID 获取职位上的招聘人员信息，如招聘负责人、用人经理
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job-manager/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_jobManager.go
func (j *jobManager) Get(ctx context.Context, req *GetJobManagerReq, options ...larkcore.RequestOptionFunc) (*GetJobManagerResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/jobs/:job_id/managers/:manager_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetJobManagerResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取招聘流程信息
//
// - 获取全部招聘流程信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/job_process/list
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/list_jobProcess.go
func (j *jobProcess) List(ctx context.Context, req *ListJobProcessReq, options ...larkcore.RequestOptionFunc) (*ListJobProcessResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/job_processes"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, j.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListJobProcessResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, j.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 创建备注
//
// - 创建备注信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/note/create
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/create_note.go
func (n *note) Create(ctx context.Context, req *CreateNoteReq, options ...larkcore.RequestOptionFunc) (*CreateNoteResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/notes"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, n.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateNoteResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, n.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取备注
//
// - 根据备注 ID 获取备注信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/note/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_note.go
func (n *note) Get(ctx context.Context, req *GetNoteReq, options ...larkcore.RequestOptionFunc) (*GetNoteResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/notes/:note_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, n.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetNoteResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, n.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取备注列表
//
// - 获取备注列表
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/note/list
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/list_note.go
func (n *note) List(ctx context.Context, req *ListNoteReq, options ...larkcore.RequestOptionFunc) (*ListNoteResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/notes"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, n.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListNoteResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, n.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 更新备注
//
// - 根据备注 ID 更新备注信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/note/patch
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/patch_note.go
func (n *note) Patch(ctx context.Context, req *PatchNoteReq, options ...larkcore.RequestOptionFunc) (*PatchNoteResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/notes/:note_id"
	apiReq.HttpMethod = http.MethodPatch
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, n.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &PatchNoteResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, n.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

//
//
// -
//
// - 官网 API 文档链接:https://open.feishu.cn/api-explorer?from=op_doc_tab&apiName=get&project=hire&resource=offer_schema&version=v1
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_offerSchema.go
func (o *offerSchema) Get(ctx context.Context, req *GetOfferSchemaReq, options ...larkcore.RequestOptionFunc) (*GetOfferSchemaResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/offer_schemas/:offer_schema_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, o.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetOfferSchemaResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, o.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取内推信息
//
// - 根据投递 ID 获取内推信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/referral/get_by_application
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/getByApplication_referral.go
func (r *referral) GetByApplication(ctx context.Context, req *GetByApplicationReferralReq, options ...larkcore.RequestOptionFunc) (*GetByApplicationReferralResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/referrals/get_by_application"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, r.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetByApplicationReferralResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, r.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取简历来源列表
//
// - 获取简历来源列表
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/resume_source/list
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/list_resumeSource.go
func (r *resumeSource) List(ctx context.Context, req *ListResumeSourceReq, options ...larkcore.RequestOptionFunc) (*ListResumeSourceResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/resume_sources"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, r.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListResumeSourceResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, r.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (r *resumeSource) ListByIterator(ctx context.Context, req *ListResumeSourceReq, options ...larkcore.RequestOptionFunc) (*ListResumeSourceIterator, error) {
	return &ListResumeSourceIterator{
		ctx:      ctx,
		req:      req,
		listFunc: r.List,
		options:  options,
		limit:    req.Limit}, nil
}

// 通过人才信息获取人才 ID
//
// - 通过人才信息获取人才 ID
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/talent/batch_get_id
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/batchGetId_talent.go
func (t *talent) BatchGetId(ctx context.Context, req *BatchGetIdTalentReq, options ...larkcore.RequestOptionFunc) (*BatchGetIdTalentResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/talents/batch_get_id"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, t.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &BatchGetIdTalentResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, t.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 获取人才信息
//
// - 根据人才 ID 获取人才信息
//
// - 官网 API 文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uMzM1YjLzMTN24yMzUjN/hire-v1/talent/get
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/hirev1/get_talent.go
func (t *talent) Get(ctx context.Context, req *GetTalentReq, options ...larkcore.RequestOptionFunc) (*GetTalentResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/hire/v1/talents/:talent_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, t.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetTalentResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, t.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
