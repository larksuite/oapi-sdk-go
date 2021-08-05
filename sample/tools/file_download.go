package main

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go"
	"io"
	"os"
)

func main() {
	f, err := os.Create("test_file_download.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	readCloser, err := lark.DownloadFileToStream(context.Background(),
		"https://sf1-ttcdn-tos.pstatp.com/obj/open-platform-opendoc/b7f456f542e8811e82657e577f126bce_WfFUj8sO1i.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer readCloser.Close()
	_, err = io.Copy(f, readCloser)
	if err != nil {
		fmt.Println(err)
		return
	}
}
