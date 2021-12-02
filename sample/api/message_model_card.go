package main

import (
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
)

func main() {
	printMessageCard1()
}

func printMessageCard1() {
	card := lark.MessageCard{
		Header: &lark.MessageCardHeader{
			Title:    &lark.MessageCardPlainText{Content: "k8s ticket"},
			Template: lark.StringPtr("indigo"),
		},
	}
	element := &lark.MessageCardAction{}
	cluster := &lark.MessageCardEmbedSelectMenuStatic{
		MessageCardEmbedSelectMenuBase: &lark.MessageCardEmbedSelectMenuBase{
			Placeholder: &lark.MessageCardPlainText{
				Content: "choose cluster",
			},
			InitialOption: "dev",
		},
	}
	element.Actions = append(element.Actions, cluster)
	card.Elements = append(card.Elements, element)
	content, err := card.JSON()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("card:%s\n", content)
}
