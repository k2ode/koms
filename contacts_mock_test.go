package main

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

	assert.Equal(t, len(idMap), 2)

	firstId := idMap["+18001112222"]
	assert.Equal(t, firstId, "0")

	secondId := idMap["+18002223333"]
	assert.Equal(t, secondId, "1")
}

func TestContactsGetContactInvalidId(t *testing.T) {
	contacts, _ := NewContactsMock()

	_, err := contacts.GetContact("invalid")
	assert.Error(t, err)
}


func TestContactsGetFromId(t *testing.T) {
	contacts, _ := NewContactsMock()

	firstContact, err := contacts.GetContact("0")
	assert.NoError(t, err)

	assert.Equal(t, firstContact.id, "0")
	assert.Equal(t, firstContact.name, "Johnny")
	assert.Equal(t, firstContact.tags, []string{"friends"})

	secondContact, err := contacts.GetContact("1")
	assert.NoError(t, err)

	assert.Equal(t, secondContact.id, "1")
	assert.Equal(t, secondContact.name, "Andrew")
	assert.Equal(t, secondContact.tags, []string{"friends"})
}
