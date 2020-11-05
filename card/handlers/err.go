package handlers

type NotFoundHandlerErr struct {
}

func newNotHandlerErr() *NotFoundHandlerErr {
	return &NotFoundHandlerErr{}
}

func (e NotFoundHandlerErr) Error() string {
	return "card, not found handler"
}

type SignatureErr struct {
}

func newSignatureErr() *SignatureErr {
	return &SignatureErr{}
}

func (e SignatureErr) Error() string {
	return "card, signature error"
}
