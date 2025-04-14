package memmock

import (
	"errors"
	"snap-rq/app"
)

// MemMockRequestsService implements the RequestsService interface
type MemMockRequestsService struct {
	StoredRequests []app.Request
}

// NewRequestsService initializes the service with mock data
func NewRequestsService() *MemMockRequestsService {
	return &MemMockRequestsService{
		StoredRequests: GenerateMockRequests(10), // You can customize count here
	}
}

func (m *MemMockRequestsService) DeleteRequest(id string) (app.Request, error) {
	panic("Not implemented")
}

func (m *MemMockRequestsService) CreateRequest(id app.Request) (error) {
	panic("Not implemented")
}

func (m *MemMockRequestsService) GetRequestListItems() ([]app.RequestListItem, error) {
	var items []app.RequestListItem
	for _, r := range m.StoredRequests {
		items = append(items, app.RequestListItem{
			Id:         r.Id,
			Name:       r.Name,
			Url:        r.Url,
			MethodType: r.MethodType,
		})
	}
	return items, nil
}

func (m *MemMockRequestsService) SaveRequest(r *app.Request) error {
	m.StoredRequests = append(m.StoredRequests, *r)
	return nil
}

func (m *MemMockRequestsService) GetRequests() ([]app.Request, error) {
	return m.StoredRequests, nil
}

func (m *MemMockRequestsService) GetRequest(id string) (app.Request, error) {
	for _, r := range m.StoredRequests {
		if r.Id == id {
			return r, nil
		}
	}
	return app.Request{}, errors.New("request not found")
}

func (m *MemMockRequestsService) UpdateRequest(updated app.Request) (app.Request, error) {
	for i, r := range m.StoredRequests {
		if r.Id == updated.Id {
			m.StoredRequests[i] = updated
			return updated, nil
		}
	}
	return app.Request{}, errors.New("request not found")
}
