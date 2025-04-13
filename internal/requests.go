package internal

import (
	"context"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/model"
	"time"
)

type RequestsViewController struct {
	*App
}

func NewRequestsViewController(app *App) *RequestsViewController {
	return &RequestsViewController{app}
}

func (r *RequestsViewController) OnAppModelsLoaded() {
	userSession := r.Models.UserSessionModel.GetUserSession()
	collection := r.Models.CollectionsModel.GetCollections()
	r.Views.RequestsView.DisplayCollection(collection)
}

func (r *RequestsViewController) OnRequestMethodSelected(request data.Request) {
	r.showPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (r *RequestsViewController) OnRequestSelected(request data.Request) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		response := http.SendRequest(ctx, request.Data)
		r.App.QueueUpdateDraw(func() {
			r.App.Views.ResponseView.SetText(response, false)
		})

		cancel()
	}()
}

func (r *RequestsViewController) OnRequestsListSelectionChanged(selectedRequest data.Request) {
	r.App.Views.UrlInput.SetUrlText(selectedRequest)
}