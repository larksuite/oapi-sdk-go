package response

import (
	"github.com/larksuite/oapi-sdk-go/core/tools"
)

const (
	ErrCodeNative                   = -1
	ErrCodeOk                       = 0
	ErrCodeAppTicketInvalid         = 10012
	ErrCodeAccessTokenInvalid       = 99991671
	ErrCodeAppAccessTokenInvalid    = 99991664
	ErrCodeTenantAccessTokenInvalid = 99991663
	ErrCodeUserAccessTokenInvalid   = 99991668
	ErrCodeUserRefreshTokenInvalid  = 99991669
)

type Detail struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type PermissionViolation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type FieldViolation struct {
	Field       string `json:"field,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}

type Help struct {
	Url         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
}

type Err struct {
	Details              []*Detail              `json:"details,omitempty"`
	PermissionViolations []*PermissionViolation `json:"permission_violations,omitempty"`
	FieldViolations      []*FieldViolation      `json:"field_violations,omitempty"`
	Helps                []*Help                `json:"helps,omitempty"`
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  *Err   `json:"error"`
}

func NewErrorOfInvalidResp(msg string) *Error {
	return &Error{
		Code: ErrCodeNative,
		Msg:  msg,
	}
}

func NewError(err error) *Error {
	return &Error{
		Code: ErrCodeNative,
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
