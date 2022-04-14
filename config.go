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
	parsePerson := func (id string) string {
		contact, err := client.GetContact(id)
		if err != nil { return "<" + id + ">" }

		return contact.name
	}

	parsePersons := func (ids []string) string {
		var res string
		for _, id := range ids {
			res = res + parsePerson(id)
		}
		return res
	}

	var res string

	if conversation.label != "" { res = conversation.label } else
	{ res = parsePersons(conversation.people) }

	if conversation.isGroupChat { res = GROUPCHAT_PREFIX + res }

	return res
}