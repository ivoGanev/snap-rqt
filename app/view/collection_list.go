package view

import (
	"snap-rq/app/entity"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CollectionListListener interface {
	OnFocusedCollectionChanged(entity.Collection)
	OnCollectionAdd(position int)
	OnCollectionRemove(collection entity.Collection, position int)
}

func (r *CollectionsList) SetListener(listener CollectionListListener) {
	r.listener = listener
}

type CollectionsList struct {
	*tview.Table
	listener    CollectionListListener
	collections map[string]int // Mapping: collection id -> collection row
}

func (r *CollectionsList) SelectCollection(collection entity.Collection) {
	collectionRow := r.collections[collection.Id]
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

	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			row, _ := r.GetSelection()
			r.listener.OnCollectionAdd(row)
			return nil
		} else if event.Key() == tcell.KeyDEL || event.Key() == tcell.KeyDelete {
			row, _ := r.GetSelection()
			cell := r.GetCell(row, 0)
			if cell != nil {
				if col, ok := cell.GetReference().(entity.Collection); ok {
					r.listener.OnCollectionRemove(col, row)
				}
			}
			return nil
		}
		return event
	})
}

func (r *CollectionsList) RenderCollections(collections []entity.Collection) {
	r.Clear()
	for i, collection := range collections {
		nameCell := tview.NewTableCell(collection.Name).
			SetReference(collection)

		r.SetCell(i, 0, nameCell)
	}
}
