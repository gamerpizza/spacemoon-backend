package message

import (
	"errors"
	"spacemoon/network/profile"
	"testing"
	"time"
)

func TestMessenger(t *testing.T) {
	var m Messenger = NewMessenger(getFakePersistence())
	var p profile.Id
	var _ ReceivedUserMessages = m.GetAllMessagesFor(Recipient(p))
}

func getFakePersistence() *fakePersistence {
	return &fakePersistence{ConversationReader: NewConversationReader(&fakePersistence{})}
}

func TestMessage(t *testing.T) {
	var m Message
	var _ string = m.String()
	var _ time.Time = m.GetPostingTime()
	var _ Author = m.GetAuthor()
	var _ Recipient = m.GetRecipient()
}

func TestMessenger_SendMessage(t *testing.T) {
	var messagePersistence Persistence = &fakePersistence{}
	var msgr Messenger = NewMessenger(messagePersistence)
	var msg Message
	var sndr profile.Id = "sender"
	var rcvr profile.Id = "receiver"
	var _ error = msgr.Send(msg).From(sndr).To(rcvr).Now()
	var receiverMessages ReceivedUserMessages = msgr.GetAllMessagesFor(Recipient(rcvr))
	_, exist := receiverMessages[Author(sndr)]
	if !exist {
		t.Fatalf("no messages found for %s", sndr)
	}
}

func TestSender_Now_ToShouldNotBeEmpty(t *testing.T) {
	var messagePersistence Persistence = &fakePersistence{}
	var msgr Messenger = NewMessenger(messagePersistence)
	var sender profile.Id = "messageSender"
	var msg Message
	var err error = msgr.Send(msg).From(sender).Now()
	if err == nil || !errors.Is(err, RecipientNotSetError) {
		t.Fatal("GetRecipient should not be empty")
	}
}

func TestSender_Now_FromShouldNotBeEmpty(t *testing.T) {
	var messagePersistence Persistence = &fakePersistence{}
	var msgr Messenger = NewMessenger(messagePersistence)
	var receiver profile.Id = "messageSender"
	var msg Message
	var err error = msgr.Send(msg).To(receiver).Now()
	if err == nil || !errors.Is(err, AuthorNotSetError) {
		t.Fatal("GetAuthor should not be empty")
	}
}

type fakePersistence struct {
	messagesByRecipient map[Recipient]ReceivedUserMessages
	ConversationReader
}

// No need for a double persistence, a nested map is rather easy to check for now
func (f *fakePersistence) GetMessagesBy(a Author) SentUserMessages {
	sentMessages := make(SentUserMessages)
	for recipient, messages := range f.messagesByRecipient {
		if conversation, exists := messages[a]; exists {
			sentMessages[recipient] = conversation
		}
	}
	return sentMessages
}

func (f *fakePersistence) GetMessagesFor(id Recipient) ReceivedUserMessages {
	return f.messagesByRecipient[id]
}

func (f *fakePersistence) Save(m Message) error {
	if f.messagesByRecipient == nil {
		f.messagesByRecipient = make(map[Recipient]ReceivedUserMessages)
	}
	if f.messagesByRecipient[m.Recipient] == nil {
		f.messagesByRecipient[m.Recipient] = make(ReceivedUserMessages)
	}
	f.messagesByRecipient[m.Recipient][m.Author] = append(f.messagesByRecipient[m.Recipient][m.Author], m)
	return nil
}
