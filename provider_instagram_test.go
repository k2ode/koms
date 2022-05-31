package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderInstagram(t *testing.T) {
	provider, err := NewProviderIG()
	assert.NoError(t, err)

	// err = provider.Sync()
	// assert.NoError(t, err)

	convos, err := provider.GetConversations()
	assert.NoError(t, err)

	fmt.Println(convos)

}