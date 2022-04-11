package koms

import (
	"errors"
	"time"
)


type ConversationData struct {
	meta Conversation
	messages []Message
}

type providerMock struct {
	conversations []ConversationData
}

func NewProviderMock() (Provider, error) {
	return &providerMock{
		conversations: []ConversationData{
			{
				meta: Conversation{
					id: "0",
					label: "Example Private Chat",
					isGroupChat: false,
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
						timestamp: time.Unix(int64(1649619517), 0),
						reactions: []Reaction{},
					},
				},
			},
			{
				meta: Conversation{
					id: "1",
					label: "Example Group Chat",
					isGroupChat: true,
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
					},
					{
						id: "1",
						from: "aitianqi",
						body: "你好世界!",
					},
				},
			},
		},
	}, nil
}

func (providerMock *providerMock) GetConversations() ([]Conversation, error) {
	var conversations []Conversation
	for _, cp := range providerMock.conversations {
		conversations = append(conversations, cp.meta)
	}
	return conversations, nil
}

func (providerMock *providerMock) GetConversationMessages(id string) ([]Message, error) {
	for _, cp := range providerMock.conversations {
		if cp.meta.id != id { continue }
		return cp.messages, nil
	}
	return nil, errors.New("invalid conversation id") 
}

func (providerMock *providerMock) SendMessage(id string, body string) error {
	for i, cp := range providerMock.conversations {
		if cp.meta.id != id { continue }
		providerMock.conversations[i].messages = append(providerMock.conversations[i].messages, Message{
			id: "0",
			from: USER,
			body: body,
			timestamp: time.Now(),
			reactions: []Reaction{},
		})
		
		return nil
	}
	return errors.New("inavlid conversation id")
}