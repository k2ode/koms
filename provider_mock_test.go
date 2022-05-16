package main

import (
	"testing"
	"time"

	"github.com/k2ode/koms/types"
	"github.com/stretchr/testify/assert"
)

func TestProviderMock(t *testing.T) {
	provider, err := NewProviderMockA()
	assert.NoError(t, err, "New mock provider should not return an error")

	assert.NotEqual(t, provider, nil, "Provider should not be nil")
}

func TestProviderMockConversations(t *testing.T) {
	provider, _ := NewProviderMockA()
	conversations, err := provider.GetConversations()
	assert.NoError(t, err, "Mock provider should not return an error when retriving conversations")

	assert.Equal(t, len(conversations), 2, "Mock Provider should return 2 conversations")

	privateChat := conversations[0]
	assert.Equal(t, privateChat.Id, "0")
	assert.False(t, privateChat.IsGroupChat)
	assert.Equal(t, privateChat.ParticipantIds, []string{"a:0"})

	groupChat := conversations[1]
	assert.Equal(t, groupChat.Id, "1")
	// assert.Equal(t, groupChat.label, "Example Group Chat")
	assert.True(t, groupChat.IsGroupChat)
	assert.Equal(t, groupChat.ParticipantIds, []string{"a:0", "a:1"})
}

func TestProviderMockConversationMessagesInvalidId(t *testing.T) {
	provider, _ := NewProviderMockA()
	_, err := provider.GetConversationMessages("LOL!")
	assert.Error(t, err)
}

func TestProviderMockConversationMessagesPrivateChat(t *testing.T) {
	provider, _ := NewProviderMockA()
	messages, err := provider.GetConversationMessages("0")
	assert.NoError(t, err, "Mock provider should not return an error when retriving conversation messages")

	assert.Equal(t, len(messages), 2)

	firstMessage := messages[0]
	assert.Equal(t, firstMessage.Id, "0")
	assert.Equal(t, firstMessage.From, USER)
	assert.Equal(t, firstMessage.Body, "hi world")
	assert.Equal(t, firstMessage.Timestamp, time.Unix(0, 0))
	assert.Equal(t, firstMessage.Reactions, []types.Reaction{})
	// assert.Equal(t, firstMessage.provider, "a")


	secondMessage := messages[1]
	assert.Equal(t, secondMessage.Id, "1")
	assert.Equal(t, secondMessage.From, "a:0")
	assert.Equal(t, secondMessage.Body, "hello there")
	assert.Equal(t, secondMessage.Timestamp, time.Unix(200, 0))
	assert.Equal(t, secondMessage.Reactions, []types.Reaction{})
	// assert.Equal(t, secondMessage.provider, "a")
}

func TestProviderMockConversationMessageGroupChat(t *testing.T) {
	provider, _ := NewProviderMockA()
	messages, err := provider.GetConversationMessages("1")
	assert.NoError(t, err, "Mock provider should not return an error when retriving conversation messages")


	assert.Equal(t, len(messages), 2)

	firstMessage := messages[0]
	assert.Equal(t, firstMessage.Id, "0")
	assert.Equal(t, firstMessage.From, USER)
	assert.Equal(t, firstMessage.Body, "hi world")
	assert.Equal(t, firstMessage.Timestamp, time.Unix(200, 0))
	assert.Equal(t, firstMessage.Reactions, []types.Reaction{})
	// assert.Equal(t, firstMessage.provider, "a")

	unicodeMessage := messages[1]
	assert.Equal(t, unicodeMessage.Id, "1")
	assert.Equal(t, unicodeMessage.From, "a:1")
	assert.Equal(t, unicodeMessage.Body, "你好世界!")
	assert.Equal(t, unicodeMessage.Timestamp, time.Unix(300, 0))
	assert.Equal(t, unicodeMessage.Reactions, []types.Reaction{})
	// assert.Equal(t, unicodeMessage.provider, "a")
}

func TestProviderMockSendMessageInvalidId(t *testing.T) {
	provider, _ := NewProviderMockA()

	err := provider.SendMessage("-1", "what is a foo bar")
	assert.Error(t, err, "should error with invalid id")
}

func TestProviderMockSendMessage(t *testing.T) {
	provider, _ := NewProviderMockA()
	
	err := provider.SendMessage("0", "what is a foo bar")
	assert.NoError(t, err, "Sending a valid message should not return an error")

	messages, err := provider.GetConversationMessages("0")
	assert.Equal(t, len(messages), 3, "Afer sending message, messages length should reflect new message")

	messageSent := messages[2]
	assert.Equal(t, messageSent.Body, "what is a foo bar")
}
