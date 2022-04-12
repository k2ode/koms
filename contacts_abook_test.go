package koms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const PATH = "mock-addressbook"

func TestContactsABook(t *testing.T) {
	contacts, err := NewContactsABook(PATH)
	assert.NoError(t, err)
	assert.NotEqual(t, contacts, nil)
}

func TestContactsABookIdMap(t *testing.T) {
	contacts, _ := NewContactsABook(PATH)

	idMap, err := contacts.GetIdMap()
	assert.NoError(t, err)

	testIds := []string{"+18005559999", "john@example.com"}

	for _, testId := range testIds {
		id, exists := idMap[testId]
		assert.True(t, exists)
		assert.Equal(t, id, "0")
	}
}

func TestContactsABookGetContact(t *testing.T) {
	contacts, _ := NewContactsABook(PATH)

	contact, err := contacts.GetContact("0")
	assert.NoError(t, err)

	assert.Equal(t, contact.id, "0")
	assert.Equal(t, contact.name, "John Woods")
	assert.Equal(t, contact.tags, []string{"dr"})
}

func TestParseABookFormat(t *testing.T) {
	content :=
			"# abook addressbook file\n\n" +
			"[format]\n" +
			"program=abook\n" +
			"version=0.6.1\n\n" +
			"[0]\n" +
			"name=joe\n" +
			"mobile=mama\n\n"
	res := ParseABookFormat(content)

	format := res["format"]
	assert.Equal(t, format["program"], "abook")
	assert.Equal(t, format["version"], "0.6.1")

	joe := res["0"]
	assert.Equal(t, joe["name"], "joe")
	assert.Equal(t, joe["mobile"], "mama")
}