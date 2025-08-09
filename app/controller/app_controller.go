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

func (c *AppController) Start() {
	logger.Debug("Controller starting")

	// Load where the user last left
	d := c.service.GetBasicFocusData()

	// Selecting a collection will automatically render the requests
	c.views.CollectionsList.RenderCollections(d.Collections)
	c.views.CollectionsList.SelectCollection(d.SelectedCollection)

	logger.Debug("Element focus: ", c.app.Application.GetFocus())
}

// App View

func (c *AppController) OnViewModeChange(mode string) {
	request := c.service.GetFocusedRequest()
	c.views.EditorView.SetTextArea(request)
}

// Edit name modal

func (c *AppController) OnEditorModalCancel() {
	c.app.HidePage(view.PAGE_EDIT_NAME)
	c.app.Focus(c.app.Views.CollectionsList)
	d := c.service.GetBasicFocusData()
	c.views.CollectionsList.SelectCollection(d.SelectedCollection)
}

func (c *AppController) OnEditorModalSave(text string, component int) {
	c.app.HidePage(view.PAGE_EDIT_NAME)

	switch component {
	case view.EDITOR_MODAL_COMPONENT_REQUESTS:
		c.service.UpdateFocusedRequest(entity.UpdateRequest{Name: &text})
		d := c.service.GetBasicFocusData()
		c.views.RequestsList.RenderRequests(d.RequestsBasic)
		c.app.Focus(c.app.Views.RequestsList)
		c.views.RequestsList.SelectRequest(d.SelectedRequest)
	case view.EDITOR_MODAL_COMPONENT_COLLETIONS:
		c.service.UpdateFocusedCollection(entity.UpdateCollection{Name: &text})
		d := c.service.GetBasicFocusData()
		c.views.CollectionsList.RenderCollections(d.Collections)
		c.app.Focus(c.app.Views.CollectionsList)
		c.views.CollectionsList.SelectCollection(d.SelectedCollection)
	}

}

func (c *AppController) OnRequestListEditName(request entity.RequestBasic) {
	c.app.ShowPage(view.PAGE_EDIT_NAME)
	c.views.NameEditorModal.Edit(view.EDITOR_MODAL_COMPONENT_REQUESTS)
	c.app.Focus(c.views.NameEditorModal.Input)
}

func (c *AppController) OnCollectionEditName(entity.Collection) {
	c.app.ShowPage(view.PAGE_EDIT_NAME)
	c.views.NameEditorModal.Edit(view.EDITOR_MODAL_COMPONENT_COLLETIONS)
	c.app.Focus(c.views.NameEditorModal.Input)
}

// Editor View

func (c *AppController) OnEditorEditTextArea(editorMode int, edit string) {
	// change the body|header of current HTTP method selected
	switch editorMode {
	case view.EDITOR_VIEW_MODE_BODY:
		c.service.UpdateFocusedRequest(entity.UpdateRequest{Body: &edit})
	case view.EDITOR_VIEW_MODE_HEADERS:
		c.service.UpdateFocusedRequest(entity.UpdateRequest{Headers: &edit})
	}
}

func (c *AppController) OnEditorModeChanged(mode int) {
	request := c.service.GetFocusedRequest()
	c.views.EditorView.SetTextArea(request)
	if mode == view.EDITOR_VIEW_MODE_HEADERS {
		c.app.Focus(c.views.EditorView.HeadersButton)
	} else {
		c.app.Focus(c.views.EditorView.BodyButton)
	}
}

func (c *AppController) OnEditorTextAreaSelected() {
	c.app.Focus(c.app.Views.EditorView.TextArea)
}

// Url Input
func (c *AppController) OnUrlInputTextChanged(urlText string) {
	c.service.UpdateFocusedRequest(entity.UpdateRequest{Url: &urlText})
}

// Landing View (Request Header Bar)

func (c *AppController) OnMethodSelection(method string) {
	c.service.UpdateFocusedRequest(entity.UpdateRequest{Method: &method})

	d := c.service.GetBasicFocusData()
	c.views.RequestsList.RenderRequests(d.RequestsBasic)
}

// Landing View (Request List)

func (c *AppController) OnRequestListRequestFocusChanged(selectedRequest entity.RequestBasic) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User selected request", selectedRequest)
	c.service.ChangeFocusedRequest(selectedRequest)
	focusedRequest := c.service.GetFocusedRequest()

	c.views.RequestHeaderBar.SetUrlText(focusedRequest.Url)
	c.views.RequestHeaderBar.SetRequestMethod(focusedRequest.Method, true)

	// Once the user changes the selection, load the historical response from memory and set it
	// TODO: Clean setting the empty response;
	go func() {
		c.service.CancelSentHttpRequest() // Cancel any requests to prevent any side-effects.
		// But at some point, we would not want to cancel the entire request but rather the side effects...Is there a better way to do this?
		c.views.ResponseWindow.SetHttpResponse(entity.HttpResponse{})
	}()
}

func (c *AppController) OnRequestListNameSelected(selected entity.RequestBasic) {
	s := fmt.Sprintf("%s %s", selected.Method, selected.Url)
	c.views.StatusBar.SetText(s)

	onHttpResult := func(result entity.HttpResult) {
		if result.Error != nil {
			c.views.ResponseWindow.SetError(result.Error)
		} else {
			c.views.ResponseWindow.SetHttpResponse(result.Response)
		}
	}
	c.service.SendHttpRequest(selected.Id, onHttpResult)
	c.views.ResponseWindow.AwaitResponse()
}

func (c *AppController) OnRequestListAdd(position int) {
	c.service.AddRequest(position)
	d := c.service.GetBasicFocusData()

	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to add a request at user position", position)

	c.views.RequestsList.RenderRequests(d.RequestsBasic)
	c.views.StatusBar.SetText("Added new request")
}

func (c *AppController) OnRequestListDuplicate(entity.RequestBasic) {
}

func (c *AppController) OnRequestListRemove(request entity.RequestBasic, position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to delete request", request, "at user position", position)
	c.service.DeleteRequest(request.Id, position)
	d := c.service.GetBasicFocusData()
	c.views.RequestsList.RenderRequests(d.RequestsBasic)

	s := fmt.Sprintf("Removed request %s", request.Name)
	c.views.StatusBar.SetText(s)
}

// Landing View (Collection list)

func (c *AppController) OnFocusedCollectionChanged(changedCollection entity.Collection) {
	d := c.service.ChangeFocusedCollection(changedCollection.Id)

	// when user selects a collection, a request item would be automatically changed
	c.views.RequestsList.RenderRequests(d.RequestsBasic)
	c.views.RequestsList.SelectRequest(d.SelectedRequest)
}

func (c *AppController) OnCollectionAdd(position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to add a collection at user position", position)
	c.service.AddCollection(position)
	d := c.service.GetBasicFocusData()
	c.views.CollectionsList.RenderCollections(d.Collections)
}

func (c *AppController) OnCollectionRemove(collection entity.Collection, position int) {
	logger.Info(APP_CONTROLLER_LOG_TAG, "User requested to remove a collection", collection, "at position", position)
	c.service.DeleteCollection(collection.Id, position)
	d := c.service.GetBasicFocusData()
	c.views.CollectionsList.RenderCollections(d.Collections)
	c.views.RequestsList.RenderRequests(d.RequestsBasic)
}
