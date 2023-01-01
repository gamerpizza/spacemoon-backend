package message

type Persistence interface {
	Save(Message) error
	GetMessagesFor(Recipient) UserMessages
}
