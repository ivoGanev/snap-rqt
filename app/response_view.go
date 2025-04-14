package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseView struct {
	*tview.TextArea
}

func NewResponseView() ResponseView {
	responseView := ResponseView{
		TextArea: tview.NewTextArea(),
	}

	return responseView
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
