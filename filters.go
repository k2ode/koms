package main

import "github.com/k2on/koms/types"

func Filter(state AppState) SearchResult {
	filters := state.search.filters

	conversations := FilterConversations(state.cache.conversations, filters)
	messages := make(map[int][]types.Message, 0)

	return AppCache{
		conversations,
		messages,
	}
}

func FilterConversations(conversations []types.Conversation, filters []SearchQueryFilter) []types.Conversation {
	if len(filters) == 0 { return conversations }

	var filtered []types.Conversation

	var ids []string
	for _, filter := range filters {
		ids = append(ids, filter.name)
	}

	Match := func(conversation types.Conversation) bool {
		HasId := func(id string) bool {
			for _, item := range conversation.ContactIds {
				if item == id { return true }
			}
			return false
		}

		for _, id := range ids {
			if !HasId(id) { return false }
		}
		return true
	}

	for _, conversation := range conversations {
		if Match(conversation) {
			filtered = append(filtered, conversation) 
		}
	}

	return filtered
}

func FilterMessages(messages []types.Message, filters []SearchQueryFilter) []types.Message {
	return messages
}