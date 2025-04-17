package app

type UrlInputController interface {
	HandleUrlTextChanged(string)
}

type UrlInputViewController struct {
	*App
	RequestsService
}

func NewUrlInputController(app *App, requestsService RequestsService) *UrlInputViewController {
	return &UrlInputViewController{app, requestsService}
}

func (c *UrlInputViewController) HandleUrlTextChanged(urlText string) {
	state := c.App.Services.StateService.GetState()

	request, err := c.RequestsService.GetRequest(state.GetSelectedRequestId())
	if err != nil {
		panic(request)
	}
	request.Url = urlText
	c.RequestsService.UpdateRequest(request)
}
