package larkim

import (
	"testing"
)

func TestGetMessageReqBuilder_CardMsgContentType(t *testing.T) {
	builder := NewGetMessageReqBuilder()
	builder.CardMsgContentType("user_card_content")
	req := builder.Build()

	if req.apiReq.QueryParams.Get("card_msg_content_type") != "user_card_content" {
		t.Errorf("expected card_msg_content_type to be user_card_content, got %s", req.apiReq.QueryParams.Get("card_msg_content_type"))
	}
}
