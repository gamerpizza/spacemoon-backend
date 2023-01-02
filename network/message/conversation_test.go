package message

import (
	"spacemoon/network/profile"
	"testing"
	"time"
)

func TestConversations(t *testing.T) {
	var m Messenger = NewMessenger(persistence)
	var messages []Message = m.GetConversationsBetween(u1).And(u2)
	for _, message := range messages {
		if author := message.Author(); author != Author(u1) && author != Author(u2) {
			t.Fatalf("unexpected author for %+v: %s", m, author)
		}
		if recipient := message.Recipient(); recipient != Recipient(u1) && recipient != Recipient(u2) {
			t.Fatalf("unexpected recipientfor %+v: %s", m, recipient)
		}
	}
	if len(messages) != expectedConversationMessages {
		t.Fatalf("expected messages: %d -- got: %d", expectedConversationMessages, len(messages))
	}
}

const u1 profile.Id = "profile1"
const u2 profile.Id = "profile2"
const u3 profile.Id = "profile3"

const expectedConversationMessages = 4

var persistence Persistence = &fakePersistence{}

func init() {
	t := time.Now().Add(-time.Hour)
	err := persistence.Save(New(Author(u1), Recipient(u2), "Hi!", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u2), Recipient(u1), "Hi! :)", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u2), Recipient(u1), "How are you doing?", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u1), Recipient(u2), "Great!", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u1), Recipient(u3), "Hi!", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u3), Recipient(u1), "Hi!", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u2), Recipient(u3), "Hi!", t))
	if err != nil {
		return
	}

	t = t.Add(time.Minute)
	err = persistence.Save(New(Author(u3), Recipient(u2), "Hi!", t))
	if err != nil {
		return
	}
}
