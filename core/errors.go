package core

type IllegalParamError struct {
	msg string
}

func (err *IllegalParamError) Error() string {
	return err.msg
}

type ClientTimeoutError struct {
	msg string
}

func (err *ClientTimeoutError) Error() string {
	return err.msg
}

type ServerTimeoutError struct {
	msg string
}

func (err *ServerTimeoutError) Error() string {
	return err.msg
}

type DialFailedError struct {
	msg string
}

func (err *DialFailedError) Error() string {
	return err.msg
}
