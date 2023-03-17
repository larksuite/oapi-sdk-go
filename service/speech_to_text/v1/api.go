// Package speech_to_text code generated by oapi sdk gen
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

package larkspeech_to_text

import (
	"context"
	"net/http"

	"github.com/larksuite/oapi-sdk-go/v3/core"
)

func NewService(config *larkcore.Config) *SpeechToTextService {
	s := &SpeechToTextService{config: config}
	s.Speech = &speech{service: s}
	return s
}

type SpeechToTextService struct {
	config *larkcore.Config
	Speech *speech // 语音识别
}

type speech struct {
	service *SpeechToTextService
}

// 语音文件识别 (ASR)
//
// - 语音文件识别接口，上传整段语音文件进行一次性识别。接口适合 60 秒以内音频识别
//
// - 单租户限流：20QPS，同租户下的应用没有限流，共享本租户的 20QPS 限流
//
// - 官网 API 文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/ai/speech_to_text-v1/speech/file_recognize
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/speech_to_textv1/fileRecognize_speech.go
func (s *speech) FileRecognize(ctx context.Context, req *FileRecognizeSpeechReq, options ...larkcore.RequestOptionFunc) (*FileRecognizeSpeechResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/speech_to_text/v1/speech/file_recognize"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, s.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &FileRecognizeSpeechResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, s.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// 语音流式识别 (ASR)
//
// - 语音流式接口，将整个音频文件分片进行传入模型。能够实时返回数据。建议每个音频分片的大小为 100-200ms
//
// - 单租户限流：20 路（一个 stream_id 称为一路会话），同租户下的应用没有限流，共享本租户的 20 路限流
//
// - 官网 API 文档链接:https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/ai/speech_to_text-v1/speech/stream_recognize
//
// - 使用 Demo 链接:https://github.com/larksuite/oapi-sdk-go/tree/v3_main/sample/apiall/speech_to_textv1/streamRecognize_speech.go
func (s *speech) StreamRecognize(ctx context.Context, req *StreamRecognizeSpeechReq, options ...larkcore.RequestOptionFunc) (*StreamRecognizeSpeechResp, error) {
	// 发起请求
	apiReq := req.apiReq
	apiReq.ApiPath = "/open-apis/speech_to_text/v1/speech/stream_recognize"
	apiReq.HttpMethod = http.MethodPost
	apiReq.SupportedAccessTokenTypes = []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant}
	apiResp, err := larkcore.Request(ctx, apiReq, s.service.config, options...)
	if err != nil {
		return nil, err
	}
	// 反序列响应结果
	resp := &StreamRecognizeSpeechResp{ApiResp: apiResp}
	err = apiResp.JSONUnmarshalBody(resp, s.service.config)
	if err != nil {
		return nil, err
	}
	return resp, err
}
