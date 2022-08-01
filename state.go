package main

import (
	"errors"

	"github.com/k2on/koms/types"
)

type AppState struct {
	cache         AppCache
	conversations map[int]ConversationState
	pos           int
	focusInput    bool
	jumpBy        int
	quit          bool
	search        SearchQuery
}

type SearchQuery struct {
	open bool
	opened bool
	filters []SearchQueryFilter
	filterPos int
	focusInput bool
	result SearchResult
}

type SearchResult = AppCache

type SearchQueryFilter struct {
	name string
}

type ConversationState struct {
	draft       string
	messagePos  int
	provider    string
	selected    []string
	carouselImageSelections map[int]int
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

func GetStateConversations(state AppState) []types.Conversation {
	isFiltered := len(state.search.filters) > 0
	if isFiltered { return state.search.result.conversations }
	return state.cache.conversations
}

func GetStateConversationsLen(state AppState) int {
	return len(GetStateConversations(state))
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

func UpdateStateConversationPos(state AppState, pos int) AppState {
	state.pos = pos
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

func UpdateStateCarouselSelectedImage(state AppState, fn func(int) int) AppState {
	return UpdateStateConversationState(state, func(convo ConversationState) ConversationState {
		if convo.carouselImageSelections == nil { convo.carouselImageSelections = make(map[int]int) }
		imageSelection := convo.carouselImageSelections[convo.messagePos]
		convo.carouselImageSelections[convo.messagePos] = fn(imageSelection)
		return convo
	})
}

func GetStateCarouselSelectedImage(state AppState) int {
	convo := GetStateConversation(state)
	return convo.carouselImageSelections[convo.messagePos]
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

func UpdateStateSearchFilterPos(state AppState, pos int) AppState {
	state.search.filterPos = pos
	return state
}

func UpdateStateSearchFilterPosFn(state AppState, fn IntMod) AppState {
	return UpdateStateSearchFilterPos(state, fn(state.search.filterPos))
}

func UpdateStateSearchOpen(state AppState) AppState {
	state.search.open = true
	return UpdateStateSearchFocus(state)
}

func UpdateStateSearchClose(state AppState) AppState {
	state.search.open = false
	state.search.opened = false
	state.search.focusInput = false
	return state
}

func UpdateStateSearchFocus(state AppState) AppState {
	state.search.focusInput = true
	return state
}

func UpdateStateSearchFilters(state AppState, filters []SearchQueryFilter) AppState {
	state.search.filters = filters
	return state
}

func UpdateStateSearchFiltersClear(state AppState) AppState {
	return UpdateStateSearchFilters(state, []SearchQueryFilter{})
}
