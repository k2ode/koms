package main

const GROUPCHAT_PREFIX = "[grey][👥][-] "

func GetClient() (Client, error) {

	providerMockA, err := NewProviderMockA()
	if err != nil { return nil, err }

	providerMockB, err := NewProviderMockB()
	if err != nil { return nil, err }

	providers := []Provider{
		providerMockA,
		providerMockB,
	}

	contacts, err := NewContactsMock()
	if err != nil { return nil, err }

	client, err := NewClient(providers, contacts)

	return client, nil
}

func ParseConversation(client Client, conversation Conversation) string {
	var res string

	if conversation.label != "" { res = conversation.label } else
	{ res = ParsePersons(client, conversation.people) }

	if conversation.isGroupChat { res = GROUPCHAT_PREFIX + res }

	return res
}

func ParsePersons(client Client, persons []string) string {
	var res string

	for _, person := range persons {
		res = res + ParsePerson(client, person)
	}

	return res
}

func ParsePerson(client Client, contactId string) string {
	contact, err := client.GetContact(contactId)
	if err != nil { return "<" + contactId + ">" }
	return contact.name
}