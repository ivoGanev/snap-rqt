package app

import (
	"snap-rq/app/style"

	"github.com/rivo/tview"
)


type selectedRow struct {
	index   int
	request RequestListItem
}

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	selectedRow   selectedRow
	styleProvider style.StyleProvider
	requests      []RequestListItem
	controller    RequestsController
}

func (r *RequestsView) RenderRequests(requests []RequestListItem) {
	r.requests = requests
	for i, request := range requests {
		methodText := r.styleProvider.GetStyledRequestMethod(string(request.MethodType))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		nameCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(i, 0, methodCell)
		r.SetCell(i, 1, nameCell)
	}
}

func NewRequestsView(styles style.StyleProvider, controller RequestsController) RequestsView {
	requestsView := RequestsView{
		Table:         tview.NewTable(),
		styleProvider: styles,
		controller: controller,
	}
	return requestsView
}

func (r *RequestsView) Init() {
	r.SetBorder(true)
	r.SetTitle("Requests")
	r.SetSelectable(true, true)
	r.SelectRequest(0)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(RequestListItem)
		if column == 0 {
			r.controller.HandleRequestMethodSelected(request)
		} else {
			r.controller.HandleRequestNameSelected(request)
		}
	})

	r.SetSelectionChangedFunc(func(row int, column int) {
		request := r.requests[row]
		r.selectedRow = selectedRow{
			index:   row,
			request: request,
		}
		r.controller.HandleSelectedRequestChanged(request)
	})
}

// Selects an item from the request list
func (r *RequestsView) SelectRequest(position int) {
	r.Select(0, position)
	r.selectedRow = selectedRow{
		index:   position,
		request: r.requests[position],
	}
}

func (r *RequestsView) ChangeMethodTypeOnSelectedRow(requestMethod string) {
	r.GetCell(r.selectedRow.index, 0).
		SetText(r.styleProvider.GetStyledRequestMethod(string(requestMethod)))
}

func (r *RequestsView) GetSelectedRequest() RequestListItem {
	return r.selectedRow.request
}
