package koms

type Provider interface {
	GetConversations() ([]Conversation, error)

	GetConversationMessages(id string) ([]Message, error)

	SendMessage(id string, body string) error
}

type provider struct {}