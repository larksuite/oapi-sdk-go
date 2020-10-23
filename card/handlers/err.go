package handlers

type NotHandlerErr struct {
}

func newNotHandlerErr() *NotHandlerErr {
	return &NotHandlerErr{}
}

func (e NotHandlerErr) Error() string {
	return "card, not find handler"
}

type SignatureErr struct {
}

func newSignatureErr() *SignatureErr {
	return &SignatureErr{}
}

func (e SignatureErr) Error() string {
	return "card, signature error"
}
