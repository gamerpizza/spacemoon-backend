package message

import (
	"spacemoon/network/profile"
	"time"
)

type ConversationThread interface {
}

func New(author Author, recipient Recipient, content string, postingTime time.Time) Message {
	return Message{
		author: author, recipient: recipient, content: content, postingTime: postingTime,
	}
}

type Message struct {
	author      Author
	recipient   Recipient
	content     string
	postingTime time.Time
}

func (m *Message) String() string {
	return m.content
}

func (m *Message) PostingTime() time.Time {
	return m.postingTime
}

func (m *Message) Author() Author {
	return m.author
}

func (m *Message) Recipient() Recipient {
	return m.recipient
}

func (m *Message) SetAuthor(from Author) {
	m.author = from
}

func (m *Message) SetRecipient(to Recipient) {
	m.recipient = to
}

type Author profile.Id
type Recipient profile.Id
type SentUserMessages map[Recipient][]Message
type ReceivedUserMessages map[Author][]Message
