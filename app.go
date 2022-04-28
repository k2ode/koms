package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeInitialState() AppState {
	return MakeEmptyState()
}

type UpdateCacheFn = func()

func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }

	state := MakeInitialState()

	conversations, conversationsUpdate := MakeConversations(client, state)
	messages,      messagesUpdate      := MakeMessages(client, state)
	preview,       previewUpdate       := MakePreview(state)
	input,         inputUpdate         := MakeInput(state)

	container                          := MakeContainer(conversations, messages, input, preview)

	update := func(newState AppState) {
		conversationsUpdate(newState)
		messagesUpdate(newState)
		previewUpdate(newState)
		
		if newState.focusInput { app.SetFocus(input) } else
		{ inputUpdate(newState) }
	}

	updateCache := MakeUpdateCacheFn(client, &state, update)

	onInputEscape := func(draft string) {
		state.focusInput = false
		state = UpdateStateDraft(state, draft)
		messagesUpdate(state)
		app.SetFocus(messages)
	}

	onInputEnter := MakeOnInputEnter(client, &state, updateCache)

	input.SetDoneFunc(
		MakeInputDoneFn(input, onInputEscape, onInputEnter),
	)

	onKeyDown := MakeOnKeyDown(app, &state, update)
	messages.SetInputCapture(onKeyDown)


	app.SetRoot(container, true)

	updateCache()


	if err := app.Run(); err != nil {
		panic(err)
	}
}

func MakeUpdateCacheFn(client Client, state *AppState, update func(state AppState)) func() {
	return func() {
		convos, err := client.GetConversations()
		if err != nil { panic(err) }

		state.cache.conversations = convos

		for i, convo := range convos {
			msgs, _ := client.GetConversationMessages(convo)
			state.cache.messages[i] = msgs

			convoState := state.conversations[i]
			convoState.messagePos = len(msgs) - 1
			state.conversations[i] = convoState
		}

		update(*state)
	}
}

type OnInputEnterFn = func(string)
func MakeOnInputEnter(client Client, state *AppState, updateCache UpdateCacheFn) OnInputEnterFn {
	return func(message string) {
		if message == "" { return }

		convo := GetCacheConversation(*state)
		// convoState := GetStateConversation(state)
		providerIds := []string{ "a" }
		err := client.SendMessage(convo, message, providerIds)
		if err != nil { panic(err) }

		msgs, _ := GetCacheMessages(*state)
		newMsgPos := len(msgs)
		*state = UpdateStateMessagePos(*state, newMsgPos)

		updateCache()
	}

}

type OnKeyDownFn = func(*tcell.EventKey) *tcell.EventKey
func MakeOnKeyDown(app *tview.Application, state *AppState, update func(AppState)) OnKeyDownFn {
	return func(event *tcell.EventKey) *tcell.EventKey {
		newState := UpdateStateFromKeyBind(*state, event.Rune())
		if newState.quit { app.Stop() }
		update(newState)
		*state = newState
		return nil
	}
}
