package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientProviderNone(t *testing.T) {
	client, err := NewClient([]Provider{}, nil)
	assert.NoError(t, err, "New client should not return an error")

	providers := client.GetProviders()
	assert.Empty(t, providers, "New client with no providers should return no providers")
}

func TestClientProviderMock(t *testing.T) {
	provider, _ := NewProviderMockA()
	client, err := NewClient([]Provider{provider}, nil)
	assert.NoError(t, err, "New client w/ mock provider should not return an error")

	providers := client.GetProviders()
	assert.Equal(t, len(providers), 1)
}

func TestClientContactsMock(t *testing.T) {
	contacts, _ := NewContactsMock()
	_, err := NewClient([]Provider{}, contacts)

	assert.NoError(t, err, "New client with mock contacts should not return an error")
}

func TestClientProviderMockGetConversations(t *testing.T) {
	provider, _ := NewProviderMockA()
	client, _ := NewClient([]Provider{provider}, nil)

	conversations, err := client.GetConversations()
	assert.NoError(t, err)

	assert.Equal(t, len(conversations), 2)
}

func TestClientContactsMockGetContact(t *testing.T) {
	contacts, _ := NewContactsMock()
	client, _ := NewClient([]Provider{}, contacts)

	contact, err := client.GetContact("0")
	assert.NoError(t, err)

	assert.Equal(t, contact.id, "0")
}