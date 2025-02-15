package ui

import (
	"snap-rq/internal/mocks"
	"snap-rq/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageName string

const (
	PAGE_REQUEST_METHOD_PICKER_MODAL PageName = "request-method-picker"
	PAGE_LANDING_VIEW                PageName = "landing-view"
)

type App struct {
	*tview.Application
	Pages *tview.Pages
	Views *Views
}

type Views struct {
	Response             *ResponseView
	Requests             *RequestsView
	Debugger             *tview.TextArea
	MethodSelectionModal *MethodSelectionModal
}

func NewApp() *App {
	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
	}

	app.Views = &Views{
		Debugger:             tview.NewTextArea(),
		Requests:             NewRequestsView(&app),
		Response:             NewResponseView(&app),
		MethodSelectionModal: NewMethodSelectionModal(&app),
	}

	return &app
}

func (app *App) Init() {
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		app.Views.Debugger.SetText(app.Views.Requests.SelectedNode.String(), false)
		return false // Allow normal drawing to continue
	})
    
	// Construct the app
	requestsModel := model.NewRequests()
	requestsModel.AddListener(app.Views.Requests)
	requestsModel.SetData(mocks.GenerateMockRequests(1000))

	// Init the layout configuration
	app.Views.Requests.Init()
	app.Views.Response.Init()
	app.Views.MethodSelectionModal.Init()
	app.Views.MethodSelectionModal.AddListener(app.Views.Requests)

	var mainWindows = tview.NewFlex()
	mainWindows.
		AddItem(app.Views.Requests, 0, 1, true).
		AddItem(app.Views.Response, 0, 1, false)

	var rootFlex = tview.NewFlex()

	rootFlex.
		SetDirection(tview.FlexRow).
		AddItem(mainWindows, 0, 10, true).
		AddItem(app.Views.Debugger, 0, 1, false)

	app.Pages.
		AddPage(string(PAGE_LANDING_VIEW), rootFlex, true, true).
		AddPage(string(PAGE_REQUEST_METHOD_PICKER_MODAL), app.Views.MethodSelectionModal, true, false)

	if err := app.
		SetFocus(app.Pages).
		SetRoot(app.Pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (app *App) ShowPage(p PageName) {
	app.Pages.ShowPage(string(p))
}

func (app *App) HidePage(p PageName) {
	app.Pages.HidePage(string(p))
}
