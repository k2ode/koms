package koms

type Contacts interface {
	GetIdMap() (map[string]string, error)

	// GetContactFromId(id string) (Contact, error)
}