package koms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderMock(t *testing.T) {
	provider, err := NewProviderMock()
	assert.NoError(t, err, "New mock provider should not return an error")

	assert.NotEqual(t, provider, nil, "Provider should not be nil")
}

func TestProviderMockConversations(t *testing.T) {
	provider, _ := NewProviderMock()
	conversations, err := provider.GetConversations()
	assert.NoError(t, err, "Mock provider should not return an error when retriving conversations")

	assert.Equal(t, len(conversations), 2, "Mock Provider should return 2 conversations")

	privateChat := conversations[0]
	assert.Equal(t, privateChat.id, "0")
	assert.Equal(t, privateChat.label, "Example Private Chat")
	assert.False(t, privateChat.isGroupChat)

	groupChat := conversations[1]
	assert.Equal(t, groupChat.id, "1")
	assert.Equal(t, groupChat.label, "Example Group Chat")
	assert.True(t, groupChat.isGroupChat)
}

func TestProviderMockConversationMessagesInvalidId(t *testing.T) {
	provider, _ := NewProviderMock()
	_, err := provider.GetConversationMessages("LOL!")
	assert.Error(t, err)
}

func TestProviderMockConversationMessagesPrivateChat(t *testing.T) {
	provider, _ := NewProviderMock()
	messages, err := provider.GetConversationMessages("0")
	assert.NoError(t, err, "Mock provider should not return an error when retriving conversation messages")

	assert.Equal(t, len(messages), 1)

	firstMessage := messages[0]
	assert.Equal(t, firstMessage.id, "0")
	assert.Equal(t, firstMessage.from, USER)
	assert.Equal(t, firstMessage.body, "hi world")
	assert.Equal(t, firstMessage.timestamp.Unix(), int64(1649619517))
	assert.Equal(t, firstMessage.reactions, []Reaction{})
}