package main

type Provider interface {
	GetConversations() ([]ConversationRaw, error)

	GetConversationMessages(id string) ([]MessageRaw, error)

	SendMessage(id string, body string) error

	GetId() string
}
