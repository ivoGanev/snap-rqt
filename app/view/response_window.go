package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseWindow struct {
	*tview.TextView
}

func NewResponseWindow() *ResponseWindow {
	responseView := ResponseWindow{
		TextView: tview.NewTextView(),
	}

	return &responseView
}

func (r *ResponseWindow) Init() {
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'p' {
			return tcell.NewEventKey(tcell.KeyCtrlV, 'v', tcell.ModNone)
		}
		if event.Rune() == 'c' {
			return tcell.NewEventKey(tcell.KeyCtrlQ, 'c', tcell.ModNone)
		}
		return event
	})

	r.SetText("No response data")
	r.SetBorder(true)
	r.SetTitle("Response")
}

