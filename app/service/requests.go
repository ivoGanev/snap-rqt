package service

import (
	"snap-rq/app/constants"
	"snap-rq/app/entity"
	"snap-rq/app/repository"
)

type RequestsService struct {
	repository repository.RequestsRepository
}

func (m RequestsService) UpdateRequest(request entity.Request) {
	m.repository.UpdateRequest(request)
}

func (m RequestsService) GetRequest(rId string) entity.Request {
	storedRequests, err := m.repository.GetRequest(rId)
	if err != nil {
		panic(err)
	}
	return storedRequests
}

func NewRequestsService(requests repository.RequestsRepository) *RequestsService {
	return &RequestsService{requests}
}

func (m *RequestsService) GetRequestsBasic(collectionId string) []entity.RequestBasic {
	storedRequests, err := m.repository.GetRequestsBasic(collectionId)
	if err != nil {
		panic(err)
	}

	return storedRequests
}

func (m *RequestsService) CreateRequest(collectionId string, position int) {
	request := entity.NewRequest(collectionId, "New Request", "", string(constants.GET), "", "", "", position)
	m.repository.ShiftRequests(collectionId, position, repository.SHIFT_UP)
	m.repository.CreateRequest(request)
}

func (m *RequestsService) DeleteRequest(collectionId, requestId string, position int) {
	m.repository.ShiftRequests(collectionId, position, repository.SHIFT_DOWN)
	m.repository.DeleteRequest(requestId)
}