package v1

import (
	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"path"
)

const serviceBasePath = "authen/v1"

type Service struct {
	config   *config.Config
	basePath string
	Authens  *AuthenService
}

func NewService(config *config.Config) *Service {
	s := &Service{
		config:   config,
		basePath: serviceBasePath,
	}
	s.Authens = newAuthenService(s)
	return s
}

type AuthenService struct {
	s *Service
}

func newAuthenService(s *Service) *AuthenService {
	return &AuthenService{
		s: s,
	}
}

type AccessTokenReqCall struct {
	ctx     *core.Context
	authens *AuthenService
	body    *AccessTokenReqBody
	optFns  []request.OptFn
}

func (rc *AccessTokenReqCall) Do() (*AccessTokenResult, error) {
	var result = &AccessTokenResult{}
	httpPath := path.Join(rc.authens.s.basePath, "access_token")
	req := request.NewRequest(httpPath, "POST", []request.AccessTokenType{request.AccessTokenTypeApp}, rc.body,
		result, rc.optFns...)
	err := api.Send(rc.ctx, rc.authens.s.config, req)
	return result, err
}

func (authens *AuthenService) AccessToken(ctx *core.Context, body *AccessTokenReqBody, optFns ...request.OptFn) *AccessTokenReqCall {
	return &AccessTokenReqCall{
		ctx:     ctx,
		authens: authens,
		body:    body,
		optFns:  optFns,
	}
}

type RefreshAccessTokenReqCall struct {
	ctx     *core.Context
	authens *AuthenService
	body    *RefreshAccessTokenReqBody
	optFns  []request.OptFn
}

func (rc *RefreshAccessTokenReqCall) Do() (*AccessTokenResult, error) {
	var result = &AccessTokenResult{}
	httpPath := path.Join(rc.authens.s.basePath, "refresh_access_token")
	req := request.NewRequest(httpPath, "POST", []request.AccessTokenType{request.AccessTokenTypeApp}, rc.body,
		result, rc.optFns...)
	err := api.Send(rc.ctx, rc.authens.s.config, req)
	return result, err
}

func (authens *AuthenService) RefreshAccessToken(ctx *core.Context, body *RefreshAccessTokenReqBody,
	optFns ...request.OptFn) *RefreshAccessTokenReqCall {
	return &RefreshAccessTokenReqCall{
		ctx:     ctx,
		authens: authens,
		body:    body,
		optFns:  optFns,
	}
}

type UserInfoReqCall struct {
	ctx     *core.Context
	authens *AuthenService
	optFns  []request.OptFn
}

func (rc *UserInfoReqCall) Do() (*UserInfoResult, error) {
	httpPath := path.Join(rc.authens.s.basePath, "user_info")
	var result = &UserInfoResult{}
	req := request.NewRequest(httpPath, "GET", []request.AccessTokenType{request.AccessTokenTypeUser}, nil,
		result, rc.optFns...)
	err := api.Send(rc.ctx, rc.authens.s.config, req)
	return result, err
}

func (authens *AuthenService) UserInfo(ctx *core.Context, optFns ...request.OptFn) *UserInfoReqCall {
	return &UserInfoReqCall{
		ctx:     ctx,
		authens: authens,
		optFns:  optFns,
	}
}
