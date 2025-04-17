package app

type StateService interface {
	GetState() State
	UpdateState(State)
}

type State struct {
	SelectedRequests map[string]SelectedRequest
	FocusedView      string
	SelectedCollection
}

// ----------> Setters
func (s *State) SetSelectedCollection(selectedCollection SelectedCollection) {
	s.SelectedCollection = selectedCollection
}

// ----------> Getters
func (s *State) GetRequestViewState(collectionID string) SelectedRequest {
	return s.SelectedRequests[collectionID]
}

func (s *State) GetSelectedCollectionId() string {
	return s.SelectedCollection.Id
}

func (s *State) GetSelectedCollectionRow() int {
	return s.SelectedCollection.Row
}

func (s *State) GetSelectedRequestId() string {
	return s.SelectedRequests[s.GetSelectedCollectionId()].Id
}

func (s *State) GetSelectedRequestRow() int {
	return s.SelectedRequests[s.GetSelectedCollectionId()].Row
}
