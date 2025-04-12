package internal

type RequestsController struct {
	*EditorController
}

func (r *RequestsView) OnRequestsModelChanged(requests *[]data.Node[http.Request], operation model.CrudOp) {
	if operation == model.UPDATE && len(*requests) > 1 {

	}
}
