package koms

import (
	"errors"
	"time"
)


type ProviderMock interface {
	GetConversations() ([]Conversation, error)

	GetConversationMessages(id string) ([]Message, error)
}

type providerMock struct {
	Conversations []Conversation
}

func NewProviderMock() (ProviderMock, error) {
	return &providerMock{
		Conversations: []Conversation{
			{
				"0",
				"Example Private Chat",
				false,
			},
			{
				"1",
				"Example Group Chat",
				true,
			},
		},
	}, nil
}

func (providerMock *providerMock) GetConversations() ([]Conversation, error) {
	return providerMock.Conversations, nil
}

func (providerMock *providerMock) GetConversationMessages(id string) ([]Message, error) {
	if id == "LOL!" { return nil, errors.New("invalid conversation id") }

	return []Message{
		{
			"0",
			USER,
			"hi world",
			time.Unix(int64(1649619517), 0),
			[]Reaction{},
		},
	}, nil
}