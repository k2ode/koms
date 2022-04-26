package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeInitialState() AppState {
	return MakeEmptyState()
}

func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }

	state := MakeInitialState()


	onInputEnter := func(string) {

	}


	conversations, conversationsUpdate := MakeConversations(client, state)
	messages,      messagesUpdate      := MakeMessages(client, state)
	input,         inputUpdate         := MakeInput(state)

	container                          := MakeContainer(conversations, messages, input)


	onInputEscape := func(draft string) {
		state.focusInput = false

		app.SetFocus(messages)

		state = UpdateStateDraft(state, draft)
	}
	onDone := MakeInputDoneFn(input, onInputEscape, onInputEnter)
	input.SetDoneFunc(onDone)

	
	onKeyDown := func(event *tcell.EventKey) *tcell.EventKey {
		newState := UpdateStateFromKeyBind(state, event.Rune())

		conversationsUpdate(newState)
		messagesUpdate(newState)
		
		if newState.focusInput { app.SetFocus(input) } else
		{ inputUpdate(state) }

		state = newState

		return nil
	}
	messages.SetInputCapture(onKeyDown)

	onLoad := func() {
		convos, err := client.GetConversations()
		if err != nil { panic(err) }

		state.cache.conversations = convos

		conversationsUpdate(state)

		msgs, _ := client.GetConversationMessages(convos[0])
		state.cache.messages[0] = msgs
		msgs1, _ := client.GetConversationMessages(convos[1])
		state.cache.messages[1] = msgs1

		messagesUpdate(state)
	}

	app.SetRoot(container, true)

	onLoad()


	if err := app.Run(); err != nil {
		panic(err)
	}
}