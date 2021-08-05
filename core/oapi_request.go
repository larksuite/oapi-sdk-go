package core

import (
	"io/ioutil"
	"net/http"
)

type OapiRequest struct {
	Header *OapiHeader
	Body   string
}

func ToOapiRequest(request *http.Request) (*OapiRequest, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	return &OapiRequest{
		Header: NewOapiHeader(request.Header),
		Body:   string(body),
	}, nil
}
