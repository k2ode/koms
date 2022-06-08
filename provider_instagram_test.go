package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderInstagram(t *testing.T) {
	provider, err := NewProviderIG("/Users/koon/.koms/ig/")
	assert.NoError(t, err)

	// convos, err := provider.GetConversations()
	// assert.NoError(t, err)

	// fmt.Println(convos)

	// // err = provider.Sync()
	// // assert.NoError(t, err)

	// // convos, err := provider.GetConversations()
	// // assert.NoError(t, err)

	id := "340282366841710300949128295624182493824"
	msgs, err := provider.GetConversationMessages(id)
	assert.NoError(t, err)
	fmt.Println(msgs)

}