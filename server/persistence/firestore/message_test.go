package firestore

import (
	"context"
	"spacemoon/network/message"
	"testing"
	"time"
)

func TestFireStorePersistence_GetMessagesFor(t *testing.T) {
	var p message.Persistence
	var err error
	p, err = GetPersistence(context.TODO())
	if err != nil {
		t.Fatal(err.Error())
	}
	var testMessage message.Message = message.New("author", "recipient", "Hello", time.Now())
	err = p.Save(testMessage)
	if err != nil {
		t.Fatal(err.Error())
	}
	messages := p.GetMessagesBy("author")
	found := false
	for _, m := range messages["recipient"] {
		if m.Content == testMessage.Content {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("message not found")
	}
}
