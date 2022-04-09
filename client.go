package koms


type Client interface {
	GetProviders() []Provider
}

type client struct {}

func NewClient() (Client, error) {
	return &client{}, nil
}

func (client *client) GetProviders() []Provider {
	return []Provider{}
}
