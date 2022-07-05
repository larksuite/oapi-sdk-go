package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	client "github.com/larksuite/oapi-sdk-go"
	"github.com/larksuite/oapi-sdk-go/card"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/service/gray_test_open_sg/v1"
	"github.com/larksuite/oapi-sdk-go/service/im/v1"
)

func uploadImage(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadFile(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func uploadImage2(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func downLoadImage(client *client.Client) {
	resp, err := client.Im.Image.Get(context.Background(), larkim.NewGetImageReqBuilder().ImageKey("img_v2_9068cbd5-71d8-4799-b29e-a01650b1328g").Build())

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

func downLoadImageV2(client *client.Client) {
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

func sendTextMsg(client *client.Client) {
	content := larkim.NewTextMsgBuilder().
		Text("加多").
		Line().
		TextLine("hello").
		TextLine("world").
		AtUser("ou_c245b0a7dff2725cfa2fb104f8b48b9d", "陆续").
		//AtAll().
		Build()

	header := make(http.Header)

	resp, err := client.Im.Message.Create(context.Background(), larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeOpenId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId("ou_c245b0a7dff2725cfa2fb104f8b48b9d").
			Content(content).
			Build()).
		Build(), core.WithHeaders(header))

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendInteractiveMsg(client *client.Client) {
	// config
	config := card.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(true).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := card.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := card.NewMessageCardHeader().
		Template("red").
		Title(card.NewMessageCardPlainText().
			Content("1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement := card.NewMessageCardDiv().
		Fields([]*card.MessageCardField{card.NewMessageCardField().
			Text(card.NewMessageCardLarkMd().
				Content("**🕐 时间：**\\n2021-02-23 20:17:51").
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 谁处理了问题
	content := "✅ " + "name" + "已处理了此告警"
	processPersonElement := card.NewMessageCardDiv().
		Fields([]*card.MessageCardField{card.NewMessageCardField().
			Text(card.NewMessageCardLarkMd().
				Content(content).
				Build()).
			IsShort(true).
			Build()}).
		Build()

	// 卡片消息体
	cardContent, err := card.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]card.MessageCardElement{divElement, processPersonElement}).
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

// 运维报警通知
//https://open.feishu.cn/tool/cardbuilder?from=cotentmodule
func sendInteractiveMonitorMsg(client *client.Client) {
	// config
	config := card.NewMessageCardConfig().
		WideScreenMode(true).
		EnableForward(false).
		UpdateMulti(false).
		Build()

	// CardUrl
	cardLink := card.NewMessageCardURL().
		PcUrl("http://www.baidu.com").
		IoSUrl("http://www.google.com").
		Url("http://open.feishu.com").
		AndroidUrl("http://www.jianshu.com").
		Build()

	// header
	header := card.NewMessageCardHeader().
		Template("red").
		Title(card.NewMessageCardPlainText().
			Content("1 级报警 - 数据平台").
			Build()).
		Build()

	// Elements
	divElement1 := card.NewMessageCardDiv().
		Fields([]*card.MessageCardField{
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("**🕐 时间：**2021-02-23 20:17:51").
					Build()).
				IsShort(true).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("**🔢 事件 ID：：**336720").
					Build()).
				IsShort(true).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("**📋 项目：**\nQA 7").
					Build()).
				IsShort(true).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("**👤 一级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("").
					Build()).
				IsShort(false).
				Build(),
			card.NewMessageCardField().
				Text(card.NewMessageCardLarkMd().
					Content("**👤 二级值班：**\n<at id=ou_c245b0a7dff2725cfa2fb104f8b48b9d>加多</at>").
					Build()).
				IsShort(true).
				Build()}).
		Build()

	divElement2 := card.NewMessageCardImage().
		Alt(card.NewMessageCardPlainText().
			Content("").
			Build()).
		ImgKey("img_v2_8b2ebeaf-c97c-411d-a4dc-4323e8cba10g").
		Title(card.NewMessageCardLarkMd().
			Content("支付方式 支付成功率低于 50%：").
			Build()).
		Build()

	divElement3 := card.NewMessageCardNote().
		Elements([]card.MessageCardNoteElement{card.NewMessageCardPlainText().
			Content("🔴 支付失败数  🔵 支付成功数").
			Build()}).
		Build()

	divElement4 := card.NewMessageCardAction().
		Actions([]card.MessageCardActionElement{card.NewMessageCardEmbedButton().
			Type(card.MessageCardButtonTypePrimary).
			Value(map[string]interface{}{"key1": "value1"}).
			Text(card.NewMessageCardPlainText().
				Content("跟进处理").
				Build()),
			card.NewMessageCardEmbedSelectMenuStatic().
				MessageCardEmbedSelectMenuStatic(card.NewMessageCardEmbedSelectMenuBase().
					Options([]*card.MessageCardEmbedSelectOption{card.NewMessageCardEmbedSelectOption().
						Value("1").
						Text(card.NewMessageCardPlainText().
							Content("屏蔽10分钟").
							Build()),
						card.NewMessageCardEmbedSelectOption().
							Value("2").
							Text(card.NewMessageCardPlainText().
								Content("屏蔽30分钟").
								Build()),
						card.NewMessageCardEmbedSelectOption().
							Value("3").
							Text(card.NewMessageCardPlainText().
								Content("屏蔽1小时").
								Build()),
						card.NewMessageCardEmbedSelectOption().
							Value("4").
							Text(card.NewMessageCardPlainText().
								Content("屏蔽24小时").
								Build()),
					}).
					Placeholder(card.NewMessageCardPlainText().
						Content("暂时屏蔽报警").
						Build()).
					Value(map[string]interface{}{"key": "value"}).
					Build()).
				Build()}).
		Build()

	divElement5 := card.NewMessageCardHr().Build()

	divElement6 := card.NewMessageCardDiv().
		Text(card.NewMessageCardLarkMd().
			Content("🙋🏼 [我要反馈误报](https://open.feishu.cn/) | 📝 [录入报警处理过程](https://open.feishu.cn/)").
			Build()).
		Build()

	// 卡片消息体
	cardContent, err := card.NewMessageCard().
		Config(config).
		Header(header).
		Elements([]card.MessageCardElement{divElement1, divElement2, divElement3, divElement4, divElement5, divElement6}).
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendImageMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendShardChatMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendShardUserMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendAudioMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendMediaMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendFileMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendStickerMsg(client *client.Client) {
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
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())

}

func sendPostMsg(client *client.Client) {
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
		Build(), core.WithRequestId("jiaduo_id"))

	if err != nil {
		switch er := err.(type) {
		case *core.IllegalParamError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ClientTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ServerTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.DialFailedError:
			fmt.Println(er.Error()) // 处理错误
		default:
			//其他处理
			fmt.Println(err)
		}
		return
	}
	fmt.Println(core.Prettify(resp))
	fmt.Println(resp.RequestId())
}

func sendPostMsgUseBuilder(client *client.Client) {
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
		case *core.IllegalParamError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ClientTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ServerTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.DialFailedError:
			fmt.Println(er.Error()) // 处理错误
		default:
			//其他处理
			fmt.Println(err)
		}
		return
	}

	if resp.Success() {
		fmt.Println(core.Prettify(resp))
	} else {
		fmt.Println(resp.RequestId(), resp.Msg, resp.Code)
	}

}

func testCreate(client *client.Client) {
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
		case *core.IllegalParamError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ClientTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.ServerTimeoutError:
			fmt.Println(er.Error()) // 处理错误
		case *core.DialFailedError:
			fmt.Println(er.Error()) // 处理错误
		default:
			//其他处理
			fmt.Println(err)
		}
		return
	}

	if resp.Success() {
		fmt.Println(core.Prettify(resp))
	} else {
		fmt.Println(resp.RequestId(), resp.Msg, resp.Code)
	}
}

func sendRawReq(cli *client.Client) {
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
	}, core.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err, content)
		return
	}

	fmt.Println(resp)
}

func sendRawImageReq(cli *client.Client) {
	img, err := os.Open("/Users/bytedance/Downloads/go-icon.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer img.Close()

	formData := core.NewFormdata().AddField("image_type", "message").AddFile("image", img)

	resp, err := cli.Post(context.Background(), "/open-apis/im/v1/images", formData, core.AccessTokenTypeTenant)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp)
}
func main() {
	var appID, appSecret = os.Getenv("APP_ID"), os.Getenv("APP_SECRET")
	var feishu_client = client.NewClient(appID, appSecret, client.WithLogLevel(core.LogLevelDebug), client.WithLogReqRespInfoAtDebugLevel(false))

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
