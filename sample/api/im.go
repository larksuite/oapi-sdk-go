package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/larksuite/oapi-sdk-go.v3"
	"github.com/larksuite/oapi-sdk-go.v3/card"
	"github.com/larksuite/oapi-sdk-go.v3/core"
	"github.com/larksuite/oapi-sdk-go.v3/service/gray_test_open_sg/v1"
	"github.com/larksuite/oapi-sdk-go.v3/service/im/v1"
)

func uploadImage(client *lark.Client) {
	pdf, err := os.Open("/Users/bytedance/Downloads/a.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	resp, err := client.Im.Image.Create(context.Background(),
		larkim.NewCreateImageReqBuilder().
			Body(larkim.NewCreateImageReqBodyBuilder().
				ImageType(larkim.ImageTypeMessage).
				Image(pdf).
				Build()).
			Build())

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadFile(client *lark.Client) {
	pdf, err := os.Open("/Users/bytedance/Downloads/redis.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pdf.Close()

	resp, err := client.Im.File.Create(context.Background(),
		larkim.NewCreateFileReqBuilder().
			Body(larkim.NewCreateFileReqBodyBuilder().
				FileType(larkim.FileTypePdf).
				FileName("redis.pdf").
				File(pdf).
				Build()).
			Build())

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadImage2(client *lark.Client) {
	body, err := larkim.NewCreateImagePathReqBodyBuilder().
		ImagePath("/Users/bytedance/Downloads/b.jpg").
		ImageType(larkim.ImageTypeMessage).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func downLoadImage(client *lark.Client) {
	resp, err := client.Im.Image.Get(context.Background(), larkim.NewGetImageReqBuilder().ImageKey("img_v2_9068cbd5-71d8-4799-b29e-a01650b1328g").Build())

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
	// 构建消息体
	content := larkim.NewTextMsgBuilder().
		Line().
		TextLine("hello").
		TextLine("world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "陆续").
		Build()

	header := make(http.Header)

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build(), larkcore.WithHeaders(header))

	if err != nil {
		fmt.Println(err)
		return
	}
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
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

// 运维报警通知
//https://open.feishu.cn/tool/cardbuilder?from=cotentmodule
func sendInteractiveMonitorMsg(client *lark.Client) {
	// config
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(false).
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

	divElement2 := larkcard.NewMessageCardImage().
		Alt(larkcard.NewMessageCardPlainText().
			Content("").
			Build()).
		ImgKey("img_v2_8b2ebeaf-c97c-411d-a4dc-4323e8cba10g").
		Title(larkcard.NewMessageCardLarkMd().
			Content("支付方式 支付成功率低于 50%：").
			Build()).
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

	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]larkcard.MessageCardElement{divElement1, divElement2, divElement3, divElement4, divElement5, divElement6}).
		CardLink(cardLink).
		String()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendImageMsg(client *lark.Client) {
	msgImage := larkim.MessageImage{ImageKey: "img_v2_63554b3a-b60f-449a-a286-0f89e353815g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendAudioMsg(client *lark.Client) {
	msgImage := larkim.MessageAudio{FileKey: "file_v2_0c7f5b4b-64ec-408a-a9eb-09aec7954a4g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendMediaMsg(client *lark.Client) {
	msgImage := larkim.MessageMedia{FileKey: "file_v2_0c7f5b4b-64ec-408a-a9eb-09aec7954a4g", ImageKey: "ssss"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendStickerMsg(client *lark.Client) {
	msgImage := larkim.MessageSticker{FileKey: "img_v2_2a0372ea-dc03-44d7-b052-255bede4d42g"}
	content, err := msgImage.String()
	if err != nil {
		fmt.Println(err)
		return
	}
	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeSticker).
			ReceiveId("121212").
			Content(content).
			Build()).
		Build())

	if err != nil {
		fmt.Println(err)
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
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypePost).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
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
	fmt.Println(larkcore.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendPostMsgUseBuilder(client *lark.Client) {
	zhCnPostText := &larkim.MessagePostText{Text: "中文内容", UnEscape: false}
	zhCnPostA := &larkim.MessagePostA{Text: "test content", Href: "http://www.baidu.com", UnEscape: false}
	enUsPostText := &larkim.MessagePostText{Text: "英文内容", UnEscape: false}
	enUsPostA := &larkim.MessagePostA{Text: "test content", Href: "http://www.baidu.com", UnEscape: false}

	zhCn := larkim.NewMessagePostContent().
		ContentTitle("title1").
		AppendContent([]larkim.MessagePostElement{zhCnPostText, zhCnPostA}).
		Build()

	enUs := larkim.NewMessagePostContent().
		ContentTitle("title1").
		AppendContent([]larkim.MessagePostElement{enUsPostA, enUsPostText}).
		Build()
	postText, err := larkim.NewMessagePost().ZhCn(zhCn).EnUs(enUs).Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypePost).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(postText).
			Build()).
		Build())

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

	if resp.Success() {
		fmt.Println(larkcore.Prettify(resp))
	} else {
		fmt.Println(resp.RequestId(), resp.Msg, resp.Code)
	}

}

func testCreate(client *lark.Client) {
	resp, err := client.GrayTestOpenSg.Moto.Create(context.Background(), larkgray_test_open_sg.
		NewCreateMotoReqBuilder().
		Level(larkgray_test_open_sg.
			NewLevelBuilder().
			Body("ss").
			Level("level").
			Type("ss").
			Build()).
		Build())
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

	if resp.Success() {
		fmt.Println(larkcore.Prettify(resp))
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

	formData := larkcore.NewFormdata().AddField("image_type", "message").AddFile("image", img)

	resp, err := cli.Post(context.Background(), "/open-apis/im/v1/images", formData, larkcore.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var feishu_client = lark.NewClient(appID, appSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithLogReqRespInfoAtDebugLevel(true),
	)

	//downLoadImageV2(feishu_client)
	//uploadImage(feishu_client)
	//uploadImage(client)
	//downLoadImage(feishu_client)
	//uploadImage2(feishu_client)
	//sendTextMsg(feishu_client)
	//sendRawReq(feishu_client)
	//sendRawImageReq(feishu_client)
	//sendImageMsg(feishu_client)
	//uploadFile(feishu_client)
	//sendFileMsg(feishu_client)
	//sendAudioMsg(client)
	//sendMediaMsg(client)
	//sendShardChatMsg(client)
	//sendShardUserMsg(client)
	//sendPostMsg(feishu_client)
	//sendPostMsgUseBuilder(feishu_client)
	//testCreate(feishu_client)
	//sendInteractiveMsg(feishu_client)
	sendInteractiveMonitorMsg(feishu_client)

}
