package controller

import (
	"fmt"
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
	"snap-rq/app/service"
	"snap-rq/app/view"
)

const APP_CONTROLLER_LOG_TAG = "[App Controller]"

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
	logger.Debug("Controller starting")

	// Load where the user last left
	d := a.service.GetBasicFocusData()
	a.views.CollectionsList.RenderCollections(d.Collections)
	a.views.CollectionsList.SelectCollection(d.SelectedCollection)
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.RequestsList.SelectRequest(d.SelectedRequest)
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

// Landing View (Request Header Bar)

func (a *AppController) OnMethodSelection(method string) {
	a.service.UpdateFocusedRequest(entity.ModRequest{Method: &method})

	d := a.service.GetBasicFocusData()
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
}

// Landing View (Request List)

func (a *AppController) OnRequestListRequestFocusChanged(selectedRequest entity.RequestBasic) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User selected request", selectedRequest)
	a.service.ChangeFocusedRequest(selectedRequest)
	focusedRequest := a.service.GetFocusedRequest()

	a.views.RequestHeaderBar.SetUrlText(focusedRequest.Url)
	a.views.RequestHeaderBar.SetRequestMethod(focusedRequest.Method)

	// Once the user changes the selection, load the historical response from memory and set it
	// TODO: Clean setting the empty response;
	go func() {
		a.service.CancelSentHttpRequest() // Cancel any requests to prevent any side-effects.
		// But at some point, we would not want to cancel the entire request but rather the side effects...Is there a better way to do this?
		a.views.ResponseWindow.SetHttpResponse(entity.HttpResponse{})
	}()
}

func (a *AppController) OnRequestListNameSelected(selected entity.RequestBasic) {
	s := fmt.Sprintf("%s %s", selected.Method, selected.Url)
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
	d := a.service.GetBasicFocusData()

	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to add a request at user position", position)

	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.StatusBar.SetText("Added new request")
}

func (a *AppController) OnRequestListDuplicate(entity.RequestBasic) {
}

func (a *AppController) OnRequestListRemove(request entity.RequestBasic, position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to delete request", request, "at user position", position)
	a.service.RemoveRequest(request.Id, position)
	d := a.service.GetBasicFocusData()
	a.views.RequestsList.RenderRequests(d.RequestsBasic)

	s := fmt.Sprintf("Removed request %s", request.Name)
	a.views.StatusBar.SetText(s)
}

// Landing View (Collection list)

func (a *AppController) OnFocusedCollectionChanged(changedCollection entity.Collection) {
	d := a.service.ChangeFocusedCollection(changedCollection.Id)

	// when user selects a collection, a request item would be automatically changed
	a.views.RequestsList.RenderRequests(d.RequestsBasic)
	a.views.RequestsList.SelectRequest(d.SelectedRequest)
}

func (a *AppController) OnCollectionAdd(position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to add a collection at user position", position)
	a.service.CreateCollection(position)
	d := a.service.GetBasicFocusData()
	a.views.CollectionsList.RenderCollections(d.Collections)
}

func (a *AppController) OnCollectionRemove(collection entity.Collection, position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to remove a collection", collection, "at position", position)
	a.service.DeleteCollection(collection.Id, position)
	d := a.service.GetBasicFocusData()
	a.views.CollectionsList.RenderCollections(d.Collections)
}
