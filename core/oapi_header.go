package core

import "strings"

type OapiHeader struct {
	header map[string][]string
}

func NewOapiHeader(m map[string][]string) *OapiHeader {
	header := make(map[string][]string, len(m))
	for k, v := range m {
		header[normalizeKey(k)] = v
	}
	return &OapiHeader{header: header}
}

func (h OapiHeader) GetNames() []string {
	var names = make([]string, 0, len(h.header))
	for k := range h.header {
		names = append(names, k)
	}
	return names
}

func (h OapiHeader) GetFirstValue(name string) string {
	values := h.header[normalizeKey(name)]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func (h OapiHeader) GetMultiValues(name string) []string {
	return h.header[normalizeKey(name)]
}

func normalizeKey(name string) string {
	return strings.ToLower(name)
}
