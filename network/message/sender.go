package message

import (
	"errors"
	"fmt"
	"spacemoon/network/profile"
	"strings"
)

type Sender interface {
	From(sender profile.Id) Sender
	To(receiver profile.Id) Sender
	Now() error
}

type messageSender struct {
	from        Author
	to          Recipient
	persistence Persistence
	message     Message
}

func (s *messageSender) From(p profile.Id) Sender {
	s.from = Author(p)
	return s
}

func (s *messageSender) To(p profile.Id) Sender {
	s.to = Recipient(p)
	return s
}

func (s *messageSender) Now() error {
	if strings.TrimSpace(string(s.from)) == "" {
		return AuthorNotSetError
	}
	if strings.TrimSpace(string(s.to)) == "" {
		return RecipientNotSetError
	}
	s.message.SetAuthor(s.from)
	s.message.SetRecipient(s.to)
	err := s.persistence.Save(s.message)
	if err != nil {
		return fmt.Errorf("could not save message: %w", err)
	}
	return nil
}

var RecipientNotSetError error = errors.New("message Recipient `To()` not set")
var AuthorNotSetError error = errors.New("message Author `From()` not set")
