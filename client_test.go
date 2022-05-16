package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientProviderNone(t *testing.T) {
	client, err := NewClient([]Provider{}, nil)
	assert.NoError(t, err, "New client should not return an error")

	providers := client.GetProviders()
	assert.Empty(t, providers, "New client with no providers should return no providers")
}

func TestClientProviderMockA(t *testing.T) {
	provider, _ := NewProviderMockA()
	client, err := NewClient([]Provider{provider}, nil)
	assert.NoError(t, err, "New client w/ mock provider should not return an error")

	providers := client.GetProviders()
	assert.Equal(t, len(providers), 1)
}

func TestClientContactsMockA(t *testing.T) {
	contacts, _ := NewContactsMock()
	_, err := NewClient([]Provider{}, contacts)

	assert.NoError(t, err, "New client with mock contacts should not return an error")
}

func TestClientProviderMockAGetConversations(t *testing.T) {
	provider, _ := NewProviderMockA()
	client, _ := NewClient([]Provider{provider}, nil)

	conversations, err := client.GetConversations()
	assert.NoError(t, err)

	assert.Equal(t, len(conversations), 2)
}

func TestClientContactsMockAGetContact(t *testing.T) {
	contacts, _ := NewContactsMock()
	client, _ := NewClient([]Provider{}, contacts)

	contact, err := client.GetContact("0")
	assert.NoError(t, err)

	assert.Equal(t, contact.Id, "0")
}

func TestClientMockContactsGetIdMap(t *testing.T) {
	contacts, _ := NewContactsMock()
	client, _ := NewClient([]Provider{}, contacts)

	idMap, err := client.GetIdMap()
	assert.NoError(t, err)
	a0, exists := idMap["a:0"]
	assert.True(t, exists)
	assert.Equal(t, a0, "0")

	a1, exists := idMap["a:1"]
	assert.True(t, exists)
	assert.Equal(t, a1, "1")

	b0, exists := idMap["b:0"]
	assert.True(t, exists)
	assert.Equal(t, b0, "0")
}

func TestClientMockABContact(t *testing.T) {
	contacts, _ := NewContactsMock()
	providerA, _ := NewProviderMockA()
	providerB, _ := NewProviderMockB()
	client, _ := NewClient([]Provider{providerA, providerB}, contacts)

	conversations, err := client.GetConversations()
	assert.NoError(t, err)
	assert.Equal(t, len(conversations), 2)

	conversation := conversations[0]
	messages, err := client.GetConversationMessages(conversation)

	assert.Equal(t, len(messages), 4)

	firstMessage := messages[0]
	assert.Equal(t, firstMessage.Provider, "a")
	assert.Equal(t, firstMessage.Timestamp, time.Unix(0, 0))

	secondMessage := messages[1]
	assert.Equal(t, secondMessage.Provider, "b")
	assert.Equal(t, secondMessage.Timestamp, time.Unix(100, 0))

	thirdMessage := messages[2]
	assert.Equal(t, thirdMessage.Provider, "a")
	assert.Equal(t, thirdMessage.Timestamp, time.Unix(200, 0))
}

func TestClientMockAGetProvider(t *testing.T) {
	providerA, _ := NewProviderMockA()
	client, _ := NewClient([]Provider{providerA}, nil)

	provider, err := client.GetProvider("a")

	assert.NoError(t, err)
	assert.Equal(t, provider.GetId(), "a")

}


func TestClientMockASendMessageEmpty(t *testing.T) {
	providerA, _ := NewProviderMockA()
	client, _ := NewClient([]Provider{providerA}, nil)

	conversations, _ := client.GetConversations()

	firstConversation := conversations[0]

	err := client.SendMessage(firstConversation, "", []string{"a"})
	assert.Error(t, err)
}

func TestClientMockASendMessage(t *testing.T) {
	providerA, _ := NewProviderMockA()
	client, _ := NewClient([]Provider{providerA}, nil)

	conversations, _ := client.GetConversations()

	firstConversation := conversations[0]

	err := client.SendMessage(firstConversation, "hello world", []string{"a"})

	assert.NoError(t, err)

	messages, err := client.GetConversationMessages(firstConversation)
	assert.NoError(t, err)

	assert.Equal(t, len(messages), 3)

}