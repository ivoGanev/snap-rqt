package app

import (
	"github.com/rivo/tview"
)

type CollectionsView struct {
	*tview.Table
	controller CollectionController
}

func (r *CollectionsView) SelectCollection(row int) {
	r.Select(row, 0)
	collectionId := r.GetCell(row, 0).GetReference().(string)
	r.controller.HandleCollectionChanged(SelectedCollection{Id: collectionId, Row: row})
}

func NewColletionsView(controller CollectionController) CollectionsView {
	return CollectionsView{
		Table: tview.NewTable(),
		controller: controller,
	}
}

func (r *CollectionsView) Init() {
	r.SetBorder(true)
	r.SetTitle("(c) Collections")
	r.SetSelectable(true, true)

	r.SetSelectionChangedFunc(func(row, column int) {
		collectionId := r.GetCell(row, 0).GetReference().(string)
		r.controller.HandleCollectionChanged(SelectedCollection{Id: collectionId, Row: row})
	})
}

func (r *CollectionsView) RenderCollections(collections []Collection) {
	for i, collection := range collections {
		nameCell := tview.NewTableCell(collection.Name).
			SetReference(collection.Id)

		r.SetCell(i, 0, nameCell)
	}
}
