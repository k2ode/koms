package koms

func NewProviderMock() (Provider, error) {
	return &provider{}, nil
}