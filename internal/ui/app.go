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
	Pages  *tview.Pages
	Views  *Views
	Models *Models
}

type Views struct {
	ResponseWindow       *ResponseView
	RequestsListWindow   *RequestsView
	Debugger             *tview.TextArea
	MethodSelectionModal *MethodSelectionModal
	UrlInput             *UrlInput
}

type Models struct {
	RequestsModel *model.Requests
}

func NewApp() *App {
	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
	}

	app.Views = &Views{
		Debugger:             tview.NewTextArea(),
		UrlInput:             NewUrlInput(&app),
		RequestsListWindow:   NewRequestsView(&app),
		ResponseWindow:       NewResponseView(&app),
		MethodSelectionModal: NewMethodSelectionModal(&app),
	}

	app.Models = &Models{RequestsModel: model.NewRequestsModel()}

	return &app
}

func (app *App) Init() {
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		app.Views.Debugger.SetText(app.Views.RequestsListWindow.SelectedNode.String(), false)
		return false // Allow normal drawing to continue
	})

	// Construct the app
	app.Models.RequestsModel.AddListener(app.Views.RequestsListWindow)
	app.Models.RequestsModel.SetAllData(mocks.GenerateMockRequests(1000))

	// Init the layout configuration
	app.Views.RequestsListWindow.AddListener(app.Views.UrlInput)
	app.Views.RequestsListWindow.Init()
	app.Views.ResponseWindow.Init()
	app.Views.UrlInput.Init()
	app.Views.MethodSelectionModal.Init()
	app.Views.MethodSelectionModal.AddListener(app.Views.RequestsListWindow)

	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(app.Views.RequestsListWindow, 0, 1, true).
		AddItem(app.Views.ResponseWindow, 0, 1, false)

	var body = tview.NewFlex()

	body.
		SetDirection(tview.FlexRow).
		AddItem(app.Views.UrlInput, 3, 0, false).
		AddItem(lrcontent, 0, 10, true).
		AddItem(app.Views.Debugger, 0, 1, false)

	app.Pages.
		AddPage(string(PAGE_LANDING_VIEW), body, true, true).
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
