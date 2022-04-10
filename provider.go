package koms

type Provider interface {
	GetConversations() ([]Conversation, error)
}

type provider struct {}