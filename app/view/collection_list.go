package view

import (
	"snap-rq/app/entity"

	"github.com/rivo/tview"
)

type CollectionListListener interface {
	OnFocusedCollectionChanged(entity.Collection)
}

func (r *CollectionsList) SetListener(listener CollectionListListener) {
	r.listener = listener
}

type CollectionsList struct {
	*tview.Table
	listener    CollectionListListener
	collections map[string]int // Mapping: collection id -> collection row
}

func (r *CollectionsList) SelectCollection(collectionId string) {
	collectionRow := r.collections[collectionId]
	r.Select(collectionRow, 0)
}

func NewColletionsList() *CollectionsList {
	return &CollectionsList{
		Table: tview.NewTable(),
	}
}

func (r *CollectionsList) Init() {
	r.SetBorder(true)
	r.SetTitle("(q) Collections")
	r.SetSelectable(true, true)

	r.SetSelectionChangedFunc(func(row, column int) {
		if row == -1 {
			return // no-op when the clicked row doesn't have any collection
		}

		collection := r.GetCell(row, 0).GetReference().(entity.Collection)
		r.listener.OnFocusedCollectionChanged(collection)
	})
}

func (r *CollectionsList) RenderCollections(collections []entity.Collection) {
	for i, collection := range collections {
		nameCell := tview.NewTableCell(collection.Name).
			SetReference(collection)

		r.SetCell(i, 0, nameCell)
	}
}
