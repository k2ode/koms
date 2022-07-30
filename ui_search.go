package main

import (
	"strconv"

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


	update := func(state AppState) {
		filtersCount := len(state.search.filters)
		height := filtersCount + 3
		search.SetRows(0, height, 0)
		updateFilters(state)

		if state.search.focusInput {
			app.SetFocus(searchInput)
		}

		if !state.search.open {
			updateParent(state)
		}
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
	searchInput := tview.NewInputField().SetPlaceholder("search...")

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
			itemText := strconv.Itoa(i + 1) + ") " + filter.name
			filters.AddItem(itemText, "", 0, func() {})
		}
		filters.SetCurrentItem(state.search.filterPos)
	}

	return filters, update
}

func MakeSearchContainer(filters ComponentSearchFilters, searchInput ComponentSearchInput) ComponentSearchContainer {

	searchContainer := tview.NewGrid().SetColumns(0).SetRows(0, 1)
	searchContainer.SetBorder(true).SetTitle("search")


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
