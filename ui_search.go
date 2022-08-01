package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// import "github.com/gdamore/tcell/v2"

type ComponentSearch = *tview.Grid
type ComponentSearchInput = *tview.InputField
type ComponentSearchFilters = *tview.List
type ComponentSearchContainer = *tview.Grid

const wordList = "Joe,Mama,Jamba"


func MakeSearch(app *tview.Application, state *AppState, updateParent UpdateStateFn) (ComponentSearch, UpdateStateFn) {
	filters, updateFilters := MakeFiltersList()
	searchInput, clear     := MakeSearchInput()


	searchContainer := MakeSearchContainer(filters, searchInput)


	search := Modal(searchContainer, 40, 0)

	update := func(newState AppState) {
		// outdatedResult := len(state.search.filters) != len(newState.search.filters)
		outdatedResult := true
		if outdatedResult {
			*state = newState

			newState.search.result = Filter(newState)
		}

		filtersCount := len(newState.search.filters)
		height := filtersCount + 3
		search.SetRows(0, height, 0)
		updateFilters(newState)

		if newState.search.focusInput {
			app.SetFocus(searchInput)
		}

		updateParent(newState)
	}


	onEscape := func(s string) {
		if len(state.search.filters) > 0 { 
			state.search.focusInput = false
			update(*state)
			app.SetFocus(filters)
			return
		}
		*state = UpdateStateSearchClose(*state)
		update(*state)
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
	filters.SetInputCapture(onKeyDownFn)

	update(*state)

	return search, update
}

type ClearFn = func()
func MakeSearchInput() (ComponentSearchInput, ClearFn) {
	searchInput := tview.NewInputField().SetPlaceholder(INPUT_PLACEHOLDER_SEARCH)

	clearFn := func() {
		searchInput.SetText("")
	}
	return searchInput, clearFn
}

func MakeFiltersList() (ComponentSearchFilters, UpdateFn) {
	filters := tview.NewList().ShowSecondaryText(false)

	update := func(state AppState) {
		UpdateSearchFiltersStyle(filters, state)

		filters.Clear()
		for i, filter := range state.search.filters {
			itemText := GetFilter(filter, i)
			filters.AddItem(itemText, "", 0, func() {})
		}
		filters.SetCurrentItem(state.search.filterPos)
	}

	return filters, update
}

func MakeSearchContainer(filters ComponentSearchFilters, searchInput ComponentSearchInput) ComponentSearchContainer {

	searchContainer := tview.NewGrid().SetColumns(0).SetRows(0, 1)
	searchContainer.SetBorder(true).SetTitle(SEARCH_TITLE)


	searchContainer.AddItem(filters, 0, 0, 1, 1, 0, 0, false)
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
