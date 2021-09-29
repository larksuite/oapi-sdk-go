package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api/core/constants"
	coreerrors "github.com/larksuite/oapi-sdk-go/api/core/errors"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/api/core/token"
	"github.com/larksuite/oapi-sdk-go/api/core/transport"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	coreconst "github.com/larksuite/oapi-sdk-go/core/constants"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const defaultMaxRetryCount = 1

var defaultHTTPRequestHeader = map[string]string{}
var defaultHTTPRequestHeaderKeysWithContext = []string{coreconst.HTTPHeaderKeyRequestID}
var Default = &Handlers{}

func init() {
	defaultHTTPRequestHeader["User-Agent"] = fmt.Sprintf("oapi-sdk-go/%s", core.SdkVersion)
	Default.init = initFunc
	Default.validate = validateFunc
	Default.build = buildFunc
	Default.sign = signFunc
	Default.unmarshalResponse = unmarshalResponseFunc
	Default.validateResponse = validateResponseFunc
	Default.retry = retryFunc
	Default.complement = complementFunc
}

type Handler func(*core.Context, *request.Request)

type Handlers struct {
	init              Handler
	validate          Handler
	build             Handler // build http request
	sign              Handler // sign token to header
	validateResponse  Handler
	unmarshalResponse Handler
	retry             Handler // when token invalid, retry
	complement        Handler
}

func Handle(ctx *core.Context, req *request.Request) {
	defer Default.complement(ctx, req)
	Default.init(ctx, req)
	if req.Err != nil {
		return
	}
	Default.validate(ctx, req)
	if req.Err != nil {
		return
	}
	i := 0
	for {
		i++
		Default.send(ctx, req)
		if !req.Retryable || i > defaultMaxRetryCount {
			return
		}
		config.ByCtx(ctx).GetLogger().Debug(ctx, fmt.Sprintf("[retry] request:%v, err: %v", req, req.Err))
		req.Err = nil
	}
}

func (hs *Handlers) send(ctx *core.Context, req *request.Request) {
	hs.build(ctx, req)
	if req.Err != nil {
		return
	}
	hs.sign(ctx, req)
	if req.Err != nil {
		return
	}
	resp, err := transport.DefaultClient.Do(req.HTTPRequest)
	if err != nil {
		req.Err = err
		return
	}
	ctx.Set(coreconst.HTTPHeader, core.NewOapiHeader(resp.Header))
	ctx.Set(coreconst.HTTPKeyStatusCode, resp.StatusCode)
	req.HTTPResponse = resp
	defer hs.retry(ctx, req)
	hs.validateResponse(ctx, req)
	if req.Err != nil {
		return
	}
	hs.unmarshalResponse(ctx, req)
}

func initFunc(ctx *core.Context, req *request.Request) {
	conf := config.ByCtx(ctx)
	req.Err = req.Init(conf.GetDomain())
}

func validateFunc(ctx *core.Context, req *request.Request) {
	if req.AccessTokenType == request.AccessTokenTypeNone {
		return
	}
	if _, ok := req.AccessibleTokenTypeSet[req.AccessTokenType]; !ok {
		req.Err = coreerrors.ErrAccessTokenTypeInvalid
	}
	if config.ByCtx(ctx).GetAppSettings().AppType == coreconst.AppTypeISV {
		if req.AccessTokenType == request.AccessTokenTypeTenant && req.TenantKey == "" {
			req.Err = coreerrors.ErrTenantKeyIsEmpty
			return
		}
	}
	if req.AccessTokenType == request.AccessTokenTypeUser && req.UserAccessToken == "" {
		req.Err = coreerrors.ErrUserAccessTokenKeyIsEmpty
		return
	}
}

func buildFunc(ctx *core.Context, req *request.Request) {
	conf := config.ByCtx(ctx)
	if !req.Retryable {
		if req.Input != nil {
			switch req.Input.(type) {
			case *request.FormData:
				reqBodyFromFormData(ctx, req)
				conf.GetLogger().Debug(ctx, fmt.Sprintf("[build]request:%v, body:formdata:%s", req, req.RequestBodyFilePath))
			default:
				reqBodyFromInput(ctx, req)
				conf.GetLogger().Debug(ctx, fmt.Sprintf("[build]request:%v, body:%s", req, string(req.RequestBody)))
			}
		} else {
			conf.GetLogger().Debug(ctx, fmt.Sprintf("[build]request:%v", req))
		}
		if req.Err != nil {
			return
		}
	}
	if req.RequestBody != nil {
		req.RequestBodyStream = bytes.NewBuffer(req.RequestBody)
	} else if req.RequestBodyFilePath != "" {
		req.RequestBodyStream, req.Err = os.Open(req.RequestBodyFilePath)
		if req.Err != nil {
			return
		}
	}
	r, err := http.NewRequestWithContext(ctx, req.HttpMethod, req.Url(), req.RequestBodyStream)
	if err != nil {
		req.Err = err
		return
	}
	for k, v := range defaultHTTPRequestHeader {
		r.Header.Set(k, v)
	}
	for _, k := range defaultHTTPRequestHeaderKeysWithContext {
		if v, ok := ctx.Get(k); ok {
			r.Header.Set(k, fmt.Sprint(v))
		}
	}
	if req.ContentType != "" {
		r.Header.Set(coreconst.ContentType, req.ContentType)
	}
	req.HTTPRequest = r
}

func signFunc(ctx *core.Context, req *request.Request) {
	var httpRequest *http.Request
	var err error
	switch req.AccessTokenType {
	case request.AccessTokenTypeApp:
		httpRequest, err = setAppAccessToken(ctx, req.HTTPRequest)
	case request.AccessTokenTypeTenant:
		httpRequest, err = setTenantAccessToken(ctx, req.HTTPRequest)
	case request.AccessTokenTypeUser:
		httpRequest, err = setUserAccessToken(ctx, req.HTTPRequest)
	default:
		httpRequest, err = req.HTTPRequest, req.Err
	}
	if req.NeedHelpDeskAuth {
		conf := config.ByCtx(ctx)
		if conf.GetHelpDeskAuthorization() == "" {
			err = errors.New("help desk API, please set the helpdesk information of config.AppSettings")
		} else if httpRequest != nil {
			httpRequest.Header.Set("X-Lark-Helpdesk-Authorization", conf.GetHelpDeskAuthorization())
		}
	}
	req.HTTPRequest = httpRequest
	req.Err = err
}

func validateResponseFunc(_ *core.Context, req *request.Request) {
	resp := req.HTTPResponse
	contentType := resp.Header.Get(coreconst.ContentType)
	if req.IsResponseStream {
		if resp.StatusCode == http.StatusOK {
			req.IsResponseStreamReal = true
			return
		}
		if strings.Contains(contentType, coreconst.ContentTypeJson) {
			req.IsResponseStreamReal = false
			return
		}
		if resp.StatusCode != http.StatusOK {
			req.Err = fmt.Errorf("response is stream, but status code:%d, contentType:%s", resp.StatusCode, contentType)
		}
		return
	}
	if !strings.Contains(contentType, coreconst.ContentTypeJson) {
		respBody, err := readResponse(resp)
		if err != nil {
			req.Err = err
			return
		}
		req.Err = response.NewErrorOfInvalidResp(fmt.Sprintf("content-type: %s, is not: %s, body:%s", contentType, coreconst.ContentTypeJson, string(respBody)))
	}
}

func unmarshalResponseFunc(ctx *core.Context, req *request.Request) {
	resp := req.HTTPResponse
	if req.IsResponseStreamReal {
		defer resp.Body.Close()
		switch output := req.Output.(type) {
		case io.Writer:
			_, err := io.Copy(output, resp.Body)
			if err != nil {
				req.Err = err
				return
			}
		default:
			req.Err = fmt.Errorf("request`s Output type must implement `io.Writer` interface")
			return
		}
		return
	}
	respBody, err := readResponse(resp)
	if err != nil {
		req.Err = err
		return
	}
	config.ByCtx(ctx).GetLogger().Debug(ctx, fmt.Sprintf("[unmarshalResponse] request:%v, response:body:%s",
		req, string(respBody)))
	if req.DataFilled() {
		err := unmarshalJSON(req.Output, req.IsNotDataField, respBody)
		if err != nil {
			req.Err = err
			return
		}
	} else {
		req.Err = fmt.Errorf("request out do not write")
		return
	}
}

func readResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func retryFunc(_ *core.Context, req *request.Request) {
	if req.Err != nil {
		if err, ok := req.Err.(*response.Error); ok {
			req.Info.Retryable = err.Retryable()
			return
		}
	}
	req.Info.Retryable = false
}

func complementFunc(ctx *core.Context, req *request.Request) {
	if req.RequestBodyFilePath != "" {
		if err := os.Remove(req.RequestBodyFilePath); err != nil {
			config.ByCtx(ctx).GetLogger().Debug(ctx, fmt.Sprintf("[complement] request:%v, "+
				"delete tmp file(%s) err: %v", req, req.RequestBodyFilePath, err))
		}
	}
	switch err := req.Err.(type) {
	case *response.Error:
		switch err.Code {
		case response.ErrCodeAppTicketInvalid:
			applyAppTicket(ctx)
		}
	default:
		if req.Err == coreerrors.ErrAppTicketIsEmpty {
			applyAppTicket(ctx)
		}
	}
}

// apply app ticket
func applyAppTicket(ctx *core.Context) {
	conf := config.ByCtx(ctx)
	req := request.NewRequestByAuth(constants.ApplyAppTicketPath, http.MethodPost,
		&token.ApplyAppTicketReq{
			AppID:     conf.GetAppSettings().AppID,
			AppSecret: conf.GetAppSettings().AppSecret,
		}, &response.NoData{})
	Handle(ctx, req)
	if req.Err != nil {
		conf.GetLogger().Error(ctx, req.Err)
	}
}

func unmarshalJSON(v interface{}, isNotDataField bool, data []byte) error {
	var e response.Error
	if isNotDataField {
		if ret, ok := v.(*map[string]interface{}); ok {
			e = response.Error{}
			err := json.Unmarshal(data, &e)
			if err != nil {
				return err
			}
			err = json.Unmarshal(data, ret)
			if err != nil {
				return err
			}
			delete(*ret, "code")
			delete(*ret, "msg")
		} else {
			typ := reflect.TypeOf(v)
			name := typ.Elem().Name()
			responseTyp := reflect.StructOf([]reflect.StructField{
				{
					Name:      "Error",
					Anonymous: true,
					Type:      reflect.TypeOf(response.Error{}),
				},
				{
					Name:      name,
					Anonymous: true,
					Type:      typ,
				},
			})
			responseV := reflect.New(responseTyp).Elem()
			responseV.Field(1).Set(reflect.ValueOf(v))
			s := responseV.Addr().Interface()
			err := json.Unmarshal(data, s)
			if err != nil {
				return err
			}
			e = responseV.Field(0).Interface().(response.Error)
		}
	} else {
		out := &response.Response{
			Data: v,
		}
		err := json.Unmarshal(data, out)
		if err != nil {
			return err
		}
		e = out.Error
	}
	if e.Code == response.ErrCodeOk {
		return nil
	}
	return &e
}

func reqBodyFromFormData(_ *core.Context, req *request.Request) {
	var reqBody io.ReadWriter
	fd := req.Input.(*request.FormData)
	hasStream := fd.HasStream()
	if hasStream {
		var reqBodyFile *os.File
		reqBodyFile, req.Err = ioutil.TempFile("", ".larksuiteoapisdk")
		if req.Err != nil {
			return
		}
		defer reqBodyFile.Close()
		req.RequestBodyFilePath = reqBodyFile.Name()
		reqBody = reqBodyFile
	} else {
		reqBody = &bytes.Buffer{}
	}
	writer := multipart.NewWriter(reqBody)
	for key, val := range fd.Params() {
		err := writer.WriteField(key, fmt.Sprint(val))
		if err != nil {
			req.Err = err
			return
		}
	}
	for _, file := range fd.Files() {
		part, err := writer.CreatePart(file.MIMEHeader())
		if err != nil {
			req.Err = err
			return
		}
		_, err = io.Copy(part, file)
		if err != nil {
			req.Err = err
			return
		}
	}
	req.ContentType = writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		req.Err = err
		return
	}
	if !hasStream {
		req.RequestBody, req.Err = ioutil.ReadAll(reqBody)
	}
}

func reqBodyFromInput(_ *core.Context, req *request.Request) {
	var bs []byte
	if input, ok := req.Input.(string); ok {
		bs = []byte(input)
	} else {
		reqBody := new(bytes.Buffer)
		err := json.NewEncoder(reqBody).Encode(req.Input)
		if err != nil {
			req.Err = err
			return
		}
		bs = reqBody.Bytes()
	}
	req.ContentType = coreconst.DefaultContentType
	req.RequestBody = bs
}
