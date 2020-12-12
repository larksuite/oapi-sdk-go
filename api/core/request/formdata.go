package request

import (
	"bytes"
	"fmt"
	"io"
	"net/textproto"
	"strings"
)

type FormData struct {
	params map[string]interface{}
	files  []*File
}

func (fd *FormData) Params() map[string]interface{} {
	return fd.params
}

func (fd *FormData) Files() []*File {
	return fd.files
}

func NewFormData() *FormData {
	return &FormData{
		params: map[string]interface{}{},
	}
}

func (fd *FormData) AddParam(field string, val interface{}) *FormData {
	fd.params[field] = val
	return fd
}

func (fd *FormData) AddFile(field string, file *File) *FormData {
	file.fieldName = field
	fd.files = append(fd.files, file)
	return fd
}

func (fd *FormData) HasStream() bool {
	for _, file := range fd.files {
		isStream := file.IsStream()
		if isStream {
			return isStream
		}
	}
	return false
}

type File struct {
	fieldName     string
	name          string
	typ           string
	content       *bytes.Buffer
	contentStream io.Reader
	isStream      bool
}

func NewFile() *File {
	return &File{}
}

func (f *File) ContentStream() io.Reader {
	return f.contentStream
}

func (f *File) SetContentStream(reader io.Reader) *File {
	f.contentStream = reader
	f.isStream = true
	switch f.contentStream.(type) {
	case *bytes.Buffer:
		f.isStream = false
	}
	return f
}

func (f *File) SetContent(content []byte) *File {
	f.content = bytes.NewBuffer(content)
	return f
}

func (f *File) Type() string {
	return f.typ
}

func (f *File) SetType(typ string) *File {
	f.typ = typ
	return f
}

func (f *File) Name() string {
	return f.name
}

func (f *File) IsStream() bool {
	return f.isStream
}

func (f *File) SetName(name string) *File {
	f.name = name
	return f
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.contentStream != nil {
		return f.contentStream.Read(p)
	}
	return f.content.Read(p)
}

func (f *File) MIMEHeader() textproto.MIMEHeader {
	fieldName := "file"
	if f.fieldName != "" {
		fieldName = f.fieldName
	}
	name := "unknown"
	if f.name != "" {
		name = f.name
	}
	typ := "application/octet-stream"
	if f.typ != "" {
		typ = f.typ
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldName), escapeQuotes(name)))
	h.Set("Content-Type", typ)
	return h
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
