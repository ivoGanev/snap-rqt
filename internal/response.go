package internal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseView struct {
	*tview.TextArea
	app *App
}

func NewResponseView(app *App) *ResponseView {
	responseView := ResponseView{
		app:      app,
		TextArea: tview.NewTextArea(),
	}

	return &responseView
}

func (r *ResponseView) Init() {
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'p' {
			return tcell.NewEventKey(tcell.KeyCtrlV, 'v', tcell.ModNone)
		}
		if event.Rune() == 'c' {
			return tcell.NewEventKey(tcell.KeyCtrlQ, 'c', tcell.ModNone)
		}
		return event
	})

	r.SetText("No response data", false)
	r.SetBorder(true)
	r.SetTitle("Response")
}
