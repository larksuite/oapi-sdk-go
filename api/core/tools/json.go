package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func MarshalJSON(value interface{}, forceSendFields []string) ([]byte, error) {
	if len(forceSendFields) == 0 {
		return json.Marshal(value)
	}
	mustInclude := make(map[string]bool)
	for _, f := range forceSendFields {
		mustInclude[f] = true
	}
	dataMap, err := structToMap(value, mustInclude)
	if err != nil {
		return nil, err
	}
	return json.Marshal(dataMap)
}

func structToMap(schema interface{}, mustInclude map[string]bool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	s := reflect.ValueOf(schema)
	st := s.Type()

	for i := 0; i < s.NumField(); i++ {
		jsonTag := st.Field(i).Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		tag, err := parseJSONTag(jsonTag)
		if err != nil {
			return nil, err
		}
		if tag.ignore {
			continue
		}
		v := s.Field(i)
		f := st.Field(i)
		if !includeField(v, f, mustInclude) {
			continue
		}

		// nil maps are treated as empty maps.
		if f.Type.Kind() == reflect.Map && v.IsNil() {
			m[tag.apiName] = map[string]string{}
			continue
		}
		if f.Type.Kind() == reflect.Slice && v.IsNil() {
			m[tag.apiName] = []string{}
			continue
		}
		if tag.stringFormat {
			m[tag.apiName] = formatAsString(v, f.Type.Kind())
		} else {
			m[tag.apiName] = v.Interface()
		}
	}
	return m, nil
}

func formatAsString(v reflect.Value, kind reflect.Kind) string {
	if kind == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return fmt.Sprintf("%v", v.Interface())
}

type jsonTag struct {
	apiName      string
	stringFormat bool
	ignore       bool
}

func parseJSONTag(val string) (jsonTag, error) {
	if val == "-" {
		return jsonTag{ignore: true}, nil
	}
	var tag jsonTag
	i := strings.Index(val, ",")
	if i == -1 || val[:i] == "" {
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	tag = jsonTag{
		apiName: val[:i],
	}
	switch val[i+1:] {
	case "omitempty":
	case "omitempty,string":
		tag.stringFormat = true
	default:
		return tag, fmt.Errorf("malformed json tag: %s", val)
	}
	return tag, nil
}

func includeField(v reflect.Value, f reflect.StructField, mustInclude map[string]bool) bool {
	if f.Type.Kind() == reflect.Ptr && v.IsNil() {
		return false
	}
	if f.Type.Kind() == reflect.Interface && v.IsNil() {
		return false
	}
	return mustInclude[f.Name] || !isEmptyValue(v)
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
