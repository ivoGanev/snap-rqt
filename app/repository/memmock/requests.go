package memmock

import (
	"errors"
	"snap-rq/app/entity"
	"snap-rq/app/repository"
	"snap-rq/app/repository/memmock/random"
	"sort"
)

type MemMockRequestsRepository struct {
	StoredRequests []entity.Request
}

func (m *MemMockRequestsRepository) ShiftRequests(collectionId string, startingPosition int, direction string) {
	switch direction {
	case repository.SHIFT_UP:
		// Sort in descending order to prevent overwriting positions
		sort.SliceStable(m.StoredRequests, func(i, j int) bool {
			return m.StoredRequests[i].RowPosition > m.StoredRequests[j].RowPosition
		})
		for i := range m.StoredRequests {
			if m.StoredRequests[i].CollectionID == collectionId &&
				m.StoredRequests[i].RowPosition >= startingPosition {
				m.StoredRequests[i].RowPosition += 1
			}
		}

	case repository.SHIFT_DOWN:
		// Sort in ascending order to prevent overwriting positions
		sort.SliceStable(m.StoredRequests, func(i, j int) bool {
			return m.StoredRequests[i].RowPosition < m.StoredRequests[j].RowPosition
		})
		for i := range m.StoredRequests {
			if m.StoredRequests[i].CollectionID == collectionId &&
				m.StoredRequests[i].RowPosition >= startingPosition {
				m.StoredRequests[i].RowPosition -= 1
			}
		}

	default:
		// Optional: handle unsupported direction
	}
}


// NewRequestsRepository initializes the service with mock data
func NewRequestsRepository(c *MemMockCollectionRepository) *MemMockRequestsRepository {
	m := &MemMockRequestsRepository{}

	collections, err := c.GetCollections()
	if err != nil {
		panic(err)
	}

	var requests []entity.Request
	for _, collection := range collections {
		requestBatch := random.Requests(100, collection.Id)
		requests = append(requests, requestBatch...)
	}
	m.StoredRequests = requests
	return m
}

func (m *MemMockRequestsRepository) DeleteRequest(id string) (entity.Request, error) {
	panic("Not implemented")
}

func (m *MemMockRequestsRepository) CreateRequest(request entity.Request) error {
	m.StoredRequests = append(m.StoredRequests, request)
	return nil
}

func (m *MemMockRequestsRepository) GetRequestsBasic(collectionId string) ([]entity.RequestBasic, error) {
	// load request data only needed for basic landing page to save CPU and memory
	var items []entity.RequestBasic
	for _, r := range m.StoredRequests {
		if r.CollectionID == collectionId {
			items = append(items, entity.NewRequestBasicFromRequest(r))
		}
	}
	return items, nil
}

func (m *MemMockRequestsRepository) SaveRequest(r *entity.Request) error {
	m.StoredRequests = append(m.StoredRequests, *r)
	return nil
}

func (m *MemMockRequestsRepository) GetRequests() ([]entity.Request, error) {
	return m.StoredRequests, nil
}

func (m *MemMockRequestsRepository) GetRequest(id string) (entity.Request, error) {
	for _, r := range m.StoredRequests {
		if r.Id == id {
			return r, nil
		}
	}
	return entity.Request{}, errors.New("request not found")
}

func (m *MemMockRequestsRepository) UpdateRequest(updated entity.Request) (entity.Request, error) {
	for i, r := range m.StoredRequests {
		if r.Id == updated.Id {
			m.StoredRequests[i] = updated
			return updated, nil
		}
	}
	return entity.Request{}, errors.New("request not found")
}
