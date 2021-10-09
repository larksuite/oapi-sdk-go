package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
)

func main() {
	printMessageText()
	printMessagePost()
	printMessageCard()
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
	fmt.Println(s)
}

func printMessagePost() {
	post := make(lark.MessagePost)
	post[lark.LanguageTypeZhCN] = &lark.MessagePostContent{
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
	}
	s, err := post.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
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
	fmt.Println(s)
}

func printMessageImage() {
	image := &lark.MessageImage{ImageKey: "img_v2_9221f258-db3e-4a40-b9cb-24decddee2bg"}
	s, err := image.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func printMessageShareChat() {
	shareChat := &lark.MessageShareChat{ChatId: "oc_xxxxxxxxxxxx"}
	s, err := shareChat.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func printMessageShareUser() {
	shareUser := &lark.MessageShareUser{UserId: "ou-xxxxxxxxxxxxxx"}
	s, err := shareUser.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func printMessageAudio() {
	audio := &lark.MessageAudio{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := audio.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
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
	fmt.Println(s)
}

func printMessageFile() {
	file := &lark.MessageFile{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := file.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}

func printMessageSticker() {
	sticker := &lark.MessageSticker{
		FileKey: "file-xxxxxxxxxx",
	}
	s, err := sticker.JSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
