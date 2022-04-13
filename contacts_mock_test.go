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

	assert.Equal(t, len(idMap), 1)

	id := idMap["911"]
	assert.Equal(t, id, "12")
}

func TestContactsGetContactInvalidId(t *testing.T) {
	contacts, _ := NewContactsMock()

	_, err := contacts.GetContact("invalid")
	assert.Error(t, err)
}


func TestContactsGetFromId(t *testing.T) {
	contacts, _ := NewContactsMock()

	contact, err := contacts.GetContact("12")
	assert.NoError(t, err)

	assert.NotNil(t, contact)

	assert.Equal(t, contact.id, "12")
	assert.Equal(t, contact.name, "The Police")
	assert.Equal(t, contact.tags, []string{"friends"})
}
