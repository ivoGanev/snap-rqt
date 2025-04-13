package internal

import (
	"snap-rq/internal/data"
	"snap-rq/internal/model"
	"snap-rq/internal/styles"
	"snap-rq/internal/view"

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
	Pages         *tview.Pages
	Views         *Views
	Models        *Models
	Controllers   *Controllers
	Store         *data.Store
	StyleProvider styles.StyleProvider
}

type Views struct {
	CollectionsView          *view.CollectionsView
	ResponseView             *view.ResponseView
	RequestsView             *view.RequestsView
	Debugger                 *tview.TextArea
	RequestMethodPickerModal *view.RequestMethodPickerModal
	UrlInput                 *view.UrlInput
	NavHelp                  *view.NavHelp
}

type Models struct {
	*model.CollectionsModel
	*model.RequestsModel
	*model.UserSessionModel
}

type Controllers struct {
	*RequestMethodPickerController
	*RequestsViewController
	*UrlInputController
}

type OnAppModelsLoadedListener interface {
	OnAppModelsLoaded()
}

func NewApp() *App {
	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
	}

	app.StyleProvider = &styles.Default{}

	app.Views = &Views{
		CollectionsView:          view.NewColletionsView(),
		NavHelp:                  view.NewNavigationHelp(),
		UrlInput:                 view.NewUrlInput(),
		RequestsView:             view.NewRequestsView(app.StyleProvider),
		ResponseView:             view.NewResponseView(),
		RequestMethodPickerModal: view.NewMethodPickerModal(),
		Debugger:                 tview.NewTextArea(),
	}

	app.Models = &Models{
		CollectionsModel: model.NewCollectionModel(*app.Store),
		RequestsModel:    model.NewRequestModel(*app.Store),
		UserSessionModel: model.NewUserSessionModel(*app.Store),
	}

	app.Controllers = &Controllers{
		RequestMethodPickerController: NewMethodPickerModalController(&app),
		RequestsViewController:        NewRequestsViewController(&app),
		UrlInputController:            NewUrlInputController(&app),
	}

	// app.Store = &data.MockStore{}
	return &app
}

func (app *App) Init() {
	app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		return false // Allow normal drawing to continue
	})

	// Bind app listeners
	var r []OnAppModelsLoadedListener
	r[0] = app.Views.RequestsView

	// Bind model listeners
	app.Models.RequestsModel.AddListener(app.Controllers.RequestsViewController)
	// app.Models.CollectionsModel.AddListener(app.Controllers.Coll)

	// Init layout
	app.Views.RequestsView.Init()
	app.Views.ResponseView.Init()
	app.Views.NavHelp.Init()
	app.Views.UrlInput.Init()
	app.Views.CollectionsView.Init()
	app.Views.RequestMethodPickerModal.Init()

	// Bind controllers
	app.Views.RequestMethodPickerModal.SetRequestMethodPickerListener(app.Controllers.RequestMethodPickerController)
	app.Views.RequestsView.SetRequestsViewListener(app.Controllers.RequestsViewController)

	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(app.Views.CollectionsView, 0, 1, false).
		AddItem(app.Views.RequestsView, 0, 3, true).
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
		AddPage(string(PAGE_REQUEST_METHOD_PICKER_MODAL), app.Views.RequestMethodPickerModal, true, false)

	// Load app data
	app.Models.RequestsModel.Load()
	app.Models.CollectionsModel.Load()
	app.Models.UserSessionModel.Load()

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
