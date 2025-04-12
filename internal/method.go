package internal

import (
	"fmt"
	"snap-rq/internal/controller"
	"snap-rq/internal/http"

	"github.com/rivo/tview"
)

// type OnMethodSelectionModalChangeListener interface {
// 	OnMethodSelectionChanged(string)
// }

type MethodSelectionModal struct {
	*tview.Grid
	controller controller.MethodController
	// listeners []OnMethodSelectionModalChangeListener
}

func NewMethodSelectionModal(controller controller.MethodController) *MethodSelectionModal {
	return &MethodSelectionModal{
		controller: controller,
		Grid:       tview.NewGrid(),
	}
}

func (m *MethodSelectionModal) Init() {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	for _, method := range http.RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				m.editor.HandleMethodSelectionChanged(mainText)
				// m.OnMethodSelectionChanged(mainText)
				// for _, l := range m.listeners {
				// 	l.OnMethodSelectionChanged(mainText)
				// }
				m.Hide()
			})
	}

	m.SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)
}

func (m *MethodSelectionModal) OnMethodSelectionChanged(method string) {
	m.SelectedNode.Data.Method = http.RequestMethod(method)

	methodText := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(http.RequestMethod(method)), method)
	m.RequestsView.GetCell(r.SelectedRow, 0).SetText(methodText)
	m.app.SetFocus(r)
}

// func (m *MethodSelectionModal) Show() {
// 	m.app.ShowPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
// }

// func (m *MethodSelectionModal) Hide() {
// 	m.app.HidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
// }

// func (r *MethodSelectionModal) AddListener(l OnMethodSelectionModalChangeListener) {
// 	r.listeners = append(r.listeners, l)
// }

// func (r *MethodSelectionModal) RemoveListener(l OnMethodSelectionModalChangeListener) {
// 	for i, lis := range r.listeners {
// 		if lis == l {
// 			r.listeners = slices.Delete(r.listeners, i, i+1)
// 			return
// 		}
// 	}
// }
