package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// func render(app *tview.Application, state AppState, client Client) {
// 	container := MakeContainer()

// 	updateState := func(s AppState) {
// 		render(app, s, client)
// 	}

// 	conversations, _ := MakeConversations(client, state)
// 	AddContainerConversations(container, conversations)

// 	preview, _ := MakePreview(state)
// 	AddContainerPreview(container, preview)

// 	messages, updateMessages := MakeMessages(client, state)

// 	setDraft := func(draft string) { state.drafts[state.conversationPos] = draft }

// 	handleEnter := func(message string) {
// 		setDraft("")

// 		client.SendMessage(state.conversations[state.conversationPos], message, []string{"a"})

// 		newState := state

// 		newMsgs, _ := client.GetConversationMessages(state.conversations[state.conversationPos])

// 		newState.messages[state.conversationPos] = newMsgs

// 		updateMessages(newState)
// 	}

// 	handleEscape := func(draft string) {
// 		setDraft(draft)
// 		app.SetFocus(messages)
// 	}

// 	input, _ := MakeInput(state, handleEscape, handleEnter)
// 	AddContainerInput(container, input)

// 	handleKeyDown := func(event *tcell.EventKey) *tcell.EventKey {

// 		newState := state

// 		lenX := len(state.conversations) - 1
// 		lenY := len(state.messages[state.conversationPos]) - 1

// 		incX := MakeInc(lenX)
// 		descX := MakeDesc(lenX)

// 		incY := MakeInc(lenY)
// 		descY := MakeDesc(lenY)

// 		if event.Rune() == BIND_KEY_TOP { newState.messagePos = 0 }
// 		if event.Rune() == BIND_KEY_BOTTOM { newState.messagePos = lenY }
// 		if event.Rune() == BIND_KEY_DOWN { newState.messagePos = incY(state.messagePos) }
// 		if event.Rune() == BIND_KEY_UP { newState.messagePos = descY(state.messagePos) }
// 		if event.Rune() == BIND_KEY_LEFT { newState.conversationPos = descX(state.conversationPos) }
// 		if event.Rune() == BIND_KEY_RIGHT { newState.conversationPos = incX(state.conversationPos) }
// 		if event.Rune() == BIND_KEY_CHAT { app.SetFocus(input) }

// 		posChange := state.messagePos != newState.messagePos || state.conversationPos != newState.conversationPos
// 		if posChange { updateState(newState) }

// 		return nil
// 	}
// 	messages.SetInputCapture(handleKeyDown)

// 	AddContainerMessages(container, messages)

// 	app.SetRoot(container, true)
// }


func MakeInitialState() AppState {
	return MakeEmptyState()
}

func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }

	state := MakeInitialState()

	onInputEscape := func(string) {}

	onInputEnter := func(string) {}


	onLoad := func() {
		
	}


	conversations, conversationsUpdate := MakeConversations(client, state)
	messages,      messagesUpdate      := MakeMessages(client, state)
	input,         inputUpdate         := MakeInput(state, onInputEscape, onInputEnter)

	container                          := MakeContainer(conversations, messages, input)


	onKeyDown := func(event *tcell.EventKey) *tcell.EventKey {

		// state.conversations[0] = ConversationState{messagePos: state.conversations[0].messagePos + 1}
		// messagesUpdate(state)
		newState := UpdateStateFromKeyBind(state, event.Rune())

		// conversationsUpdate(*state)

		messagesUpdate(newState)
		
		// if newState.focusInput { app.SetFocus(input) }

		state = newState

		return nil
		// if event.Rune() == BIND_KEY_TOP    { messagePos = 0 }
		// if event.Rune() == BIND_KEY_BOTTOM { messagePos = lenY }
		// if event.Rune() == BIND_KEY_DOWN   { messagePos = incY(messagePosStart) }
		// if event.Rune() == BIND_KEY_UP     { messagePos = descY(messagePosStart) }
		// if event.Rune() == BIND_KEY_LEFT   { pos = descX(posStart) }
		// if event.Rune() == BIND_KEY_RIGHT  { pos = incX(posStart) }
		// if event.Rune() == BIND_KEY_CHAT   { app.SetFocus(input); return nil }
	}



	messages.SetInputCapture(onKeyDown)

	app.SetRoot(container, true)


	convos, err := client.GetConversations()
	if err != nil { panic(err) }

	state.cache.conversations = convos

	conversationsUpdate(state)

	msgs, _ := client.GetConversationMessages(convos[0])
	state.cache.messages[0] = msgs
	msgs1, _ := client.GetConversationMessages(convos[1])
	state.cache.messages[1] = msgs1

	messagesUpdate(state)

	inputUpdate(state)

	onLoad()


	if err := app.Run(); err != nil {
		panic(err)
	}
}