package main

import "errors"

type contactsMock struct {}


func NewContactsMock() (Contacts, error) {
	return &contactsMock{}, nil
}

func (contacts *contactsMock) GetIdMap() (map[string]string, error) {
	idMap := make(map[string]string)
	idMap["+18001112222"] = "0"
	idMap["+18002223333"] = "1"
	return idMap, nil
}

func (contacts *contactsMock) GetContact(id string) (Contact, error) {
	if id == "0" {
		return Contact{
			id: "0",
			name: "Johnny",
			tags: []string{"friends"},
		}, nil
	}
	if id == "1" {
		return Contact{
			id: "1",
			name: "Andrew",
			tags: []string{"friends"},
		}, nil
	}
	return Contact{}, errors.New("invalid contact id")
}