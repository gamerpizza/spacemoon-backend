package message

import "testing"

func TestMessenger_GetConversationsWith(t *testing.T) {
	var m Messenger = NewMessenger(persistence, stubLoginPersistence{})
	conversations := m.GetConversationsWith(u1)
	if len(conversations[u2]) != 4 || len(conversations[u3]) != 2 {
		t.Fatalf("did not get the expected conversations: %+v", conversations)
	}
}
