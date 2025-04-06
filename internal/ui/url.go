package ui

import (
	"github.com/rivo/tview"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
)

type UrlInput struct {
	*tview.InputField
	app *App
	selectedNode *data.Node[http.Request]
}

func (r *UrlInput) OnRequestsListSelectionChanged(selection *data.Node[http.Request]) {
	r.SetText(selection.Data.Url)
	r.selectedNode = selection
}

func NewUrlInput(app *App) *UrlInput {
	return &UrlInput{
		InputField: tview.NewInputField(),
		app: app,
	}
}

func (r *UrlInput) Init() {
	r.SetTitle("Url")
	r.SetBorder(true)
	r.SetChangedFunc(func(text string){
		if r.selectedNode != nil {
			r.selectedNode.Data.Url = text
			r.app.Models.RequestsModel.SetData(r.selectedNode.Id, r.selectedNode)
		}
	})
}
