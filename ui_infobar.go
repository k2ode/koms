package main

import "github.com/rivo/tview"

type ComponentInfobar = *tview.TextView

func MakeInfobar(state AppState) (ComponentInfobar, UpdateStateFn) {
	display := tview.NewTextView()
	updateInfobar := MakeUpdateFnInfoBar(display)
	updateInfobar(state)
	return display, updateInfobar
}

func MakeUpdateFnInfoBar(display ComponentInfobar) UpdateStateFn {
	return func(state AppState) {
		provider := GetInfobar(state)
		display.SetText(provider)
	}
}