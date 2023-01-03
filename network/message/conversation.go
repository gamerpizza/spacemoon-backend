package message

import "spacemoon/network/profile"

type ConversationGetter interface {
	//And allows not only to have a cleaner more readable code but also to let us set up how we want to use
	//our persistence to get messages between to profiles
	And(secondProfileId profile.Id) []Message
}

type conversationGetter struct {
	firstProfileId profile.Id
	persistence    Persistence
}

func (c conversationGetter) And(secondProfileId profile.Id) []Message {
	var messages []Message
	messages = append(messages, c.persistence.GetMessagesFor(Recipient(c.firstProfileId))[Author(secondProfileId)]...)
	messages = append(messages, c.persistence.GetMessagesFor(Recipient(secondProfileId))[Author(c.firstProfileId)]...)
	return messages
}

type UserConversations map[profile.Id]Conversation
type Conversation []Message
