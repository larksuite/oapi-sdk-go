package lark

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"reflect"
	"strings"
)

import "time"

// StringPtr returns a pointer to the string value passed in.
func StringPtr(v string) *string {
	return &v
}

// StringValue returns the value of the string pointer passed in or
// "" if the pointer is nil.
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

// BoolPtr returns a pointer to the bool value passed in.
func BoolPtr(v bool) *bool {
	return &v
}

// BoolValue returns the value of the bool pointer passed in or
// false if the pointer is nil.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// IntPtr returns a pointer to the int value passed in.
func IntPtr(v int) *int {
	return &v
}

// IntValue returns the value of the int pointer passed in or
// 0 if the pointer is nil.
func IntValue(v *int) int {
	if v != nil {
		return *v
	}
	return 0
}

// Int8Ptr returns a pointer to the int8 value passed in.
func Int8Ptr(v int8) *int8 {
	return &v
}

// Int8Value returns the value of the int8 pointer passed in or
// 0 if the pointer is nil.
func Int8Value(v *int8) int8 {
	if v != nil {
		return *v
	}
	return 0
}

// Int16Ptr returns a pointer to the int16 value passed in.
func Int16Ptr(v int16) *int16 {
	return &v
}

// Int16Value returns the value of the int16 pointer passed in or
// 0 if the pointer is nil.
func Int16Value(v *int16) int16 {
	if v != nil {
		return *v
	}
	return 0
}

// Int32Ptr returns a pointer to the int32 value passed in.
func Int32Ptr(v int32) *int32 {
	return &v
}

// Int32Value returns the value of the int32 pointer passed in or
// 0 if the pointer is nil.
func Int32Value(v *int32) int32 {
	if v != nil {
		return *v
	}
	return 0
}

// Int64Ptr returns a pointer to the int64 value passed in.
func Int64Ptr(v int64) *int64 {
	return &v
}

// Int64Value returns the value of the int64 pointer passed in or
// 0 if the pointer is nil.
func Int64Value(v *int64) int64 {
	if v != nil {
		return *v
	}
	return 0
}

// Float32Ptr returns a pointer to the float32 value passed in.
func Float32Ptr(v float32) *float32 {
	return &v
}

// Float32Value returns the value of the float32 pointer passed in or
// 0 if the pointer is nil.
func Float32Value(v *float32) float32 {
	if v != nil {
		return *v
	}
	return 0
}

// Float64Ptr returns a pointer to the float64 value passed in.
func Float64Ptr(v float64) *float64 {
	return &v
}

// Float64Value returns the value of the float64 pointer passed in or
// 0 if the pointer is nil.
func Float64Value(v *float64) float64 {
	if v != nil {
		return *v
	}
	return 0
}

// TimePtr returns a pointer to the time.Time value passed in.
func TimePtr(v time.Time) *time.Time {
	return &v
}

// TimeValue returns the value of the time.Time pointer passed in or
// time.Time{} if the pointer is nil.
func TimeValue(v *time.Time) time.Time {
	if v != nil {
		return *v
	}
	return time.Time{}
}

// Prettify returns the string representation of a value.
func Prettify(i interface{}) string {
	var buf bytes.Buffer
	prettify(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// DownloadFile returns the url of resource
func DownloadFile(ctx context.Context, url string) ([]byte, error) {
	r, err := downloadFileToStream(ctx, url)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func FileNameByHeader(header http.Header) string {
	filename := ""
	_, media, _ := mime.ParseMediaType(header.Get("Content-Disposition"))
	if len(media) > 0 {
		filename = media["filename"]
	}
	return filename
}

type DecryptErr struct {
	Message string
}

func newDecryptErr(message string) *DecryptErr {
	return &DecryptErr{Message: message}
}

func (e DecryptErr) Error() string {
	return e.Message
}

func downloadFileToStream(ctx context.Context, url string) (io.ReadCloser, error) {
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

// prettify will recursively walk value v to build a textual
// representation of the value.
func prettify(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		strType := v.Type().String()
		if strType == "time.Time" {
			fmt.Fprintf(buf, "%s", v.Interface())
			break
		} else if strings.HasPrefix(strType, "io.") {
			buf.WriteString("<buffer>")
			break
		}

		buf.WriteString("{\n")

		var names []string
		for i := 0; i < v.Type().NumField(); i++ {
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if name[0:1] == strings.ToLower(name[0:1]) {
				continue // ignore unexported fields
			}
			if (f.Kind() == reflect.Ptr || f.Kind() == reflect.Slice || f.Kind() == reflect.Map) && f.IsNil() {
				continue // ignore unset fields
			}
			names = append(names, name)
		}

		for i, n := range names {
			val := v.FieldByName(n)
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(n + ": ")
			prettify(val, indent+2, buf)

			if i < len(names)-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	case reflect.Slice:
		strType := v.Type().String()
		if strType == "[]uint8" {
			fmt.Fprintf(buf, "<binary> len %d", v.Len())
			break
		}

		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			prettify(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}

		buf.WriteString(nl + id + "]")
	case reflect.Map:
		buf.WriteString("{\n")

		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			prettify(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}

		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")
	default:
		if !v.IsValid() {
			fmt.Fprint(buf, "<invalid value>")
			return
		}
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		case io.ReadSeeker, io.Reader:
			format = "buffer(%p)"
		}
		fmt.Fprintf(buf, format, v.Interface())
	}
}
