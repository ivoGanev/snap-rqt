package memmock

import "snap-rq/app"

type MemMockStateService struct {
	state app.State
}

func NewStateService(c app.CollectionService, r app.RequestsService) *MemMockStateService {
	collections, err := c.GetCollections()
	if err != nil {
		panic(err)
	}

	selectedRequests := make(map[string]app.SelectedRequest)
	for _, collection := range  collections {
		req, _ := r.GetRequestListItems(collection.Id)
		selectedRequests[collection.Id] = app.SelectedRequest{ Id: req[0].Id, Row: 0, Column: 1 }
	}

	state := app.State{
		SelectedRequests: selectedRequests,
		SelectedCollection: app.SelectedCollection{
			Row: 0,
			Id:  collections[0].Id,
		},
		FocusedView: app.VIEW_NAME_REQUESTS,
	}
	return &MemMockStateService{state}
}

func (m *MemMockStateService) GetState() app.State {
	return m.state
}

func (m *MemMockStateService) UpdateState(stateAction app.State) {
	state := m.GetState()
	m.state = state
}
