package app

import (
	"snap-rq/app/style"
	"snap-rq/app/view"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	VIEW_NAME_REQUESTS               = "requests"
	VIEW_NAME_RESPONSE               = "response"
	VIEW_NAME_COLLECTIONS            = "collections"
	PAGE_REQUEST_METHOD_PICKER_MODAL = "request-method-picker"
	PAGE_LANDING_VIEW                = "landing-view"
	ENABLE_DEBUG                     = false
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
	view.HotkeysHelp
	Debugger *tview.TextArea
}

type Controllers struct {
	*RequestMethodPickerViewController
	*RequestsViewController
	*UrlInputViewController
	*CollectionViewController
}

type Services struct {
	RequestsService
	CollectionService
	StateService
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
		CollectionViewController:          NewCollectionViewController(&app),
	}

	app.Views = &Views{
		CollectionsView:              NewColletionsView(app.Controllers.CollectionViewController),
		HotkeysHelp:                  view.NewHotkeysHelp(),
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

	// Load app data
	collections, err := app.Services.CollectionService.GetCollections()
	if err != nil {
		panic(err)
	}

	requests, err := app.Services.RequestsService.GetRequestListItems(collections[0].Id)
	if err != nil {
		panic(err)
	}
	app.Views.RequestsView.RenderRequests(requests)
	app.Views.CollectionsView.RenderCollections(collections)

	// Load and set app state
	appState := app.Services.StateService.GetState()
	app.Views.CollectionsView.SelectCollection(appState.GetSelectedCollectionRow())
	app.Views.RequestsView.SelectRequest(appState.GetSelectedRequestRow())

	// Init layout
	app.Views.RequestsView.Init()
	app.Views.ResponseView.Init()
	app.Views.HotkeysHelp.Init()
	app.Views.UrlInputView.Init()
	app.Views.CollectionsView.Init()
	app.Views.RequestMethodPickerModalView.Init()

	// Build Editor Layout
	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(app.Views.CollectionsView, 0, 1, true).
		AddItem(app.Views.RequestsView, 0, 3, false).
		AddItem(app.Views.ResponseView, 0, 3, false)

	var body = tview.NewFlex()

	body.
		SetDirection(tview.FlexRow).
		AddItem(app.Views.HotkeysHelp, 5, 0, false).
		AddItem(app.Views.UrlInputView, 3, 0, false).
		AddItem(lrcontent, 0, 10, true)

	if ENABLE_DEBUG {
		body.AddItem(app.Views.Debugger, 0, 1, false)
	}

	app.Pages.
		AddPage(string(PAGE_LANDING_VIEW), body, true, true).
		AddPage(string(PAGE_REQUEST_METHOD_PICKER_MODAL), app.Views.RequestMethodPickerModalView, true, false)

	// set hotkeys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'c' {
			// Collections focus hotkey
			app.focus(app.Views.CollectionsView)
			return nil
		}
		if event.Rune() == 'r' {
			// Requests focus hotkey
			app.focus(app.Views.RequestsView)
			return nil
		}
		if event.Rune() == 'q' {
			// Quit
			app.Stop()
			return nil
		}
		return event
	})

	// Start the app
	if err := app.
		SetFocus(app.Pages).
		SetRoot(app.Pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func (app *App) showPage(p string) {
	app.Pages.ShowPage(string(p))
}

func (app *App) hidePage(p string) {
	app.Pages.HidePage(string(p))
}

func (app *App) focus(p tview.Primitive) {
	app.SetFocus(p)
}
