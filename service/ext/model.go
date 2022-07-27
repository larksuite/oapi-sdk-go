/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice, shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package larkext

import larkcore "github.com/larksuite/oapi-sdk-go/v3/core"

const (
	FileTypeDoc     = "doc"
	FileTypeSheet   = "sheet"
	FileTypeBitable = "bitable"
)

type CreateFileReq struct {
	apiReq *larkcore.ApiReq
	Body   *CreateFileReqBody `body:""`
}

type CreateFileResp struct {
	*larkcore.ApiResp `json:"-"`
	larkcore.CodeError
	Data *CreateFileRespData `json:"data"`
}

func (c *CreateFileResp) Success() bool {
	return c.Code == 0
}

type CreateFileRespData struct {
	Url      string `json:"url,omitempty"`
	Token    string `json:"token,omitempty"`
	Revision int64  `json:"revision,omitempty"`
}

type CreateFileReqBody struct {
	Title string `json:"title,omitempty"`
	Type_ string `json:"type,omitempty"`
}

type CreateFileReqBodyBuilder struct {
	title string `json:"title,omitempty"`
	type_ string `json:"type,omitempty"`
}

func NewCreateFileReqBodyBuilder() *CreateFileReqBodyBuilder {
	return &CreateFileReqBodyBuilder{}
}

func (c *CreateFileReqBodyBuilder) Title(title string) *CreateFileReqBodyBuilder {
	c.title = title
	return c
}

func (c *CreateFileReqBodyBuilder) Type(type_ string) *CreateFileReqBodyBuilder {
	c.type_ = type_
	return c
}

func (c *CreateFileReqBodyBuilder) Build() *CreateFileReqBody {
	body := &CreateFileReqBody{}
	body.Type_ = c.type_
	body.Title = c.title
	return body
}

type CreateFileReqBuilder struct {
	apiReq *larkcore.ApiReq
	body   *CreateFileReqBody `body:""`
}

func NewCreateFileReqBuilder() *CreateFileReqBuilder {
	builder := &CreateFileReqBuilder{}
	builder.apiReq = &larkcore.ApiReq{
		PathParams:  larkcore.PathParams{},
		QueryParams: larkcore.QueryParams{},
	}
	return builder
}

func (c *CreateFileReqBuilder) FolderToken(folderToken string) *CreateFileReqBuilder {
	c.apiReq.PathParams.Set("folderToken", folderToken)
	return c
}

func (c *CreateFileReqBuilder) Body(body *CreateFileReqBody) *CreateFileReqBuilder {
	c.body = body
	return c
}

func (c *CreateFileReqBuilder) Build() *CreateFileReq {
	req := &CreateFileReq{}
	req.apiReq = &larkcore.ApiReq{}
	req.apiReq.Body = c.body
	req.apiReq.PathParams = c.apiReq.PathParams
	return req
}
