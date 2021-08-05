package core

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/core/constants"
	"net/http"
)

const ResponseFormat = `{"codemsg":"%s"}`
const ChallengeResponseFormat = `{"challenge":"%s"}`

type OapiResponse struct {
	StatusCode  int
	ContentType string
	Header      map[string][]string
	Body        string
}

func NewOapiResponseOfErr(err error) *OapiResponse {
	return &OapiResponse{
		StatusCode:  http.StatusInternalServerError,
		ContentType: constants.DefaultContentType,
		Header:      nil,
		Body:        fmt.Sprintf(ResponseFormat, err.Error()),
	}
}

func (r *OapiResponse) Write(statusCode int, contentType string, body string) {
	r.StatusCode = statusCode
	r.ContentType = contentType
	r.Body = body
}

func (r OapiResponse) WriteTo(response http.ResponseWriter) error {
	for k, vs := range r.Header {
		for _, v := range vs {
			response.Header().Set(k, v)
		}
	}
	response.Header().Set(constants.ContentType, r.ContentType)
	response.WriteHeader(r.StatusCode)
	_, err := fmt.Fprint(response, r.Body)
	return err
}
