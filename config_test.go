package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigGetClient(t *testing.T) {
	client, err := GetClient()

	assert.NoError(t, err)
	assert.NotEqual(t, client, nil)
}