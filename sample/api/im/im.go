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

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func uploadImage(client *lark.Client) {
	pdf, err := os.Open("/Users/bytedance/Downloads/封面.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	resp, err := client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.CreateImageImageTypeMessage).
				Image(pdf).
				Build()).
			Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadFile(client *lark.Client) {
	file, err := os.Open("/Users/bytedance/Downloads/领域特定语言.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	resp, err := client.Im.File.Create(context.Background(),
		larkim.NewCreateFileReqBuilder().
			Body(larkim.NewCreateFileReqBodyBuilder().
				FileType(larkim.CreateFileFileTypePdf).
				FileName("open-redis.pdf").
				File(file).
				Build()).
			Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func downFile(client *lark.Client) {

	resp, err := client.Im.File.Get(context.Background(), larkim.NewGetFileReqBuilder().FileKey("file_v2_21d3ce59-5d45-4421-9e09-a1733fdd67fg").Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	resp.WriteFile("a.mp4")

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadImage2(client *lark.Client) {
	body, err := larkim.NewCreateImagePathReqBodyBuilder().
		ImagePath("/Users/bytedance/Downloads/b.jpg").
		ImageType(larkim.CreateImageImageTypeMessage).
		Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Image.Create(context.Background(), larkim.NewCreateImageReqBuilder().Body(body).Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func downLoadImage(client *lark.Client) {
	resp, err := client.Im.Image.Get(context.Background(), larkim.NewGetImageReqBuilder().ImageKey("img_v2_b9f85d3e-a4e5-4ae7-9b8a-d72880189b1g").Build())

	if err != nil {
		fmt.Println(larkcore.Prettify(err))
		return
	}

	if resp.Code != 0 {
		fmt.Println(larkcore.Prettify(resp))
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

func downLoadImageV2(client *lark.Client) {
	resp, err := client.Im.Image.Get(context.Background(), larkim.NewGetImageReqBuilder().ImageKey("img_v112_cd2657c7-ad1e-410a-8e76-942c89203bfg").Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Code != 0 {
		fmt.Println(resp)
		return
	}
	fmt.Println(resp.FileName)
	fmt.Println(resp.RequestId())

	resp.WriteFile("a.jpg")
}

func sendTextMsg(client *lark.Client) {

	// Build the text message content
	//content := "{\"text\":\"hello,world\\n<at user_id=\\\"ou_c245b0a7dff2725cfa2fb104f8b48b9d\\\">加多</at>\"}"
	content := larkim.NewTextMsgBuilder().
		Text("hello,world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "加多").
		Build()

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId("ou_e8daec8c7bd6269852c84239ac85db3e").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(resp.Data.MessageId)
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendInteractiveMsg(client *lark.Client) {
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := larkcard.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateRed).
		Title(larkcard.NewMessageCardPlainText().
			Content("1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content("**🕐 时间：**\\n2021-02-23 20:17:51").
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 谁处理了问题
	content := "✅ " + "name" + "已处理了此告警"
	processPersonElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement, processPersonElement}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(cardContent).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

// 运维报警通知
// https://open.feishu.cn/tool/cardbuilder?from=cotentmodule
func sendInteractiveMonitorMsg(client *lark.Client) {
	// config
	config := larkcard.NewMessageCardConfig().
		EnableForward(true).
		UpdateMulti(true).
		Build()

	// header
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateRed).
		Title(larkcard.NewMessageCardPlainText().
			Content("1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement1 := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**🕐 时间：**2021-02-23 20:17:51").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**🔢 事件 ID：：**336720").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),

			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**📋 项目：**\nQA 7").
					Build()).
				IsShort(true).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**👤 一级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build(),

			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),
			larkcard.NewMessageCardField().
				Text(larkcard.NewMessageCardLarkMd().
					Content("**👤 二级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build()}).
		Build()

	divElement3 := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content("🔴 支付失败数  🔵 支付成功数").
			Build()}).
		Build()

	divElement4 := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{larkcard.NewMessageCardEmbedButton().
			Type(larkcard.MessageCardButtonTypePrimary).
			Value(map[string]interface{}{"key1": "value1"}).
			Text(larkcard.NewMessageCardPlainText().
				Content("跟进处理").
				Build()),
			larkcard.NewMessageCardEmbedSelectMenuStatic().
				MessageCardEmbedSelectMenuStatic(larkcard.NewMessageCardEmbedSelectMenuBase().
					Options([]*larkcard.MessageCardEmbedSelectOption{larkcard.NewMessageCardEmbedSelectOption().
						Value("1").
						Text(larkcard.NewMessageCardPlainText().
							Content("屏蔽10分钟").
							Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("2").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽30分钟").
								Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("3").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽1小时").
								Build()),
						larkcard.NewMessageCardEmbedSelectOption().
							Value("4").
							Text(larkcard.NewMessageCardPlainText().
								Content("屏蔽24小时").
								Build()),
					}).
					Placeholder(larkcard.NewMessageCardPlainText().
						Content("暂时屏蔽报警").
						Build()).
					Value(map[string]interface{}{"key": "value"}).
					Build()).
				Build()}).
		Build()

	divElement5 := larkcard.NewMessageCardHr().Build()

	divElement6 := larkcard.NewMessageCardDiv().
		Text(larkcard.NewMessageCardLarkMd().
			Content("🙋🏼 [我要反馈误报](https://open.feishu.cn/) | 📝 [录入报警处理过程](https://open.feishu.cn/)").
			Build()).
		Build()

	// CardUrl
	cardLink := larkcard.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	low := "low"
	priority := larkcard.NewMessageCardMarkdown().
		Content(fmt.Sprintf(`**Priority**: (~~*%s*~~)  **%s**`, low, "high")).
		Build()
	fmt.Println(priority)
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement1, divElement3, divElement4, divElement5, divElement6, priority}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return
	}

	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(cardContent).
			Build()).
		Build()
	resp, err := client.Im.Message.Create(context.Background(), req)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendImageMsg(client *lark.Client) {
	msgImage := larkim.MessageImage{ImageKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeImage).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendShardChatMsg(client *lark.Client) {
	msgImage := larkim.MessageShareChat{ChatId: "oc_4ea14cc15e39ef80a579ca74895a33ff"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeShareChat).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendShardUserMsg(client *lark.Client) {
	msgImage := larkim.MessageShareUser{UserId: "ou_487f709a942d16edafe57fd6fbc4bcf5"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeShareUser).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendAudioMsg(client *lark.Client) {
	msgImage := larkim.MessageAudio{FileKey: "file_v2_d37fd916-1b78-4dc7-b3e9-f5a57eb40ddg"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeAudio).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendMediaMsg(client *lark.Client) {
	msgImage := larkim.MessageMedia{FileKey: "file_v2_cd6a059f-f143-491a-a0b3-ec2823784dbg", ImageKey: "img_v2_8ebfff45-629a-4e99-a7c3-c5f89a3d0d1g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeMedia).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendFileMsg(client *lark.Client) {
	msgImage := larkim.MessageFile{FileKey: "file_v2_4fa17cda-01f3-4aac-927a-7833ab482fcg"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeFile).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendStickerMsg(client *lark.Client) {
	msgImage := larkim.MessageSticker{FileKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeSticker).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendPostMsg(client *lark.Client) {
	// 2.1 创建text与href元素
	zhCnPostText := &larkim.MessagePostText{Text: "中文内容", UnEscape: false}
	zhCnPostA := &larkim.MessagePostA{Text: "test content", Href: "http://www.baidu.com", UnEscape: false}
	enUsPostText := &larkim.MessagePostText{Text: "英文内容", UnEscape: false}
	enUsPostA := &larkim.MessagePostA{Text: "test content", Href: "http://www.baidu.com", UnEscape: false}

	// 2.2 构建消息content
	zhCnMessagePostContent := &larkim.MessagePostContent{Title: "title1", Content: [][]larkim.MessagePostElement{{zhCnPostText, zhCnPostA}}}
	enUsMessagePostContent := &larkim.MessagePostContent{Title: "title2", Content: [][]larkim.MessagePostElement{{enUsPostText, enUsPostA}}}
	messagePostText := &larkim.MessagePost{ZhCN: zhCnMessagePostContent, EnUS: enUsMessagePostContent}
	content, err := messagePostText.String()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypePost).
			ReceiveId("ou_e8daec8c7bd6269852c84239ac85db3e").
			Content(content).
			Build()).
		Build(), larkcore.WithRequestId("jiaduo_id"))

	if err != nil {
		switch er := err.(type) {
		case *larkcore.IllegalParamError:
			fmt.Println(er.Error()) // 处理错误
		case *larkcore.ClientTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *larkcore.ServerTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *larkcore.DialFailedError:
			fmt.Println(er.Error()) // 处理错误
		default:
			//其他处理
			fmt.Println(err)
		}
		return
	}

	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendPostMsgUseBuilder(client *lark.Client) {
	// 第一行
	// 文本 &超链接
	zhCnPostText := &larkim.MessagePostText{Text: "第一行:", UnEscape: false}
	enUsPostText := &larkim.MessagePostText{Text: "英文内容", UnEscape: false}

	zhCnPostA := &larkim.MessagePostA{Text: "超链接", Href: "http://www.baidu.com", UnEscape: false}
	enUsPostA := &larkim.MessagePostA{Text: "link", Href: "http://www.baidu.com", UnEscape: false}

	// At人
	zhCnPostAt := &larkim.MessagePostAt{UserId: "ou_c245b0a7dff2725cfa2fb104f8b48b9d", UserName: "加多"}
	enCnPostAt := &larkim.MessagePostAt{UserId: "ou_c245b0a7dff2725cfa2fb104f8b48b9d", UserName: "jiaduo"}

	// 图片
	//zhCnPostImage := &larkim.MessagePostImage{ImageKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}
	//enCnPostImage := &larkim.MessagePostImage{ImageKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}

	// 第二行
	// 文本 &超链接
	zhCnPostText21 := &larkim.MessagePostText{Text: "第二行:", UnEscape: false}
	enUsPostText21 := &larkim.MessagePostText{Text: "英文内容", UnEscape: false}

	zhCnPostText22 := &larkim.MessagePostText{Text: "文本测试", UnEscape: false}
	enUsPostText22 := &larkim.MessagePostText{Text: "英文内容", UnEscape: false}

	// 图片
	//zhCnPostImage2 := &larkim.MessagePostImage{ImageKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}
	//enCnPostImage2 := &larkim.MessagePostImage{ImageKey: "img_v2_a66c4f79-c7b5-4899-b5e3-622766c4f82g"}

	// 中文
	zhCn := larkim.NewMessagePostContent().
		ContentTitle("我是一个标题").
		AppendContent([]larkim.MessagePostElement{zhCnPostText, zhCnPostA, zhCnPostAt}).
		//AppendContent([]larkim.MessagePostElement{zhCnPostImage}).
		AppendContent([]larkim.MessagePostElement{zhCnPostText21, zhCnPostText22}).
		//AppendContent([]larkim.MessagePostElement{zhCnPostImage2}).
		Build()

	// 英文
	enUs := larkim.NewMessagePostContent().
		ContentTitle("im a title").
		AppendContent([]larkim.MessagePostElement{enUsPostA, enUsPostText, enCnPostAt}).
		//AppendContent([]larkim.MessagePostElement{enCnPostImage}).
		AppendContent([]larkim.MessagePostElement{enUsPostText21, enUsPostText22}).
		//AppendContent([]larkim.MessagePostElement{enCnPostImage2}).
		Build()

	// 构建消息体
	postText, err := larkim.NewMessagePost().ZhCn(zhCn).EnUs(enUs).Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypePost).
			ReceiveId("ou_e3f3fca5204cdf7552531c84a32f60d1").
			Content(postText).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.Success() {
		fmt.Println(larkcore.Prettify(resp))
		fmt.Println(*resp.Data.MessageId)
	} else {
		fmt.Println(resp.RequestId(), resp.Msg, resp.Code)
	}

}

func sendRawReq(cli *lark.Client) {
	content := larkim.NewTextMsgBuilder().
		Text("加多").
		Line().
		TextLine("hello").
		TextLine("world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "陆续").
		Build()

	// 放到client里面
	resp, err := cli.Post(context.Background(), "/open-apis/im/v1/messages?receive_id_type=open_id", map[string]interface{}{
		"receive_id": "ou_c245b0a7dff2725cfa2fb104f8b48b9d",
		"msg_type":   "text",
		"content":    "{\"text\":\"<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at> test content\"}",
	}, larkcore.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err, content)
		return
	}

	fmt.Println(resp)
}

func sendRawImageReq(cli *lark.Client) {
	img, err := os.Open("/Users/bytedance/Downloads/go-icon.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer img.Close()

	formData := larkcore.NewFormdata().
		AddField("image_type", "message").
		AddFile("image", img)

	resp, err := cli.Post(context.Background(), "/open-apis/im/v1/images", formData, larkcore.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}

type CustomLogger struct {
}

func (c *CustomLogger) Debug(context.Context, ...interface{}) {
	// 实现自己的日志逻辑
}

func (c *CustomLogger) Info(context.Context, ...interface{}) {
	// 实现自己的日志逻辑
}

func (c *CustomLogger) Warn(context.Context, ...interface{}) {
	// 实现自己的日志逻辑
}

func (c *CustomLogger) Error(context.Context, ...interface{}) {
	// 实现自己的日志逻辑
}

type CustomCache struct {
}

func (c *CustomCache) Set(ctx context.Context, key string, value string, expireTime time.Duration) error {
	// 缓存token
	return nil
}

func (c *CustomCache) Get(ctx context.Context, key string) (string, error) {
	// 获取token
	token := ""
	return token, nil
}

type CustomHttpClient struct {
	httpClient *http.Client
}

func (c *CustomHttpClient) Do(req *http.Request) (*http.Response, error) {
	// 请求前做点事情
	// 发起请求
	resp, err := c.httpClient.Do(req)
	// 请求后做点事情
	return resp, err
}

func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	client := lark.NewClient(appID, appSecret)

	// 发送文本消息
	//sendTextMsg(client)

	// 发送富文本消息
	sendPostMsgUseBuilder(client)

	// 发送图片消息
	//uploadImage(client)
	//sendImageMsg(client)

	// 发送文件消息
	//uploadFile(client)
	//sendFileMsg(client)

	// 发送交互卡片消息
	//sendInteractiveMonitorMsg(client)

	// 发送群名片消息
	//sendShardChatMsg(client)

	// 发送个人名片消息
	//sendShardUserMsg(client)

	// 发送语音 audio
	//uploadFile(client)
	//sendAudioMsg(client)

	// 发送视频消息
	//uploadFile(client)
	//uploadImage(client)
	//sendMediaMsg(client)

	// 发送表情
	//sendStickerMsg(client)

	//downFile(client)
	//downLoadImageV2(feishu_client)
	//uploadImage(client)
	//downLoadImage(feishu_client)
	//uploadImage2(feishu_client)
	//sendRawReq(feishu_client)
	//sendRawImageReq(feishu_client)
	//uploadFile(feishu_client)
	//sendAudioMsg(client)
	//sendShardUserMsg(client)
	//sendPostMsg(feishu_client)
	//testCreate(feishu_client)
	//sendInteractiveMonitorMsg(feishu_client)
	//sendInteractiveMonitorMsg(feishu_client)
	//sendPostMsgUseBuilder(feishu_client)

}
