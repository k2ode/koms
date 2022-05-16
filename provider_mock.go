package main

import (
	"errors"
	"time"

	. "github.com/k2ode/koms/types"
)


type ConversationData struct {
	meta     ConversationRaw
	messages []MessageRaw
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
				meta: ConversationRaw{
					Id: "0",
					IsGroupChat: false,
					ParticipantIds: []string{"a:0"},
					Provider: "a",
				},
				messages: []MessageRaw{
					{
						Id: "0",
						From: USER,
						Body: "hi world",
						Timestamp: time.Unix(0, 0),
						Reactions: []Reaction{},
					},
					{
						Id: "1",
						From: "a:0",
						Body: "hello there",
						Timestamp: time.Unix(200, 0),
						Reactions: []Reaction{},
					},
				},
			},
			{
				meta: ConversationRaw{
					Id: "1",
					IsGroupChat: true,
					ParticipantIds: []string{"a:0", "a:1"},
					Provider: "a",
				},
				messages: []MessageRaw{
					{
						Id: "0",
						From: USER,
						Body: "hi world",
						Timestamp: time.Unix(200, 0),
						Reactions: []Reaction{},
					},
					{
						Id: "1",
						From: "a:1",
						Body: "你好世界!",
						Timestamp: time.Unix(300, 0),
						Reactions: []Reaction{},
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
				meta: ConversationRaw{
					Id: "0",
					IsGroupChat: false,
					ParticipantIds: []string{"b:0"},
					Provider: "b",
				},
				messages: []MessageRaw{
					{
						Id: "0",
						From: USER,
						Body: "hi world",
						Timestamp: time.Unix(100, 0),
						Reactions: []Reaction{},
					},
					{
						Id: "1",
						From: "0",
						Body: "ay look at this",
						Timestamp: time.Unix(300, 0),
						Reactions: []Reaction{},
					},
				},
			},
		},
	}, nil
}

func (providerMock *providerMock) GetId() string {
	return providerMock.id
}

func (providerMock *providerMock) GetConversations() ([]ConversationRaw, error) {
	var conversations []ConversationRaw
	for _, cp := range providerMock.conversations {
		conversations = append(conversations, cp.meta)
	}
	return conversations, nil
}

func (providerMock *providerMock) GetConversationMessages(id string) ([]MessageRaw, error) {
	for _, cp := range providerMock.conversations {
		if cp.meta.Id != id { continue }
		return cp.messages, nil
	}
	return nil, errors.New("invalid conversation id") 
}

func (providerMock *providerMock) SendMessage(id string, body string) error {
	for i, cp := range providerMock.conversations {
		if cp.meta.Id != id { continue }
		providerMock.conversations[i].messages = append(providerMock.conversations[i].messages, MessageRaw{
			Id: "0",
			From: USER,
			Body: body,
			Timestamp: time.Now(),
			Reactions: []Reaction{},
		})
		
		return nil
	}
	return errors.New("inavlid conversation id")
}