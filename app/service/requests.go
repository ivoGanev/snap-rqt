package service

import (
	"snap-rq/app/constants"
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
)

const REQUEST_SERVICE_LOG_TAG = "[Request Service]"

type RequestsService struct {
	repository repository.RequestsRepository
}

func (m RequestsService) UpdateRequest(request entity.Request) entity.Request {
	request, err := m.repository.UpdateRequest(request)
	if err != nil {
		logger.Println(REQUEST_SERVICE_LOG_TAG, err)
	}
	return request
}

func (m RequestsService) GetRequest(rId string) entity.Request  {
	storedRequest, err := m.repository.GetRequest(rId)
	if err != nil {
		logger.Println(REQUEST_SERVICE_LOG_TAG, err)
	}
	return storedRequest
}

func NewRequestsService(requests repository.RequestsRepository) *RequestsService {
	return &RequestsService{requests}
}

func (m *RequestsService) GetRequestsBasic(collectionId string) []entity.RequestBasic {
	storedRequests, err := m.repository.GetRequestsBasic(collectionId)
	if err != nil {
		logger.Println(REQUEST_SERVICE_LOG_TAG, err)
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
