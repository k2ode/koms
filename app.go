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

func MakeInput(state AppState) *tview.InputField {
	conversationDraft := state.drafts[state.conversationPos]
	input := tview.NewInputField().
		SetText(conversationDraft)

	return input
}

type AppState struct {
	conversations   []Conversation
	messages        map[int][]Message
	drafts          map[int]string
	conversationPos int
	messagePos      int
}

func MakeMessages(state AppState, handleKeyDown func(e *tcell.EventKey) *tcell.EventKey) *tview.List {
	list := tview.NewList()
	list.SetInputCapture(handleKeyDown)

	messages, exists := state.messages[state.conversationPos]
	if !exists { list.AddItem("NO MESSAGES IN STATE", "", 0, nil); return list }

	for _, message := range messages {
		list.AddItem(message.body, "", 0, nil)
	}

	list.SetCurrentItem(state.messagePos)


	return list
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

	input := MakeInput(state)
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

		updateState(newState)

		return nil
	}

	messages := MakeMessages(state, handleKeyDown)
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