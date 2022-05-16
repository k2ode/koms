package main

import (
	"testing"

	"github.com/k2ode/koms/types"
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

	MockConvo := func(id string) types.Conversation {
		return types.Conversation{
			Conversations: []types.ConversationRaw{
				{ Id: id, ParticipantIds: []string{"0"}, IsGroupChat: false, Label: "", Provider: "a" },
			},
		}
	}

	convo1 := MockConvo("6")
	convo2 := MockConvo("1")
	state.cache.conversations = append(state.cache.conversations, convo1, convo2)
	state.cache.messages[0] = []types.Message{
		{ Id: "9" },
		{ Id: "10" },
		{ Id: "11" },
		{ Id: "12" },
	}

	state.conversations[0] = ConversationState{
		messagePos: 2,
		draft: "hmm",
		provider: "a",
		selected: []string{},
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
	assert.Equal(t, messages[0].Id, "9")
}

func TestStateGetCacheConversation(t *testing.T) {
	state := MakeMockState()

	conversation := GetCacheConversation(state)
	assert.Equal(t, conversation.Conversations[0].Id, "6")
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

func TestStateGetStateMessage(t *testing.T) {
	state := MakeMockState()

	msg, err := GetStateMessage(state)
	assert.NoError(t, err)
	assert.Equal(t, msg.Id, "11")
}

func TestStateUpdateStateDraft(t *testing.T) {
	state := MakeMockState()

	newState := UpdateStateDraft(state, "thinking")

	assert.Equal(t, newState.conversations[newState.pos].draft, "thinking")
}

func TestStateUpdateStateMessagePos(t *testing.T) {
	state := MakeMockState()

	newState := UpdateStateMessagePos(state, 0)

	assert.Equal(t, newState.conversations[newState.pos].messagePos, 0)
}


func TestStateUpdateStateMessagePosFn(t *testing.T) {
	state := MakeMockState()

	newState := UpdateStateMessagePosFn(state, func(i int) int { return i - 1 })

	assert.Equal(t, newState.conversations[newState.pos].messagePos, 1)
}

func TestStateUpdateStateProvider(t *testing.T) {
	state := MakeMockState()

	assert.Equal(t, GetStateProvider(state), "a")

	state = UpdateStateProvider(state, "b")

	assert.Equal(t, GetStateProvider(state), "b")
}

func TestStateUpdateStateSelected(t *testing.T) {
	state := MakeMockState()

	assert.Equal(t, state.conversations[0].selected, []string{})

	state = UpdateStateSelectedToggle(state, "0")

	assert.Equal(t, state.conversations[0].selected, []string{"0"})

	state = UpdateStateSelectedToggle(state, "1")

	assert.Equal(t, state.conversations[0].selected, []string{"0", "1"})

	state = UpdateStateSelectedToggle(state, "0")

	assert.Equal(t, state.conversations[0].selected, []string{"1"})
}
