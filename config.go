package main

import (
	"strings"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/k2ode/koms/types"
)

// Bindings

const BIND_KEY_LEFT   = 'h'
const BIND_KEY_DOWN   = 'j'
const BIND_KEY_UP     = 'k'
const BIND_KEY_RIGHT  = 'l'
const BIND_KEY_TOP    = 'g'
const BIND_KEY_BOTTOM = 'G'
const BIND_KEY_SELECT = 'v'
const BIND_KEY_QUIT   = 'q'
const BIND_KEY_CHAT   = '/'


// Colors

const FOCUS_BACKGROUND_NORMAL = tcell.ColorWhite
const FOCUS_BACKGROUND_INSERT = tcell.ColorGray

const FOCUS_FOREGROUND_NORMAL = tcell.ColorDefault
const FOCUS_FOREGROUND_INSERT = tcell.ColorWhite

// Layout

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

func ParseConversation(client Client, conversation types.Conversation) string {
	parseIds := func (ids []string) string {
		return strings.Join(ids, ", ")
	}

	var result string

	if conversation.Label != "" { result = conversation.Label } else
	{ result = parseIds(conversation.ContactIds) }

	return result
}

func ParseMessage(client Client, conversation ConversationState, message types.Message) string {
	var prefix string
	id := message.Provider + message.Id
	isSelected := Contains(conversation.selected, id)
	if isSelected { prefix = "S " }
	return prefix + message.Provider + ": " + message.Body
}

func GetMessagePreview(message types.Message) string {
	return message.Body
}

func GetProviderDisplay(state AppState) string {
	provider := state.conversations[state.pos].provider
	return "provider: " + provider
}

func UpdateStateFromKeyBind(state AppState, key rune) AppState {
	switch {
		case key == BIND_KEY_TOP || key == BIND_KEY_BOTTOM:
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

			jumpBy := state.jumpBy
			if jumpBy == -1 { jumpBy = 1 }
			state.jumpBy = -1

			var fn func(int) int
			if key == BIND_KEY_DOWN { fn = MakeIncBy(maxMsgs, jumpBy) } else
			{ fn = MakeDescBy(maxMsgs, jumpBy) }

			state = UpdateStateMessagePosFn(state, fn)

			break
		case key == BIND_KEY_CHAT:
			state.focusInput = true
			break
		case key == BIND_KEY_QUIT:
			state.quit = true
			break
		case unicode.IsDigit(key):
			var jumpBy int
			numb := int(key - '0')

			if state.jumpBy == -1 { jumpBy = numb } else
			{ jumpBy = state.jumpBy * 10 + numb }

			state.jumpBy = jumpBy
			break
		case key == BIND_KEY_SELECT:
			msg, err := GetStateMessage(state)
			if err != nil { return state }

			id := msg.Provider + msg.Id
			state = UpdateStateSelectedToggle(state, id)
			break
	}

	msg, err := GetStateMessage(state)
	if err == nil { state = UpdateStateProvider(state, msg.Provider) }

	return state
}

func UpdateMessagesStyle(messages MessagesComponent, state AppState) {
	isFocus := state.focusInput
	colorBackground := GetFocusBackgroundColor(isFocus)
	messages.SetSelectedBackgroundColor(colorBackground)
	// colorForeground := GetFocusForegroundColor(isFocus)
	// messages.SetSelectedTextColor(colorForeground)
}

func UpdateConversationsStyle(conversations ConversationsComponent, state AppState) {
	isFocus := state.focusInput
	colorBackground := GetFocusBackgroundColor(isFocus)
	conversations.SetSelectedBackgroundColor(colorBackground)
	// colorForeground := GetFocusForegroundColor(isFocus)
	// conversations.SetSelectedTextColor(colorForeground)
}

func GetFocusBackgroundColor(focusInput bool) tcell.Color {
	if focusInput { return FOCUS_BACKGROUND_INSERT }
	return FOCUS_BACKGROUND_NORMAL
}

func GetFocusForegroundColor(focusInput bool) tcell.Color {
	if focusInput { return FOCUS_FOREGROUND_INSERT }
	return FOCUS_FOREGROUND_NORMAL
}

type AutoCompleteFn = func(draft string) (entries []string)
func MakeAutoCompleteFn() AutoCompleteFn {
	return func(draft string) (entries []string) {
		return
	}
}
