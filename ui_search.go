package main

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// import "github.com/gdamore/tcell/v2"

type ComponentSearch = *tview.Grid
type ComponentSearchInput = *tview.InputField
type ComponentSearchParticipants = *tview.List
type ComponentSearchContainer = *tview.Grid

const wordList = "Joe,Mama,Jamba"


func MakeSearch(app *tview.Application, state *AppState) (ComponentSearch, UpdateStateFn) {
	participants, updateParticipants := MakeParticipantsList()
	searchInput, clear               := MakeSearchInput()


	searchContainer := MakeSearchContainer(participants, searchInput)


	search := Modal(searchContainer, 40, 0)


	update := func(state AppState) {
		participantLen := len(state.search.filters)
		height := participantLen + 3
		search.SetRows(0, height, 0)
		updateParticipants(state)
	}


	onEscape := func(s string) {
		app.SetFocus(participants)
	}

	onEnter := func(text string) {
		if text == "" { return }
		clear()
		state.search.filters = append(state.search.filters, SearchQueryFilter{ text })
		update(*state)
	}

	doneFn := MakeInputDoneFn(searchInput, onEscape, onEnter)
	searchInput.SetDoneFunc(doneFn)

	autocomplete := MakeInputAutocompleteFn()
	searchInput.SetAutocompleteFunc(autocomplete)


	onKeyDownFn := MakeSearchOnKeyDown(state, update)
	participants.SetInputCapture(onKeyDownFn)

	update(*state)

	return search, update
}

type ClearFn = func()
func MakeSearchInput() (ComponentSearchInput, ClearFn) {
	searchInput := tview.NewInputField().SetPlaceholder("search...")

	clearFn := func() {
		searchInput.SetText("")
	}
	return searchInput, clearFn
}

func MakeParticipantsList() (ComponentSearchParticipants, UpdateFn) {
	participants := tview.NewList().ShowSecondaryText(false)

	update := func(state AppState) {
		participants.Clear()
		for i, participant := range state.search.filters {
			itemText := strconv.Itoa(i + 1) + ") " + participant.name
			participants.AddItem(itemText, "", 0, func() {})
		}
		participants.SetCurrentItem(state.search.filterPos)
	}

	return participants, update
}

func MakeSearchContainer(participants ComponentSearchParticipants, searchInput ComponentSearchInput) ComponentSearchContainer {

	searchContainer := tview.NewGrid().SetColumns(0).SetRows(0, 1)
	searchContainer.SetBorder(true).SetTitle("search")


	searchContainer.AddItem(participants, 0, 0, 1, 1, 0, 0, false)
	searchContainer.AddItem(searchInput, 1, 0, 1, 1, 0, 0, true)

	return searchContainer
}

func MakeSearchOnKeyDown(state *AppState, update UpdateFn) OnKeyDownFn {
	return func(event *tcell.EventKey) *tcell.EventKey {
		newState := UpdateStateSearchFromKeyBind(*state, event.Rune())
		update(newState)
		*state = newState
		return nil
	}
}
