package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppState struct {
	conversations   []Conversation
	messages        map[int][]Message
	drafts          map[int]string
	conversationPos int
	messagePos      int
}

func render(app *tview.Application, state AppState, client Client) {
	container := MakeContainer()

	updateState := func(s AppState) {
		render(app, s, client)
	}


	conversations, _ := MakeConversations(client, state)
	AddContainerConversations(container, conversations)

	preview, _ := MakePreview(state)
	AddContainerPreview(container, preview)

	messages, updateMessages := MakeMessages(client, state)

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

	input, _ := MakeInput(state, handleEscape, handleEnter)
	AddContainerInput(container, input)
	
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

	AddContainerMessages(container, messages)




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