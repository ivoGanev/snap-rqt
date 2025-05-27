package memmock

import (
	"errors"
	"snap-rq/app/entity"
	"snap-rq/app/repository/memmock/random"
)

type MemMockRequestsRepository struct {
	StoredRequests []entity.Request
}

// NewRequestsRepository initializes the service with mock data
func NewRequestsRepository(c MemMockCollectionRepository) MemMockRequestsRepository {
	m := MemMockRequestsRepository{}

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

func (m *MemMockRequestsRepository) CreateRequest(id entity.Request) error {
	panic("Not implemented")
}

func (m *MemMockRequestsRepository) GetRequestsBasic(collectionId string) ([]entity.RequestBasic, error) {
	// load request data only needed for basic landing page to save CPU and memory
	var items []entity.RequestBasic
	for _, r := range m.StoredRequests {
		if r.CollectionID == collectionId {
			items = append(items, entity.RequestBasic{
				Id:         r.Id,
				Name:       r.Name,
				Url:        r.Url,
				MethodType: r.MethodType,
			})
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
