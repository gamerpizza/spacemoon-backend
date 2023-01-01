package message

import "spacemoon/network/profile"

func NewMessenger(p Persistence) Messenger {
	return messenger{persistence: p}
}

type Messenger interface {
	Get(p profile.Id) UserMessages
	Send(Message) Sender
}

type messenger struct {
	persistence Persistence
}

func (m messenger) Get(p profile.Id) UserMessages {
	return m.persistence.GetMessagesFor(Recipient(p))
}

func (m messenger) Send(msg Message) Sender {
	return &messageSender{persistence: m.persistence, message: msg}
}
