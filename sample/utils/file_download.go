package main

import (
	"context"
	"fmt"
	"github.com/larksuite/oapi-sdk-go/v2"
	"io/ioutil"
)

func main() {
	bs, err := lark.DownloadFile(context.Background(), "https://sf1-ttcdn-tos.pstatp.com/obj/open-platform-opendoc/b7f456f542e8811e82657e577f126bce_WfFUj8sO1i.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ioutil.WriteFile("test_file_download.png", bs, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
}
