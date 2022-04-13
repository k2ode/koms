package main


type Client interface {
	GetProviders() []Provider

	GetConversations() ([]Conversation, error)
}

type client struct {
	providers []Provider
}

func NewClient(providers []Provider, contacts Contacts) (Client, error) {
	return &client{ providers }, nil
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