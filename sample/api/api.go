package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/api/core/response"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/test"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	"io/ioutil"
	"os"
)

var conf = test.GetInternalConf("online")

func main() {
	//testSendMessage()
	testUploadFile()
	//testDownloadFile()
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
	req := request.NewRequestWithNative("message/v4/send", "POST",
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
	err = api.Send(coreCtx, conf, request.NewRequestWithNative("image/v4/put", "POST",
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
	req := request.NewRequestWithNative("image/v4/get", "GET",
		request.AccessTokenTypeTenant, nil, ret,
		request.SetQueryParams(map[string]interface{}{"image_key": "[image key]"}), request.SetResponseStream())
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
