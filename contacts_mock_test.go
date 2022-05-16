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

	assert.Equal(t, len(idMap), 3)

	a0 := idMap["a:0"]
	assert.Equal(t, a0, "0")

	a1 := idMap["a:1"]
	assert.Equal(t, a1, "1")

	b0 := idMap["b:0"]
	assert.Equal(t, b0, "0")
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

	assert.Equal(t, firstContact.Id, "0")
	assert.Equal(t, firstContact.Name, "Johnny")
	assert.Equal(t, firstContact.Tags, []string{"friends"})

	secondContact, err := contacts.GetContact("1")
	assert.NoError(t, err)

	assert.Equal(t, secondContact.Id, "1")
	assert.Equal(t, secondContact.Name, "Andrew")
	assert.Equal(t, secondContact.Tags, []string{"friends"})
}
