package view

import (
	"snap-rq/app/constants"

	"github.com/rivo/tview"
)

type MethodPickerListener interface {
	OnRequestMethodPickerSelected(method string)
}

func (r *MethodPickerModal) SetListener(listener MethodPickerListener) {
	r.listener = listener
}

// MethodPickerModal: A modal dialog that allows the user to pick a request method type, e.g. GET, PUT, PATCH, etc.
type MethodPickerModal struct {
	*tview.Grid
	listener MethodPickerListener
}

func NewMethodPickerModal() *MethodPickerModal {
	return &MethodPickerModal{
		Grid:     tview.NewGrid(),
	}
}

func (m *MethodPickerModal) Init() {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	// Create GET, PUT, PATCH...
	for _, method := range constants.RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				m.listener.OnRequestMethodPickerSelected(mainText)
			})
	}

	m.SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)
}
