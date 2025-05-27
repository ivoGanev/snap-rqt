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

type RequestsList struct {
	*tview.Table
	styles   style.StyleProvider
	requests []entity.RequestBasic
	listener RequestListListener
}

func (r *RequestsList) RenderRequests(requests []entity.RequestBasic) {
	r.requests = requests
	for i, request := range requests {
		methodText := r.styles.GetStyledRequestMethod(string(request.MethodType))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		requestCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(i, METHOD_COLUMN, methodCell)
		r.SetCell(i, REQUEST_COLUMN, requestCell)
	}
}

func NewRequestsList(styles style.StyleProvider) *RequestsList {
	requestsView := RequestsList{
		Table:  tview.NewTable(),
		styles: styles,
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
		request := r.requests[row]
		r.listener.OnSelectedRequestChanged(request)
	})
}

// Selects an item from the request list
func (r *RequestsList) SelectRequest(row int) {
	request := r.requests[row]
	r.Select(row, REQUEST_COLUMN)
	r.listener.OnSelectedRequestChanged(request)
}

// Selects the request method table item on a specific row of the requests list
func (r *RequestsList) SelectMethod(row int) {
	request := r.requests[row]
	r.Select(row, METHOD_COLUMN)
	r.listener.OnSelectedRequestChanged(request)
}

func (r *RequestsList) ChangeMethodTypeOnSelectedRow(row int, requestMethod string) {
	r.GetCell(row, 0).
		SetText(r.styles.GetStyledRequestMethod(string(requestMethod)))
}
