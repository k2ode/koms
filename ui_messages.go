package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MessagesComponent = *tview.List

func MakeMessages(client Client, state AppState) (MessagesComponent, UpdateStateFn) {
	messages := tview.NewList()
	UpdateMessageStyle(messages, state)

	updateMessages := MakeMessagesUpdateFn(client, messages)

	updateMessages(state)

	return messages, updateMessages
}

func UpdateMessageStyle(messages MessagesComponent, state AppState) {
	isFocus := state.focusInput
	colorBackground := GetMessageFocusBackgroundColor(isFocus)
	messages.SetSelectedBackgroundColor(colorBackground)
	colorForeground := GetMessageFocusForegroundColor(isFocus)
	messages.SetSelectedTextColor(colorForeground)
}

func GetMessageFocusBackgroundColor(focusInput bool) tcell.Color {
	if focusInput { return MESSAGE_FOCUS_BACKGROUND_INSERT }
	return MESSAGE_FOCUS_BACKGROUND_NORMAL 
}

func GetMessageFocusForegroundColor(focusInput bool) tcell.Color {
	if focusInput { return MESSAGE_FOCUS_FOREGROUND_INSERT }
	return MESSAGE_FOCUS_FOREGROUND_NORMAL 
}

func MakeMessagesUpdateFn(client Client, messages MessagesComponent) UpdateStateFn {
	return func(state AppState) {
		UpdateMessageStyle(messages, state)
		messages.Clear()

		conversationMessages, exists := GetCacheMessages(state)
		if !exists { return }

		for _, message := range conversationMessages {
			parsedMessage := ParseMessage(client, message)
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