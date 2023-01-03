package message

import "spacemoon/network/profile"

// Persistence is used to save and retrieve messages
type Persistence interface {
	ConversationReader
	// Save saves a message and returns and error if one happens when trying to do so
	Save(Message) error
	// GetMessagesFor retrieves all the messages for a specific user profile.Id as a Recipient
	GetMessagesFor(Recipient) ReceivedUserMessages
	GetMessagesBy(Author) SentUserMessages
}

type ConversationReader interface {
	//GetMessagesBetween returns a ConversationGetter that returns the conversations between the two profile.Id
	//it works like this: `GetMessagesBetween(firstProfileId).And(secondProfileId)`
	GetMessagesBetween(firstProfileId profile.Id) ConversationGetter
}

func NewConversationReader(p Persistence) ConversationReader {
	return conversationReader{persistence: p}
}

type conversationReader struct {
	persistence Persistence
}

func (c conversationReader) GetMessagesBetween(firstProfileId profile.Id) ConversationGetter {
	return conversationGetter{persistence: c.persistence, firstProfileId: firstProfileId}
}
