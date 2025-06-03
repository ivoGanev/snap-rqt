package view

import (
	"snap-rq/app/entity"
	"snap-rq/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	METHOD_COLUMN  = 0
	REQUEST_COLUMN = 1
)

type RequestListListener interface {
	OnRequestListMethodSelected(request entity.RequestBasic)
	OnRequestListNameSelected(request entity.RequestBasic)
	OnRequestListRequestFocusChanged(request entity.RequestBasic)
	OnRequestListAdd(position int) //  'position' indicates the position of the request currently in focus (i.e. not the position where the user expects the next request to be added)
	OnRequestListRemove(request entity.RequestBasic, position int)
	OnRequestListDuplicate(request entity.RequestBasic)
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
	for _, request := range requests {
		row := request.RowPosition
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
		Table:    tview.NewTable(),
		styles:   styles,
		requests: make(map[string]int),
	}
	return &requestsView
}

func (r *RequestsList) Init() {
	r.SetBorder(true)
	r.SetTitle("(w) Requests")
	r.SetSelectable(true, true)

	r.SetSelectedFunc(func(row int, column int) {
		request := r.getRequest(row, column)
		if column == METHOD_COLUMN {
			r.listener.OnRequestListMethodSelected(request)
		} else {
			r.listener.OnRequestListNameSelected(request)
		}
	})

	r.SetSelectionChangedFunc(func(row, column int) {
		r.listener.OnRequestListRequestFocusChanged(r.getRequest(row, column))
	})

	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			row, _ := r.GetSelection()
			r.listener.OnRequestListAdd(row)
			return nil
		} else if event.Key() == tcell.KeyDEL || event.Key() == tcell.KeyDelete {
			row, column := r.GetSelection()
			r.listener.OnRequestListRemove(r.getRequest(row, column), row)
			return nil
		} else if event.Rune() == 'd' {
			row, column := r.GetSelection()
			r.listener.OnRequestListDuplicate(r.getRequest(row, column))
			return nil
		}
		return event
	})
}

// Selects an item from the request list
func (r *RequestsList) SelectRequest(requestId string) {
	requestRow := r.requests[requestId]
	r.Select(requestRow, REQUEST_COLUMN)
}

// Selects the request method table item on a specific row of the requests list
func (r *RequestsList) SelectMethod(requestId string) {
	requestRow := r.requests[requestId]
	r.Select(requestRow, METHOD_COLUMN)
}

func (r *RequestsList) ChangeRequestMethod(requestId string, requestMethod string) {
	requestRow := r.requests[requestId]
	r.GetCell(requestRow, 0).
		SetText(r.styles.GetStyledRequestMethod(string(requestMethod)))
}

func (r *RequestsList) getRequest(row int, column int) entity.RequestBasic {
	ref := r.GetCell(row, column).GetReference()
	request, _ := ref.(entity.RequestBasic)
	return request
}
