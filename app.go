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



	conversations, conversationsUpdate := MakeConversations(client, state)
	messages,      messagesUpdate      := MakeMessages(client, state)
	input,         inputUpdate         := MakeInput(state)

	container                          := MakeContainer(conversations, messages, input)

	update := func(newState AppState) {
		conversationsUpdate(newState)
		messagesUpdate(newState)
		
		if newState.focusInput { app.SetFocus(input) } else
		{ inputUpdate(newState) }
	}

	onLoad := func() {
		convos, err := client.GetConversations()
		if err != nil { panic(err) }

		state.cache.conversations = convos

		msgs, _ := client.GetConversationMessages(convos[0])
		state.cache.messages[0] = msgs
		msgs1, _ := client.GetConversationMessages(convos[1])
		state.cache.messages[1] = msgs1

		update(state)
	}

	onInputEscape := func(draft string) {
		state.focusInput = false

		app.SetFocus(messages)

		state = UpdateStateDraft(state, draft)
	}

	onInputEnter := func(message string) {
		if message == "" { return }

		convo := GetCacheConversation(state)
		// convoState := GetStateConversation(state)
		providerIds := []string{ "a" }
		err := client.SendMessage(convo, message, providerIds)
		if err != nil { panic(err) }

		msgs, _ := GetCacheMessages(state)
		newMsgPos := len(msgs)
		state = UpdateStateMessagePos(state, newMsgPos)

		onLoad()
	}

	onDone := MakeInputDoneFn(input, onInputEscape, onInputEnter)
	input.SetDoneFunc(onDone)

	
	onKeyDown := func(event *tcell.EventKey) *tcell.EventKey {
		newState := UpdateStateFromKeyBind(state, event.Rune())
		update(newState)
		state = newState
		return nil
	}
	messages.SetInputCapture(onKeyDown)


	app.SetRoot(container, true)

	onLoad()


	if err := app.Run(); err != nil {
		panic(err)
	}
}