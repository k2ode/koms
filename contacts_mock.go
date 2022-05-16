package main

import (
	"errors"

	"github.com/k2ode/koms/types"
)

type contactsMock struct {}


func NewContactsMock() (Contacts, error) {
	return &contactsMock{}, nil
}

func (contacts *contactsMock) GetIdMap() (IdMap, error) {
	idMap := make(map[string]string)
	idMap["a:0"] = "0"
	idMap["a:1"] = "1"
	idMap["b:0"] = "0"
	return idMap, nil
}

func (contacts *contactsMock) GetContact(id string) (types.Contact, error) {
	if id == "0" {
		return types.Contact{
			Id: "0",
			Name: "Johnny",
			Tags: []string{"friends"},
		}, nil
	}
	if id == "1" {
		return types.Contact{
			Id: "1",
			Name: "Andrew",
			Tags: []string{"friends"},
		}, nil
	}
	return types.Contact{}, errors.New("invalid contact id")
}