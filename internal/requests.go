package internal

import (
	"context"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"time"
)

type RequestsViewController struct {
	*App
}

func (r *RequestsViewController) OnRequestModelChanged(request *data.RequestSimple) {
	// when a request changes, we need to update the text
	panic("unimplemented")
}

func NewRequestsViewController(app *App) *RequestsViewController {
	return &RequestsViewController{app}
}

// func (r *RequestsViewController) OnAppModelsLoaded() {
// 	userSession := r.Models.UserSessionModel.GetUserSession()
// 	collection := r.Models.CollectionsModel.GetCollection(userSession.CollectionId)
// 	r.Views.RequestsView.DisplayCollection(collection)
// }

func (r *RequestsViewController) RenderRequests(requests []data.RequestSimple) {
	r.App.Views.RequestsView.RenderRequests(requests)
}

func (r *RequestsViewController) OnRequestMethodSelected(request data.RequestSimple) {
	r.showPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (r *RequestsViewController) OnRequestSelected(request data.RequestSimple) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		request := r.Models.ProjectModel.GetRequest(request.Id)
		response := http.SendRequest(ctx, request.Data)
		r.App.QueueUpdateDraw(func() {
			r.App.Views.ResponseView.SetText(response, false)
		})

		cancel()
	}()
}

func (r *RequestsViewController) OnRequestsListSelectionChanged(selectedRequest data.RequestSimple) {
	r.App.Views.UrlInput.SetUrlText(selectedRequest.Name)
}
