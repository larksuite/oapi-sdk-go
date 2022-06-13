package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	client2 "github.com/feishu/oapi-sdk-go"
	"github.com/feishu/oapi-sdk-go/core"
	"github.com/feishu/oapi-sdk-go/service/im/v1"
)

func uploadImage(client *client2.Client) {

	pdf, err := os.Open("/Users/bytedance/Downloads/a.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	resp, err := client.Im.Images.Create(context.Background(),
		im.NewCreateImageReqBuilder().
			Body(im.NewCreateImageReqBodyBuilder().
				ImageType(im.IMAGE_TYPE_MESSAGE).
				Image(pdf).
				Build()).
			Build())

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func uploadImage2(client *client2.Client) {
	body, err := im.NewCreateImagePathReqBodyBuilder().ImagePath("/Users/bytedance/Downloads/a.jpg").ImageType(im.IMAGE_TYPE_MESSAGE).Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Images.Create(context.Background(), im.NewCreateImageReqBuilder().Body(body).Build())

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func downLoadImage(client *client2.Client) {
	resp, err := client.Im.Images.Get(context.Background(), im.NewGetImageReqBuilder().ImageKey("img_v2_cd2657c7-ad1e-410a-8e76-942c89203bfg").Build())

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}

	if resp.Code != 0 {
		fmt.Println(core.Prettify(resp))
		return
	}
	fmt.Println(resp.FileName)
	fmt.Println(resp.RequestId())

	bs, err := ioutil.ReadAll(resp.File)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("test_download_v2.jpg", bs, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func downLoadImageV2(client *client2.Client) {
	resp, err := client.Im.Images.Get(context.Background(), im.NewGetImageReqBuilder().ImageKey("img_v2_cd2657c7-ad1e-410a-8e76-942c89203bfg").Build())

	if err != nil {
		fmt.Println(core.Prettify(err))
		return
	}

	if resp.Code != 0 {
		fmt.Println(core.Prettify(resp))
		return
	}
	fmt.Println(resp.FileName)
	fmt.Println(resp.RequestId())

	resp.WriteFile("a.jpg")

}
func main() {

	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var client = client2.NewClient(appID, appSecret)
	//downLoadImageV2(client)
	//uploadImage(client)
	//uploadImage(client)
	//downLoadImage(client)
	uploadImage2(client)

}
