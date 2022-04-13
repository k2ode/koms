package main

import (
	"errors"
	"os"
	"strings"
)

type ABookFormat map[string]map[string]string

type contactsABook struct {
	parsed ABookFormat
}

func NewContactsABook(path string) (Contacts, error) {
	file, err := os.ReadFile(path)
	if err != nil { return nil, err }

	content := string(file)
	parsed := ParseABookFormat(content)

	return &contactsABook{ parsed }, nil
}

func ParseABookFormat(content string) ABookFormat {
	lines := Lines(content)

	var id string
	result := make(ABookFormat)
	for _, line := range(lines) {
		isEmpty := line == ""
		if isEmpty { continue }

		isComment := line[0] == '#'
		if isComment { continue }

		isNewId := line[0] == '['
		if isNewId {
			id = Between(line, "[", "]")
			result[id] = make(map[string]string)
			continue
		}

		lineParts := strings.Split(line, "=")
		key := lineParts[0]
		val := lineParts[1]

		result[id][key] = val
	}

	return result
}

func (contacts *contactsABook) GetIdMap() (map[string]string, error) {
	idMap := make(map[string]string)
	for id, attrs := range contacts.parsed {

		addIfExists := func (attr string) {
			val, exists := attrs[attr]
			if !exists { return }
			idMap[val] = id
		}

		addIfExists("email")
		addIfExists("mobile")
	}
	return idMap, nil
}

func (contacts *contactsABook) GetContact(id string) (Contact, error) {
	for contactId, attrs := range contacts.parsed {
		if contactId != id { continue }

		var name string
		var tags []string

		name = attrs["name"]
		groups, hasGroups := attrs["groups"]
		if hasGroups {
			tags = strings.Split(groups, ",")
		}

		contact := Contact{
			id,
			name,
			tags,
		}
		return contact, nil
	}
	return Contact{}, errors.New("unknown contact")
}
