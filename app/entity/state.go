package entity

type State struct {
	FocusedRequests   map[string]FocusedRequest
	FocusedView       string
	FocusedCollection FocusedCollection
}

type FocusedRequest struct {
	Row    int
	Column int
	Id     string
}

type FocusedCollection struct {
	Id  string
	Row int
}
