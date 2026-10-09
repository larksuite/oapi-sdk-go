package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	core "github.com/larksuite/oapi-sdk-go/v3/core"
	agent "github.com/larksuite/oapi-sdk-go/v3/service/agent/v1"
)

func (e *example) upload(ctx context.Context) (result error) {
	var reader io.Reader = bytes.NewReader([]byte("SDK Agent upload example\n"))
	if path := os.Getenv("AGENT_FILE_PATH"); path != "" {
		for _, part := range strings.Split(filepath.ToSlash(path), "/") {
			if part == ".." {
				return fmt.Errorf("AGENT_FILE_PATH must not contain parent traversal")
			}
		}
		file, err := os.Open(filepath.Clean(path))
		if err != nil {
			return fmt.Errorf("open upload file: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				if result == nil {
					result = fmt.Errorf("close upload file: %w", err)
				} else {
					log.Printf("file cleanup failed: %T", err)
				}
			}
		}()
		reader = file
	}
	filename := os.Getenv("AGENT_FILENAME")
	if filename == "" {
		filename = "report.txt"
	}
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	body := agent.NewCreateFileReqBodyBuilder().
		File(reader).
		Filename(filename).
		Uuid(uuid).
		Build()
	request := agent.NewCreateFileReqBuilder().
		Body(body).
		Build()
	response, err := e.client.Agent.V1.File.Create(
		ctx,
		request,
		core.WithUserAccessToken(e.token),
	)
	if err != nil {
		return fmt.Errorf("upload: %w", err)
	}
	if !response.Success() {
		return fmt.Errorf("upload: HTTP %d, code %d", response.StatusCode, response.Code)
	}
	if response.Data == nil || response.Data.Uri == nil {
		return fmt.Errorf("upload: missing URI")
	}
	log.Printf("AGENT_FILE_URI=%s", *response.Data.Uri)
	log.Print("upload succeeded; set AGENT_FILE_URI and run chat to reference it")
	return nil
}
