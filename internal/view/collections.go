package view

import (
	"snap-rq/internal/data"
	"github.com/rivo/tview"
)

type CollectionsView struct {
	*tview.Table
	listener CollectionsViewChangedListener
}

type CollectionsViewChangedListener interface {
	OnCollectionListSelectionChanged(selection *data.Collection)
}

func (r *CollectionsView) SetListener(l CollectionsViewChangedListener) {
	r.listener = l
}

func NewColletionsView() *CollectionsView {
	return &CollectionsView{
		Table: tview.NewTable(),
	}
}

func (r *CollectionsView) Init() {
	r.SetBorder(true)
	r.SetTitle("Collections")
	r.SetSelectable(true, true)
}
