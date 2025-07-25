package controller

import (
	"fmt"
	"snap-rq/app/entity"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

type AppController struct {
	app     *view.AppView
	views   *view.Views
	service *service.AppService
}


func NewAppController(app view.AppView, appService *service.AppService) AppController {
	var controller = AppController{
		&app,
		&app.Views,
		appService,
	}

	return controller
}

func (a *AppController) Start() {
	// Load and render app data on load
	d := a.service.FetchLandingData()
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.RequestsList.SelectRequest(d.SelectedRequestId)
	a.views.CollectionsList.RenderCollections(d.Collections)
	a.views.CollectionsList.SelectCollection(d.SelectedCollectionId)
}

// App View


func (a *AppController) OnViewModeChange(mode string) {
	request := a.service.GetFocusedRequest()
	a.views.EditorView.SetTextArea(request)
}


// Editor View

func (a *AppController) OnEditTextArea(editorMode int, edit string) {
	// change the body|header of current HTTP method selected
	switch editorMode {
	case view.EDITOR_VIEW_MODE_BODY:
		a.service.UpdateFocusedRequest(entity.ModRequest{Body: &edit})
	case view.EDITOR_VIEW_MODE_HEADERS:
		a.service.UpdateFocusedRequest(entity.ModRequest{Headers: &edit})
	}
}

func (a *AppController) OnEditorModeChanged() {
	request := a.service.GetFocusedRequest()
	a.views.EditorView.SetTextArea(request)
}


func (a *AppController) OnUrlInputTextChanged(urlText string) {
	a.service.UpdateFocusedRequest(entity.ModRequest{Url: &urlText})
}

// Landing View (Request List)


func (a *AppController) OnRequestMethodPickerSelected(method string) {
	d := a.service.FetchBasicFocusData()
	a.views.RequestsList.ChangeRequestMethod(d.SelectedRequestId, method)
	a.app.HidePage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
	a.app.Focus(a.views.RequestsList)
}

func (a *AppController) OnRequestListMethodSelected(entity.RequestBasic) {
	a.app.ShowPage(view.PAGE_REQUEST_METHOD_PICKER_MODAL)
}

func (a *AppController) OnRequestListRequestFocusChanged(selectedRequest entity.RequestBasic) {
	a.service.ChangeFocusedRequest(selectedRequest)
	a.views.UrlInputView.SetUrlText(selectedRequest.Url)

	// Once the user changes the selection, load the historical response from memory and set it
	// TODO: Clean setting the empty response;
	go func() {
		a.service.CancelSentHttpRequest() // Cancel any requests to prevent any side-effects.
		// But at some point, we would not want to cancel the entire request but rather the side effects...Is there a better way to do this?
		a.views.ResponseWindow.SetHttpResponse(entity.HttpResponse{})
	}()
}

func (a *AppController) OnRequestListNameSelected(selected entity.RequestBasic) {
	s := fmt.Sprintf("%s %s", selected.MethodType, selected.Url)
	a.views.StatusBar.SetText(s)

	onHttpResult := func(result entity.HttpResult) {
		if result.Error != nil {
			a.views.ResponseWindow.SetError(result.Error)
		} else {
			a.views.ResponseWindow.SetHttpResponse(result.Response)
		}
	}
	a.service.SendHttpRequest(selected.Id, onHttpResult)
	a.views.ResponseWindow.AwaitResponse()
}

func (a *AppController) OnRequestListAdd(position int) {
	a.service.AddRequest(position)
	d := a.service.FetchBasicFocusData()
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.StatusBar.SetText("Added new request")
}

func (a *AppController) OnRequestListDuplicate(entity.RequestBasic) {
}

func (a *AppController) OnRequestListRemove(request entity.RequestBasic, position int) {
	a.service.RemoveRequest(request.Id, position)
	d := a.service.FetchBasicFocusData()
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	
	s := fmt.Sprintf("Removed request %s", request.Name)
	a.views.StatusBar.SetText(s)
}


// Collection list (Request List)


func (a *AppController) OnFocusedCollectionChanged(changedCollection entity.Collection) {
	d := a.service.ChangeFocusedCollection(changedCollection.Id)

	// when user selects a collection, a request item would be automatically changed
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.RequestsList.SelectRequest(d.SelectedRequestId)
}
