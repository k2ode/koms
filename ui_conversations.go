package main

import (
	"github.com/rivo/tview"
)

type ComponentConversations = *tview.List

func MakeConversations(client Client, state AppState) (ComponentConversations, UpdateStateFn) {
	conversations := tview.NewList()
	UpdateConversationsStyle(conversations, state)

	updateConversations := MakeConversationsUpdateFn(client, conversations)

	updateConversations(state)

	return conversations, updateConversations
}

func MakeConversationsUpdateFn(client Client, component ComponentConversations) UpdateStateFn {
	return func(state AppState) {
		UpdateConversationsStyle(component, state)

		conversations := GetStateConversations(state)

		component.Clear()
		for _, conversation := range conversations {
			label := ParseConversation(client, conversation)
			component.AddItem(label, "", 0, nil)
		}
		component.SetCurrentItem(state.pos)
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