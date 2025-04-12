package internal

import (
	"slices"
	"snap-rq/internal/controller"
	"github.com/rivo/tview"
)

type CollectionsView struct {
	*tview.Table
	app       *App
	listeners []CollectionsViewChangedListener
}

type CollectionsViewChangedListener interface {
	OnCollectionListSelectionChanged(selection *Collection)
}

func (u *CollectionsView) AddListener(l CollectionsViewChangedListener) {
	u.listeners = append(u.listeners, l)
}

func (u *CollectionsView) RemoveListener(l CollectionsViewChangedListener) {
	for i, lis := range u.listeners {
		if lis == l {
			u.listeners = slices.Delete(u.listeners, i, i+1)
			return
		}
	}
}

// Initialisation
func NewColletionsView(c *controller.EditorController) *CollectionsView {
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
