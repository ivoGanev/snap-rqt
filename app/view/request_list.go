package view

import (
	"snap-rq/app/entity"
	"snap-rq/app/style"

	"github.com/rivo/tview"
)

const (
	METHOD_COLUMN  = 0
	REQUEST_COLUMN = 1
)

type RequestListListener interface {
	OnRequestMethodSelected(entity.RequestBasic)
	OnRequestNameSelected(entity.RequestBasic)
	OnSelectedRequestChanged(entity.RequestBasic)
}

func (r *RequestsList) SetListener(listener RequestListListener) {
	r.listener = listener
}

type RequestRow struct {
	Row     int
	Request entity.RequestBasic
}

type RequestsList struct {
	*tview.Table
	styles   style.StyleProvider
	requests map[string]int // Mapping: request id -> request row
	listener RequestListListener
}

func (r *RequestsList) RenderRequests(requests []entity.RequestBasic) {
	for row, request := range requests {
		r.requests[request.Id] = row

		methodText := r.styles.GetStyledRequestMethod(string(request.MethodType))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		requestCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(row, METHOD_COLUMN, methodCell)
		r.SetCell(row, REQUEST_COLUMN, requestCell)
	}
}

func NewRequestsList(styles style.StyleProvider) *RequestsList {
	requestsView := RequestsList{
		Table:  tview.NewTable(),
		styles: styles,
		requests: make(map[string]int),
	}
	return &requestsView
}

func (r *RequestsList) Init() {
	r.SetBorder(true)
	r.SetTitle("(r) Requests")
	r.SetSelectable(true, true)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(entity.RequestBasic)
		if column == METHOD_COLUMN {
			r.listener.OnRequestMethodSelected(request)
		} else {
			r.listener.OnRequestNameSelected(request)
		}
	})

	r.SetSelectionChangedFunc(func(row, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(entity.RequestBasic)
		r.listener.OnSelectedRequestChanged(request)
	})
}

// Selects an item from the request list
func (r *RequestsList) SelectRequest(requestId string) {
	requestRow := r.requests[requestId]
	r.Select(requestRow, REQUEST_COLUMN)
	// r.listener.OnSelectedRequestChanged(requestRow.Request)
}

// Selects the request method table item on a specific row of the requests list
func (r *RequestsList) SelectMethod(requestId string) {
	requestRow := r.requests[requestId]
	r.Select(requestRow, METHOD_COLUMN)
	// r.listener.OnSelectedRequestChanged(requestRow.Request)
}

func (r *RequestsList) ChangeMethodTypeOnSelectedRow(requestId string, requestMethod string) {
	requestRow := r.requests[requestId]
	r.GetCell(requestRow, 0).
		SetText(r.styles.GetStyledRequestMethod(string(requestMethod)))
}
