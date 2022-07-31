package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ComponentInput = *tview.InputField

func MakeInput(state AppState) (ComponentInput, UpdateStateFn) {
	input := tview.NewInputField()
	input.SetFieldBackgroundColor(0x000000)

	updateInput := MakeInputUpdateFn(input)
	updateInput(state)

	return input, updateInput
}

type InputDoneFn = func(tcell.Key)
type HandleInputFn = func(string)
func MakeInputDoneFn(input *tview.InputField, handleEscape, handleEnter HandleInputFn) InputDoneFn {
	return func(key tcell.Key) {
		text := input.GetText()
		if key == tcell.KeyEscape { handleEscape(text) }
		if key == tcell.KeyEnter { input.SetText(""); handleEnter(text) }
	}
}

type AutocompleteFn = func(currentText string) (entries []string)
func MakeInputAutocompleteFn() AutocompleteFn {
	return func(draft string) (entries []string) {
		d := strings.ToLower(draft)
		r := []string{}
		if d == "" { return r }
		for _, name := range strings.Split(wordList, ",") {
			if strings.HasPrefix(strings.ToLower(name), d) { r = append(r, name)}
		}
		return r
	}
}

func MakeInputUpdateFn(input *tview.InputField) UpdateStateFn {
	return func(state AppState) {
		draft := GetStateDraft(state)
		input.SetText(draft)
	}
}

func AddContainerInput(container *tview.Grid, input ComponentInput) {
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