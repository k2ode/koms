package main

import "github.com/k2ode/koms/types"

type IdMap map[string]string

type Contacts interface {
	GetIdMap() (IdMap, error)

	GetContact(id string) (types.Contact, error)
}