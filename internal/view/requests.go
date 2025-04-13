package view

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/styles"

	"github.com/rivo/tview"
)

type RequestsViewListener interface {
	OnRequestSelected(data.Request)
	OnRequestMethodSelected(data.Request)
	OnRequestsListSelectionChanged(data.Request)
}

func (r *RequestsView) SetRequestsViewListener(l RequestsViewListener) {
	r.eventListener = l
}

type selectedRow struct {
	index   int
	request data.Request
}

// Displays a collection of requests
type RequestsView struct {
	*tview.Table
	collection    *data.Collection
	selectedRow   selectedRow
	eventListener RequestsViewListener
	styleProvider styles.StyleProvider
}

func (r *RequestsView) DisplayCollection(collection *data.Collection) {
	r.collection = collection
	requests := *collection.Data
	for i, request := range requests {
		methodText := r.styleProvider.GetStyledRequestMethod(string(request.Data.Method))
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
		request, _ := ref.(data.Request)
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
		request := (*r.collection.Data)[row]
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
		request: (*r.collection.Data)[position],
	}
}

func (r *RequestsView) ChangeMethodTypeOnSelectedRow(method http.RequestMethod) {
	r.GetCell(r.selectedRow.index, 0).
		SetText(r.styleProvider.GetStyledRequestMethod(string(method)))
}

func (r *RequestsView) GetSelectedRequest() *data.Request {
	return &r.selectedRow.request
}