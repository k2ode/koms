package main

import "github.com/rivo/tview"

type ComponentProviderDisplay = *tview.TextView

func MakeProviderDisplay(state AppState) (ComponentProviderDisplay, UpdateStateFn) {
	display := tview.NewTextView()
	updateProviderDisplay := MakeProviderDisplayUpdateFn(display)
	updateProviderDisplay(state)
	return display, updateProviderDisplay
}


func MakeProviderDisplayUpdateFn(display ComponentProviderDisplay) UpdateStateFn {
	return func(state AppState) {
		provider := GetProviderDisplay(state)
		display.SetText(provider)
	}
}