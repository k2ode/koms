package main

import (
	"errors"
	"sort"

	"github.com/k2ode/koms/types"
)


type Client interface {
	GetProviders() []Provider

	GetProvider(id string) (Provider, error)
	
	GetConversations() ([]types.Conversation, error)

	GetContact(id string) (types.Contact, error)

	GetIdMap() (IdMap, error)

	GetConversationMessages(conversation types.Conversation) ([]types.Message, error)

	SendMessage(conversation types.Conversation, message string, providerIds []string) error
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

func (client *client) GetConversations() ([]types.Conversation, error) {
	var all []types.ConversationRaw

	for _, provider := range client.GetProviders() {
		providerConversations, err := provider.GetConversations()
		if err != nil { return nil, err }

		all = append(all, providerConversations...)
	}

	var conversations []types.Conversation

	if client.contacts == nil {
		for _, conversation := range all {
			personOrGroupChat := types.Conversation{
				Conversations: []types.ConversationRaw{ conversation },
				ContactIds: conversation.ParticipantIds,
				IsGroupChat: conversation.IsGroupChat,
				Label: conversation.Label,
			}
			conversations = append(conversations, personOrGroupChat)
		}
		return conversations, nil
	}

	// vvvvvvv    move all this to contacts      vvvvvv
	idMap, err := client.GetIdMap()
	if err != nil { return []types.Conversation{}, err }

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
		for _, id := range conversation.ParticipantIds {
			contactIds = append(contactIds, matchId(id))
		}

		if conversation.IsGroupChat {
			groupChat := types.Conversation{
				Conversations: []types.ConversationRaw{ conversation },
				ContactIds: contactIds,
				IsGroupChat: true,
				Label: conversation.Label,
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
			person := types.Conversation{
				Conversations: []types.ConversationRaw{},
				ContactIds: contactIds,
				IsGroupChat: false,
			}
			conversations = append(conversations, person)
			convPos = position
			position++
		} 

		// conversation.provider = 

		conversations[convPos].Conversations = append(conversations[convPos].Conversations, conversation)
	}

	return conversations, nil
}

func (client *client) GetContact(id string) (types.Contact, error) {
	return client.contacts.GetContact(id)
}

func (client *client) GetIdMap() (IdMap, error) {
	return client.contacts.GetIdMap()
}

func (client *client) GetConversationMessages(conversation types.Conversation) ([]types.Message, error) {
	var messages []types.Message



	for _, convo := range conversation.Conversations {
		provider, exists := client.providers[convo.Provider]
		if !exists { return messages, errors.New("invalid provider") }
		messagesRaw, err := provider.GetConversationMessages(convo.Id)
		if err != nil { panic(err) }


		var conversationMessages []types.Message

		for _, messageRaw := range messagesRaw {
			conversationMessages = append(conversationMessages, types.Message{
				Id: messageRaw.Id,
				From: types.Contact{},
				Body: messageRaw.Body,
				Provider: provider.GetId(),
				Timestamp: messageRaw.Timestamp,
				Reactions: messageRaw.Reactions,
			})
		}


		messages = append(messages, conversationMessages...)
	}

	sort.Slice(messages, func(p, q int) bool {
		return messages[p].Timestamp.Before(messages[q].Timestamp)
	})

	return messages, nil
}

func (client *client) SendMessage(conversation types.Conversation, message string, providerIds []string) error {
	if message == "" { return errors.New("empty message") }

	for _, providerId := range providerIds {
		for _, convo := range conversation.Conversations {
			if convo.Provider == providerId {
				err := client.providers[providerId].SendMessage(convo.Id, message)
				if err != nil { return err }
			}
		}

	}

	return nil
}