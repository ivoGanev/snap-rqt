package app

type CollectionController interface {
	HandleCollectionChanged(newData SelectedCollection)
}

type CollectionViewController struct {
	*App
}

func NewCollectionViewController(app *App) *CollectionViewController {
	return &CollectionViewController{app}
}

func (c *CollectionViewController) HandleCollectionChanged(newData SelectedCollection) {
	state := c.App.StateService.GetState()
	service := c.App.Services.RequestsService

	state.SetSelectedCollection(newData)
	requests, _ := service.GetRequestListItems(state.GetSelectedCollectionId())
	// when user selects a collection, a request item would be automatically changed
	// load the collection
	c.App.Views.RequestsView.RenderRequests(requests)
	c.App.Views.RequestsView.SelectRequest(state.GetSelectedRequestRow())
}
