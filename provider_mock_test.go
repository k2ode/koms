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