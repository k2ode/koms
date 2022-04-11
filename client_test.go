package koms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientNoProviders(t *testing.T) {
	client, err := NewClient()
	assert.NoError(t, err, "New client should not return an error")

	providers := client.GetProviders()
	assert.Empty(t, providers, "New client with no providers should return no providers")
}

func TestClientMockProvider(t *testing.T) {
	provider, _ := NewProviderMock()
	client, err := NewClient(provider)
	assert.NoError(t, err, "New client w/ mock provider should not return an error")

	providers := client.GetProviders()
	assert.Equal(t, len(providers), 1)
}