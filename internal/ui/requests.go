package ui

import (
	"context"
	"fmt"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"time"

	"github.com/rivo/tview"
)

type RequestsView struct {
	*tview.Table
	app          *App
	SelectedNode *data.Node[http.Request]
	SelectedRow  int
	data         *[]data.Node[http.Request]
}

// OnSelectionChanged implements OnMethodSelectionModalChangeListener.
func (r *RequestsView) OnSelectionChanged(method string) {
	r.SelectedNode.Data.Method = http.RequestMethod(method)

	methodText := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(http.RequestMethod(method)), method)
	r.GetCell(r.SelectedRow, 0).SetText(methodText)
	r.app.SetFocus(r)
}

// OnRequestAdded implements model.ReqestsListener.
func (r *RequestsView) OnRequestAdded(string) {
	panic("unimplemented")
}

// OnRequestChanged implements model.ReqestsListener.
func (r *RequestsView) OnRequestChanged(string) {
	panic("unimplemented")
}

// OnRequestRemoved implements model.ReqestsListener.
func (r *RequestsView) OnRequestRemoved(string) {
	panic("unimplemented")
}

// OnRequestsSet implements model.ReqestsListener.
func (r *RequestsView) OnRequestsSet(data []data.Node[http.Request]) {
	r.data = &data

	size := len(data)
	count := 0

	for count < size {
		request := data[count]

		method := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(request.Data.Method), string(request.Data.Method))

		methodCell := tview.NewTableCell(method)
		methodCell.SetReference(request)

		nameCell := tview.NewTableCell(string(request.Name))
		nameCell.SetReference(request)
		r.SetCell(count, 0, methodCell)
		r.SetCell(count, 1, nameCell)
		count++
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
	r.SelectRequest(1)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, ok := ref.(data.Node[http.Request])
		if !ok {
			panic("Failed to cast reference to *http.Request")
		} else {
			if column == 0 {
				r.app.Views.MethodSelectionModal.Show()
			} else {
				go func() {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

					response := http.SendRequest(ctx, *request.Data)
					r.app.QueueUpdateDraw(func() {
						r.app.Views.Response.SetText(response)
					})

					cancel()
				}()
			}
		}
	})

	r.SetSelectionChangedFunc(func(row int, column int) {
		data := *r.data
		r.SelectedNode = &data[row]
		r.SelectedRow = row
	})
}

func (r *RequestsView) SelectRequest(row int) {
	data := *r.data
	r.SelectedNode = &data[row]
	r.SelectedRow = row
	r.Select(0, row)
}
