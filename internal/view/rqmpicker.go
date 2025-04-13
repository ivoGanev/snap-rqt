package view

import (
	"snap-rq/internal/http"
	"github.com/rivo/tview"
)

type RequestMethodPickerListener interface {
	OnRequestMethodSelected(http.RequestMethod)
}

func (m *RequestMethodPickerModal) SetRequestMethodPickerListener(l RequestMethodPickerListener) {
	m.listener = l
}

// RequestMethodPickerModal: A modal dialog that allows the user to pick a request method type, e.g. GET, PUT, PATCH, etc.
type RequestMethodPickerModal struct {
	*tview.Grid
	listener RequestMethodPickerListener
}

func NewMethodPickerModal() *RequestMethodPickerModal {
	return &RequestMethodPickerModal{
		Grid: tview.NewGrid(),
	}
}

func (m *RequestMethodPickerModal) Init() {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	// Create GET, PUT, PATCH...
	for _, method := range http.RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				if m.listener != nil {
					m.listener.OnRequestMethodSelected(http.RequestMethod(mainText))
				}
			})
	}

	m.SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)
}
