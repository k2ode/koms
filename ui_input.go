package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func MakeInput(state AppState, handleEscape func(string), handleEnter func(string)) (*tview.InputField, func(AppState)) {
	input := tview.NewInputField()

	doneFunc := MakeInputDoneFunc(input, handleEscape, handleEnter)
	input.SetDoneFunc(doneFunc)

	updateInput := MakeInputUpdateFunc(input)
	updateInput(state)

	return input, updateInput
}

func MakeInputDoneFunc(input *tview.InputField, handleEscape func(string), handleEnter func(string)) func(tcell.Key) {
	return func(key tcell.Key) {
		text := input.GetText()
		if key == tcell.KeyEscape { handleEscape(text) }
		if key == tcell.KeyEnter { input.SetText(""); handleEnter(text) }
	}
}

func MakeInputUpdateFunc(input *tview.InputField) func(AppState) {
	return func(state AppState) {
		draft := state.drafts[state.conversationPos]
		input.SetText(draft)
	}
}

func AddContainerInput(container *tview.Grid, input *tview.InputField) {
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