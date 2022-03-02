package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
)

func main() {
	printMessageCardWithAction()
	//printMessageCardWithI18n()
}

func printMessageCardWithAction() {
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
