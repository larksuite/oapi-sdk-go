package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

var (
	ErrAppTicketIsEmpty = errors.New("app ticket is empty")
)

func (app *App) SendRequest(ctx context.Context, httpMethod string, httpPath string,
	accessTokenType AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	return app.SendRequestWithAccessTokenTypes(ctx, httpMethod, httpPath, []AccessTokenType{accessTokenType}, input, options...)
}

func (app *App) SendRequestWithAccessTokenTypes(ctx context.Context, httpMethod string, httpPath string,
	accessTokenTypes []AccessTokenType, input interface{}, options ...RequestOptionFunc) (*RawResponse, error) {
	option := &requestOption{}
	for _, optionFunc := range options {
		optionFunc(option)
	}
	accessibleTokenTypeSet := make(map[AccessTokenType]struct{})
	accessTokenType := accessTokenTypes[0]
	for _, t := range accessTokenTypes {
		if t == AccessTokenTypeTenant {
			accessTokenType = t // default
		}
		accessibleTokenTypeSet[t] = struct{}{}
	}
	if option.tenantKey != "" {
		if _, ok := accessibleTokenTypeSet[AccessTokenTypeTenant]; ok {
			accessTokenType = AccessTokenTypeTenant
		}
	}
	if option.userAccessToken != "" {
		if _, ok := accessibleTokenTypeSet[AccessTokenTypeUser]; ok {
			accessTokenType = AccessTokenTypeUser
		}
	}
	paths, queries, body := parseInput(input, option)
	if _, ok := body.(*Formdata); ok {
		option.fileUpload = true
	}
	contentType, rawBody, err := payload(body)
	if err != nil {
		return nil, err
	}
	fullURL, err := jointURL(app.domain, httpPath, paths, queries)
	if err != nil {
		return nil, err
	}
	req := &request{
		method:          httpMethod,
		url:             fullURL,
		contentType:     contentType,
		body:            rawBody,
		accessTokenType: accessTokenType,
		option:          option,
	}
	return req.do(ctx, app)
}

type RequestOptionFunc func(option *requestOption)

func WithNeedHelpDeskAuth() RequestOptionFunc {
	return func(option *requestOption) {
		option.needHelpDeskAuth = true
	}
}

func WithTenantKey(tenantKey string) RequestOptionFunc {
	return func(option *requestOption) {
		option.tenantKey = tenantKey
	}
}

func WithFileUpload() RequestOptionFunc {
	return func(option *requestOption) {
		option.fileUpload = true
	}
}

func WithHTTPHeader(header http.Header) RequestOptionFunc {
	return func(option *requestOption) {
		option.header = header
	}
}

func WithUserAccessToken(userAccessToken string) RequestOptionFunc {
	return func(option *requestOption) {
		option.userAccessToken = userAccessToken
	}
}

func parseInput(input interface{}, option *requestOption) (map[string]interface{}, map[string]interface{}, interface{}) {
	if input == nil {
		return nil, nil, nil
	}
	if _, ok := input.(*Formdata); ok {
		return nil, nil, input
	}
	var hasHTTPTag bool
	paths, queries := map[string]interface{}{}, map[string]interface{}{}
	vv := reflect.ValueOf(input)
	vt := reflect.TypeOf(input)
	if vt.Kind() == reflect.Ptr {
		vv = vv.Elem()
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return nil, nil, input
	}
	var body interface{}
	for i := 0; i < vt.NumField(); i++ {
		fieldValue := vv.Field(i)
		fieldType := vt.Field(i)
		if path, ok := fieldType.Tag.Lookup("path"); ok {
			hasHTTPTag = true
			if path != "" && !isEmpty(fieldValue) {
				paths[path] = reflect.Indirect(fieldValue).Interface()
			}
			continue
		}
		if query, ok := fieldType.Tag.Lookup("query"); ok {
			hasHTTPTag = true
			if query != "" && !isEmpty(fieldValue) {
				queries[query] = reflect.Indirect(fieldValue).Interface()
			}
			continue
		}
		if _, ok := fieldType.Tag.Lookup("body"); ok {
			hasHTTPTag = true
			body = fieldValue.Interface()
		}
	}
	if !hasHTTPTag {
		body = input
		if option.fileUpload {
			body = toFormdata(input)
		}
		return nil, nil, body
	}
	if body != nil {
		if option.fileUpload {
			body = toFormdata(body)
		}
	}
	return paths, queries, body
}

func toFormdata(body interface{}) *Formdata {
	formdata := &Formdata{}
	v := reflect.ValueOf(body)
	t := reflect.TypeOf(body)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		if isEmpty(fieldValue) {
			continue
		}
		if fieldName := fieldType.Tag.Get("json"); fieldName != "" {
			if strings.HasSuffix(fieldName, ",omitempty") {
				fieldName = fieldName[:len(fieldName)-10]
			}
			formdata.AddField(fieldName, reflect.Indirect(fieldValue).Interface())
		}
	}
	return formdata
}

func isEmpty(value reflect.Value) bool {
	if (value.Kind() == reflect.Ptr || value.Kind() == reflect.Slice || value.Kind() == reflect.Map) && value.IsNil() {
		return true
	}
	if (value.Kind() == reflect.Slice || value.Kind() == reflect.Map) && value.Len() == 0 {
		return true
	}
	return false
}

func jointURL(domain Domain, httpPath string, paths, queries map[string]interface{}) (string, error) {
	// path
	var pathSegs []string
	for _, p := range strings.Split(httpPath, "/") {
		if strings.Index(p, ":") == 0 {
			varName := p[1:]
			v, ok := paths[varName]
			if !ok {
				return "", fmt.Errorf("http path:%s, name: %s, not found value", httpPath, varName)
			}
			val := fmt.Sprint(v)
			if val == "" {
				return "", fmt.Errorf("http path:%s, name: %s, value is empty", httpPath, varName)
			}
			pathSegs = append(pathSegs, val)
			continue
		}
		pathSegs = append(pathSegs, p)
	}
	newPath := strings.Join(pathSegs, "/")
	if strings.Index(newPath, "http") != 0 {
		if strings.Index(newPath, "/open-apis") == 0 {
			newPath = fmt.Sprintf("%s%s", domain, newPath)
		} else {
			newPath = fmt.Sprintf("%s/open-apis/%s", domain, newPath)
		}
	}
	// query
	query := make(url.Values)
	for k, v := range queries {
		sv := reflect.ValueOf(v)
		if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
			for i := 0; i < sv.Len(); i++ {
				query.Add(k, fmt.Sprint(sv.Index(i)))
			}
		} else {
			query.Set(k, fmt.Sprint(v))
		}
	}
	if len(query) > 0 {
		newPath = fmt.Sprintf("%s?%s", newPath, query.Encode())
	}
	return newPath, nil
}

func payload(body interface{}) (string, []byte, error) {
	if fd, ok := body.(*Formdata); ok {
		return fd.content()
	}
	contentType := defaultContentType
	if body == nil {
		return contentType, nil, nil
	}
	bs, err := json.Marshal(body)
	return contentType, bs, err
}

type requestOption struct {
	tenantKey        string
	userAccessToken  string
	needHelpDeskAuth bool
	fileUpload       bool
	header           http.Header
}

type request struct {
	method          string
	url             string
	contentType     string
	body            []byte
	accessTokenType AccessTokenType
	option          *requestOption
	retryCount      int
}

func (r *request) do(ctx context.Context, app *App) (*RawResponse, error) {
	err := r.validate(app)
	if err != nil {
		return nil, err
	}
	rawResp, code, err := r.send(ctx, app)
	if code == errCodeAppTicketInvalid || err == ErrAppTicketIsEmpty {
		app.logger.Warn(ctx, fmt.Sprintf("app_ticket invalid, send apply app_ticket request"))
		r.applyAppTicket(ctx, app)
	}
	return rawResp, err
}

func (r *request) validate(app *App) error {
	if app.settings.type_ == AppTypeMarketplace {
		if r.accessTokenType == AccessTokenTypeTenant && r.option.tenantKey == "" {
			return errors.New("tenant key is empty")
		}
	}
	if r.accessTokenType == AccessTokenTypeUser && r.option.userAccessToken == "" {
		return errors.New("user access token is empty")
	}
	return nil
}

func (r *request) send(ctx context.Context, app *App) (*RawResponse, int, error) {
	var rawResp *RawResponse
	var code int
	for r.retryCount < 2 {
		app.logger.Debug(ctx, fmt.Sprintf("send request %v", r))
		rawRequest, err := http.NewRequestWithContext(ctx, r.method, r.url, bytes.NewBuffer(r.body))
		if err != nil {
			return nil, 0, err
		}
		for k, vs := range r.option.header {
			for _, v := range vs {
				rawRequest.Header.Add(k, v)
			}
		}
		rawRequest.Header.Set("User-Agent", fmt.Sprintf("oapi-sdk-go-v2/%s", version))
		if r.contentType != "" {
			rawRequest.Header.Set(contentTypeHeader, r.contentType)
		}
		switch r.accessTokenType {
		case AccessTokenTypeApp:
			err = r.signAppAccessToken(ctx, rawRequest, app)
		case AccessTokenTypeTenant:
			err = r.signTenantAccessToken(ctx, rawRequest, app)
		case AccessTokenTypeUser:
			err = r.signUserAccessToken(rawRequest)
		}
		if err != nil {
			return nil, 0, err
		}
		err = r.signHelpdeskAuthToken(rawRequest, app)
		if err != nil {
			return nil, 0, err
		}
		resp, err := http.DefaultClient.Do(rawRequest)
		if err != nil {
			return nil, 0, err
		}
		body, err := r.readResponse(resp)
		if err != nil {
			return nil, 0, err
		}
		rawResp = &RawResponse{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			RawBody:    body,
		}
		app.logger.Debug(ctx, fmt.Sprintf("send request %v, response %v", r, rawResp))
		if r.retryCount == 1 || !strings.Contains(resp.Header.Get(contentTypeHeader), contentTypeJson) {
			break
		}
		codeError := &CodeError{}
		err = json.Unmarshal(body, codeError)
		if err != nil {
			return nil, 0, err
		}
		code = codeError.Code
		if code != errCodeAccessTokenInvalid && code != errCodeAppAccessTokenInvalid &&
			code != errCodeTenantAccessTokenInvalid {
			break
		}
		r.retryCount++
	}
	return rawResp, code, nil
}

func (r *request) applyAppTicket(ctx context.Context, app *App) {
	rawResp, err := app.SendRequest(ctx, http.MethodPost, applyAppTicketPath, accessTokenTypeNone, &applyAppTicketReq{
		AppID:     app.settings.id,
		AppSecret: app.settings.secret,
	})
	if err != nil {
		app.logger.Error(ctx, fmt.Sprintf("apply app_ticket, error: %v", err))
		return
	}
	if !strings.Contains(rawResp.Header.Get(contentTypeHeader), contentTypeJson) {
		app.logger.Error(ctx, fmt.Sprintf("apply app_ticket, response content-type not json, response: %v", rawResp))
		return
	}
	codeError := &CodeError{}
	err = json.Unmarshal(rawResp.RawBody, codeError)
	if err != nil {
		app.logger.Error(ctx, fmt.Sprintf("apply app_ticket, json unmarshal error: %v", err))
		return
	}
	if codeError.Code != 0 {
		app.logger.Error(ctx, fmt.Sprintf("apply app_ticket, response error: %v", codeError))
		return
	}
}

func (r *request) String() string {
	bodyStr := ""
	if r.option.fileUpload {
		bodyStr = fmt.Sprintf("<binary> len %d", len(r.body))
	} else {
		bodyStr = string(r.body)
	}
	return fmt.Sprintf("%s %s, body: %s", r.method, r.url, bodyStr)
}

const expiryDelta = 3 * time.Minute

// internal app access token
func (r *request) customAppAccessToken(ctx context.Context, app *App) (string, error) {
	rawResp, err := app.SendRequest(ctx, http.MethodPost, appAccessTokenInternalUrlPath, accessTokenTypeNone, &internalAccessTokenReq{
		AppID:     app.settings.id,
		AppSecret: app.settings.secret,
	})
	if err != nil {
		return "", err
	}
	appAccessTokenResp := &appAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, appAccessTokenResp)
	if err != nil {
		return "", err
	}
	if appAccessTokenResp.Code != 0 {
		return "", appAccessTokenResp.CodeError
	}
	expire := time.Duration(appAccessTokenResp.Expire)*time.Second - expiryDelta
	err = app.store.Put(ctx, appAccessTokenKey(app.settings.id), appAccessTokenResp.AppAccessToken, expire)
	if err != nil {
		app.logger.Warn(ctx, fmt.Sprintf("custom app appAccessToken store, err:%v", err))
	}
	return appAccessTokenResp.AppAccessToken, err
}

// get internal tenant access token
func (r *request) customTenantAccessToken(ctx context.Context, app *App) (string, error) {
	rawResp, err := app.SendRequest(ctx, http.MethodPost, tenantAccessTokenInternalUrlPath, accessTokenTypeNone, &internalAccessTokenReq{
		AppID:     app.settings.id,
		AppSecret: app.settings.secret,
	})
	if err != nil {
		return "", err
	}
	tenantAccessTokenResp := &tenantAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, tenantAccessTokenResp)
	if err != nil {
		return "", err
	}
	if tenantAccessTokenResp.Code != 0 {
		return "", tenantAccessTokenResp.CodeError
	}
	expire := time.Duration(tenantAccessTokenResp.Expire)*time.Second - expiryDelta
	err = app.store.Put(ctx, tenantAccessTokenKey(app.settings.id, r.option.tenantKey), tenantAccessTokenResp.TenantAccessToken, expire)
	if err != nil {
		app.logger.Warn(ctx, fmt.Sprintf("custom app tenantAccessToken store, err:%v", err))
	}
	return tenantAccessTokenResp.TenantAccessToken, err
}

// get marketplace app access token
func (r *request) marketplaceAppAccessToken(ctx context.Context, app *App) (string, error) {
	appTicket, err := r.appTicket(ctx, app)
	if err != nil {
		return "", err
	}
	if appTicket == "" {
		return "", ErrAppTicketIsEmpty
	}
	rawResp, err := app.SendRequest(ctx, http.MethodPost, appAccessTokenUrlPath, accessTokenTypeNone, &marketplaceAppAccessTokenReq{
		AppID:     app.settings.id,
		AppSecret: app.settings.secret,
		AppTicket: appTicket,
	})
	if err != nil {
		return "", err
	}
	appAccessTokenResp := &appAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, appAccessTokenResp)
	if err != nil {
		return "", err
	}
	if appAccessTokenResp.Code != 0 {
		return "", appAccessTokenResp.CodeError
	}
	expire := time.Duration(appAccessTokenResp.Expire)*time.Second - expiryDelta
	err = app.store.Put(ctx, appAccessTokenKey(app.settings.id), appAccessTokenResp.AppAccessToken, expire)
	if err != nil {
		app.logger.Warn(ctx, fmt.Sprintf("marketplace app appAccessToken store, err:%v", err))
	}
	return appAccessTokenResp.AppAccessToken, err
}

// get marketplace tenant access token
func (r *request) marketplaceTenantAccessToken(ctx context.Context, app *App) (string, error) {
	appAccessToken, err := r.marketplaceAppAccessToken(ctx, app)
	if err != nil {
		return "", err
	}
	rawResp, err := app.SendRequest(ctx, http.MethodPost, tenantAccessTokenUrlPath, accessTokenTypeNone, &marketplaceTenantAccessTokenReq{
		AppAccessToken: appAccessToken,
		TenantKey:      r.option.tenantKey,
	})
	if err != nil {
		return "", err
	}
	tenantAccessTokenResp := &tenantAccessTokenResp{}
	err = json.Unmarshal(rawResp.RawBody, tenantAccessTokenResp)
	if err != nil {
		return "", err
	}
	if tenantAccessTokenResp.Code != 0 {
		return "", tenantAccessTokenResp.CodeError
	}
	expire := time.Duration(tenantAccessTokenResp.Expire)*time.Second - expiryDelta
	err = app.store.Put(ctx, tenantAccessTokenKey(app.settings.id, r.option.tenantKey), tenantAccessTokenResp.TenantAccessToken, expire)
	if err != nil {
		app.logger.Warn(ctx, fmt.Sprintf("custom app tenantAccessToken store, err:%v", err))
	}
	return tenantAccessTokenResp.TenantAccessToken, err
}

func (r *request) authorizationToHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (r *request) appTicket(ctx context.Context, app *App) (string, error) {
	return app.store.Get(ctx, appTicketKey(app.settings.id))
}

func (r *request) signAppAccessToken(ctx context.Context, httpRequest *http.Request, app *App) error {
	// from store get app access token
	if r.retryCount == 0 {
		tok, err := app.store.Get(ctx, appAccessTokenKey(app.settings.id))
		if err != nil {
			return err
		}
		if tok != "" {
			r.authorizationToHeader(httpRequest, tok)
			return nil
		}
	}
	// from api get app access token
	var appAccessToken string
	var err error
	if app.settings.type_ == AppTypeCustom {
		appAccessToken, err = r.customAppAccessToken(ctx, app)
	} else {
		appAccessToken, err = r.marketplaceAppAccessToken(ctx, app)
	}
	if err != nil {
		return err
	}
	r.authorizationToHeader(httpRequest, appAccessToken)
	return nil
}

func (r *request) signTenantAccessToken(ctx context.Context, httpRequest *http.Request, app *App) error {
	// from store get tenant access token
	if r.retryCount == 0 {
		tok, err := app.store.Get(ctx, tenantAccessTokenKey(app.settings.id, r.option.tenantKey))
		if err != nil {
			return err
		}
		if tok != "" {
			r.authorizationToHeader(httpRequest, tok)
			return nil
		}
	}
	// from api get tenant access token
	var tenantAccessToken string
	var err error
	if app.settings.type_ == AppTypeCustom {
		tenantAccessToken, err = r.customTenantAccessToken(ctx, app)
	} else {
		tenantAccessToken, err = r.marketplaceTenantAccessToken(ctx, app)
	}
	if err != nil {
		return err
	}
	r.authorizationToHeader(httpRequest, tenantAccessToken)
	return nil
}

func (r *request) signUserAccessToken(httpRequest *http.Request) error {
	r.authorizationToHeader(httpRequest, r.option.userAccessToken)
	return nil
}

func (r *request) readResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func (r *request) signHelpdeskAuthToken(rawRequest *http.Request, app *App) error {
	if r.option.needHelpDeskAuth {
		if app.settings.helpdeskAuthToken == "" {
			return errors.New("help desk API, please set the helpdesk information of lark.App")
		}
		rawRequest.Header.Set("X-Lark-Helpdesk-Authorization", app.settings.helpdeskAuthToken)
	}
	return nil
}

type tenantAccessTokenResp struct {
	CodeError
	Expire            int    `json:"expire"`
	TenantAccessToken string `json:"tenant_access_token"`
}

type internalAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type marketplaceAppAccessTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	AppTicket string `json:"app_ticket"`
}

type marketplaceTenantAccessTokenReq struct {
	AppAccessToken string `json:"app_access_token"`
	TenantKey      string `json:"tenant_key"`
}

type appAccessTokenResp struct {
	CodeError
	Expire         int    `json:"expire"`
	AppAccessToken string `json:"app_access_token"`
}

type applyAppTicketReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

const (
	appTicketKeyPrefix         = "app_ticket"
	appAccessTokenKeyPrefix    = "app_access_token"
	tenantAccessTokenKeyPrefix = "tenant_access_token"
)

func appTicketKey(appID string) string {
	return fmt.Sprintf("%s-%s", appTicketKeyPrefix, appID)
}

func appAccessTokenKey(appID string) string {
	return fmt.Sprintf("%s-%s", appAccessTokenKeyPrefix, appID)
}

func tenantAccessTokenKey(appID, tenantKey string) string {
	return fmt.Sprintf("%s-%s-%s", tenantAccessTokenKeyPrefix, appID, tenantKey)
}

type Formdata struct {
	fields map[string]interface{}
	data   *struct {
		content     []byte
		contentType string
	}
}

func NewFormdata() *Formdata {
	return &Formdata{}
}

func (fd *Formdata) AddField(field string, val interface{}) *Formdata {
	if fd.fields == nil {
		fd.fields = map[string]interface{}{}
	}
	fd.fields[field] = val
	return fd
}

func (fd *Formdata) AddFile(field string, r io.Reader) *Formdata {
	return fd.AddField(field, r)
}

func (fd *Formdata) content() (string, []byte, error) {
	if fd.data != nil {
		return fd.data.contentType, fd.data.content, nil
	}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	for key, val := range fd.fields {
		if r, ok := val.(io.Reader); ok {
			part, err := writer.CreateFormFile(key, "unknown-file")
			if err != nil {
				return "", nil, err
			}
			_, err = io.Copy(part, r)
			if err != nil {
				return "", nil, err
			}
			continue
		}
		err := writer.WriteField(key, fmt.Sprint(val))
		if err != nil {
			return "", nil, err
		}
	}
	contentType := writer.FormDataContentType()
	err := writer.Close()
	if err != nil {
		return "", nil, err
	}
	fd.data = &struct {
		content     []byte
		contentType string
	}{content: buf.Bytes(), contentType: contentType}
	return fd.data.contentType, fd.data.content, nil
}

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  *struct {
		Details              []*CodeErrorDetail              `json:"details,omitempty"`
		PermissionViolations []*CodeErrorPermissionViolation `json:"permission_violations,omitempty"`
		FieldViolations      []*CodeErrorFieldViolation      `json:"field_violations,omitempty"`
	} `json:"error"`
}

func (ce CodeError) Error() string {
	return ce.String()
}

func (ce CodeError) String() string {
	return Prettify(ce)
}

type CodeErrorDetail struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type CodeErrorPermissionViolation struct {
	Type        string `json:"type,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Description string `json:"description,omitempty"`
}

type CodeErrorFieldViolation struct {
	Field       string `json:"field,omitempty"`
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
}
