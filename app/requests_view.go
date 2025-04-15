package app

import (
	"snap-rq/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	METHOD_COLUMN  = 0
	REQUEST_COLUMN = 1
)

type SessionState struct {
	RowIndex    int
	ColumnIndex int
	request     RequestListItem
}

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	sessionState  SessionState
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
		r.setSessionState(row, column, request)
		r.controller.HandleSelectedRequestChanged(request)
	})

	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			r.controller.HandleColumnSwitch()
			return nil
		}
		return event
	})
}

// Selects an item from the request list
func (r *RequestsView) SelectRequest(row int) {
	request := r.requests[row]
	r.Select(row, REQUEST_COLUMN)
	r.setSessionState(row, REQUEST_COLUMN, request)
	r.controller.HandleSelectedRequestChanged(request)
}

// Selects the request method table item on a specific row of the requests list
func (r *RequestsView) SelectMethod(row int) {
	request := r.requests[row]
	r.Select(row, METHOD_COLUMN)
	r.setSessionState(row, METHOD_COLUMN, request)
	r.controller.HandleSelectedRequestChanged(request)
}

func (r *RequestsView) SwitchColumns() {
	row := r.sessionState.RowIndex
	if r.sessionState.ColumnIndex == METHOD_COLUMN {
		r.SelectRequest(row)
	} else {
		r.SelectMethod(row)
	}
}

func (r *RequestsView) ChangeMethodTypeOnSelectedRow(requestMethod string) {
	r.GetCell(r.sessionState.RowIndex, 0).
		SetText(r.styleProvider.GetStyledRequestMethod(string(requestMethod)))
}

func (r *RequestsView) GetSelectedRequest() RequestListItem {
	return r.sessionState.request
}

func (r *RequestsView) setSessionState(rowIndex int, columnIndex int, request RequestListItem) {
	r.sessionState = SessionState{
		ColumnIndex: columnIndex,
		RowIndex:    rowIndex,
		request:     request,
	}
}
