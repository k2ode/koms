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
	conversations []ConversationData
}

func NewProviderMockA() (Provider, error) {
	return &providerMock{
		conversations: []ConversationData{
			{
				meta: Conversation{
					id: "0",
					isGroupChat: false,
					people: []string{"a:0"},
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
						timestamp: time.Unix(int64(1649619517), 0),
						reactions: []Reaction{},
					},
					{
						id: "1",
						from: "a:0",
						body: "hello there",
						timestamp: time.Unix(int64(1649619617), 0),
						reactions: []Reaction{},
					},
				},
			},
			{
				meta: Conversation{
					id: "1",
					isGroupChat: true,
					people: []string{"a:0", "a:1"},
				},
				messages: []Message{
					{
						id: "0",
						from: USER,
						body: "hi world",
					},
					{
						id: "1",
						from: "a:1",
						body: "你好世界!",
					},
				},
			},
		},
	}, nil
}

func NewProviderMockB() (Provider, error) {
	return &providerMock{
		conversations: []ConversationData{
			{
				meta: Conversation{
					id: "0",
					isGroupChat: false,
					people: []string{"b:0"},
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