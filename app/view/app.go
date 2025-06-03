package view

import (
	"snap-rq/app/style"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Views struct {
	CollectionsList   *CollectionsList
	ResponseWindow    *ResponseWindow
	RequestsList      *RequestsList
	MethodPickerModal *MethodPickerModal
	UrlInputView      *UrlInputView
	HotkeysHelp       *HotkeysHelp
	StatusBar         *StatusBar
	Debugger          *tview.TextArea
}

type App struct {
	*tview.Application
	Pages  *tview.Pages
	Styles style.StyleProvider
	Views  Views
}

const (
	VIEW_NAME_REQUESTS               = "requests"
	VIEW_NAME_RESPONSE               = "response"
	VIEW_NAME_COLLECTIONS            = "collections"
	PAGE_REQUEST_METHOD_PICKER_MODAL = "request-method-picker"
	PAGE_LANDING_VIEW                = "landing-view"
	ENABLE_DEBUG                     = false
)

func NewApp() App {
	var styleProvider = style.DefaultStylesProvider{}

	var collectionListView = NewColletionsList()
	var hotkeyHelpView = NewHotkeysHelp()
	var urlInputView = NewUrlInput()
	var requestsListView = NewRequestsList(&styleProvider)
	var responseWindowView = NewResponseWindow()
	var methodPickerView = NewMethodPickerModal()
	var statusBar = NewStatusBar()

	var views = Views{
		CollectionsList:   collectionListView,
		HotkeysHelp:       hotkeyHelpView,
		UrlInputView:      urlInputView,
		RequestsList:      requestsListView,
		ResponseWindow:    responseWindowView,
		MethodPickerModal: methodPickerView,
		StatusBar:         statusBar,
		Debugger:          tview.NewTextArea(),
	}

	app := App{
		Application: tview.NewApplication(),
		Pages:       tview.NewPages(),
		Views:       views,
	}
	app.Styles = &styleProvider

	// Render UI
	collectionListView.Init()
	hotkeyHelpView.Init()
	urlInputView.Init()
	requestsListView.Init()
	responseWindowView.Init()
	methodPickerView.Init()

	return app
}

func (app *App) Init() {
	views := app.Views

	// Build Editor Layout
	var lrcontent = tview.NewFlex()
	lrcontent.
		AddItem(views.CollectionsList, 0, 1, true).
		AddItem(views.RequestsList, 0, 3, false).
		AddItem(views.ResponseWindow, 0, 3, false)

	var body = tview.NewFlex()

	body.
		SetDirection(tview.FlexRow).
		AddItem(views.HotkeysHelp, 5, 0, false).
		AddItem(views.UrlInputView, 3, 0, false).
		AddItem(lrcontent, 0, 10, true)

	body.AddItem(views.StatusBar, 2, 0, false)

	if ENABLE_DEBUG {
		body.AddItem(views.Debugger, 0, 1, false)
	}

	app.Pages.
		AddPage(string(PAGE_LANDING_VIEW), body, true, true).
		AddPage(string(PAGE_REQUEST_METHOD_PICKER_MODAL), views.MethodPickerModal, true, false)

	// set hotkeys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Swap collecti0n/request views
		focus := app.GetFocus()
		if focus != nil && event.Key() == tcell.KeyTAB {
			if focus == views.CollectionsList {
				app.Focus(views.RequestsList)
				return nil
			} else if focus == views.RequestsList {
				app.Focus(views.CollectionsList)
				return nil
			}
		}
		if event.Rune() == 'c' {
			// Collections focus hotkey
			app.Focus(views.CollectionsList)
			return nil
		}
		if event.Rune() == 'r' {
			// Requests focus hotkey
			app.Focus(views.RequestsList)
			return nil
		}
		if event.Key() == tcell.KeyEsc || event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyESC {
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

func (app *App) ShowPage(p string) {
	app.Pages.ShowPage(string(p))
}

func (app *App) HidePage(p string) {
	app.Pages.HidePage(string(p))
}

func (app *App) Focus(p tview.Primitive) {
	app.SetFocus(p)
}
