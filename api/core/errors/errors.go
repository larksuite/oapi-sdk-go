package errors

import "errors"

var (
	ErrAccessTokenTypeIsInValid  = errors.New("access token type is invalid")
	ErrTenantKeyIsEmpty          = errors.New("tenant key is empty")
	ErrUserAccessTokenKeyIsEmpty = errors.New("user access token is empty")
	ErrAppTicketIsEmpty          = errors.New("app ticket is empty")
)
