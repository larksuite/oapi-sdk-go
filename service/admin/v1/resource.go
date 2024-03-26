// Code generated by Lark OpenAPI.

package larkadmin

import (
	"context"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"net/http"
)

type V1 struct {
	AdminDeptStat *adminDeptStat // 部门维度的数据报表
	AdminUserStat *adminUserStat // 用户维度的数据报表
	AuditInfo     *auditInfo     // 行为审计日志（灰度租户可见）
	Badge         *badge         // 勋章
	BadgeGrant    *badgeGrant    // 勋章授予名单
	BadgeImage    *badgeImage    // 勋章图片
	Password      *password      // 登录密码管理
}

func New(config *larkcore.Config) *V1 {
	return &V1{
		AdminDeptStat: &adminDeptStat{config: config},
		AdminUserStat: &adminUserStat{config: config},
		AuditInfo:     &auditInfo{config: config},
		Badge:         &badge{config: config},
		BadgeGrant:    &badgeGrant{config: config},
		BadgeImage:    &badgeImage{config: config},
		Password:      &password{config: config},
	}
}

type adminDeptStat struct {
	config *larkcore.Config
}
type adminUserStat struct {
	config *larkcore.Config
}
type auditInfo struct {
	config *larkcore.Config
}
type badge struct {
	config *larkcore.Config
}
type badgeGrant struct {
	config *larkcore.Config
}
type badgeImage struct {
	config *larkcore.Config
}
type password struct {
	config *larkcore.Config
}

// List 获取部门维度的用户活跃和功能使用数据
//
// - 该接口用于获取部门维度的用户活跃和功能使用数据，即IM（即时通讯）、日历、云文档、音视频会议功能的使用数据。
//
// - - 只有企业自建应用才有权限调用此接口;;- 当天的数据会在第二天的早上九点半产出（UTC+8）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/admin_dept_stat/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/list_adminDeptStat.go
func (a *adminDeptStat) List(ctx context.Context, req *ListAdminDeptStatReq, options ...larkcore.RequestOptionFunc) (*ListAdminDeptStatResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/admin_dept_stats"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListAdminDeptStatResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// List 获取用户维度的用户活跃和功能使用数据
//
// - 用于获取用户维度的用户活跃和功能使用数据，即IM（即时通讯）、日历、云文档、音视频会议功能的使用数据。
//
// - - 只有企业自建应用才有权限调用此接口;;- 当天的数据会在第二天的早上九点半产出（UTC+8）
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/admin_user_stat/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/list_adminUserStat.go
func (a *adminUserStat) List(ctx context.Context, req *ListAdminUserStatReq, options ...larkcore.RequestOptionFunc) (*ListAdminUserStatResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/admin_user_stats"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListAdminUserStatResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// List
//
// -
//
// - 官网API文档链接:https://open.feishu.cn/document/ukTMukTMukTM/uQjM5YjL0ITO24CNykjN/audit_log/audit_data_get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/list_auditInfo.go
func (a *auditInfo) List(ctx context.Context, req *ListAuditInfoReq, options ...larkcore.RequestOptionFunc) (*ListAuditInfoResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/audit_infos"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, a.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListAuditInfoResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, a.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (a *auditInfo) ListByIterator(ctx context.Context, req *ListAuditInfoReq, options ...larkcore.RequestOptionFunc) (*ListAuditInfoIterator, error) {
	return &ListAuditInfoIterator{
		ctx:      ctx,
		req:      req,
		listFunc: a.List,
		options:  options,
		limit:    req.Limit}, nil
}

// Create 创建勋章
//
// - 使用该接口可以创建一枚完整的勋章信息，一个租户下最多可创建1000枚勋章。
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge/create
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/create_badge.go
func (b *badge) Create(ctx context.Context, req *CreateBadgeReq, options ...larkcore.RequestOptionFunc) (*CreateBadgeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateBadgeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Get 获取勋章详情
//
// - 可以通过该接口查询勋章的详情
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge/get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/get_badge.go
func (b *badge) Get(ctx context.Context, req *GetBadgeReq, options ...larkcore.RequestOptionFunc) (*GetBadgeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetBadgeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// List 获取勋章列表
//
// - 可以通过该接口列出租户下所有的勋章，勋章的排列顺序是按照创建时间倒序排列。
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/list_badge.go
func (b *badge) List(ctx context.Context, req *ListBadgeReq, options ...larkcore.RequestOptionFunc) (*ListBadgeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListBadgeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (b *badge) ListByIterator(ctx context.Context, req *ListBadgeReq, options ...larkcore.RequestOptionFunc) (*ListBadgeIterator, error) {
	return &ListBadgeIterator{
		ctx:      ctx,
		req:      req,
		listFunc: b.List,
		options:  options,
		limit:    req.Limit}, nil
}

// Update 修改勋章信息
//
// - 通过该接口可以修改勋章的信息
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge/update
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/update_badge.go
func (b *badge) Update(ctx context.Context, req *UpdateBadgeReq, options ...larkcore.RequestOptionFunc) (*UpdateBadgeResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id"
	apiReq.HttpMethod = http.MethodPut
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &UpdateBadgeResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Create 创建勋章的授予名单
//
// - 通过该接口可以为特定勋章创建一份授予名单，一枚勋章下最多可创建1000份授予名单。
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge-grant/create
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/create_badgeGrant.go
func (b *badgeGrant) Create(ctx context.Context, req *CreateBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*CreateBadgeGrantResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id/grants"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateBadgeGrantResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Delete 删除授予名单
//
// - 通过该接口可以删除特定授予名单的信息
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge-grant/delete
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/delete_badgeGrant.go
func (b *badgeGrant) Delete(ctx context.Context, req *DeleteBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*DeleteBadgeGrantResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id/grants/:grant_id"
	apiReq.HttpMethod = http.MethodDelete
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &DeleteBadgeGrantResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Get 获取授予名单的信息
//
// - 通过该接口可以获取特定授予名单的信息
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge-grant/get
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/get_badgeGrant.go
func (b *badgeGrant) Get(ctx context.Context, req *GetBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*GetBadgeGrantResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id/grants/:grant_id"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &GetBadgeGrantResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// List 获取勋章的授予名单列表
//
// - 通过该接口可以获取特定勋章下的授予名单列表，授予名单的排列顺序按照创建时间倒序排列。
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge-grant/list
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/list_badgeGrant.go
func (b *badgeGrant) List(ctx context.Context, req *ListBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*ListBadgeGrantResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id/grants"
	apiReq.HttpMethod = http.MethodGet
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ListBadgeGrantResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (b *badgeGrant) ListByIterator(ctx context.Context, req *ListBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*ListBadgeGrantIterator, error) {
	return &ListBadgeGrantIterator{
		ctx:      ctx,
		req:      req,
		listFunc: b.List,
		options:  options,
		limit:    req.Limit}, nil
}

// Update 修改授予名单
//
// - 通过该接口可以修改特定授予名单的相关信息
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge-grant/update
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/update_badgeGrant.go
func (b *badgeGrant) Update(ctx context.Context, req *UpdateBadgeGrantReq, options ...larkcore.RequestOptionFunc) (*UpdateBadgeGrantResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badges/:badge_id/grants/:grant_id"
	apiReq.HttpMethod = http.MethodPut
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &UpdateBadgeGrantResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Create 上传勋章图片
//
// - 通过该接口可以上传勋章详情图、挂饰图的文件，获取对应的文件key
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/badge_image/create
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/create_badgeImage.go
func (b *badgeImage) Create(ctx context.Context, req *CreateBadgeImageReq, options ...larkcore.RequestOptionFunc) (*CreateBadgeImageResp, error) {
	options = append(options, larkcore.WithFileUpload())
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/badge_images"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, b.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &CreateBadgeImageResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, b.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Reset 重置密码
//
// - 重置用户的企业邮箱密码，仅当用户的邮箱和企业邮箱(别名)一致时生效，可用于处理飞书企业邮箱登录死锁的问题。;;邮箱死锁：当用户的登录凭证与飞书企业邮箱一致时，目前飞书登录流程要求用户输入验证码，由于飞书邮箱无单独的帐号体系，则未登录时无法收取邮箱验证码，即陷入死锁
//
// - 官网API文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/admin-v1/password/reset
//
// - 使用Demo链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/adminv1/reset_password.go
func (p *password) Reset(ctx context.Context, req *ResetPasswordReq, options ...larkcore.RequestOptionFunc) (*ResetPasswordResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/admin/v1/password/reset"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, p.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &ResetPasswordResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, p.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}