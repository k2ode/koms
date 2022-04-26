package main

import (
	"strings"
)

const BIND_KEY_LEFT   = 'h'
const BIND_KEY_DOWN   = 'j'
const BIND_KEY_UP     = 'k'
const BIND_KEY_RIGHT  = 'l'
const BIND_KEY_TOP    = 'g'
const BIND_KEY_BOTTOM = 'G'
const BIND_KEY_CHAT   = '/'

// 0 is treated as auto
const ROWS_CONTENT = 0
const ROWS_INPUT = 1

const COLUMNS_CONVERSATIONS = 30
const COLUMNS_MESSAGES = 0
const COLUMNS_PREVIEW = 0

const ROW_POS_INPUT = 1 // the second row in the container grid
const COLUMN_POS_INPUT = 0

const ROW_SPAN_INPUT = 1
const COLUMN_SPAN_INPUT = 3 // span all 3 columns in the container

const HEIGHT_MIN_INPUT = 0
const WIDTH_MIN_INPUT = 0

const ROW_POS_MSGS = 0
const COLUMN_POS_MSGS = 1 // the second column in the container

const ROW_SPAN_MSGS = 1
const COLUMN_SPAN_MSGS = 1

const HEIGHT_MIN_MSGS = 0
const WIDTH_MIN_MSGS = 0


const ROW_POS_CONVOS = 0
const COLUMN_POS_CONVOS = 0

const ROW_SPAN_CONVOS = 1
const COLUMN_SPAN_CONVOS = 1

const HEIGHT_MIN_CONVOS = 0
const WIDTH_MIN_CONVOS = 0


const ROW_POS_PREVIEW = 0
const COLUMN_POS_PREVIEW = 2

const ROW_SPAN_PREVIEW = 1
const COLUMN_SPAN_PREVIEW = 1

const HEIGHT_MIN_PREVIEW = 0
const WIDTH_MIN_PREVIEW = 0


func GetClient() (Client, error) {

	providerMockA, err := NewProviderMockA()
	if err != nil { return nil, err }

	providerMockB, err := NewProviderMockB()
	if err != nil { return nil, err }

	providers := []Provider{
		providerMockA,
		providerMockB,
	}

	contacts, err := NewContactsMock()
	if err != nil { return nil, err }

	client, err := NewClient(providers, contacts)

	return client, nil
}

func ParseConversation(client Client, conversation Conversation) string {
	parseIds := func (ids []string) string {
		return strings.Join(ids, ", ")
	}

	var result string

	if conversation.label != "" { result  = conversation.label } else
	{ result = parseIds(conversation.contactIds) }

	return result
}

func ParseMessage(client Client, message Message) string {
	messagePrefix := message.from.name
	if message.fromUser { messagePrefix = "[blue]" }
	return messagePrefix + message.provider + ": " + message.body
}

func GetMessagePreview(message Message) string {
	return message.body
}

func UpdateStateFromKeyBind(state AppState, key rune) AppState {
	switch {
		case key == BIND_KEY_TOP || key == BIND_KEY_BOTTOM:
			// conversation := GetStateConversation(state)
			var messagePos int

			if key == BIND_KEY_BOTTOM {
				msgs, exists := GetCacheMessages(state)
				if !exists { return state }
				maxMsgs := len(msgs) - 1
				messagePos = maxMsgs 
			}

			state = UpdateStateMessagePos(state, messagePos)

			break
		case key == BIND_KEY_LEFT || key == BIND_KEY_RIGHT:
			maxConvos := len(state.cache.conversations) - 1

			var fn func(int) int
			if key == BIND_KEY_RIGHT { fn = MakeInc(maxConvos) } else
			{ fn = MakeDesc(maxConvos) }

			state.pos = fn(state.pos) 

			break
		case key == BIND_KEY_UP || key == BIND_KEY_DOWN:
			msgs, exists := GetCacheMessages(state) 
			if !exists { return state }
			maxMsgs := len(msgs) - 1

			var fn func(int) int
			if key == BIND_KEY_DOWN { fn = MakeInc(maxMsgs) } else
			{ fn = MakeDesc(maxMsgs) }

			state = UpdateStateMessagePosFn(state, fn)

			break
		case key == BIND_KEY_CHAT:
			state.focusInput = true
			break
	}
	return state
}