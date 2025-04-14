package app

import (
	"github.com/rivo/tview"
)

// RequestMethodPickerModalView: A modal dialog that allows the user to pick a request method type, e.g. GET, PUT, PATCH, etc.
type RequestMethodPickerModalView struct {
	*tview.Grid
	controller RequestMethodPickerController
}

func NewMethodPickerModal(controller RequestMethodPickerController) RequestMethodPickerModalView {
	return RequestMethodPickerModalView{
		Grid:       tview.NewGrid(),
		controller: controller,
	}
}

func (m *RequestMethodPickerModalView) Init() {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	// Create GET, PUT, PATCH...
	for _, method := range RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				m.controller.HandleRequestMethodSelected(mainText)
			})
	}

	m.SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)
}
