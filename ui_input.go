package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InputComponent = *tview.InputField

func MakeInput(state AppState) (InputComponent, UpdateStateFn) {
	input := tview.NewInputField()
	input.SetFieldBackgroundColor(0x000000)

	updateInput := MakeInputUpdateFn(input)
	updateInput(state)

	return input, updateInput
}

func MakeInputDoneFn(input *tview.InputField, handleEscape func(string), handleEnter func(string)) func(tcell.Key) {
	return func(key tcell.Key) {
		text := input.GetText()
		if key == tcell.KeyEscape { handleEscape(text) }
		if key == tcell.KeyEnter { input.SetText(""); handleEnter(text) }
	}
}

func MakeInputUpdateFn(input *tview.InputField) UpdateStateFn {
	return func(state AppState) {
		draft := GetStateDraft(state)
		input.SetText(draft)
	}
}

func AddContainerInput(container *tview.Grid, input InputComponent) {
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