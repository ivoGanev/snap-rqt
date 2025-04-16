package app

import (
	"context"
	"snap-rq/app/http"
	"time"
)

type RequestsController interface {
	HandleRequestNameSelected(RequestListItem)
	HandleRequestMethodSelected(RequestListItem)
	HandleSelectedRequestChanged(RequestListItem)
}

type RequestsViewController struct {
	*App
	RequestsService
}


func NewRequestsViewController(app *App, requestsService RequestsService) *RequestsViewController {
	return &RequestsViewController{app, requestsService}
}

func (r *RequestsViewController) HandleRequestMethodSelected(request RequestListItem) {
	r.showPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (r *RequestsViewController) HandleRequestNameSelected(request RequestListItem) {

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		request, err := r.RequestsService.GetRequest(request.Id)
		if err != nil {
			panic(err)
		}

		response := http.SendRequest(ctx, request.AsHttpRequest())
		r.App.QueueUpdateDraw(func() {
			r.App.Views.ResponseView.SetText(response, false)
		})

		cancel()
	}()
}

func (r *RequestsViewController) HandleSelectedRequestChanged(selectedRequest RequestListItem) {
	r.App.Views.UrlInputView.SetUrlText(selectedRequest.Url)
}
