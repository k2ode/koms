package main

import (
	"strings"
)

func GetClient() (Client, error) {

	providerMockA, err := NewProviderMockA()
	if err != nil { return nil, err }

	providerMockB, err := NewProviderMockB()
	if err != nil { return nil, err }

	providers := []Provider{
		providerMockA,
		providerMockB,
	}

	contacts, err := NewContactsMock()
	if err != nil { return nil, err }

	client, err := NewClient(providers, contacts)

	return client, nil
}

func ParseConversation(client Client, conversation Conversation) string {
	parseIds := func (ids []string) string {
		return strings.Join(ids, ", ")
	}

	var result string

	if conversation.label != "" { result  = conversation.label } else
	{ result = parseIds(conversation.contactIds) }

	return result
}

func ParseMessage(client Client, message Message) string {
	messagePrefix := message.from.name
	if message.fromUser { messagePrefix = "[blue]" }
	return messagePrefix + message.body
}