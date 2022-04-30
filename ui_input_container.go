package main

import "github.com/rivo/tview"

type ComponentInputContainer = *tview.Grid

func MakeContainerInput(state AppState, input InputComponent, provider ComponentProviderDisplay) (ComponentInputContainer, UpdateStateFn) {
	container := tview.NewGrid().
		SetRows(0).
		SetColumns(COLUMNS_CONVERSATIONS, COLUMNS_MESSAGES)

	container.AddItem(input, 0, 1, 1, 1, 0, 0, false)
	container.AddItem(provider, 0, 0, 1, 1, 0, 0, false)

	return container, nil
}

func AddContainerInputContainer(container *tview.Grid, input ComponentInputContainer) {
	container.AddItem(
		input,
		ROW_POS_INPUT,
		COLUMN_POS_INPUT,
		ROW_SPAN_INPUT,
		COLUMN_SPAN_INPUT,
		HEIGHT_MIN_INPUT,
		WIDTH_MIN_INPUT,
		false,
	)
}