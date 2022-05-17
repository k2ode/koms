package main

import "github.com/k2on/koms/types"

type Provider interface {
	GetConversations() ([]types.ConversationRaw, error)

	GetConversationMessages(id string) ([]types.MessageRaw, error)

	SendMessage(id string, body string) error

	GetId() string
}
