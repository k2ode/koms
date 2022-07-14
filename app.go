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

	providerDisplay, providerUpdate := MakeProviderDisplay(state)
	inputContainer, _               := MakeContainerInput(state, input, providerDisplay)

	container                          := MakeContainer(conversations, messages, inputContainer, preview)

	pages := tview.NewPages()

	UIUpdate := func(newState AppState) {
		conversationsUpdate(newState)
		messagesUpdate(newState)
		previewUpdate(newState)
		providerUpdate(newState)

		if newState.search.open { pages.ShowPage("search") }

		if newState.focusInput { app.SetFocus(input) } else
		{ inputUpdate(newState) }
	}

	cacheAndUIUpdate := MakeUpdateCacheFn(client, &state, UIUpdate)

	onInputEscape := func(draft string) {
		state.focusInput = false
		state = UpdateStateDraft(state, draft)
		messagesUpdate(state)
		conversationsUpdate(state)
		app.SetFocus(messages)
	}

	onInputEnter := MakeOnInputEnter(client, &state, cacheAndUIUpdate)

	doneFn := MakeInputDoneFn(input, onInputEscape, onInputEnter)
	input.SetDoneFunc(doneFn)

	autoCompleteFn := MakeAutoCompleteFn()
	input.SetAutocompleteFunc(autoCompleteFn)

	onKeyDown := MakeOnKeyDown(app, &state, UIUpdate)
	messages.SetInputCapture(onKeyDown)

	search, _ := MakeSearch(app, &state)

	pages.
		AddPage("container", container, true, true).
		AddPage("search", search, true, false)


	app.SetRoot(pages, true)

	cacheAndUIUpdate()


	if err := app.Run(); err != nil {
		panic(err)
	}
}

func MakeUpdateCacheFn(client Client, state *AppState, updateUI func(state AppState)) UpdateCacheFn {
	return func() {
		convos, err := client.GetConversations()
		if err != nil { panic(err) }

		state.cache.conversations = convos

		for i, convo := range convos {
			msgs, _ := client.GetConversationMessages(convo)
			state.cache.messages[i] = msgs

			convoState := state.conversations[i]
			messagePos := len(msgs) - 1
			convoState.messagePos = messagePos
			if messagePos != -1 {
				lastMsg := msgs[messagePos]
				lastMsgProvider := lastMsg.Provider
				convoState.provider = lastMsgProvider
			} else {
				// this should be the default provider for the (app > contact) w that heiarchy
				convoState.provider = "ADD DEFAULT"
			}
			state.conversations[i] = convoState
		}

		updateUI(*state)
	}
}

type OnInputEnterFn = func(string)
func MakeOnInputEnter(client Client, state *AppState, updateCache UpdateCacheFn) OnInputEnterFn {
	return func(message string) {
		if message == "" { return }

		convo := GetCacheConversation(*state)
		provider := GetStateProvider(*state)
		providerIds := []string{ provider }
		err := client.SendMessage(convo, message, providerIds)
		if err != nil { panic(err) }

		msgs, _ := GetCacheMessages(*state)
		newMsgPos := len(msgs)
		*state = UpdateStateConversationPos(*state, 0)
		*state = UpdateStateMessagePos(*state, newMsgPos)

		updateCache()
	}

}

type OnKeyDownFn = func(*tcell.EventKey) *tcell.EventKey
type UpdateFn = func(AppState)
func MakeOnKeyDown(app *tview.Application, state *AppState, update UpdateFn) OnKeyDownFn {
	return func(event *tcell.EventKey) *tcell.EventKey {
		newState := UpdateStateFromKeyBind(*state, event.Rune())
		if newState.quit { app.Stop() }
		update(newState)
		*state = newState
		return nil
	}
}
