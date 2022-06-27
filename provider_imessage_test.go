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

	fmt.Println(convo)
	fmt.Println(err)
}