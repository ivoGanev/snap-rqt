package view

import (
	"snap-rq/app/input"
	logger "snap-rq/app/log"
	"snap-rq/app/style"

	"github.com/rivo/tview"
)

type Views struct {
	CollectionsList  *CollectionsList
	ResponseWindow   *ResponseWindow
	RequestsList     *RequestsList
	EditorView       *EditorView
	RequestHeaderBar *RequestHeaderBar
	HotkeysHelp      *HotkeysHelp
	StatusBar        *StatusBar
	Debugger         *tview.TextArea
	NameEditorModal  *NameEditorModal
	puppetRow        *tview.Flex  // this is the middle row that can be morphed into catalog and edit view
}

type AppView struct {
	*tview.Application
	Pages        *tview.Pages
	Styles       style.StyleProvider
	Views        Views
	ViewMode     string
	listener     AppViewListener
	inputHandler *input.Handler
}

const (
	VIEW_NAME_REQUESTS    = "requests"
	VIEW_NAME_RESPONSE    = "response"
	VIEW_NAME_COLLECTIONS = "collections"
	PAGE_LANDING_VIEW     = "landing-view"
	PAGE_EDIT_NAME        = "edit-name"
	ENABLE_DEBUG          = false
	MODE_CATALOG_VIEW     = "lv"
	MODE_EDITOR_VIEW      = "ev"
)

type AppViewListener interface {
	OnViewModeChange(mode string)
}

func NewAppView(listener AppViewListener) AppView {

	var application = tview.NewApplication()

	var styleProvider = style.DefaultStylesProvider{}
	var inputHandler = input.NewHandler()

	var collectionListView = NewColletionsList(inputHandler)
	var hotkeyHelpView = NewHotkeysHelp()
	var editorView = NewEditorView(inputHandler)
	var requestHeaderBar = NewRequestHeaderBar(&styleProvider, inputHandler)
	var requestsListView = NewRequestsList(&styleProvider, inputHandler)
	var responseWindowView = NewResponseWindow(application)
	var statusBar = NewStatusBar()
	var nameEditor = NewNameEditorModal(inputHandler)
	var puppetRow = tview.NewFlex()

	var views = Views{
		CollectionsList:  collectionListView,
		HotkeysHelp:      hotkeyHelpView,
		RequestHeaderBar: requestHeaderBar,
		RequestsList:     requestsListView,
		ResponseWindow:   responseWindowView,
		StatusBar:        statusBar,
		EditorView:       editorView,
		Debugger:         tview.NewTextArea(),
		NameEditorModal:  nameEditor,
		puppetRow:        puppetRow,
	}

	appView := AppView{
		Application:  application,
		Pages:        tview.NewPages(),
		Views:        views,
		inputHandler: inputHandler,
		listener:     listener,
	}
	appView.Styles = &styleProvider

	return appView
}

func (app *AppView) Init() {
	views := app.Views

	// Build landing page
	puppetRow := views.puppetRow

	puppetRow.AddItem(app.Views.CollectionsList, 0, 1, false)
	puppetRow.AddItem(app.Views.RequestsList, 0, 3, false)
	puppetRow.AddItem(app.Views.ResponseWindow, 0, 3, false)

	app.ViewMode = MODE_CATALOG_VIEW

	var body = tview.NewFlex()

	body.SetDirection(tview.FlexRow)
	body.AddItem(views.HotkeysHelp, 5, 0, false)
	body.AddItem(views.RequestHeaderBar, 3, 0, false)
	body.AddItem(puppetRow, 0, 10, false)

	body.AddItem(views.StatusBar, 1, 0, false)

	if ENABLE_DEBUG {
		body.AddItem(views.Debugger, 0, 1, false)
	}
	
	app.Pages.AddPage(string(PAGE_LANDING_VIEW), body, true, true)
	app.Pages.AddPage(string(PAGE_EDIT_NAME), views.NameEditorModal, true, false)
	
	app.inputHandler.SetInputCapture(app.Application, input.SourceApp, func(action input.Action) {
		focus := app.GetFocus()
		switch action {
		case input.ActionSwapFocus:
			switch focus {
			case views.CollectionsList:
				app.Focus(views.RequestsList)
			case views.RequestsList:
				app.Focus(views.CollectionsList)
			}
		case input.ActionFocusCollections:
			app.Focus(views.CollectionsList)
		case input.ActionFocusRequests:
			app.Focus(views.RequestsList)
		case input.ActionSwapPuppetModes:
			if app.ViewMode == MODE_CATALOG_VIEW {
				app.MorphToEdit()
			} else {
				app.MorphToCatalog()
			}
		case input.ActionHeaderBarEditUrl:
			app.Focus(views.RequestHeaderBar.UrlInput)
		case input.ActionHeaderBarSelectMethod:
			app.Focus(views.RequestHeaderBar.RequestMethodDD)
		case input.ActionQuit:
			app.Stop()
		}
	})

	// Init UI
	views.CollectionsList.Init()
	views.HotkeysHelp.Init()
	views.RequestHeaderBar.Init()
	views.RequestsList.Init()
	views.ResponseWindow.Init()
	views.EditorView.Init(app)
}

func (app *AppView) Start() {
	// Start the app
	if err := app.
		SetRoot(app.Pages, true).
		SetFocus(app.Views.CollectionsList).
		EnableMouse(true).
		Run(); err != nil {
		logger.Println(err)
		panic(err)
	}
}

func (app *AppView) MorphToCatalog() {
	if app.ViewMode != MODE_CATALOG_VIEW {
		puppetRow := app.Views.puppetRow
		puppetRow.RemoveItem(app.Views.EditorView)
		puppetRow.RemoveItem(app.Views.ResponseWindow)

		puppetRow.
			AddItem(app.Views.CollectionsList, 0, 1, true).
			AddItem(app.Views.RequestsList, 0, 3, false).
			AddItem(app.Views.ResponseWindow, 0, 3, false)

		app.Focus(app.Views.RequestsList)
		app.ViewMode = MODE_CATALOG_VIEW
		app.listener.OnViewModeChange(app.ViewMode)
	}
}

func (app *AppView) MorphToEdit() {
	if app.ViewMode != MODE_EDITOR_VIEW {
		puppetRow := app.Views.puppetRow
		puppetRow.RemoveItem(app.Views.CollectionsList)
		puppetRow.RemoveItem(app.Views.RequestsList)
		puppetRow.RemoveItem(app.Views.ResponseWindow)

		puppetRow.
			AddItem(app.Views.EditorView, 0, 4, true).
			AddItem(app.Views.ResponseWindow, 0, 3, false)

		app.Focus(app.Views.EditorView)
		app.ViewMode = MODE_EDITOR_VIEW
		app.listener.OnViewModeChange(app.ViewMode)

		// Select body/header button based on previous mode
		lastEditorMode := app.Views.EditorView.GetCurrentMode()
		if lastEditorMode == EDITOR_VIEW_FOCUS_HEADERS {
			app.Focus(app.Views.EditorView.HeadersButton)
		} else {
			app.Focus(app.Views.EditorView.BodyButton)
		}
	}
}

func (app *AppView) ShowPage(p string) {
	app.Pages.ShowPage(string(p))
}

func (app *AppView) HidePage(p string) {
	app.Pages.HidePage(string(p))
}

func (app *AppView) Focus(p tview.Primitive) {
	app.SetFocus(p)
}
