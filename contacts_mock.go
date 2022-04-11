package koms

type contactsMock struct {}


func NewContactsMock() (Contacts, error) {
	return &contactsMock{}, nil
}

func (contacts *contactsMock) GetIdMap() (map[string]string, error) {
	idMap := make(map[string]string)
	idMap["911"] = "12"
	return idMap, nil
}