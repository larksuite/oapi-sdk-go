package larkim

import (
	"context"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestListChatIteratorStopsWhenSinglePageResponseHasPageToken(t *testing.T) {
	pageToken := "stale-token"
	hasMore := false
	chatID := "chat-id"
	calls := 0

	req := NewListChatReqBuilder().Build()
	iterator := &ListChatIterator{
		ctx: context.Background(),
		req: req,
		listFunc: func(ctx context.Context, req *ListChatReq, options ...larkcore.RequestOptionFunc) (*ListChatResp, error) {
			calls++
			return &ListChatResp{
				CodeError: larkcore.CodeError{Code: 0},
				Data: &ListChatRespData{
					Items:     []*ListChat{{ChatId: &chatID}},
					PageToken: &pageToken,
					HasMore:   &hasMore,
				},
			}, nil
		},
	}

	hasNext, item, err := iterator.Next()
	if err != nil {
		t.Fatalf("first Next returned error: %v", err)
	}
	if !hasNext || item == nil || item.ChatId == nil || *item.ChatId != chatID {
		t.Fatalf("first Next returned unexpected item: hasNext=%v item=%#v", hasNext, item)
	}

	hasNext, item, err = iterator.Next()
	if err != nil {
		t.Fatalf("second Next returned error: %v", err)
	}
	if hasNext || item != nil {
		t.Fatalf("second Next should stop, got hasNext=%v item=%#v", hasNext, item)
	}
	if calls != 1 {
		t.Fatalf("iterator should not request another page, got calls=%d", calls)
	}
}

func TestListChatIteratorClearsPageTokenOnLastPage(t *testing.T) {
	firstToken := "page-2"
	firstHasMore := true
	lastHasMore := false
	firstChatID := "first-chat"
	lastChatID := "last-chat"
	calls := 0

	req := NewListChatReqBuilder().Build()
	iterator := &ListChatIterator{
		ctx: context.Background(),
		req: req,
		listFunc: func(ctx context.Context, req *ListChatReq, options ...larkcore.RequestOptionFunc) (*ListChatResp, error) {
			calls++
			switch calls {
			case 1:
				return &ListChatResp{
					CodeError: larkcore.CodeError{Code: 0},
					Data: &ListChatRespData{
						Items:     []*ListChat{{ChatId: &firstChatID}},
						PageToken: &firstToken,
						HasMore:   &firstHasMore,
					},
				}, nil
			case 2:
				return &ListChatResp{
					CodeError: larkcore.CodeError{Code: 0},
					Data: &ListChatRespData{
						Items:   []*ListChat{{ChatId: &lastChatID}},
						HasMore: &lastHasMore,
					},
				}, nil
			default:
				t.Fatalf("iterator requested an unexpected extra page, calls=%d", calls)
				return nil, nil
			}
		},
	}

	for _, wantChatID := range []string{firstChatID, lastChatID} {
		hasNext, item, err := iterator.Next()
		if err != nil {
			t.Fatalf("Next returned error: %v", err)
		}
		if !hasNext || item == nil || item.ChatId == nil || *item.ChatId != wantChatID {
			t.Fatalf("Next returned unexpected item: hasNext=%v item=%#v wantChatID=%s", hasNext, item, wantChatID)
		}
	}

	hasNext, item, err := iterator.Next()
	if err != nil {
		t.Fatalf("final Next returned error: %v", err)
	}
	if hasNext || item != nil {
		t.Fatalf("final Next should stop, got hasNext=%v item=%#v", hasNext, item)
	}
	if calls != 2 {
		t.Fatalf("iterator should request exactly two pages, got calls=%d", calls)
	}
}
