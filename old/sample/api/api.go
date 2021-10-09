package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"github.com/larksuite/oapi-sdk-go/sample/configs"
	"io/ioutil"
	"os"
)

// for redis store and logrus
var conf = configs.TestConfigWithLogrusAndRedisStore(core.DomainFeiShu)

// var conf = configs.TestConfig("https://open.feishu.cn")
// var conf = configs.TestConfig(core.DomainFeiShu)

func main() {
	//testSendMessage()
	//testSendCardMessage()
	//testUploadFile()
	testDownloadFile()
}

// send card message
func testSendCardMessage() {
	coreCtx := core.WrapContext(context.Background())
	cardContent := "{\"config\":{\"wide_screen_mode\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"div\",\"text\":{\"tag\":\"lark_md\",\"content\":\"[飞书](https://www.feishu.cn)整合即时沟通、日历、音视频会议、云文档、云盘、工作台等功能于一体，成就组织和个人，更高效、更愉悦。\"}},{\"tag\":\"action\",\"actions\":[{\"tag\":\"button\",\"text\":{\"tag\":\"plain_text\",\"content\":\"主按钮\"},\"type\":\"primary\",\"value\":{\"key\":\"primary\"}},{\"tag\":\"button\",\"text\":{\"tag\":\"plain_text\",\"content\":\"次按钮\"},\"type\":\"default\",\"value\":{\"key\":\"default\"}}]}]}}"
	card := map[string]interface{}{}
	err := json.Unmarshal([]byte(cardContent), &card)
	if err != nil {
		panic(err)
	}
	body := map[string]interface{}{
		"user_id":  "77bbc392",
		"msg_type": "interactive",
		"card":     card,
	}
	ret := make(map[string]interface{})
	req := request.NewRequestWithNative("/open-apis/message/v4/send", "POST",
		request.AccessTokenTypeTenant, body, &ret,
		//应用市场应用 request.SetTenantKey("TenantKey"),
	)
	err = api.Send(coreCtx, conf, req)
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(ret))
}

// send message
func testSendMessage() {
	coreCtx := core.WrapContext(context.Background())
	body := map[string]interface{}{
		"open_id":  "[open_id]",
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": "test",
		},
	}
	ret := make(map[string]interface{})
	req := request.NewRequestWithNative("/open-apis/message/v4/send", "POST",
		request.AccessTokenTypeTenant, body, &ret,
		//应用市场应用 request.SetTenantKey("TenantKey"),
	)
	err := api.Send(coreCtx, conf, req)
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(ret))
}

type UploadImage struct {
	ImageKey string `json:"image_key"`
}

// upload image
func testUploadFile() {
	coreCtx := core.WrapContext(context.Background())
	// coreCtx.Set(constants.HTTPHeaderKeyRequestID, "2020122212081301001702714534518-xxxxx")
	var formData = request.NewFormData()
	formData.AddParam("image_type", "message")
	bs, err := ioutil.ReadFile("test.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	formData.AddFile("image", request.NewFile().SetContent(bs))
	/*
		// stream upload, file implement io.Reader
		file, err := os.Open("test.png")
		if err != nil {
			fmt.Println(err)
			return
		}
		formData.AddFile("image", request.NewFile().SetContentStream(file))
	*/
	ret := &UploadImage{}
	err = api.Send(coreCtx, conf, request.NewRequestWithNative("/open-apis/image/v4/put", "POST",
		request.AccessTokenTypeTenant, formData, ret))
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	fmt.Println(tools.Prettify(ret))
}

// download image
func testDownloadFile() {
	coreCtx := core.WrapContext(context.Background())
	ret := &bytes.Buffer{}
	/*
		// stream download: ret implement io.Writer
		ret, err := os.Create("[file path]")
		if err != nil {
			fmt.Println(err)
			return
		}
	*/
	req := request.NewRequestWithNative("/open-apis/image/v4/get", "GET",
		request.AccessTokenTypeTenant, nil, ret,
		request.SetQueryParams(map[string]interface{}{"image_key": "img_v2_f6203671-41d6-46ed-adc9-c50aa840330g"}), request.SetResponseStream())
	err := api.Send(coreCtx, conf, req)
	fmt.Println(coreCtx.GetRequestID())
	fmt.Println(coreCtx.GetHTTPStatusCode())
	if err != nil {
		fmt.Println(tools.Prettify(err))
		e := err.(*response.Error)
		fmt.Println(e.Code)
		fmt.Println(e.Msg)
		return
	}
	err = ioutil.WriteFile("test_download.png", ret.Bytes(), os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
}
