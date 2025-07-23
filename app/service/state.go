package service

import (
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
)

type StateService struct {
	state repository.StateRepository
}

func NewStateService(stateRepository repository.StateRepository) *StateService {
	return &StateService{stateRepository}
}

func (s *StateService) SetFocusedCollection(collectionId string) {
	state := s.state.GetState()
	state.FocusedCollectionId = collectionId
	s.state.SetState(state)
}

func (s *StateService) SetFocusedRequest(collectionID string, requestId string) {
	state := s.state.GetState()
	state.FocusedRequestIds[collectionID] = requestId
	logger.Println("[State Service]", "Setting focused request to", requestId)
	s.state.SetState(state)
}

func (s *StateService) GetFocusedRequestByCollection(collectionID string) string {
	return s.state.GetState().FocusedRequestIds[collectionID]
}

func (s *StateService) GetFocusedCollectionId() string {
	return s.state.GetState().FocusedCollectionId
}

func (s *StateService) GetFocusedRequestId() string {
	return s.state.GetState().FocusedRequestIds[s.GetFocusedCollectionId()]
}
