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
		msgs, exists := state.cache.messages[state.pos]
		if !exists { preview.SetText(""); return }
		if len(msgs) == 0 { preview.SetText(""); return }
		x, y, width, height := preview.GetInnerRect()
		size := Size{
			x, y, width, height,
		}
		previewText := GetMessagePreview(state, size)
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
