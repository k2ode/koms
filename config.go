package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/k2on/koms/types"
	"github.com/rivo/tview"
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
const BIND_KEY_SEARCH = 's'
const BIND_KEY_NEXT   = '.'
const BIND_KEY_PREV   = ','

const BIND_KEY_SEARCH_CLOSE = RUNE_ESCAPE
const BIND_KEY_SEARCH_FILTER_REMOVE = 'D'
const BIND_KEY_SEARCH_FILTER_CLEAR = 'C'
const BIND_KEY_SEARCH_FOCUS = '/'


// Text

const INPUT_PLACEHOLDER_SEARCH = "search..."
const SEARCH_TITLE = "search"

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

	// imess, _ := NewProviderIMessage()

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

	result += " - " + GetLastActivity(conversation).Local().String()

	return result
}

func ParseMessage(client Client, conversation ConversationState, message types.Message) string {
	var prefix string
	id := message.Provider + message.Raw.Id
	isSelected := Contains(conversation.selected, id)
	if isSelected { prefix = "S " }
	return prefix + message.Provider + ": " + message.Raw.Body
}

type Size struct {
	x int
	y int
	width int
	height int
}

func GetMessagePreview(state AppState, size Size) string {
	msg, err := GetStateMessage(state)
	if err != nil { return err.Error() }
	return msg.Raw.Body
}

func GetInfobar(state AppState) string {
	var infobar string

	provider := state.conversations[state.pos].provider
	infobar += "provider: " + provider

	filterLen := len(state.search.filters)
	if filterLen > 0 { infobar += fmt.Sprintf(", %s filter(s)", strconv.Itoa(filterLen)) }

	return infobar
}

func GetFilter(filter SearchQueryFilter, pos int) string {
	return strconv.Itoa(pos + 1) + ") " + filter.name
}

func UpdateStateFromKeyBind(state AppState, key rune) AppState {
	getSize := func() int { 
		msgs, exists := GetCacheMessages(state)
		if !exists { return -1 }
		return len(msgs) - 1
	}
	setPosition := func(state AppState, fn IntMod) AppState {
		return UpdateStateMessagePosFn(state, fn)
	}
	fallback := func(state AppState, key rune) AppState {
		switch {
			case key == BIND_KEY_LEFT || key == BIND_KEY_RIGHT:
				maxConvos := GetStateConversationsLen(state) - 1

				var fn func(int) int
				if key == BIND_KEY_RIGHT { fn = MakeInc(maxConvos) } else
				{ fn = MakeDesc(maxConvos) }

				state.pos = fn(state.pos)

				break
			case key == BIND_KEY_NEXT || key == BIND_KEY_PREV:
				msg, err := GetStateMessage(state)
				if err != nil { break }
				lenImg := len(msg.Raw.Images)
				if lenImg == 0 { break }

				var fn func(int) int
				if key == BIND_KEY_NEXT { fn = MakeInc(lenImg - 1) } else
				{ fn = MakeDesc(lenImg - 1) }

				state = UpdateStateCarouselSelectedImage(state, fn)
				break
			case key == BIND_KEY_CHAT:
				state.focusInput = true
				break
			case key == BIND_KEY_QUIT:
				state.quit = true
				break
			case key == BIND_KEY_SELECT:
				msg, err := GetStateMessage(state)
				if err != nil { return state }

				id := msg.Provider + msg.Raw.Id
				state = UpdateStateSelectedToggle(state, id)
				break
			case key == BIND_KEY_SEARCH:
				state = UpdateStateSearchOpen(state)
				break
		}
		return state
	}
	state = VerticleListKeyBinds(state, key, getSize, setPosition, fallback)

	msg, err := GetStateMessage(state)
	if err == nil { state = UpdateStateProvider(state, msg.Provider) }

	return state
}

func UpdateStateSearchFromKeyBind(state AppState, key rune) AppState {
	getSize := func() int {
		return len(state.search.filters) - 1
	}
	setPosition := func(state AppState, fn IntMod) AppState {
		return UpdateStateSearchFilterPosFn(state, fn)
	}
	fallback := func(state AppState, key rune) AppState {
		switch key {
			case BIND_KEY_SEARCH_FILTER_REMOVE:
				filters := state.search.filters
				pos := state.search.filterPos

				filtersUpdated := RemoveSearchQueryFilter(filters, pos)
				state = UpdateStateSearchFilters(state, filtersUpdated)

				filtersUpdatedLen := len(filtersUpdated)
				if filtersUpdatedLen == pos { state = UpdateStateSearchFilterPos(state, pos - 1) }
				if filtersUpdatedLen == 0 { state = UpdateStateSearchFocus(state) }
				break
			case BIND_KEY_SEARCH_FILTER_CLEAR:
				state = UpdateStateSearchFiltersClear(state)
				state = UpdateStateSearchFocus(state)
				state = UpdateStateSearchFilterPos(state, 0)
				break
			case BIND_KEY_SEARCH_CLOSE:
				state = UpdateStateSearchClose(state)
				break
			case BIND_KEY_SEARCH_FOCUS:
 				state = UpdateStateSearchFocus(state) 
				break
		}

		return state
	}
	state = VerticleListKeyBinds(state, key, getSize, setPosition, fallback)
	
	return state
}

func UpdateMessagesStyle(messages MessagesComponent, state AppState) {
	isFocus := state.focusInput
	UpdateListStyle(messages, isFocus)
}

func UpdateConversationsStyle(conversations ComponentConversations, state AppState) {
	isFocus := state.focusInput
	UpdateListStyle(conversations, isFocus)
}

func UpdateSearchFiltersStyle(filters ComponentSearchFilters, state AppState) {
	isFocus := state.search.focusInput
	UpdateListStyle(filters, isFocus)
}

func UpdateListStyle(component *tview.List, isFocus bool) {
	colorBackground := GetFocusBackgroundColor(isFocus)
	component.SetSelectedBackgroundColor(colorBackground)
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
