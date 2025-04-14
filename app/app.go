package app

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"snap-rq/app/style"
)

type PageName string

const (
	PAGE_REQUEST_METHOD_PICKER_MODAL PageName = "request-method-picker"
	PAGE_LANDING_VIEW                PageName = "landing-view"
	ENABLE_DEBUG                     bool     = false
)

type App struct {
	*tview.Application
	*tview.Pages
	*Views
	*Controllers
	*Services
	style.StyleProvider
}

type Views struct {
	CollectionsView
	ResponseView
	RequestsView
	tview.TextArea
	RequestMethodPickerModalView
	UrlInputView
	NavHelpView
	Debugger *tview.TextArea
}

type Controllers struct {
	*RequestMethodPickerViewController
	*RequestsViewController
	*UrlInputViewController
}

type Services struct {
	RequestsService
	CollectionService
}

type OnAppModelsLoadedListener interface {
	OnAppModelsLoaded()
}

func NewApp(services *Services) App {
	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
		Services:    services,
	}

	app.StyleProvider = &style.DefaultStylesProvider{}

	app.Controllers = &Controllers{
		RequestMethodPickerViewController: NewMethodPickerModalController(&app),
		RequestsViewController:            NewRequestsViewController(&app, services),
		UrlInputViewController:            NewUrlInputController(&app, services),
	}

	app.Views = &Views{
		CollectionsView:              NewColletionsView(),
		NavHelpView:                  NewNavigationHelpView(),
		UrlInputView:                 NewUrlInputView(app.Controllers.UrlInputViewController),
		RequestsView:                 NewRequestsView(app.StyleProvider, app.Controllers.RequestsViewController),
		ResponseView:                 NewResponseView(),
		RequestMethodPickerModalView: NewMethodPickerModal(app.Controllers.RequestMethodPickerViewController),
		Debugger:                     tview.NewTextArea(),
	}

	return app
}

func (app *App) Init() {
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		return false // Allow normal drawing to continue
	})

	// Load initial data
	requests, err := app.Services.RequestsService.GetRequestListItems()
	if err != nil {
		panic(err)
	}
	app.Views.RequestsView.RenderRequests(requests)
	app.Views.RequestsView.SelectRequest(0)

	collections, err := app.Services.CollectionService.GetCollections()
	if err != nil {
		panic(err)
	}
	app.Views.CollectionsView.RenderCollections(collections)

	// Init layout
	app.Views.RequestsView.Init()
	app.Views.ResponseView.Init()
	app.Views.NavHelpView.Init()
	app.Views.UrlInputView.Init()
	app.Views.CollectionsView.Init()
	app.Views.RequestMethodPickerModalView.Init()

	// Build Editor Layout
	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(app.Views.CollectionsView, 0, 1, false).
		AddItem(app.Views.RequestsView, 0, 3, true).
		AddItem(app.Views.ResponseView, 0, 3, false)

	var body = tview.NewFlex()

	body.
		SetDirection(tview.FlexRow).
		AddItem(app.Views.NavHelpView, 5, 0, false).
		AddItem(app.Views.UrlInputView, 3, 0, false).
		AddItem(lrcontent, 0, 10, true)

	if ENABLE_DEBUG {
		body.AddItem(app.Views.Debugger, 0, 1, false)
	}

	app.Pages.
		AddPage(string(PAGE_LANDING_VIEW), body, true, true).
		AddPage(string(PAGE_REQUEST_METHOD_PICKER_MODAL), app.Views.RequestMethodPickerModalView, true, false)

	// Start the app
	if err := app.
		SetFocus(app.Pages).
		SetRoot(app.Pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (app *App) showPage(p PageName) {
	app.Pages.ShowPage(string(p))
}

func (app *App) hidePage(p PageName) {
	app.Pages.HidePage(string(p))
}

func (app *App) focus(p tview.Primitive) {
	app.SetFocus(p)
}
