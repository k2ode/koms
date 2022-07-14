package main

import (
	"strings"
	"time"
	"unicode"

	"github.com/k2on/koms/types"
	"github.com/rivo/tview"
)

func Lines(str string) []string {
	return strings.Split(str, "\n")
}

func Between(str string, start string, end string) string {
	return strings.Split(strings.Split(str, start)[1], end)[0]
}

type IntMod = func(int) int

func MakeIncBy(max int, by int) IntMod {
	return func(i int) int {
		if i == max { return 0 }
		if i + by > max { return max }
		return i + by
	}
}

func MakeInc(max int) IntMod {
	return MakeIncBy(max, 1)
}

func MakeDescBy(max int, by int) IntMod {
	return func(i int) int {
		if i == 0 { return max }
		if i - by < 0 { return 0 }
		return i - by
	}
}

func MakeDesc(max int) IntMod {
	return MakeDescBy(max, 1)
}

func Contains(haystack []string, needle string) bool { 
	return Find(haystack, needle) != -1
}

func Find(haystack []string, needle string) int {
	for index, item := range haystack {
		if item != needle { continue } 
		return index
	}
	return -1
}

func GetLastActivity(conversation types.Conversation) time.Time {
	base := time.Unix(0, 0)
	if len(conversation.Conversations) == 0 { return base }
	last := base

	for _, convo := range conversation.Conversations {
		if convo.LastActivity.After(last) { last = convo.LastActivity }
	}

	return last
}

func Modal(primative tview.Primitive, width, height int) *tview.Grid {
	grid := tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(primative, 1, 1, 1, 1, 0, 0, true)
	return grid
}

type SetPositionFn = func(state AppState, fn IntMod) AppState
type FallbackFn = func(state AppState, key rune) AppState
func VerticleListKeyBinds(state AppState, key rune, getSize func() int, setPosition SetPositionFn, fallback FallbackFn) AppState {
	switch {
		case key == BIND_KEY_TOP || key == BIND_KEY_BOTTOM:
			var messagePos int
			if key == BIND_KEY_BOTTOM { messagePos = getSize() }

			state = setPosition(state, func(_ int) int { return messagePos})
			break
		case key == BIND_KEY_UP || key == BIND_KEY_DOWN:
			size := getSize()

			jumpBy := state.jumpBy
			if jumpBy == -1 { jumpBy = 1 }
			state.jumpBy = -1

			var fn func(int) int
			if key == BIND_KEY_DOWN { fn = MakeIncBy(size, jumpBy) } else
			{ fn = MakeDescBy(size, jumpBy) }

			state = setPosition(state, fn)
			break
		case unicode.IsDigit(key):
			var jumpBy int
			numb := int(key - '0')

			if state.jumpBy == -1 { jumpBy = numb } else
			{ jumpBy = state.jumpBy * 10 + numb }

			state.jumpBy = jumpBy
			break
		default:
			state = fallback(state, key)
	}

	return state
}

func RemoveSearchQueryFilter(filters []SearchQueryFilter, pos int) []SearchQueryFilter {
	return append(filters[:pos], filters[pos + 1:]...)
}
