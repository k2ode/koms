package koms

import "errors"

type contactsMock struct {}


func NewContactsMock() (Contacts, error) {
	return &contactsMock{}, nil
}

func (contacts *contactsMock) GetIdMap() (map[string]string, error) {
	idMap := make(map[string]string)
	idMap["911"] = "12"
	return idMap, nil
}

func (contacts *contactsMock) GetContact(id string) (Contact, error) {
	if id == "12" {
		return Contact{
			id: "12",
			name: "The Police",
			tags: []string{"friends"},
		}, nil
	}
	return Contact{}, errors.New("invalid contact id")
}