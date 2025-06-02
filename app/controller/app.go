package controller

import (
	"snap-rq/app/entity"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

type AppController struct {
	rootView   *view.App
	appService *service.AppService
}


func NewAppController(rootView view.App, appService service.AppService) AppController {
	var controller = AppController{
		&rootView,
		&appService,
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
	a.appService.UpdateFocusedRequest(entity.PatchRequest{Url: &urlText})
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

func (a *AppController) OnRequestListNameSelected(selected entity.RequestBasic) {
	go func() {
		response := a.appService.SendHttpRequestById(selected.Id)
		a.rootView.QueueUpdateDraw(func() {
			a.rootView.Views.ResponseWindow.SetText(response, false)
		})
	}()
}

func (a *AppController) OnRequestListAdd(position int) {
	a.appService.AddRequest(position)
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
}

func (a *AppController) OnRequestListDuplicate(entity.RequestBasic) {
}

func (a *AppController) OnRequestListRemove(request entity.RequestBasic, position int) {
	a.appService.RemoveRequest(request.Id, position)
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
}

func (a *AppController) OnRequestListRequestFocusChanged(selectedRequest entity.RequestBasic) {
	a.appService.ChangeFocusedRequest(selectedRequest)
	a.rootView.Views.UrlInputView.SetUrlText(selectedRequest.Url)
}

func (a *AppController) OnFocusedCollectionChanged(changedCollection entity.Collection) {
	d := a.appService.ChangeFocusedCollection(changedCollection.Id)

	// when user selects a collection, a request item would be automatically changed
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.RequestsList.SelectRequest(d.SelectedRequestId)
}
