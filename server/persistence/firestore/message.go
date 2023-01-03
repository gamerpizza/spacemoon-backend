package firestore

import (
	"fmt"
	"google.golang.org/api/iterator"
	"spacemoon/network/message"
	"spacemoon/network/profile"
)

type conversationGetter struct {
	firstProfileId profile.Id
	persistence    *fireStorePersistence
}

func (c conversationGetter) And(secondProfileId profile.Id) []message.Message {
	coll := c.persistence.storage.Collection(messagesCollection)
	docs1 := coll.Where("Recipient", "==", c.firstProfileId).Where("Author", "==", secondProfileId).Documents(c.persistence.ctx)
	docs2 := coll.Where("Recipient", "==", secondProfileId).Where("Author", "==", c.firstProfileId).Documents(c.persistence.ctx)
	var messages []message.Message
	for {
		doc, err := docs1.Next()
		if err != nil {
			return nil
		}
		var m message.Message
		err = doc.DataTo(&m)
		if err != nil {
			return nil
		}
		messages = append(messages, m)
		break
	}
	for {
		doc, err := docs2.Next()
		if err != nil {
			return nil
		}
		var m message.Message
		err = doc.DataTo(&m)
		if err != nil {
			return nil
		}
		messages = append(messages, m)
	}
}

func (p *fireStorePersistence) GetMessagesBetween(firstProfileId profile.Id) message.ConversationGetter {
	return conversationGetter{firstProfileId: firstProfileId, persistence: p}
}

func (p *fireStorePersistence) Save(m message.Message) error {
	coll := p.storage.Collection(messagesCollection)
	_, err := coll.NewDoc().Set(p.ctx, m)
	if err != nil {
		return fmt.Errorf("could not save to firestore: %w", err)
	}
	return nil
}

func (p *fireStorePersistence) GetMessagesFor(recipient message.Recipient) message.ReceivedUserMessages {
	coll := p.storage.Collection(messagesCollection)
	docs := coll.Where("Recipient", "==", recipient).Documents(p.ctx)
	messages := message.ReceivedUserMessages{}
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		msg := message.Message{}
		err = doc.DataTo(&msg)
		if err != nil {
			return nil
		}
		messages[msg.Author] = append(messages[msg.Author], msg)
	}
	return messages
}

func (p *fireStorePersistence) GetMessagesBy(author message.Author) message.SentUserMessages {
	coll := p.storage.Collection(messagesCollection)
	docs := coll.Where("Author", "==", author).Documents(p.ctx)
	messages := message.SentUserMessages{}
	for {
		doc, err := docs.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}
		msg := message.Message{}
		err = doc.DataTo(&msg)
		if err != nil {
			return nil
		}
		messages[msg.Recipient] = append(messages[msg.Recipient], msg)
	}
	return messages
}
