package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


func MakeContainer() *tview.Grid {
	containerRows := []int{ROWS_CONTENT, ROWS_INPUT}
	containerColumns := []int{COLUMNS_CONVERSATIONS, COLUMNS_MESSAGES, COLUMNS_PREVIEW}

	container := tview.NewGrid().
		SetRows(containerRows...).
		SetColumns(containerColumns...).
		SetBorders(true)

	return container
}

func MakeInput(state AppState, handleEscape func(string), handleEnter func(string)) *tview.InputField {
	conversationDraft := state.drafts[state.conversationPos]

	input := tview.NewInputField()

	doneFunc := func(key tcell.Key) {
		text := input.GetText()
		if key == tcell.KeyEscape { handleEscape(text) }
		if key == tcell.KeyEnter { input.SetText(""); handleEnter(text) }
	}

	input.SetText(conversationDraft).
		SetDoneFunc(doneFunc)

	return input
}

type AppState struct {
	conversations   []Conversation
	messages        map[int][]Message
	drafts          map[int]string
	conversationPos int
	messagePos      int
}

func MakeMessages(state AppState) (*tview.List, func(s AppState)) {
	list := tview.NewList()

	updateMessages := func(s AppState) {
		messages, exists := s.messages[s.conversationPos]
		if !exists { return }

		list.Clear()

		for _, message := range messages {
			list.AddItem(message.body, "", 0, nil)
		}
		list.SetCurrentItem(state.messagePos)
	}

	updateMessages(state)

	return list, updateMessages
}

func MakeConversations(state AppState, client Client) *tview.List {
	list := tview.NewList()

	for _, convo := range state.conversations {
		label := ParseConversation(client, convo)
		list.AddItem(label, "", 0, nil)
	}

	list.SetCurrentItem(state.conversationPos)

	return list
}

func render(app *tview.Application, state AppState, client Client) {
	container := MakeContainer()

	updateState := func(s AppState) {
		render(app, s, client)
	}




	messages, updateMessages := MakeMessages(state)

	setDraft := func(draft string) { state.drafts[state.conversationPos] = draft }

	handleEnter := func(message string) {
		setDraft("")

		client.SendMessage(state.conversations[state.conversationPos], message, []string{"a"})

		newState := state

		newMsgs, _ := client.GetConversationMessages(state.conversations[state.conversationPos])
	
		newState.messages[state.conversationPos] = newMsgs

		updateMessages(newState)
	}

	handleEscape := func(draft string) {
		setDraft(draft)
		app.SetFocus(messages)
	}

	input := MakeInput(state, handleEscape, handleEnter)
	container.AddItem(
		input,
		ROW_POS_INPUT,
		COLUMN_POS_INPUT,
		ROW_SPAN_INPUT,
		COLUMN_SPAN_INPUT,
		HEIGHT_MIN_INPUT,
		WIDTH_MIN_INPUT,
		false,
	)
	
	handleKeyDown := func(event *tcell.EventKey) *tcell.EventKey {

		newState := state

		lenX := len(state.conversations) - 1 
		lenY := len(state.messages[state.conversationPos]) - 1 

		incX := MakeInc(lenX)
		descX := MakeDesc(lenX)

		incY := MakeInc(lenY)
		descY := MakeDesc(lenY)


		if event.Rune() == BIND_KEY_TOP { newState.messagePos = 0 }
		if event.Rune() == BIND_KEY_BOTTOM { newState.messagePos = lenY }
		if event.Rune() == BIND_KEY_DOWN { newState.messagePos = incY(state.messagePos) }
		if event.Rune() == BIND_KEY_UP { newState.messagePos = descY(state.messagePos) }
		if event.Rune() == BIND_KEY_LEFT { newState.conversationPos = descX(state.conversationPos) }
		if event.Rune() == BIND_KEY_RIGHT { newState.conversationPos = incX(state.conversationPos) }
		if event.Rune() == BIND_KEY_CHAT { app.SetFocus(input) }

		posChange := state.messagePos != newState.messagePos || state.conversationPos != newState.conversationPos
		if posChange { updateState(newState) }

		return nil
	}
	messages.SetInputCapture(handleKeyDown)

	focusMessages := true
	container.AddItem(
		messages,
		ROW_POS_MSGS,
		COLUMN_POS_MSGS,
		ROW_SPAN_MSGS,
		COLUMN_SPAN_MSGS,
		HEIGHT_MIN_MSGS,
		WIDTH_MIN_MSGS,
		focusMessages,
	)

	conversations := MakeConversations(state, client)
	container.AddItem(
		conversations,
		ROW_POS_CONVOS,
		COLUMN_POS_CONVOS,
		ROW_SPAN_CONVOS,
		COLUMN_SPAN_CONVOS,
		HEIGHT_MIN_CONVOS,
		WIDTH_MIN_CONVOS,
		false,
	)




	app.SetRoot(container, true)
}


func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }

	conversations, err := client.GetConversations()
	if err != nil { panic(err) }

	state := AppState{
		conversations: conversations,
		messages: make(map[int][]Message),
		drafts: make(map[int]string),
		conversationPos: 0,
		messagePos: 0,
	}

	msgs, _ := client.GetConversationMessages(conversations[0])
	state.messages[0] = msgs
	msgs1, _ := client.GetConversationMessages(conversations[1])
	state.messages[1] = msgs1

	render(app, state, client)

	if err := app.Run(); err != nil {
		panic(err)
	}
}