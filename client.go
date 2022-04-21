package main

import (
	"errors"
	"sort"
)


type Client interface {
	GetProviders() []Provider

	GetProvider(id string) (Provider, error)
	
	GetConversations() ([]Conversation, error)

	GetContact(id string) (Contact, error)

	GetIdMap() (IdMap, error)

	GetConversationMessages(conversation Conversation) ([]Message, error)
}

type client struct {
	providers map[string]Provider
	contacts  Contacts
	idMap     IdMap
}

func NewClient(providers []Provider, contacts Contacts) (Client, error) {
	providerMap := make(map[string]Provider)
	for _, provider := range providers {
		providerMap[provider.GetId()] = provider
	}

	var idMap IdMap
	if contacts != nil { 
		var err error
		idMap, err = contacts.GetIdMap()
		if err != nil { return nil, err }
	}

	return &client{ providerMap, contacts, idMap }, nil
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

func (client *client) GetConversations() ([]Conversation, error) {
	var all []ConversationRaw

	for _, provider := range client.GetProviders() {
		providerConversations, err := provider.GetConversations()
		if err != nil { return nil, err }

		all = append(all, providerConversations...)
	}

	var conversations []Conversation

	if client.contacts == nil {
		for _, conversation := range all {
			personOrGroupChat := Conversation{
				conversations: []ConversationRaw{ conversation },
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
	if err != nil { return []Conversation{}, err }

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
			groupChat := Conversation{
				conversations: []ConversationRaw{ conversation },
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
			person := Conversation{
				conversations: []ConversationRaw{},
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

func (client *client) GetConversationMessages(conversation Conversation) ([]Message, error) {
	var messages []Message



	for _, convo := range conversation.conversations {
		provider, exists := client.providers[convo.provider]
		if !exists { return messages, errors.New("invalid provider") }
		messagesRaw, err := provider.GetConversationMessages(convo.id)
		if err != nil { panic(err) }


		var conversationMessages []Message

		for _, messageRaw := range messagesRaw {
			conversationMessages = append(conversationMessages, Message{
				id: messageRaw.id,
				from: Contact{},
				body: messageRaw.body,
				provider: provider.GetId(),
				timestamp: messageRaw.timestamp,
				reactions: messageRaw.reactions,
			})
		}


		messages = append(messages, conversationMessages...)
	}

	sort.Slice(messages, func(p, q int) bool {
		return messages[p].timestamp.Before(messages[q].timestamp)
	})

	return messages, nil
}