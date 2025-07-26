package service

import (
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
)

type StateService struct {
	state repository.StateRepository
}

const STATE_SERVICE_LOG_TAG = "[State Service]"


func NewStateService(stateRepository repository.StateRepository) *StateService {
	return &StateService{stateRepository}
}

func (s *StateService) SetFocusedCollection(collectionId string) {
    state, err := s.state.GetState()
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to get state:", err)
        return
    }
    state.FocusedCollectionId = collectionId
    err = s.state.SetState(state)
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to set focused collection:", err)
    }
}

func (s *StateService) SetFocusedRequest(collectionID string, requestId string) {
    state, err := s.state.GetState()
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to get state:", err)
        return
    }
    state.FocusedRequestIds[collectionID] = requestId
    logger.Println(STATE_SERVICE_LOG_TAG, "Setting focused request to", requestId)
    err = s.state.SetState(state)
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to set focused request:", err)
    }
}

func (s *StateService) GetFocusedRequestByCollection(collectionID string) string {
    state, err := s.state.GetState()
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to get state:", err)
        return ""
    }
    if state.FocusedRequestIds == nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "FocusedRequestIds is nil for collection:", collectionID)
        return ""
    }
    return state.FocusedRequestIds[collectionID]
}

func (s *StateService) GetFocusedCollectionId() string {
    state, err := s.state.GetState()
    if err != nil {
        logger.Println(STATE_SERVICE_LOG_TAG, "Failed to get state:", err)
        return ""
    }
    if state.FocusedCollectionId == "" {
        logger.Println(STATE_SERVICE_LOG_TAG, "FocusedCollectionId is empty")
    }
    return state.FocusedCollectionId
}

func (s *StateService) GetFocusedRequestId() string {
    collectionID := s.GetFocusedCollectionId()
    if collectionID == "" {
        logger.Println(STATE_SERVICE_LOG_TAG, "No focused request found for collection:", collectionID)
        return ""
    }
    requestID := s.GetFocusedRequestByCollection(collectionID)
    if requestID == "" {
        logger.Println(STATE_SERVICE_LOG_TAG, "No focused request found for collection:", collectionID)
    }
    return requestID
}
