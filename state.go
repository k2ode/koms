package main

type AppState struct {
	cache         AppCache
	conversations map[int]ConversationState
	pos           int
}

type ConversationState struct {
	draft      string
	messagePos int
	provider   string
}

type AppCache struct {
	conversations []Conversation
	messages      map[int][]Message
}

func MakeEmptyState() AppState {
	return AppState{
		cache: AppCache{
			conversations: []Conversation{},
			messages: make(map[int][]Message),
		},
		conversations: make(map[int]ConversationState),
		pos: 0,
	}
}

func GetStateConversation(state AppState) ConversationState {
	return state.conversations[state.pos]
}

func GetCacheConversation(state AppState) Conversation {
	return state.cache.conversations[state.pos]
}

func GetCacheMessages(state AppState) ([]Message, bool) {
	messages, exists := state.cache.messages[state.pos]
	return messages, exists
}

func GetStateMessagePos(state AppState) int {
	return state.conversations[state.pos].messagePos
}

func GetStateDraft(state AppState) string {
	return state.conversations[state.pos].draft
}

func GetStateProvider(state AppState) string {
	return state.conversations[state.pos].provider
}