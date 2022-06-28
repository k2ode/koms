package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderIMessage(t *testing.T) {
	provider, err := NewProviderIMessage()
	assert.NoError(t, err)

	convo, err := provider.GetConversations()

	msg, err := provider.GetConversationMessages(convo[3].Id)

	fmt.Println(msg)
	fmt.Println(err)
}