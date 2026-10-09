package larkcore

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"testing"
)

func TestUploadFilenameOverride(t *testing.T) {
	for _, name := range []string{"report.txt", "中文报告.txt"} {
		t.Run(name, func(t *testing.T) {
			body := struct {
				File     io.Reader `json:"file,omitempty" filename:"Filename"`
				Filename *string   `json:"filename,omitempty"`
			}{File: bytes.NewBufferString("contents"), Filename: &name}
			contentType, content, err := toFormdata(&body).content()
			if err != nil {
				t.Fatal(err)
			}
			_, params, err := mime.ParseMediaType(contentType)
			if err != nil {
				t.Fatal(err)
			}
			reader := multipart.NewReader(bytes.NewReader(content), params["boundary"])
			parts := map[string]string{}
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				data, err := io.ReadAll(part)
				if err != nil {
					t.Fatal(err)
				}
				parts[part.FormName()] = string(data)
				if part.FormName() == "file" && part.FileName() != name {
					t.Fatalf("filename = %q", part.FileName())
				}
				if err := part.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if parts["file"] != "contents" || parts["filename"] != name {
				t.Fatalf("multipart fields mismatch")
			}
		})
	}
}

func TestUploadFilenameRejectsPathsAndControlCharacters(t *testing.T) {
	for _, name := range []string{"../report.txt", "sub/report.txt", "sub\\report.txt", "bad\r\nname", "bad\x00name"} {
		body := struct {
			File     io.Reader `json:"file,omitempty" filename:"Filename"`
			Filename *string   `json:"filename,omitempty"`
		}{File: bytes.NewBufferString("contents"), Filename: &name}
		if _, _, err := toFormdata(&body).content(); err == nil {
			t.Fatal("expected unsafe upload filename rejection")
		}
	}
}

func TestUnmappedFilenameDoesNotOverrideFile(t *testing.T) {
	name := "business-name.txt"
	body := struct {
		File     io.Reader `json:"file,omitempty"`
		Filename *string   `json:"filename,omitempty"`
	}{File: bytes.NewBufferString("contents"), Filename: &name}
	fields := toFormdata(&body).fields
	if formdataFilename(fields["file"]) != "unknown-file" {
		t.Fatal("unmapped filename changed file name")
	}
}

func TestUploadFilenameUnsetKeepsReaderName(t *testing.T) {
	empty := ""
	for _, name := range []*string{nil, &empty} {
		body := struct {
			File     io.Reader `json:"file,omitempty" filename:"Filename"`
			Filename *string   `json:"filename,omitempty"`
		}{File: &testNamedReader{Reader: bytes.NewReader([]byte("data")), name: "original.txt"}, Filename: name}
		if got := formdataFilename(toFormdata(&body).fields["file"]); got != "original.txt" {
			t.Fatalf("filename = %q", got)
		}
	}
}
