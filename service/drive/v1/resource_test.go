package larkdrive

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

type streamDownloadBody struct {
	reader *strings.Reader
	reads  int32
	closed int32
}

func (b *streamDownloadBody) Read(p []byte) (int, error) {
	atomic.AddInt32(&b.reads, 1)
	return b.reader.Read(p)
}

func (b *streamDownloadBody) Close() error {
	atomic.StoreInt32(&b.closed, 1)
	return nil
}

type streamDownloadClient struct {
	t            *testing.T
	body         *streamDownloadBody
	expectedPath string
	statusCode   int
}

func (c *streamDownloadClient) Do(req *http.Request) (*http.Response, error) {
	c.t.Helper()
	if req.URL.Path != c.expectedPath {
		c.t.Fatalf("unexpected request path: %s", req.URL.Path)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer user-token" {
		c.t.Fatalf("unexpected authorization header: %s", got)
	}
	return &http.Response{
		StatusCode: c.statusCode,
		Header: http.Header{
			"Content-Type":        []string{"application/octet-stream"},
			"Content-Disposition": []string{`attachment; filename="large.bin"`},
		},
		Body:    c.body,
		Request: req,
	}, nil
}

func TestDownloadStreamDoesNotBufferResponse(t *testing.T) {
	tests := []struct {
		name         string
		expectedPath string
		statusCode   int
		download     func(context.Context, *V1) (io.Reader, *larkcore.ApiResp, string, error)
	}{
		{
			name:         "export task",
			expectedPath: "/open-apis/drive/v1/export_tasks/file/file-token/download",
			statusCode:   http.StatusOK,
			download: func(ctx context.Context, client *V1) (io.Reader, *larkcore.ApiResp, string, error) {
				resp, err := client.ExportTask.DownloadStream(
					ctx,
					NewDownloadExportTaskReqBuilder().FileToken("file-token").Build(),
					larkcore.WithUserAccessToken("user-token"),
				)
				if err != nil {
					return nil, nil, "", err
				}
				return resp.File, resp.ApiResp, resp.FileName, nil
			},
		},
		{
			name:         "file",
			expectedPath: "/open-apis/drive/v1/files/file-token/download",
			statusCode:   http.StatusOK,
			download: func(ctx context.Context, client *V1) (io.Reader, *larkcore.ApiResp, string, error) {
				resp, err := client.File.DownloadStream(
					ctx,
					NewDownloadFileReqBuilder().FileToken("file-token").Build(),
					larkcore.WithUserAccessToken("user-token"),
				)
				if err != nil {
					return nil, nil, "", err
				}
				return resp.File, resp.ApiResp, resp.FileName, nil
			},
		},
		{
			name:         "media",
			expectedPath: "/open-apis/drive/v1/medias/file-token/download",
			statusCode:   http.StatusOK,
			download: func(ctx context.Context, client *V1) (io.Reader, *larkcore.ApiResp, string, error) {
				resp, err := client.Media.DownloadStream(
					ctx,
					NewDownloadMediaReqBuilder().FileToken("file-token").Build(),
					larkcore.WithUserAccessToken("user-token"),
				)
				if err != nil {
					return nil, nil, "", err
				}
				return resp.File, resp.ApiResp, resp.FileName, nil
			},
		},
		{
			name:         "media range",
			expectedPath: "/open-apis/drive/v1/medias/file-token/download",
			statusCode:   http.StatusPartialContent,
			download: func(ctx context.Context, client *V1) (io.Reader, *larkcore.ApiResp, string, error) {
				resp, err := client.Media.DownloadStream(
					ctx,
					NewDownloadMediaReqBuilder().FileToken("file-token").Build(),
					larkcore.WithUserAccessToken("user-token"),
				)
				if err != nil {
					return nil, nil, "", err
				}
				return resp.File, resp.ApiResp, resp.FileName, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &streamDownloadBody{reader: strings.NewReader("large-file")}
			config := &larkcore.Config{
				BaseUrl:      "https://open.feishu.cn",
				AppId:        "app-id",
				AppSecret:    "app-secret",
				HttpClient:   &streamDownloadClient{t: t, body: body, expectedPath: tt.expectedPath, statusCode: tt.statusCode},
				Logger:       larkcore.NewDefaultLogger(larkcore.LogLevelError),
				Serializable: &larkcore.DefaultSerialization{},
			}

			file, apiResp, fileName, err := tt.download(context.Background(), New(config))
			if err != nil {
				t.Fatalf("download stream failed: %v", err)
			}
			if fileName != "large.bin" {
				t.Fatalf("unexpected file name: %s", fileName)
			}
			if atomic.LoadInt32(&body.reads) != 0 {
				t.Fatalf("expected download body unread before caller consumes it")
			}

			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("read streamed file failed: %v", err)
			}
			if string(data) != "large-file" {
				t.Fatalf("unexpected streamed file body: %s", data)
			}
			if err := apiResp.Body.Close(); err != nil {
				t.Fatalf("close streamed response failed: %v", err)
			}
			if atomic.LoadInt32(&body.closed) != 1 {
				t.Fatalf("expected streamed response body closed")
			}
		})
	}
}
