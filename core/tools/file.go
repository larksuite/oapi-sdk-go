package tools

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// download file return ReadCloser
//  defer r.Close()
//
func DownloadFileToStream(ctx context.Context, url string) (io.ReadCloser, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code:%d", resp.StatusCode)
	}
	return resp.Body, nil
}

func DownloadFile(ctx context.Context, url string) ([]byte, error) {
	r, err := DownloadFileToStream(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}
