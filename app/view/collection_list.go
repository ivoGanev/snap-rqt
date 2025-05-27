package view

import (
	"snap-rq/app/entity"

	"github.com/rivo/tview"
)

type CollectionListListener interface {
	OnCollectionChanged(entity.FocusedCollection)
}

func (r *CollectionsList) SetListener(listener CollectionListListener) {
	r.listener = listener
}

type CollectionsList struct {
	*tview.Table
	listener CollectionListListener
}

func (r *CollectionsList) SelectCollection(row int) {
	r.Select(row, 0)
	collectionId := r.GetCell(row, 0).GetReference().(string)
	r.listener.OnCollectionChanged(entity.FocusedCollection{Id: collectionId, Row: row})
}

func NewColletionsList() *CollectionsList {
	return &CollectionsList{
		Table: tview.NewTable(),
	}
}

func (r *CollectionsList) Init() {
	r.SetBorder(true)
	r.SetTitle("(c) Collections")
	r.SetSelectable(true, true)

	r.SetSelectionChangedFunc(func(row, column int) {
		collectionId := r.GetCell(row, 0).GetReference().(string)
		r.listener.OnCollectionChanged(entity.FocusedCollection{Id: collectionId, Row: row})
	})
}

func (r *CollectionsList) RenderCollections(collections []entity.Collection) {
	for i, collection := range collections {
		nameCell := tview.NewTableCell(collection.Name).
			SetReference(collection.Id)

		r.SetCell(i, 0, nameCell)
	}
}
