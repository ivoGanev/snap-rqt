package memmock

import "snap-rq/app/entity"

type MemMockStateRepository struct {
	state entity.State
}

func NewStateService(c MemMockCollectionRepository, r MemMockRequestsRepository) MemMockStateRepository {
	collections, err := c.GetCollections()
	if err != nil {
		panic(err)
	}

	selectedRequests := make(map[string]entity.FocusedRequest)
	for _, collection := range collections {
		req, _ := r.GetRequestsBasic(collection.Id)
		selectedRequests[collection.Id] = entity.FocusedRequest{Id: req[0].Id, Row: 0, Column: 1}
	}

	state := entity.State{
		FocusedRequests: selectedRequests,
		FocusedCollection: entity.FocusedCollection{
			Row: 0,
			Id:  collections[0].Id,
		},
		FocusedView: "requests", // TODO: Make it a constant
	}
	return MemMockStateRepository{state}
}

func (m *MemMockStateRepository) GetState() entity.State {
	return m.state
}

func (m *MemMockStateRepository) SetState(state entity.State) {
	m.state = state
}
