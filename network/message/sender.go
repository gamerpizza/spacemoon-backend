package message

import (
	"errors"
	"fmt"
	"spacemoon/login"
	"spacemoon/network/profile"
	"strings"
	"time"
)

type Sender interface {
	From(sender profile.Id) Sender
	To(receiver profile.Id) Sender
	Now() error
}

type messageSender struct {
	from             Author
	to               Recipient
	persistence      Persistence
	loginPersistence login.Persistence
	message          Message
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
	s.message.PostingTime = time.Now()
	check, err := s.loginPersistence.Check(login.UserName(s.to))
	if err != nil {
		return fmt.Errorf("could not check for recipient: %w", err)
	}
	if !check {
		return RecipientNotFoundError
	}
	err = s.persistence.Save(s.message)
	if err != nil {
		return fmt.Errorf("could not save message: %w", err)
	}
	return nil
}

var RecipientNotSetError error = errors.New("message GetRecipient `To()` not set")
var RecipientNotFoundError error = errors.New("invalid recipient, user not found")
var AuthorNotSetError error = errors.New("message GetAuthor `From()` not set")
