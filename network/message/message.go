package message

import (
	"spacemoon/network/profile"
	"time"
)

type Message struct {
	author    profile.Id
	recipient profile.Id
}

func (m *Message) String() string {
	return ""
}

func (m *Message) Time() time.Time {
	return time.Time{}
}

func (m *Message) Author() Author {
	return Author(m.author)
}

func (m *Message) Recipient() Recipient {
	return Recipient(m.recipient)
}

func (m *Message) SetAuthor(from profile.Id) {
	m.author = from
}

func (m *Message) SetRecipient(to profile.Id) {
	m.recipient = to
}

type Author profile.Id
type Recipient profile.Id

type UserMessages map[Author][]Message
type Messages map[Recipient]UserMessages
