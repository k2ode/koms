package main

import (
	"errors"
	"sort"
)


type Client interface {
	GetProviders() []Provider

	GetProvider(id string) (Provider, error)
	
	GetConversations() ([]PersonOrGroupChat, error)

	GetContact(id string) (Contact, error)

	GetIdMap() (IdMap, error)

	GetConversationMessages(conversation PersonOrGroupChat) ([]Message, error)
}

type client struct {
	providers map[string]Provider
	contacts    Contacts
}

func NewClient(providers []Provider, contacts Contacts) (Client, error) {
	providerMap := make(map[string]Provider)
	for _, provider := range providers {
		providerMap[provider.GetId()] = provider
	}

	return &client{ providerMap, contacts }, nil
}

func (client *client) GetProviders() []Provider {
	var providers []Provider
	for _, provider := range client.providers {
		providers = append(providers, provider)
	}
	return providers
}

func (client *client) GetProvider(id string) (Provider, error) {
	for _, provider := range client.providers {
		if provider.GetId() == id { return provider, nil }
	}
	return nil, errors.New("invalid provider")
}

func (client *client) GetConversations() ([]PersonOrGroupChat, error) {
	var all []Conversation

	for _, provider := range client.GetProviders() {
		providerConversations, err := provider.GetConversations()
		if err != nil { return nil, err }

		all = append(all, providerConversations...)
	}

	var conversations []PersonOrGroupChat

	if client.contacts == nil {
		for _, conversation := range all {
			personOrGroupChat := PersonOrGroupChat{
				conversations: []Conversation{ conversation },
				contactIds: conversation.participantIds,
				isGroupChat: conversation.isGroupChat,
				label: conversation.label,
			}
			conversations = append(conversations, personOrGroupChat)
		}
		return conversations, nil
	}

	// vvvvvvv    move all this to contacts      vvvvvv
	idMap, err := client.GetIdMap()
	if err != nil { return []PersonOrGroupChat{}, err }

	matchId := func (id string) string {
		match, exists := idMap[id]
		if !exists { return id }
		return match
	} 


	// map a contact id to []conversations position
	contactConversations := make(map[string]int)
	position := 0

	for _, conversation := range all {
		var contactIds []string
		for _, id := range conversation.participantIds {
			contactIds = append(contactIds, matchId(id))
		}

		if conversation.isGroupChat {
			groupChat := PersonOrGroupChat{
				conversations: []Conversation{ conversation },
				contactIds: contactIds,
				isGroupChat: true,
				label: conversation.label,
			}
			conversations = append(conversations, groupChat)
			position++

			continue
		}

		contactId := contactIds[0]

		var convPos int
		convPos, exists := contactConversations[contactId]

		if !exists {
			contactConversations[contactId] = position
			person := PersonOrGroupChat{
				conversations: []Conversation{},
				contactIds: contactIds,
				isGroupChat: false,
			}
			conversations = append(conversations, person)
			convPos = position
			position++
		} 

		// conversation.provider = 

		conversations[convPos].conversations = append(conversations[convPos].conversations, conversation)
	}

	return conversations, nil
}

func (client *client) GetContact(id string) (Contact, error) {
	return client.contacts.GetContact(id)
}

func (client *client) GetIdMap() (IdMap, error) {
	return client.contacts.GetIdMap()
}

func (client *client) GetConversationMessages(conversation PersonOrGroupChat) ([]Message, error) {
	var messages []Message



	for _, convo := range conversation.conversations {
		provider, exists := client.providers[convo.provider]
		if !exists { return messages, errors.New("invalid provider") }
		conversationMessages, err := provider.GetConversationMessages(convo.id)
		if err != nil { panic(err) }
		messages = append(messages, conversationMessages...)
	}

	sort.Slice(messages, func(p, q int) bool {
		return messages[p].timestamp.Before(messages[q].timestamp)
	})

	return messages, nil
}