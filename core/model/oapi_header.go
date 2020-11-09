package model

import "strings"

type OapiHeader struct {
	m map[string][]string
}

func NewOapiHeader(m map[string][]string) *OapiHeader {
	return &OapiHeader{m: m}
}

func (h OapiHeader) GetNames() []string {
	var names = make([]string, 0, len(h.m))
	for k := range h.m {
		names = append(names, k)
	}
	return names
}

func (h OapiHeader) GetFirstValues(name string) string {
	values := h.m[h.normalizeKey(name)]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func (h OapiHeader) GetMultiValues(name string) []string {
	return h.m[h.normalizeKey(name)]
}

func (h OapiHeader) normalizeKey(name string) string {
	return strings.ToLower(name)
}
