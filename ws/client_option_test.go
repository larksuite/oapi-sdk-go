package ws

import "testing"

func TestWithReconnectCount_SetsValue(t *testing.T) {
	cli := NewClient("appID", "appSecret", WithReconnectCount(3))
	if cli.reconnectCount != 3 {
		t.Fatalf("expected reconnectCount=3, got %d", cli.reconnectCount)
	}
}

func TestNewClient_DefaultReconnectCount(t *testing.T) {
	cli := NewClient("appID", "appSecret")
	if cli.reconnectCount != -1 {
		t.Fatalf("expected default reconnectCount=-1, got %d", cli.reconnectCount)
	}
}
