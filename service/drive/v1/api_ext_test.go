/*
 * MIT License
 *
 * Copyright (c) 2022 Lark Technologies Pte. Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package larkdrive

import (
	"context"
	"testing"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
)

func TestListFileIteratorStopsAfterEmptyNextPageToken(t *testing.T) {
	emptyToken := ""
	hasMore := false
	requestCount := 0
	iterator := &ListFileIterator{
		ctx: context.Background(),
		req: NewListFileReqBuilder().Build(),
		listFunc: func(context.Context, *ListFileReq, ...larkcore.RequestOptionFunc) (*ListFileResp, error) {
			requestCount++
			return &ListFileResp{
				Data: &ListFileRespData{
					Files:         []*File{{}},
					NextPageToken: &emptyToken,
					HasMore:       &hasMore,
				},
			}, nil
		},
	}

	hasNext, file, err := iterator.Next()
	if err != nil {
		t.Fatalf("first Next() returned error: %v", err)
	}
	if !hasNext || file == nil {
		t.Fatalf("first Next() = (%v, %v), want an item", hasNext, file)
	}

	hasNext, file, err = iterator.Next()
	if err != nil {
		t.Fatalf("second Next() returned error: %v", err)
	}
	if hasNext || file != nil {
		t.Fatalf("second Next() = (%v, %v), want end of iteration", hasNext, file)
	}
	if requestCount != 1 {
		t.Fatalf("list request count = %d, want 1", requestCount)
	}
}

func TestListFileIteratorContinuesWithNonEmptyNextPageToken(t *testing.T) {
	nextToken := "page-2"
	requestCount := 0
	iterator := &ListFileIterator{
		ctx: context.Background(),
		req: NewListFileReqBuilder().Build(),
		listFunc: func(_ context.Context, req *ListFileReq, _ ...larkcore.RequestOptionFunc) (*ListFileResp, error) {
			requestCount++
			switch requestCount {
			case 1:
				return &ListFileResp{
					Data: &ListFileRespData{
						Files:         []*File{{}},
						NextPageToken: &nextToken,
					},
				}, nil
			case 2:
				if got := req.apiReq.QueryParams.Get("page_token"); got != nextToken {
					t.Fatalf("page_token = %q, want %q", got, nextToken)
				}
				return &ListFileResp{
					Data: &ListFileRespData{Files: []*File{{}}},
				}, nil
			default:
				t.Fatalf("unexpected list request %d", requestCount)
				return nil, nil
			}
		},
	}

	for call := 1; call <= 2; call++ {
		hasNext, file, err := iterator.Next()
		if err != nil {
			t.Fatalf("Next() call %d returned error: %v", call, err)
		}
		if !hasNext || file == nil {
			t.Fatalf("Next() call %d = (%v, %v), want an item", call, hasNext, file)
		}
	}

	hasNext, file, err := iterator.Next()
	if err != nil {
		t.Fatalf("final Next() returned error: %v", err)
	}
	if hasNext || file != nil {
		t.Fatalf("final Next() = (%v, %v), want end of iteration", hasNext, file)
	}
	if requestCount != 2 {
		t.Fatalf("list request count = %d, want 2", requestCount)
	}
}
