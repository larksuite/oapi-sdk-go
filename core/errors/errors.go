package errors

import "fmt"

type TokenInvalidErr struct {
	token string
}

func NewTokenInvalidErr(token string) *TokenInvalidErr {
	return &TokenInvalidErr{token: token}
}

func (e TokenInvalidErr) Error() string {
	return fmt.Sprintf("AppSettings.verificationToken not equal token(%s)", e.token)
}
