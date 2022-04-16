package main

import (
	"errors"
	"time"
)


type ConversationData struct {
	meta Conversation
	messages []Message
}

type providerMock struct {
	id            string
	conversations []ConversationData
}

func NewProviderMockA() (Provider, error) {
	return &providerMock{
		id: "a",
		conversations: []ConversationData{
			{
				meta: Conversation{
					id: "0",
					isGroupChat: false,
					participantIds: []string{"a:0"},
					provider: "a",
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
						provider: "a",
						timestamp: time.Unix(0, 0),
						reactions: []Reaction{},
					},
					{
						id: "1",
						from: "a:0",
						body: "hello there",
						provider: "a",
						timestamp: time.Unix(200, 0),
						reactions: []Reaction{},
					},
				},
			},
			{
				meta: Conversation{
					id: "1",
					isGroupChat: true,
					participantIds: []string{"a:0", "a:1"},
					provider: "a",
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
						provider: "a",
						timestamp: time.Unix(200, 0),
						reactions: []Reaction{},
					},
					{
						id: "1",
						from: "a:1",
						body: "你好世界!",
						provider: "a",
						timestamp: time.Unix(300, 0),
						reactions: []Reaction{},
					},
				},
			},
		},
	}, nil
}

func NewProviderMockB() (Provider, error) {
	return &providerMock{
		id: "b",
		conversations: []ConversationData{
			{
				meta: Conversation{
					id: "0",
					isGroupChat: false,
					participantIds: []string{"b:0"},
					provider: "b",
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
						provider: "b",
						timestamp: time.Unix(100, 0),
						reactions: []Reaction{},
					},
				},
			},
		},
	}, nil
}

func (providerMock *providerMock) GetId() string {
	return providerMock.id
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
			provider: providerMock.id,
			timestamp: time.Now(),
			reactions: []Reaction{},
		})
		
		return nil
	}
	return errors.New("inavlid conversation id")
}