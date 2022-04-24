package main

import (
	"github.com/rivo/tview"
)

type UpdateStateFn = func(AppState)

func MakeContainer() *tview.Grid {
	containerRows := []int{ROWS_CONTENT, ROWS_INPUT}
	containerColumns := []int{COLUMNS_CONVERSATIONS, COLUMNS_MESSAGES, COLUMNS_PREVIEW}

	container := tview.NewGrid().
		SetRows(containerRows...).
		SetColumns(containerColumns...).
		SetBorders(true)

	return container
}