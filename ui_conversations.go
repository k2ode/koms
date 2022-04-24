package main

import "github.com/rivo/tview"

type ConversationsComponent = *tview.List

func MakeConversations(client Client, state AppState) (ConversationsComponent, UpdateStateFn) {
	conversations := tview.NewList()

	updateConversations := MakeConversationsUpdateFn(client, conversations)

	updateConversations(state)

	return conversations, updateConversations
}

func MakeConversationsUpdateFn(client Client, conversations ConversationsComponent) UpdateStateFn {
	return func(state AppState) {
		for _, conversation := range state.conversations {
			label := ParseConversation(client, conversation)
			conversations.AddItem(label, "", 0, nil)
		}
		conversations.SetCurrentItem(state.conversationPos)
	}
}

func AddContainerConversations(container *tview.Grid, conversations *tview.List) {
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
}