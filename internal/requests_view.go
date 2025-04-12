package internal

import (
	"context"
	"fmt"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/model"
	"time"
	"github.com/rivo/tview"
)

type selectedRow struct {
	index   int
	request Request
}

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	app         *App
	collection  *Collection
	selectedRow *selectedRow
}


func (r *RequestsView) DisplayCollection(collection *Collection) {
	r.collection = collection
	requests := *collection.Data
	for i, request := range requests {
		methodText := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(request.Data.Method), request.Data.Method)

		methodCell := tview.NewTableCell(methodText).SetReference(request)
		nameCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(i, 0, methodCell)
		r.SetCell(i, 1, nameCell)
	}
}

func NewRequestsView(app *App) *RequestsView {
	requestsView := RequestsView{
		app:   app,
		Table: tview.NewTable(),
	}

	return &requestsView
}

func (r *RequestsView) Init() {
	r.SetBorder(true)
	r.SetTitle("Requests")
	r.SetSelectable(true, true)
	r.SelectRequestRow(0)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(data.Node[http.Request])

		if column == 0 {
			// c.Show("method-selection")
			r.app.Views.MethodSelectionModal.Show()
		} else {
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

				response := http.SendRequest(ctx, *request.Data)
				r.app.QueueUpdateDraw(func() {
					r.app.Views.ResponseView.SetText(response, false)
				})

				cancel()
			}()
		}
	})

	r.SetSelectionChangedFunc(func(row int, column int) {
		r.processSelectionChanged(row)
	})
}

func (r *RequestsView) SelectRequestRow(index int) {
	r.Select(0, index)
	r.selectedRow = &selectedRow{
		index:   index,
		request: (*r.collection.Data)[index],
	}
}
