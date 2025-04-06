package ui

import (
	"context"
	"fmt"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/model"
	"time"

	"github.com/rivo/tview"
	"slices"
)

type RequestsListChangedListener interface {
	OnRequestsListSelectionChanged(selection *data.Node[http.Request])
}

type RequestsView struct {
	*tview.Table
	app             *App
	SelectedNode    *data.Node[http.Request]
	SelectedRow     int
	data            *[]data.Node[http.Request]
	changeListeners []RequestsListChangedListener
}

func (r *RequestsView) OnMethodSelectionChanged(method string) {
	r.SelectedNode.Data.Method = http.RequestMethod(method)

	methodText := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(http.RequestMethod(method)), method)
	r.GetCell(r.SelectedRow, 0).SetText(methodText)
	r.app.SetFocus(r)
}

func (r *RequestsView) OnRequestsModelChanged(requests *[]data.Node[http.Request], operation model.CrudOp, multiplicity model.Multiplicity) {
	if operation == model.UPDATE && multiplicity == model.MANY {
		r.setAllRequests(requests)
	}
}

func (r *RequestsView) setAllRequests(requests *[]data.Node[http.Request]) {
	r.data = requests

	for i, request := range *requests {
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

	r.Select(0, 1)
	r.processSelectionChanged(0)

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
						r.app.Views.ResponseWindow.SetText(response, false)
					})

					cancel()
				}()
			}
		}
	})

	r.SetSelectionChangedFunc(func(row int, column int) {
		r.processSelectionChanged(row)
	})
}


func (r *RequestsView) processSelectionChanged(row int) {
	data := *r.data
	r.SelectedNode = &data[row]
	r.SelectedRow = row

	for _, l := range r.changeListeners {
		l.OnRequestsListSelectionChanged(&data[row])
	}
}

func (u *RequestsView) AddListener(l RequestsListChangedListener) {
	u.changeListeners = append(u.changeListeners, l)
}

func (u *RequestsView) RemoveListener(l RequestsListChangedListener) {
	for i, lis := range u.changeListeners {
		if lis == l {
			u.changeListeners = slices.Delete(u.changeListeners, i, i+1)
			return
		}
	}
}
