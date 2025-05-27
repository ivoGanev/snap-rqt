package service

import (
	"snap-rq/app/entity"
	"snap-rq/app/repository"
)

type StateService struct {
	state repository.StateRepository
}

func NewStateService(stateRepository repository.StateRepository) *StateService {
	return &StateService{stateRepository}
}

func (s *StateService) SetFocusedCollection(selectedCollection entity.FocusedCollection) {
	state := s.state.GetState()
	state.FocusedCollection = selectedCollection
	s.state.SetState(state)
}

func (s *StateService) SetFocusedRequest(collectionID string, request entity.FocusedRequest) {
	state := s.state.GetState()
	state.FocusedRequests[collectionID] = request
	s.state.SetState(state)
}

func (s *StateService) GetFocusedRequest(collectionID string) entity.FocusedRequest {
	return s.state.GetState().FocusedRequests[collectionID]
}

func (s *StateService) GetFocusedCollectionId() string {
	return s.state.GetState().FocusedCollection.Id
}

func (s *StateService) GetFocusedCollectionRow() int {
	return s.state.GetState().FocusedCollection.Row
}

func (s *StateService) GetFocusedRequestId() string {
	return s.state.GetState().FocusedRequests[s.GetFocusedCollectionId()].Id
}

func (s *StateService) GetFocusedRequestRow() int {
	return s.state.GetState().FocusedRequests[s.GetFocusedCollectionId()].Row
}
