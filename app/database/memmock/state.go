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

	requests, err := r.GetRequestListItems(collections[0].Id)
	if err != nil {
		panic(err)
	}

	state := app.State{
		RequestsViewState: map[string]app.RequestsViewState{
			collections[0].Id: {
				RowIndex:    0,
				ColumnIndex: 0,
				RequestId:   requests[0].Id,
			},
		},
		AppViewState: app.AppViewState{
			FocusedView:           app.VIEW_NAME_REQUESTS,
			SelectedCollectionRow: 0,
			SelectedCollectionId:  collections[0].Id,
		},
	}
	return &MemMockStateService{state}
}

func (m *MemMockStateService) GetState() app.State {
	return m.state
}

func (m *MemMockStateService) UpdateState(stateAction app.StateAction) {
	state := m.GetState()
	stateAction.Apply(&state)
	m.state = state
}
