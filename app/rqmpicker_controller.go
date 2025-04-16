package app

type RequestMethodPickerController interface {
	HandleRequestMethodSelected(string)
}

type RequestMethodPickerViewController struct {
	*App
}

func NewMethodPickerModalController(app *App) *RequestMethodPickerViewController {
	return &RequestMethodPickerViewController{app}
}

func (c *RequestMethodPickerViewController) HandleRequestMethodSelected(method string) {
	state := c.App.Services.StateService.GetState()
	rstate := state.GetRequestViewState(state.AppViewState.SelectedCollectionId)

	c.App.Views.RequestsView.ChangeMethodTypeOnSelectedRow(rstate.RowIndex, method)
	c.App.hidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
	c.App.focus(c.App.Views.RequestsView)
}
