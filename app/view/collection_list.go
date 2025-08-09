package view

import (
	"snap-rq/app/constants"
	"snap-rq/app/entity"
	"snap-rq/app/input"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CollectionListListener interface {
	OnFocusedCollectionChanged(entity.Collection)
	OnCollectionAdd(position int)
	OnCollectionRemove(collection entity.Collection, position int)
	OnCollectionEditName(entity.Collection)
}

func (r *CollectionsList) SetListener(listener CollectionListListener) {
	r.listener = listener
}

type CollectionsList struct {
	*tview.Table
	input    *input.Handler
	listener CollectionListListener
}

func (r *CollectionsList) SelectCollection(collection entity.Collection) {
	r.Select(collection.RowPosition, 0)
}

func NewColletionsList(input *input.Handler) *CollectionsList {
	return &CollectionsList{
		Table: tview.NewTable(),
		input: input,
	}
}

func (r *CollectionsList) Init() {
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return r.input.SetInputCapture(constants.ViewCollections, event)
	})

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


	r.input.AddListener(func(action input.Action) {
		row, _ := r.GetSelection()
		switch action {
		case input.ActionAddCollection:
			r.listener.OnCollectionAdd(row)
			r.Select(row, 0)
		case input.ActionEditCollectionName:
			cell := r.GetCell(row, 0)
			if col, ok := cell.GetReference().(entity.Collection); ok {
				r.listener.OnCollectionEditName(col)
			}
		case input.ActionRemoveCollection:
			cell := r.GetCell(row, 0)
			if cell != nil {
				if col, ok := cell.GetReference().(entity.Collection); ok {
					r.listener.OnCollectionRemove(col, row)
					r.Select(row-1, 0)
				}
			}
		}
	})

}

func (r *CollectionsList) RenderCollections(collections []entity.Collection) {
	r.Clear()

	// Sort collections by RowPosition ascending
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].RowPosition < collections[j].RowPosition
	})

	for i, collection := range collections {
		nameCell := tview.NewTableCell(collection.Name).
			SetReference(collection)

		r.SetCell(i, 0, nameCell)
	}
}
