package main


type Client interface {
	GetProviders() []Provider

	GetConversations() ([]Conversation, error)

	GetContact(id string) (Contact, error)
}

type client struct {
	providers []Provider
	contacts  Contacts
}

func NewClient(providers []Provider, contacts Contacts) (Client, error) {
	return &client{ providers, contacts }, nil
}

func (client *client) GetProviders() []Provider {
	return client.providers
}

func (client *client) GetConversations() ([]Conversation, error) {
	conversations := []Conversation{}

	for _, provider := range client.GetProviders() {
		providerConversations, err := provider.GetConversations()
		if err != nil { return nil, err }
		conversations = append(conversations, providerConversations...)
	}
	return conversations, nil
}

func (client *client) GetContact(id string) (Contact, error) {
	return client.contacts.GetContact(id)
}