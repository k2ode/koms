package koms

type Contacts interface {
	GetIdMap() (map[string]string, error)

	GetContact(id string) (Contact, error)
}