package main

import (
	"github.com/rivo/tview"
)

type MessagesComponent = *tview.List

func MakeMessages(client Client, state AppState) (MessagesComponent, UpdateStateFn) {
	messages := tview.NewList()
	UpdateMessagesStyle(messages, state)

	updateMessages := MakeMessagesUpdateFn(client, messages)

	updateMessages(state)

	return messages, updateMessages
}

func MakeMessagesUpdateFn(client Client, messages MessagesComponent) UpdateStateFn {
	return func(state AppState) {
		UpdateMessagesStyle(messages, state)
		messages.Clear()

		conversationMessages, exists := GetCacheMessages(state)
		if !exists { return }

		messagesFiltered := FilterMessages(conversationMessages, state.search.filters)

		conversation := GetStateConversation(state)
		for _, message := range messagesFiltered {
			parsedMessage := ParseMessage(client, conversation, message)
			messages.AddItem(parsedMessage, "", 0, nil)
		}

		messagePos := GetStateMessagePos(state)
		messages.SetCurrentItem(messagePos)
	}
}

func AddContainerMessages(container *tview.Grid, messages *tview.List) {
	isFocused := true
	container.AddItem(
		messages,
		ROW_POS_MSGS,
		COLUMN_POS_MSGS,
		ROW_SPAN_MSGS,
		COLUMN_SPAN_MSGS,
		HEIGHT_MIN_MSGS,
		WIDTH_MIN_MSGS,
		isFocused,
	)
}
