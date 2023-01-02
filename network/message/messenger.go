package message

import "spacemoon/network/profile"

func NewMessenger(p Persistence) Messenger {
	return messenger{persistence: p}
}

type Messenger interface {
	GetAllMessagesFor(p Recipient) ReceivedUserMessages
	Send(Message) Sender
	GetConversationsBetween(firstProfileId profile.Id) ConversationGetter
}

type messenger struct {
	persistence Persistence
}

func (m messenger) GetConversationsBetween(p profile.Id) ConversationGetter {
	return conversationGetter{persistence: m.persistence, firstProfileId: p}
}

func (m messenger) GetAllMessagesFor(p Recipient) ReceivedUserMessages {
	return m.persistence.GetMessagesFor(p)
}

func (m messenger) Send(msg Message) Sender {
	return &messageSender{persistence: m.persistence, message: msg}
}
