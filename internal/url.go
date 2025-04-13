package internal

type UrlInputController struct {
	*App
}

func NewUrlInputController(c *App) *UrlInputController {
	return &UrlInputController{c}
}

func (c *UrlInputController) OnUrlTextChanged(text string) {
	selectedRequest := c.App.Views.RequestsView.GetSelectedRequest()
	c.App.Models.ProjectModel.UpdateRequestUrl(selectedRequest.Id, text)
}
