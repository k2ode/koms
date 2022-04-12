package koms


type Client interface {
	GetProviders() []Provider
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
