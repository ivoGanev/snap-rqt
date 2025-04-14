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
	c.App.Views.RequestsView.ChangeMethodTypeOnSelectedRow(method)
	c.App.hidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
	c.App.focus(c.App.Views.RequestsView)
}
