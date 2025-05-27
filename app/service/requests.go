package service

import (
	"snap-rq/app/entity"
	"snap-rq/app/repository"
)

type RequestsService struct {
	requests repository.RequestsRepository
}

func (m RequestsService) UpdateRequest(request entity.Request) {
	m.requests.UpdateRequest(request)
}

func (m RequestsService) GetRequest(rId string) entity.Request {
	storedRequests, err := m.requests.GetRequest(rId)
	if err != nil {
		panic(err)
	}
	return storedRequests
}

func NewRequestsService(requests repository.RequestsRepository) *RequestsService {
	return &RequestsService{requests}
}

func (m *RequestsService) GetRequestsBasic(collectionId string) []entity.RequestBasic {
	storedRequests, err := m.requests.GetRequestsBasic(collectionId)
	if err != nil {
		panic(err)
	}

	return storedRequests
}
