package main

type Provider interface {
	GetConversations() ([]Conversation, error)

	GetConversationMessages(id string) ([]Message, error)

	SendMessage(id string, body string) error
}
