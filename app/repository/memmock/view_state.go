package memmock

import (
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
)

type MemMockViewStateRepository struct {
	state entity.AppViewState
}

const VIEW_STATE_LOG_TAG = "[MEMMOCK State Repository]"

func NewStateRepository(c *MemMockCollectionRepository, r *MemMockRequestsRepository) *MemMockViewStateRepository {
	collections, err := c.GetCollections()
	if err != nil {
		panic(err)
	}

	selectedRequests := make(map[string]string)
	for _, collection := range collections {
		req, _ := r.GetRequestsBasic(collection.Id)
		selectedRequests[collection.Id] = req[0].Id
	}

	state := entity.AppViewState{
		FocusedRequestIds:   selectedRequests,
		FocusedCollectionId: collections[0].Id,
		FocusedView:         "requests", // TODO: Make it a constant
	}
	return &MemMockViewStateRepository{state}
}

func (m *MemMockViewStateRepository) GetState() entity.AppViewState {
	return m.state
}

func (m *MemMockViewStateRepository) SetState(state entity.AppViewState) {
	m.state = state
	logger.Println(VIEW_STATE_LOG_TAG, "Setting state: ", m.state)
}
