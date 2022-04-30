package main

type AppState struct {
	cache         AppCache
	conversations map[int]ConversationState
	pos           int
	focusInput    bool
	jumpBy        int
	quit          bool
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
		jumpBy: -1,
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

func GetStateMessage(state AppState) Message {
	return state.cache.messages[state.pos][state.conversations[state.pos].messagePos]
}

func UpdateStateConversationState(state AppState, fn func(ConversationState) ConversationState) AppState {
	conversation := GetStateConversation(state)
	state.conversations[state.pos] = fn(conversation)
	return state
}

func UpdateStateDraft(state AppState, draft string) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		convo.draft = draft
		return convo
	})
}

func UpdateStateMessagePos(state AppState, pos int) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		convo.messagePos = pos
		return convo
	})
}

func UpdateStateMessagePosFn(state AppState, fn func(int) int) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		convo.messagePos = fn(convo.messagePos)
		return convo
	})
}

func UpdateStateProvider(state AppState, provider string) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		convo.provider = provider
		return convo
	})
}