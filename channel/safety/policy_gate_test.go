package safety

import (
	"github.com/larksuite/oapi-sdk-go/v3/channel/types"
	"testing"
)

func ptrBool(b bool) *bool {
	return &b
}

func TestPolicyGate(t *testing.T) {
	cfg := types.PolicyConfig{
		GroupAllowlist:      []string{"group1"},
		RequireMention:      ptrBool(true),
		RespondToMentionAll: ptrBool(false),
		DMMode:              "allowlist",
		DMAllowlist:         []string{"user1"},
	}

	pg := NewPolicyGate(&cfg, nil)

	// Test Group Allowed
	msg1 := &types.NormalizedMessage{
		ChatType:     "group",
		ChatID:       "group1",
		MentionedBot: true,
	}
	dec1 := pg.Evaluate(msg1)
	if !dec1.Allowed {
		t.Errorf("Expected group message to be allowed")
	}

	// Test Group Not Allowed
	msg2 := &types.NormalizedMessage{
		ChatType:     "group",
		ChatID:       "group2",
		MentionedBot: true,
	}
	dec2 := pg.Evaluate(msg2)
	if dec2.Allowed || dec2.Reason != types.RejectReasonGroupNotAllowed {
		t.Errorf("Expected group_not_allowed")
	}

	// Test Group No Mention
	msg3 := &types.NormalizedMessage{
		ChatType:     "group",
		ChatID:       "group1",
		MentionedBot: false,
	}
	dec3 := pg.Evaluate(msg3)
	if dec3.Allowed || dec3.Reason != types.RejectReasonNoMention {
		t.Errorf("Expected no_mention")
	}

	// Test Group Mention All Blocked
	msg4 := &types.NormalizedMessage{
		ChatType:     "group",
		ChatID:       "group1",
		MentionedBot: true,
		MentionAll:   true,
	}
	dec4 := pg.Evaluate(msg4)
	if dec4.Allowed || dec4.Reason != types.RejectReasonMentionAll {
		t.Errorf("Expected mention_all_blocked")
	}

	// Test DM Allowed
	msg5 := &types.NormalizedMessage{
		ChatType: "p2p",
		UserID:   "user1",
	}
	dec5 := pg.Evaluate(msg5)
	if !dec5.Allowed {
		t.Errorf("Expected DM to be allowed")
	}

	// Test DM Not Allowed
	msg6 := &types.NormalizedMessage{
		ChatType: "p2p",
		UserID:   "user2",
	}
	dec6 := pg.Evaluate(msg6)
	if dec6.Allowed || dec6.Reason != types.RejectReasonSenderNotAllowed {
		t.Errorf("Expected sender_not_allowed")
	}

	// Update Config
	pg.UpdateConfig(types.PolicyConfig{
		DMMode: "open",
	})
	dec7 := pg.Evaluate(msg6)
	if !dec7.Allowed {
		t.Errorf("Expected DM to be allowed after config update")
	}
}
