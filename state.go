package main

import (
	"errors"

	"github.com/k2ode/koms/types"
)

type AppState struct {
	cache         AppCache
	conversations map[int]ConversationState
	pos           int
	focusInput    bool
	jumpBy        int
	quit          bool
}

type ConversationState struct {
	draft       string
	messagePos  int
	provider    string
	selected    []string
}

type AppCache struct {
	conversations []types.Conversation
	messages      map[int][]types.Message
}

func MakeEmptyState() AppState {
	return AppState{
		cache: AppCache{
			conversations: []types.Conversation{},
			messages: make(map[int][]types.Message),
		},
		conversations: make(map[int]ConversationState),
		pos: 0,
		jumpBy: -1,
	}
}

func GetStateConversation(state AppState) ConversationState {
	return state.conversations[state.pos]
}

func GetCacheConversation(state AppState) types.Conversation {
	return state.cache.conversations[state.pos]
}

func GetCacheMessages(state AppState) ([]types.Message, bool) {
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

func GetStateMessage(state AppState) (types.Message, error) {
	msgs, exists := GetCacheMessages(state)
	if !exists { return types.Message{}, errors.New("no cached messages for convo") }
	if len(msgs) == 0 { return types.Message{}, errors.New("no messages in convo") }
	messagePos := GetStateMessagePos(state)
	return msgs[messagePos], nil
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

func UpdateStateSelected(state AppState, fn func([]string) []string) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		convo.selected = fn(convo.selected)
		return convo
	})
}

func UpdateStateSelectedToggle(state AppState, toggledId string) AppState {
	return UpdateStateSelected(state, func(ids []string) []string {
		result := []string{}
		removed := false
		for _, id := range ids {
			if id == toggledId { removed = true; continue }
			result = append(result, id)
		}
		if !removed { result = append(result, toggledId) }
		return result 
	})
}
