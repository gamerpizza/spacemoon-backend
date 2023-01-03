package message

import "spacemoon/network/profile"

func NewMessenger(p Persistence) Messenger {
	return messenger{persistence: p}
}

type Messenger interface {
	GetAllMessagesFor(p Recipient) ReceivedUserMessages
	Send(Message) Sender
	GetConversationsBetween(firstProfileId profile.Id) ConversationGetter
	//GetConversationsWith retrieves all the conversations that a profile.Id is referenced in
	//either as an Author or as a Recipient
	GetConversationsWith(profile.Id) UserConversations
}

type messenger struct {
	persistence Persistence
}

func (m messenger) ListConversationsFor(id profile.Id) []profile.Id {
	//TODO implement me
	panic("implement me")
}

func (m messenger) GetConversationsWith(p profile.Id) UserConversations {
	receivedMessages := m.persistence.GetMessagesFor(Recipient(p))
	sentMessages := m.persistence.GetMessagesBy(Author(p))

	conversations := make(UserConversations)
	for author, messages := range receivedMessages {
		conversations[profile.Id(author)] = append(conversations[profile.Id(author)], messages...)
	}
	for recipient, messages := range sentMessages {
		conversations[profile.Id(recipient)] = append(conversations[profile.Id(recipient)], messages...)
	}

	return conversations
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
