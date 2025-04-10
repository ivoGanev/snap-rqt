package view

import (
	"snap-rq/internal/data"
	"snap-rq/internal/mocks"
	"snap-rq/internal/model"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type PageName string

const (
	PAGE_REQUEST_METHOD_PICKER_MODAL PageName = "request-method-picker"
	PAGE_LANDING_VIEW                PageName = "landing-view"
	ENABLE_DEBUG                     bool     = false
)

type App struct {
	*tview.Application
	Pages  *tview.Pages
	Views  *Views
	Models *Models
}

type Views struct {
	CollectionsView      *CollectionsView
	ResponseView         *ResponseView
	RequestsListView     *RequestsView
	Debugger             *tview.TextArea
	MethodSelectionModal *MethodSelectionModal
	UrlInput             *UrlInput
	NavHelp              *NavHelp
}

type Models struct {
	RequestsModel    *model.Requests
	CollectionsModel *model.Collections
}

func NewApp() *App {
	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
	}

	app.Views = &Views{
		CollectionsView:      NewColletionsView(&app),
		NavHelp:              NewNavigationHelp(&app),
		UrlInput:             NewUrlInput(&app),
		RequestsListView:     NewRequestsView(&app),
		ResponseView:         NewResponseView(&app),
		MethodSelectionModal: NewMethodSelectionModal(&app),
		Debugger:             tview.NewTextArea(),
	}

	app.Models = &Models{RequestsModel: model.NewRequestsModel()}

	return &app
}

func (app *App) Init() {
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		app.Views.Debugger.SetText(app.Views.RequestsListView.SelectedNode.String(), false)
		return false // Allow normal drawing to continue
	})

	// Load app data
	var store data.Store = &data.MockStore{}
	store.LoadAllCollections()
	requests := store.LoadAllRequests()

	loadedRequests := mocks.GenerateMockRequests(1000)

	app.Models.RequestsModel.SetAllData(loadedRequests)
	app.Models.CollectionsModel.SetCollections()

	// Handle model listeners
	app.Models.RequestsModel.AddListener(app.Views.RequestsListView)
	app.Models.CollectionsModel.AddListener(app.Views.CollectionsView)

	// Init layout and bind controllers
	app.Views.RequestsListView.AddListener(app.Views.UrlInput)
	app.Views.RequestsListView.Init()
	app.Views.ResponseView.Init()
	app.Views.NavHelp.Init()
	app.Views.UrlInput.Init()
	app.Views.CollectionsView.Init()
	app.Views.MethodSelectionModal.Init()
	app.Views.MethodSelectionModal.AddListener(app.Views.RequestsListView)

	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(app.Views.CollectionsView, 0, 1, false).
		AddItem(app.Views.RequestsListView, 0, 3, true).
		AddItem(app.Views.ResponseView, 0, 3, false)

	var body = tview.NewFlex()

	body.
		SetDirection(tview.FlexRow).
		AddItem(app.Views.NavHelp, 5, 0, false).
		AddItem(app.Views.UrlInput, 3, 0, false).
		AddItem(lrcontent, 0, 10, true)

	if ENABLE_DEBUG {
		body.AddItem(app.Views.Debugger, 0, 1, false)
	}

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
