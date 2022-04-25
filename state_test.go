package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateMakeEmptyState(t *testing.T) {
	state := MakeEmptyState()
	assert.Empty(t, state.cache.conversations)
	assert.Empty(t, state.cache.messages)
	assert.Equal(t, state.pos, 0)
}

func MakeMockState() AppState {
	state := MakeEmptyState()
	convo := Conversation{
		conversations: []ConversationRaw{
			{ id: "6", participantIds: []string{"0"}, isGroupChat: false, label: "", provider: "a" },
		},
	}
	state.cache.conversations = append(state.cache.conversations, convo)
	state.cache.messages[0] = []Message{
		{ id: "9" },
	}

	state.conversations[0] = ConversationState{
		messagePos: 2,
		draft: "hmm",
		provider: "a",
	}

	return state
}

func TestStateGetConversationState(t *testing.T) {
	state := MakeMockState()

	conversationState := GetStateConversation(state)
	assert.NotNil(t, conversationState)
	assert.Equal(t, conversationState.messagePos, 2)
	assert.Equal(t, conversationState.draft, "hmm")
}

func TestStateGetCacheMessages(t *testing.T) {
	state := MakeMockState()

	messages, exists := GetCacheMessages(state)
	assert.True(t, exists)
	assert.NotEmpty(t, messages)
	assert.Equal(t, messages[0].id, "9")
}

func TestStateGetCacheConversation(t *testing.T) {
	state := MakeMockState()

	conversation := GetCacheConversation(state)
	assert.Equal(t, conversation.conversations[0].id, "6")
}

func TestStateGetStateMessagePos(t *testing.T) {
	state := MakeMockState()

	pos := GetStateMessagePos(state)
	assert.Equal(t, pos, 2)
}

func TestStateGetStateDraft(t *testing.T) {
	state := MakeMockState()

	draft := GetStateDraft(state)
	assert.Equal(t, draft, "hmm")
}

func TestStateGetStateProvider(t *testing.T) {
	state := MakeMockState()

	provider := GetStateProvider(state)
	assert.Equal(t, provider , "a")
}