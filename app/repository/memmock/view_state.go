package memmock

import "snap-rq/app/entity"

type MemMockViewStateRepository struct {
	state entity.AppViewState
}

func NewStateService(c *MemMockCollectionRepository, r *MemMockRequestsRepository) *MemMockViewStateRepository {
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
		FocusedRequestIds: selectedRequests,
		FocusedCollectionId: collections[0].Id,
		FocusedView: "requests", // TODO: Make it a constant
	}
	return &MemMockViewStateRepository{state}
}

func (m *MemMockViewStateRepository) GetState() entity.AppViewState {
	return m.state
}

func (m *MemMockViewStateRepository) SetState(state entity.AppViewState) {
	m.state = state
}
