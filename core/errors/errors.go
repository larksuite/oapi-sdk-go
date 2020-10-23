package errors

type TokenInvalidErr struct {
}

func NewTokenInvalidErr() *TokenInvalidErr {
	return &TokenInvalidErr{}
}

func (e TokenInvalidErr) Error() string {
	return "token invalid"
}
