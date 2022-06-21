package core

import (
	"context"
	"testing"
)

func TestEncryptedEventMsg(t *testing.T) {
	en, err := EncryptedEventMsg(context.Background(), map[string]string{"key1": "value1", "key2": "value2"}, "encrypteKey")
	if err != nil {
		t.Errorf("TestEncryptedEventMsg failed ,%v", err)
	}

	if en == "" {
		t.Errorf("TestEncryptedEventMsg failed ,%v", err)
	}
}
