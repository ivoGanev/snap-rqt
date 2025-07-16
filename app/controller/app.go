package controller

import (
	"fmt"
	"snap-rq/app/entity"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

type AppController struct {
	rootView   *view.AppView
	appService *service.AppService
}

func NewAppController(rootView view.AppView, appService *service.AppService) AppController {
	var controller = AppController{
		&rootView,
		appService,
	}

	return controller
}

func (a *AppController) Start() {
	// Load and render app data on load
	d := a.appService.FetchLandingData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.RequestsList.SelectRequest(d.SelectedRequestId)
	a.rootView.Views.CollectionsList.RenderCollections(d.Collections)
	a.rootView.Views.CollectionsList.SelectCollection(d.SelectedCollectionId)
}

func (a *AppController) OnUrlInputTextChanged(urlText string) {
	a.appService.UpdateFocusedRequest(entity.ModRequest{Url: &urlText})
}

func (a *AppController) OnRequestMethodPickerSelected(method string) {
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.ChangeRequestMethod(d.SelectedRequestId, method)
	a.rootView.HidePage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
	a.rootView.Focus(a.rootView.Views.RequestsList)
}

func (a *AppController) OnRequestListMethodSelected(entity.RequestBasic) {
	a.rootView.ShowPage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (a *AppController) OnRequestListRequestFocusChanged(selectedRequest entity.RequestBasic) {
	a.appService.ChangeFocusedRequest(selectedRequest)
	a.rootView.Views.UrlInputView.SetUrlText(selectedRequest.Url)

	// Once the user changes the selection, load the historical response from memory and set it
	// TODO: Clean setting the empty response;
	go func() {
		a.appService.CancelSentHttpRequest() // Cancel any requests to prevent any side-effects.
		// But at some point, we would not want to cancel the entire request but rather the side effects...Is there a better way to do this?
		a.rootView.Views.ResponseWindow.SetHttpResponse(entity.HttpResponse{})
	}()
}

func (a *AppController) OnRequestListNameSelected(selected entity.RequestBasic) {
	s := fmt.Sprintf("%s %s", selected.MethodType, selected.Url)
	a.rootView.Views.StatusBar.SetText(s)

	onHttpResult := func(result entity.HttpResult) {
		if result.Error != nil {
			a.rootView.Views.ResponseWindow.SetError(result.Error)
		} else {
			a.rootView.Views.ResponseWindow.SetHttpResponse(result.Response)
		}
	}
	a.appService.SendHttpRequest(selected.Id, onHttpResult)
	a.rootView.Views.ResponseWindow.AwaitResponse()
}

func (a *AppController) OnRequestListAdd(position int) {
	a.appService.AddRequest(position)
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.StatusBar.SetText("Added new request")
}

func (a *AppController) OnRequestListDuplicate(entity.RequestBasic) {
}

func (a *AppController) OnRequestListRemove(request entity.RequestBasic, position int) {
	a.appService.RemoveRequest(request.Id, position)
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)

	s := fmt.Sprintf("Removed request %s", request.Name)
	a.rootView.Views.StatusBar.SetText(s)
}

func (a *AppController) OnFocusedCollectionChanged(changedCollection entity.Collection) {
	d := a.appService.ChangeFocusedCollection(changedCollection.Id)

	// when user selects a collection, a request item would be automatically changed
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.RequestsList.SelectRequest(d.SelectedRequestId)
}
