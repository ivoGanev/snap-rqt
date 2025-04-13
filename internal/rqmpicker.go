package internal

import (
	"snap-rq/internal/http"
)

type RequestMethodPickerController struct {
	*App
}

func NewMethodPickerModalController(app *App) *RequestMethodPickerController {
	return &RequestMethodPickerController{app}
}

func (c *RequestMethodPickerController) OnRequestMethodSelected(method http.RequestMethod) {
	c.App.Views.RequestsView.ChangeMethodTypeOnSelectedRow(method)
	c.App.hidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
	c.App.focus(c.App.Views.RequestsView)
}
