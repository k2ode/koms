package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
const RUNE_LEFT   = 'h'
const RUNE_DOWN   = 'j'
const RUNE_UP     = 'k'
const RUNE_RIGHT  = 'l'
const RUNE_TOP    = 'g'
const RUNE_BOTTOM = 'G'

func AddBindings(list *tview.List, handleHover func(int), handleLeft func(int), handleRight func(int)) {
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		endPosition := list.GetItemCount() - 1
		position := list.GetCurrentItem()
		var fn func(i int) int
		switch {
		case event.Rune() == RUNE_DOWN:
			fn = func(i int) int {
				if i == endPosition { return 0 }
				return i + 1
			}
		case event.Rune() == RUNE_UP:
			fn = func(i int) int {
				if i == 0 { return endPosition }
				return i - 1
			}
		case event.Rune() == RUNE_RIGHT:
			handleRight(position)
			return nil
		
		case event.Rune() == RUNE_LEFT:
			handleLeft(position)
			return nil
		
		case event.Rune() == RUNE_TOP:
			fn = func (_ int) int { return 0 }
		
		case event.Rune() == RUNE_BOTTOM:
			fn = func (_ int) int { return endPosition }

		default:
			return nil
		}
		
		newPos := fn(position)
		list.SetCurrentItem(newPos)
		handleHover(newPos)
		
		return event
	})
}

func UIListConversations(client Client, messagePreview func(int), exit func(int), focusMessage func(int)) *tview.List {
	listConversations := tview.NewList()

	AddBindings(
		listConversations,
		messagePreview,
		exit,
		focusMessage,
	)

	return listConversations
}
	
func UIListMessages() *tview.List {
	listMessages := tview.NewList().SetSelectedFocusOnly(true)

	return listMessages
}

func UIInput() *tview.InputField {
	input := tview.NewInputField()

	return input
}

func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }

	var conversations []Conversation

	listMessages := UIListMessages()

	updateMessages := func (messages []Message) {
		listMessages.Clear()
		for _, message := range messages {
			msg := ParseMessage(client, message)
			listMessages.AddItem(msg, "", 0, nil)
		}
	}

	messagePreview := func(newPos int) {

		conversation := conversations[newPos]

		messages, err := client.GetConversationMessages(conversation)

		if err != nil { panic(err) }


		updateMessages(messages)

	}


	exit := func(_ int) {
		app.Stop()	
	}

	messageFocus := func(_ int) {

	}


	listConversations := UIListConversations(client, messagePreview, exit, messageFocus)

	updateConversations := func () {
		conversations, err = client.GetConversations()
		if err != nil { panic(err) }

		for _, conversation := range conversations {
			labelConversation := ParseConversation(client, conversation)
			listConversations.AddItem(labelConversation, "", 0, nil)
		}
	}

	updateConversations()

	input := UIInput()

	gridConversation := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0).
		SetBorders(true).
		AddItem(listMessages, 0, 0, 1, 1, 0, 0, false).
		AddItem(input, 1, 0, 1, 1, 0, 0, false)

	gridContainer := tview.NewGrid().
		SetRows(0).
		SetColumns(30, 0).
		SetBorders(true).
		AddItem(listConversations, 0, 0, 1, 1, 0, 0, true).
		AddItem(gridConversation, 0, 1, 1, 1, 0, 0, false)

	if err := app.SetRoot(gridContainer, true).Run(); err != nil {
		panic(err)
	}
}
