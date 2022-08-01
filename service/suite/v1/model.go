// Code generated by lark suite oapi sdk gen
package v1

import (
	"github.com/larksuite/oapi-sdk-go/api/core/tools"
)

type DocsEntity struct {
	DocsToken       string   `json:"docs_token,omitempty"`
	DocsType        string   `json:"docs_type,omitempty"`
	Title           string   `json:"title,omitempty"`
	DocsOwner       string   `json:"docs_owner,omitempty"`
	OwnerId         string   `json:"owner_id,omitempty"`
	ForceSendFields []string `json:"-"`
}

func (s *DocsEntity) MarshalJSON() ([]byte, error) {
	type cp DocsEntity
	raw := cp(*s)
	return tools.MarshalJSON(raw, s.ForceSendFields)
}

type DocsMeta struct {
	DocsToken        string   `json:"docs_token,omitempty"`
	DocsType         string   `json:"docs_type,omitempty"`
	Title            string   `json:"title,omitempty"`
	OwnerId          string   `json:"owner_id,omitempty"`
	CreateTime       int      `json:"create_time,omitempty"`
	LatestModifyUser string   `json:"latest_modify_user,omitempty"`
	LatestModifyTime int      `json:"latest_modify_time,omitempty"`
	ForceSendFields  []string `json:"-"`
}

func (s *DocsMeta) MarshalJSON() ([]byte, error) {
	type cp DocsMeta
	raw := cp(*s)
	return tools.MarshalJSON(raw, s.ForceSendFields)
}

type RequestDoc struct {
	DocsToken       string   `json:"docs_token,omitempty"`
	DocsType        string   `json:"docs_type,omitempty"`
	ForceSendFields []string `json:"-"`
}

func (s *RequestDoc) MarshalJSON() ([]byte, error) {
	type cp RequestDoc
	raw := cp(*s)
	return tools.MarshalJSON(raw, s.ForceSendFields)
}

type DocsApiMetaReqBody struct {
	RequestDocs     []*RequestDoc `json:"request_docs,omitempty"`
	ForceSendFields []string      `json:"-"`
}

func (s *DocsApiMetaReqBody) MarshalJSON() ([]byte, error) {
	type cp DocsApiMetaReqBody
	raw := cp(*s)
	return tools.MarshalJSON(raw, s.ForceSendFields)
}

type DocsApiMetaResult struct {
	DocsMetas []*DocsMeta `json:"docs_metas,omitempty"`
}

type DocsApiSearchReqBody struct {
	SearchKey       string   `json:"search_key,omitempty"`
	Count           int      `json:"count,omitempty"`
	Offset          int      `json:"offset,omitempty"`
	OwnerIds        []string `json:"owner_ids,omitempty"`
	ChatIds         []string `json:"chat_ids,omitempty"`
	DocsTypes       []string `json:"docs_types,omitempty"`
	ForceSendFields []string `json:"-"`
}

func (s *DocsApiSearchReqBody) MarshalJSON() ([]byte, error) {
	type cp DocsApiSearchReqBody
	raw := cp(*s)
	return tools.MarshalJSON(raw, s.ForceSendFields)
}

type DocsApiSearchResult struct {
	DocsEntities []*DocsEntity `json:"docs_entities,omitempty"`
	HasMore      bool          `json:"has_more,omitempty"`
	Total        int           `json:"total,omitempty"`
}
