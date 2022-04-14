package main

type IdMap map[string]string

type Contacts interface {
	GetIdMap() (IdMap, error)

	GetContact(id string) (Contact, error)
}