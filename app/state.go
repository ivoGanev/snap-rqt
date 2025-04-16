package app

type State struct {
	RequestsViewState map[string]RequestsViewState
	AppViewState      AppViewState
}

type RequestsViewState struct {
	RowIndex    int
	ColumnIndex int
	RequestId   string
}

type AppViewState struct {
	FocusedView           string
	SelectedCollectionRow int
	SelectedCollectionId  string
}

type StateAction interface {
	Apply(state *State)
}

type SetRowIndexAction struct {
	CollectionID string
	RowIndex     int
}

func (a SetRowIndexAction) Apply(state *State) {
	viewState := state.RequestsViewState[a.CollectionID]
	viewState.RowIndex = a.RowIndex
	state.RequestsViewState[a.CollectionID] = viewState
}

type SetRequestViewStateAction struct {
	CollectionID string
	ViewState    RequestsViewState
}

func (a SetRequestViewStateAction) Apply(state *State) {
	state.RequestsViewState[a.CollectionID] = a.ViewState
}

type SetActiveRequestViewStateAction struct {
	ViewState RequestsViewState
}

func (a SetActiveRequestViewStateAction) Apply(state *State) {
	id := state.AppViewState.SelectedCollectionId
	state.RequestsViewState[id] = a.ViewState
}

func (s *State) Apply(action StateAction) {
	action.Apply(s)
}

func (s *State) GetRequestViewState(collectionID string) RequestsViewState {
	return s.RequestsViewState[collectionID]
}

type StateService interface {
	GetState() State
	UpdateState(StateAction)
}
