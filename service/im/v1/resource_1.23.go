//go:build go1.23

package larkim

import (
	"context"
	"fmt"
	"iter"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func (c *chat) ListByIter(ctx context.Context, req *ListChatReq, options ...larkcore.RequestOptionFunc) (iter.Seq2[*ListChat, error], error) {
	firstResp, err := c.List(ctx, req, options...)
	if err != nil {
		return nil, err
	}
	if !firstResp.Success() {
		return nil, fmt.Errorf("Code:%d,Msg:%s", firstResp.Code, firstResp.Msg) // follow old iterator behavior
	}

	return func(yield func(*ListChat, error) bool) {
		resp := firstResp
		for {
			if resp.Data != nil {
				for _, item := range resp.Data.Items {
					if !yield(item, nil) {
						return
					}
				}
			}

			if resp.Data == nil || resp.Data.HasMore == nil || !*resp.Data.HasMore || resp.Data.PageToken == nil {
				return
			}
			req.apiReq.QueryParams.Set("page_token", *resp.Data.PageToken)

			nextResp, err := c.List(ctx, req, options...)
			if err != nil {
				yield(nil, err)
				return
			}
			if !nextResp.Success() {
				yield(nil, fmt.Errorf("Code:%d,Msg:%s", nextResp.Code, nextResp.Msg))
				return
			}
			resp = nextResp
		}
	}, nil
}
