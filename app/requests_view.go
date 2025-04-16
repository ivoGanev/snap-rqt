package app

import (
	"github.com/rivo/tview"
	"snap-rq/app/style"
)

const (
	METHOD_COLUMN  = 0
	REQUEST_COLUMN = 1
)

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	styleProvider style.StyleProvider
	requests      []RequestListItem
	controller    RequestsController
}

func (r *RequestsView) RenderRequests(requests []RequestListItem) {
	r.requests = requests
	for i, request := range requests {
		methodText := r.styleProvider.GetStyledRequestMethod(string(request.MethodType))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		requestCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(i, METHOD_COLUMN, methodCell)
		r.SetCell(i, REQUEST_COLUMN, requestCell)
	}
}

func NewRequestsView(styles style.StyleProvider, controller RequestsController) RequestsView {
	requestsView := RequestsView{
		Table:         tview.NewTable(),
		styleProvider: styles,
		controller:    controller,
	}
	return requestsView
}

func (r *RequestsView) Init() {
	r.SetBorder(true)
	r.SetTitle("(r) Requests")
	r.SetSelectable(true, true)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(RequestListItem)
		if column == METHOD_COLUMN {
			r.controller.HandleRequestMethodSelected(request)
		} else {
			r.controller.HandleRequestNameSelected(request)
		}
	})

	r.SetSelectionChangedFunc(func(row, column int) {
		request := r.requests[row]
		r.controller.HandleSelectedRequestChanged(request)
	})
}

// Selects an item from the request list
func (r *RequestsView) SelectRequest(row int) {
	request := r.requests[row]
	r.Select(row, REQUEST_COLUMN)
	r.controller.HandleSelectedRequestChanged(request)
}

// Selects the request method table item on a specific row of the requests list
func (r *RequestsView) SelectMethod(row int) {
	request := r.requests[row]
	r.Select(row, METHOD_COLUMN)
	r.controller.HandleSelectedRequestChanged(request)
}

func (r *RequestsView) ChangeMethodTypeOnSelectedRow(row int, requestMethod string) {
	r.GetCell(row, 0).
		SetText(r.styleProvider.GetStyledRequestMethod(string(requestMethod)))
}
