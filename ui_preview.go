package main

import "github.com/rivo/tview"

type PreviewComponent = *tview.TextView

func MakePreview(state AppState) (PreviewComponent, UpdateStateFn) {
	preview := tview.NewTextView()

	updatePreview := MakePreviewUpdateFunc(preview)

	updatePreview(state)

	return preview, updatePreview
}

func MakePreviewUpdateFunc(preview PreviewComponent) UpdateStateFn {
	return func(state AppState) {
		message := state.messages[state.conversationPos][state.messagePos]
		previewText := GetMessagePreview(message)
		preview.SetText(previewText)
	}
}

func AddContainerPreview(container *tview.Grid, preview PreviewComponent) {
	container.AddItem(
		preview,
		ROW_POS_PREVIEW,
		COLUMN_POS_PREVIEW,
		ROW_SPAN_PREVIEW,
		COLUMN_SPAN_PREVIEW,
		HEIGHT_MIN_PREVIEW,
		WIDTH_MIN_PREVIEW,
		false,
	)
}