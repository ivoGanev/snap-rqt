package ui

import (
	"github.com/rivo/tview"
)

type ResponseView struct {
	*tview.TextView
	app *App
}

func NewResponseView(app *App) *ResponseView {
	responseView := ResponseView{
		app: app,
		TextView: tview.NewTextView(),
	}

	return &responseView
}

func (r *ResponseView) Init() {
	r.SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			r.app.Draw()
		})

	r.SetText("No response data")
	r.SetBorder(true)
	r.SetTitle("Response")
}
