package view

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/styles"

	"github.com/rivo/tview"
)

type RequestsViewListener interface {
	OnRequestSelected(data.RequestSimple)
	OnRequestMethodSelected(data.RequestSimple)
	OnRequestsListSelectionChanged(data.RequestSimple)
}

func (r *RequestsView) SetRequestsViewListener(l RequestsViewListener) {
	r.eventListener = l
}

type selectedRow struct {
	index   int
	request data.RequestSimple
}

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	selectedRow   selectedRow
	eventListener RequestsViewListener
	styleProvider styles.StyleProvider
	requests      []data.RequestSimple
}


func (r *RequestsView) RenderRequests(requests []data.RequestSimple) {
	r.requests = requests
	for i, request := range requests {
		methodText := r.styleProvider.GetStyledRequestMethod(string(request.MethodType))
		methodCell := tview.NewTableCell(methodText).SetReference(request)
		nameCell := tview.NewTableCell(request.Name).SetReference(request)

		r.SetCell(i, 0, methodCell)
		r.SetCell(i, 1, nameCell)
	}
}

func NewRequestsView(styles styles.StyleProvider) *RequestsView {
	requestsView := RequestsView{
		Table:         tview.NewTable(),
		styleProvider: styles,
	}
	return &requestsView
}

func (r *RequestsView) Init() {
	r.SetBorder(true)
	r.SetTitle("Requests")
	r.SetSelectable(true, true)
	r.SelectRequest(0)

	r.SetSelectedFunc(func(row int, column int) {
		ref := r.GetCell(row, column).GetReference()
		request, _ := ref.(data.RequestSimple)
		if column == 0 {
			if r.eventListener != nil {
				r.eventListener.OnRequestMethodSelected(request)
			}
		} else {
			if r.eventListener != nil {
				r.eventListener.OnRequestSelected(request)
			}
		}
	})

	r.SetSelectionChangedFunc(func(row int, column int) {
		request := r.requests[row]
		r.selectedRow = selectedRow{
			index:   row,
			request: request,
		}
		if r.eventListener != nil {
			r.eventListener.OnRequestsListSelectionChanged(request)
		}
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

func (r *RequestsView) ChangeMethodTypeOnSelectedRow(method http.RequestMethod) {
	r.GetCell(r.selectedRow.index, 0).
		SetText(r.styleProvider.GetStyledRequestMethod(string(method)))
}

func (r *RequestsView) GetSelectedRequest() *data.RequestSimple {
	return &r.selectedRow.request
}
