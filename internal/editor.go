package internal

import (
	"snap-rq/internal/http"

	"github.com/rivo/tview"
)

type Editor struct {
	controller *EditorController
	state      *EditorSessionState
}

type EditorSessionState struct {
}

type EditorController struct {
	app    *App
	models *Models
	views  *Views
	pages  *tview.Pages
}

func NewEditor(app *App) Editor {
	controller := EditorController{
		app:    app,
		models: app.Models,
		views:  app.Views,
		pages:  app.Pages,
	}
	editor := Editor{
		controller: &controller,
	}
	return editor
}

func (e *EditorController) OnMethodSelectionChanged(method string) {
	e.SelectedNode.Data.Method = http.RequestMethod(method)

	methodText := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(http.RequestMethod(method)), method)
	e.RequestsView.GetCell(r.SelectedRow, 0).SetText(methodText)
	e.app.SetFocus(r)
}

func (e *EditorController) showPage(p PageName) {
	e.pages.ShowPage(string(p))
}

func (e *EditorController) hidePage(p PageName) {
	e.pages.HidePage(string(p))
}

func (e *EditorController) focus(p tview.Primitive) {
	e.app.SetFocus(p)
}
