package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
)

func main() {
	printMessageText()
	printMessagePost()
	printMessageCard()
	printMessageCardI18n()
	printMessageImage()
	printMessageShareChat()
	printMessageShareUser()
	printMessageAudio()
	printMessageVideo()
	printMessageFile()
	printMessageSticker()
}

func printMessageText() {
	text := &lark.MessageText{Text: "Text"}
	s, err := text.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageText: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessagePost() {
	post := &lark.MessagePost{
		ZhCN: &lark.MessagePostContent{
			Title: "Title",
			Content: [][]lark.MessagePostElement{
				{
					&lark.MessagePostText{
						Text:     "Text",
						UnEscape: false,
					},
					&lark.MessagePostA{
						Text:     "Fei Shu",
						Href:     "https://www.feishu.com",
						UnEscape: false,
					},
				},
				{
					&lark.MessagePostA{
						Text:     "Fei Shu",
						Href:     "https://www.feishu.com",
						UnEscape: false,
					},
				},
			},
		},
	}
	s, err := post.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessagePost: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageCard() {
	card := &lark.MessageCard{
		CardLink: &lark.MessageCardURL{
			URL:        "https://www.feishu.cn",
			AndroidURL: "https://www.feishu.cn",
			IOSURL:     "https://www.feishu.cn",
			PCURL:      "https://www.feishu.cn",
		},
		Config: &lark.MessageCardConfig{EnableForward: lark.BoolPtr(true)},
		Header: &lark.MessageCardHeader{
			Template: lark.StringPtr("blue"),
			Title: &lark.MessageCardPlainText{
				Content: "Header title",
			},
		},
		Elements: []lark.MessageCardElement{
			&lark.MessageCardImage{
				Alt: &lark.MessageCardPlainText{
					Content: "img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg",
					Lines:   nil,
				},
				Title: &lark.MessageCardPlainText{
					Content: "img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg",
				},
				ImgKey:       "img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg",
				CustomWidth:  lark.IntPtr(300),
				CompactWidth: lark.BoolPtr(false),
			},
			&lark.MessageCardAction{
				Actions: []lark.MessageCardActionElement{
					&lark.MessageCardEmbedButton{
						Text: &lark.MessageCardPlainText{
							Content: "button",
						},
						Type:  lark.MessageCardButtonTypeDanger.Ptr(),
						Value: map[string]interface{}{"value": "1"},
						Confirm: &lark.MessageCardActionConfirm{
							Title: &lark.MessageCardPlainText{
								Content: "Title",
							},
							Text: &lark.MessageCardPlainText{
								Content: "Text",
							},
						},
					},
				},
				Layout: lark.MessageCardActionLayoutFlow.Ptr(),
			},
			&lark.MessageCardMarkdown{
				Content: "**Markdown**",
			},
			&lark.MessageCardDiv{
				Text: &lark.MessageCardPlainText{
					Content: "text",
				},
				Extra: &lark.MessageCardEmbedButton{
					Text: &lark.MessageCardPlainText{
						Content: "button",
					},
					Type:  lark.MessageCardButtonTypeDanger.Ptr(),
					Value: map[string]interface{}{"value": "1"},
					Confirm: &lark.MessageCardActionConfirm{
						Title: &lark.MessageCardPlainText{
							Content: "Title",
						},
						Text: &lark.MessageCardPlainText{
							Content: "Text",
						},
					},
				},
			},
		},
	}
	s, err := card.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageCard: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageCardI18n() {
	card := &lark.MessageCard{
		CardLink: &lark.MessageCardURL{
			URL:        "https://www.feishu.cn",
			AndroidURL: "https://www.feishu.cn",
			IOSURL:     "https://www.feishu.cn",
			PCURL:      "https://www.feishu.cn",
		},
		Config: &lark.MessageCardConfig{EnableForward: lark.BoolPtr(true)},
		Header: &lark.MessageCardHeader{
			Template: lark.StringPtr("blue"),
			Title: &lark.MessageCardPlainText{
				I18n: &lark.MessageCardPlainTextI18n{
					ZhCN: "ZhCN Header title",
					EnUS: "",
					JaJP: "JaJP Header title",
				},
			},
		},
		I18nElements: &lark.MessageCardI18nElements{
			ZhCN: []lark.MessageCardElement{
				&lark.MessageCardMarkdown{
					Content: "**ZhCN**",
				},
				&lark.MessageCardMarkdown{
					Content: "**ZhCN-2**",
				},
			},
			EnUS: []lark.MessageCardElement{
				&lark.MessageCardMarkdown{
					Content: "**EnUS**",
				},
				&lark.MessageCardMarkdown{
					Content: "**EnUS-2**",
				},
			},
			JaJP: []lark.MessageCardElement{
				&lark.MessageCardMarkdown{
					Content: "**JaJP**",
				},
				&lark.MessageCardMarkdown{
					Content: "**JaJP-2**",
				},
			},
		},
	}
	s, err := card.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func printMessageImage() {
	image := &lark.MessageImage{ImageKey: "img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg"}
	s, err := image.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageImage: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageShareChat() {
	shareChat := &lark.MessageShareChat{ChatId: "oc_xxxxxxxxxxxx"}
	s, err := shareChat.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageShareChat: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageShareUser() {
	shareUser := &lark.MessageShareUser{UserId: "ou-xxxxxxxxxxxxxx"}
	s, err := shareUser.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageShareUser: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageAudio() {
	audio := &lark.MessageAudio{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := audio.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageAudio: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageVideo() {
	audio := &lark.MessageVideo{
		FileKey:  "file-xxxxxxxxxx",
		ImageKey: "img-xxxxxxxxx",
	}
	s, err := audio.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageVideo: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageFile() {
	file := &lark.MessageFile{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := file.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageFile: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}

func printMessageSticker() {
	sticker := &lark.MessageSticker{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := sticker.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println("MessageSticker: ----------------------------")
	fmt.Println(s)
	fmt.Println("-----------------------------------------")
}
