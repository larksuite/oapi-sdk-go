package response

import (
	"github.com/larksuite/oapi-sdk-go/core/tools"
)

const (
	ErrCodeIO                       = -1
	ErrCodeOk                       = 0
	ErrCodeAppTicketInvalid         = 10012
	ErrCodeAccessTokenInvalid       = 99991671
	ErrCodeAppAccessTokenInvalid    = 99991664
	ErrCodeTenantAccessTokenInvalid = 99991663
	ErrCodeUserAccessTokenInvalid   = 99991668
	ErrCodeUserRefreshTokenInvalid  = 99991669
)

type Detail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PermissionViolation struct {
	Type        string `json:"type"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

type FieldViolation struct {
	Field       string `json:"field"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Help struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Error struct {
	Code                 int                    `json:"code"`
	Msg                  string                 `json:"msg"`
	Details              []*Detail              `json:"details,omitempty"`
	PermissionViolations []*PermissionViolation `json:"permission_violations,omitempty"`
	FieldViolations      []*FieldViolation      `json:"field_violations,omitempty"`
	Helps                []*Help                `json:"helps,omitempty"`
}

func NewErrorOfInvalidResp(msg string) *Error {
	return &Error{
		Code: ErrCodeIO,
		Msg:  msg,
	}
}

func NewError(err error) *Error {
	return &Error{
		Code: ErrCodeIO,
		Msg:  err.Error(),
	}
}

func (e Error) Retryable() bool {
	switch e.Code {
	case ErrCodeAccessTokenInvalid,
		ErrCodeAppAccessTokenInvalid,
		ErrCodeTenantAccessTokenInvalid:
		return true
	default:
		return false
	}
}

func ToError(err error) *Error {
	switch realErr := err.(type) {
	case *Error:
		return realErr
	default:
		return NewError(err)
	}
}

func (e Error) Error() string {
	return tools.Prettify(e)
}
