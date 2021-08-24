package core

import (
	"context"
	"io/ioutil"
	"net/http"
)

type OapiRequest struct {
	Ctx    context.Context
	Uri    string
	Header *OapiHeader
	Body   string
}

func ToOapiRequest(request *http.Request) (*OapiRequest, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	return &OapiRequest{
		Ctx:    request.Context(),
		Uri:    request.RequestURI,
		Header: NewOapiHeader(request.Header),
		Body:   string(body),
	}, nil
}
