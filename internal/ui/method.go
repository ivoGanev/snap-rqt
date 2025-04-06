package ui

import (
	"snap-rq/internal/http"

	"github.com/rivo/tview"
	"slices"
)

type OnMethodSelectionModalChangeListener interface {
	OnMethodSelectionChanged(string)
}

type MethodSelectionModal struct {
	*tview.Grid
	app       *App
	listeners []OnMethodSelectionModalChangeListener
}

func NewMethodSelectionModal(app *App) *MethodSelectionModal {
	return &MethodSelectionModal{
		app:  app,
		Grid: tview.NewGrid(),
	}
}

func (m *MethodSelectionModal) Init() {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	for _, method := range http.RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				for _, l := range m.listeners {
					l.OnMethodSelectionChanged(mainText)
				}
				m.Hide()
			})
	}

	m.SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)
}

func (m *MethodSelectionModal) Show() {
	m.app.ShowPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (m *MethodSelectionModal) Hide() {
	m.app.HidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (r *MethodSelectionModal) AddListener(l OnMethodSelectionModalChangeListener) {
	r.listeners = append(r.listeners, l)
}

func (r *MethodSelectionModal) RemoveListener(l OnMethodSelectionModalChangeListener) {
	for i, lis := range r.listeners {
		if lis == l {
			r.listeners = slices.Delete(r.listeners, i, i+1)
			return
		}
	}
}
