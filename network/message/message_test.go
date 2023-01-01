package message

import (
	"errors"
	"spacemoon/network/profile"
	"testing"
	"time"
)

func TestMessenger(t *testing.T) {
	var m Messenger = NewMessenger(&fakePersistence{})
	var p profile.Id
	var _ UserMessages = m.Get(p)
}

func TestMessage(t *testing.T) {
	var m Message
	var _ string = m.String()
	var _ time.Time = m.Time()
	var _ Author = m.Author()
	var _ Recipient = m.Recipient()
}

func TestMessenger_SendMessage(t *testing.T) {
	var messagePersistence Persistence = &fakePersistence{}
	var msgr Messenger = NewMessenger(messagePersistence)
	var msg Message
	var sndr profile.Id = "sender"
	var rcvr profile.Id = "receiver"
	var _ error = msgr.Send(msg).From(sndr).To(rcvr).Now()
	var receiverMessages UserMessages = msgr.Get(rcvr)
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
		t.Fatal("Recipient should not be empty")
	}
}

func TestSender_Now_FromShouldNotBeEmpty(t *testing.T) {
	var messagePersistence Persistence = &fakePersistence{}
	var msgr Messenger = NewMessenger(messagePersistence)
	var receiver profile.Id = "messageSender"
	var msg Message
	var err error = msgr.Send(msg).To(receiver).Now()
	if err == nil || !errors.Is(err, AuthorNotSetError) {
		t.Fatal("Author should not be empty")
	}
}

type fakePersistence struct {
	messages Messages
}

func (f *fakePersistence) GetMessagesFor(r Recipient) UserMessages {
	return f.messages[r]
}

func (f *fakePersistence) Save(m Message) error {
	if f.messages == nil {
		f.messages = make(Messages)
	}
	if f.messages[m.Recipient()] == nil {
		f.messages[m.Recipient()] = make(UserMessages)
	}
	f.messages[m.Recipient()][m.Author()] = append(f.messages[m.Recipient()][m.Author()], m)
	return nil
}
