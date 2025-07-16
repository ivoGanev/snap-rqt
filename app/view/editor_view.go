package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	EDITOR_VIEW_MODE_HEADERS = 1
	EDITOR_VIEW_MODE_BODY    = 2
)

type EditorView struct {
	*tview.Flex
	app           *tview.Application
	headersButton tview.Button
	bodyButton    tview.Button
	currentMode   int
}

func NewEditorView(app *tview.Application) *EditorView {
	editorView := EditorView{
		Flex:          tview.NewFlex(),
		app:           app,
		headersButton: *tview.NewButton("(h) Headers"),
		bodyButton:    *tview.NewButton("(b) Body"),
		currentMode:   EDITOR_VIEW_MODE_HEADERS,
	}
	return &editorView
}

func (r EditorView) Init() {
	r.SetBorder(true)
	r.SetTitle("Edit Request")
	r.SetDirection(tview.FlexRow)

	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'b' {
			r.app.SetFocus(&r.bodyButton)
			return nil
		}
		if event.Rune() == 'h' {
			r.app.SetFocus(&r.headersButton)
			return nil
		}
		return event
	})

	top := tview.NewFlex().
		AddItem(&r.headersButton, 0, 1, false).
		AddItem(&r.bodyButton, 0, 1, false)

	textArea := tview.NewTextArea()
	textArea.SetBorder(true)
	r.AddItem(top, 3, 0, false)
	r.AddItem(textArea, 0, 1, true)
	r.app.SetFocus(&r.headersButton)
}
