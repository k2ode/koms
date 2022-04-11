package koms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactsMock(t *testing.T) {
	_, err := NewContactsMock()
	assert.NoError(t, err, "New mock contacts should not return an error")
}

func TestContactsIdMap(t *testing.T) {
	contacts, _ := NewContactsMock()

	idMap, err := contacts.GetIdMap()
	assert.NoError(t, err)

	assert.Equal(t, len(idMap), 1)

	id := idMap["911"]
	assert.Equal(t, id, "12")
}
