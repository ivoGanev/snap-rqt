package controller

import (
	"snap-rq/app/entity"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

type AppController struct {
	rootView   view.App
	appService service.AppService
}

func NewAppController(rootView view.App, appService service.AppService) AppController {
	var controller = AppController{
		rootView,
		appService,
	}

	return controller
}

func (a *AppController) Start() {
	// Load and render app data on load
	d := a.appService.FetchLandingData()
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.RequestsList.SelectRequest(d.SelectedRequestRow)
	a.rootView.Views.CollectionsList.RenderCollections(d.Collections)
	a.rootView.Views.CollectionsList.SelectCollection(d.SelectedCollectionRow)
}

func (a *AppController) OnUrlInputTextChanged(urlText string) {
	a.appService.UpdateFocusedRequest(entity.PatchRequest{Url: &urlText})
}

func (a *AppController) OnRequestMethodPickerSelected(method string) {
	d := a.appService.FetchBasicFocusData()
	a.rootView.Views.RequestsList.ChangeMethodTypeOnSelectedRow(d.SelectedRequestRow, method)
	a.rootView.HidePage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
	a.rootView.Focus(a.rootView.Views.RequestsList)
}

func (a *AppController) OnRequestMethodSelected(entity.RequestBasic) {
	a.rootView.ShowPage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (a *AppController) OnRequestNameSelected(selected entity.RequestBasic) {
	go func() {
		response := a.appService.SendHttpRequestById(selected.Id)
		a.rootView.QueueUpdateDraw(func() {
			a.rootView.Views.ResponseWindow.SetText(response, false)
		})
	}()
}

func (a *AppController) OnSelectedRequestChanged(selectedRequest entity.RequestBasic) {
	a.appService.ChangeFocusedRequest(selectedRequest)
	a.rootView.Views.UrlInputView.SetUrlText(selectedRequest.Url)
}

func (a *AppController) OnCollectionChanged(changedCollection entity.FocusedCollection) {
	d := a.appService.ChangeFocusedCollection(changedCollection)

	// when user selects a collection, a request item would be automatically changed
	a.rootView.Views.RequestsList.RenderRequests(d.RequestsBasic)
	a.rootView.Views.RequestsList.SelectRequest(d.SelectedRequestRow)
}
