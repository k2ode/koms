package main

const GROUPCHAT_PREFIX = "[grey][👥][-] "

func GetClient() (Client, error) {

	providerMock, err := NewProviderMock()
	if err != nil { return nil, err }

	providers := []Provider{
		providerMock,
	}

	contacts, err := NewContactsMock()
	if err != nil { return nil, err }

	client, err := NewClient(providers, contacts)

	return client, nil
}

func ParseConversation(conversation Conversation) string {
	res := conversation.label

	if conversation.isGroupChat { res = GROUPCHAT_PREFIX + res }

	return res
}