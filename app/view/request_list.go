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
	OnRequestListNameSelected(request entity.RequestBasic)
	OnRequestListRequestFocusChanged(request entity.RequestBasic)
	OnRequestListAdd(position int) //  'position' indicates the position of the request currently in focus (i.e. not the position where the user expects the next request to be added)
	OnRequestListRemove(request entity.RequestBasic, position int)
	OnRequestListDuplicate(request entity.RequestBasic)
	OnRequestListEditName(request entity.RequestBasic)
}

func (r *RequestsList) SetListener(listener RequestListListener) {
	r.listener = listener
}

type RequestsList struct {
	*tview.Table
	styles   style.StyleProvider
	listener RequestListListener
}

func (r *RequestsList) RenderRequests(requests []entity.RequestBasic) {
	r.Clear()
	for _, request := range requests {
		row := request.RowPosition

		methodText := r.styles.GetStyledRequestMethod(string(request.Method))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		requestCell := tview.NewTableCell(request.Name).SetReference(request)

		methodCell.SetSelectable(false)
		r.SetCell(row, METHOD_COLUMN, methodCell)
		r.SetCell(row, REQUEST_COLUMN, requestCell)
	}
}

func NewRequestsList(styles style.StyleProvider) *RequestsList {
	requestsView := RequestsList{
		Table:    tview.NewTable(),
		styles:   styles,
	}
	return &requestsView
}

func (r *RequestsList) Init() {
	r.SetBorder(true)
	r.SetTitle("(w) Requests")
	r.SetSelectable(true, true)

	r.SetSelectedFunc(func(row int, column int) {
		request := r.getRequest(row, column)
		if column == REQUEST_COLUMN {
			r.listener.OnRequestListNameSelected(request)
		}
	})

	r.SetSelectionChangedFunc(func(row, column int) {
		if row == -1 || column == -1 {
			return // no op when selection doesn't contain any requests
		}
		r.listener.OnRequestListRequestFocusChanged(r.getRequest(row, column))
	})

	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			row, _ := r.GetSelection()
			r.listener.OnRequestListAdd(row)
			r.Select(row, REQUEST_COLUMN)
			return nil
		} else if event.Key() == tcell.KeyDEL || event.Key() == tcell.KeyDelete {
			row, column := r.GetSelection()
			r.listener.OnRequestListRemove(r.getRequest(row, column), row)
			r.Select(row-1, REQUEST_COLUMN)
			return nil
		} else if event.Rune() == 'd' {
			row, column := r.GetSelection()
			r.listener.OnRequestListDuplicate(r.getRequest(row, column))
			return nil
		} else if event.Rune() == 'n' {
			row, column := r.GetSelection()
			r.listener.OnRequestListEditName(r.getRequest(row, column))
			return nil
		}
		return event
	})
}

// Selects an item from the request list
func (r *RequestsList) SelectRequest(request entity.Request) {
	r.Select(request.RowPosition, REQUEST_COLUMN)
	r.GetCell(request.RowPosition, METHOD_COLUMN).SetText(r.styles.GetStyledRequestMethod(string(request.Method)))
}

func (r *RequestsList) getRequest(row int, column int) entity.RequestBasic {
	ref := r.GetCell(row, column).GetReference()
	request, _ := ref.(entity.RequestBasic)
	return request
}
