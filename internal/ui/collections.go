package ui

import (
	"github.com/rivo/tview"
)

type CollectionsView struct {
	*tview.Table
	app *App
}

func NewColletionsView(app *App) *CollectionsView {
	collectionsView := CollectionsView{
		app:   app,
		Table: tview.NewTable(),
	}

	return &collectionsView
}


func (r *CollectionsView) Init() {
	r.SetBorder(true)
	r.SetTitle("Collections")

	r.SetSelectable(true, true)
}