package message

import (
	"spacemoon/network/profile"
	"time"
)

type ConversationThread interface {
}

func New(author Author, recipient Recipient, content string, postingTime time.Time) Message {
	return Message{
		Author: author, Recipient: recipient, Content: content, PostingTime: postingTime,
	}
}

type Message struct {
	Author      Author    `json:"author"`
	Recipient   Recipient `json:"recipient"`
	Content     string    `json:"content"`
	PostingTime time.Time `json:"posting_time"`
}

func (m *Message) String() string {
	return m.Content
}

func (m *Message) GetPostingTime() time.Time {
	return m.PostingTime
}

func (m *Message) GetAuthor() Author {
	return m.Author
}

func (m *Message) GetRecipient() Recipient {
	return m.Recipient
}

func (m *Message) SetAuthor(from Author) {
	m.Author = from
}

func (m *Message) SetRecipient(to Recipient) {
	m.Recipient = to
}

type Author profile.Id
type Recipient profile.Id
type SentUserMessages map[Recipient][]Message
type ReceivedUserMessages map[Author][]Message
