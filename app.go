package main

import "github.com/rivo/tview"

func run() {
	app := tview.NewApplication()

	client, err := GetClient()
	if err != nil { panic(err) }


	listConversations := tview.NewList()

	conversations, err := client.GetConversations()
	if err != nil { panic(err) }

	for _, conversation := range conversations {
		listConversations.AddItem(ParseConversation(conversation), "", 0, nil)
	}


	listMessages := tview.NewList()

	input := tview.NewInputField()

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