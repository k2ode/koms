package main

import (
	"github.com/rivo/tview"
)

type UpdateStateFn = func(AppState)

func MakeContainer(conversations ConversationsComponent, messages MessagesComponent, containerInput ComponentInputContainer, preview PreviewComponent) *tview.Grid {
	containerRows := []int{ROWS_CONTENT, ROWS_INPUT}
	containerColumns := []int{COLUMNS_CONVERSATIONS, COLUMNS_MESSAGES, COLUMNS_PREVIEW}

	container := tview.NewGrid().
		SetRows(containerRows...).
		SetColumns(containerColumns...).
		SetBorders(true)

	AddContainerConversations(container, conversations)
	AddContainerMessages(container, messages)
	AddContainerInputContainer(container, containerInput)
	AddContainerPreview(container, preview)

	return container
}