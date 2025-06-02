package service

import (
	"context"
	"snap-rq/app/entity"
	"snap-rq/app/http"
	"snap-rq/app/repository/memmock"
	"time"
)

type AppService struct {
	stateService       *StateService
	collectionsService *CollectionsService
	requestsService    *RequestsService
}



func NewAppService() AppService {
	collectionsRepository := memmock.NewCollectionRepository()
	requestsRepository := memmock.NewRequestsRepository(collectionsRepository)
	stateRepository := memmock.NewStateService(collectionsRepository, requestsRepository)

	stateService := NewStateService(stateRepository)
	collectionService := NewCollectionService(collectionsRepository)
	requestsService := NewRequestsService(requestsRepository)

	appService := AppService{
		stateService,
		collectionService,
		requestsService,
	}
	return appService
}

func (a *AppService) FetchLandingData() entity.BasicFocusData {
	collections := a.collectionsService.GetCollections()
	cId := collections[0].Id
	requests := a.requestsService.GetRequestsBasic(cId)

	return entity.BasicFocusData{
		Collections:          collections,
		RequestsBasic:        requests,
		SelectedCollectionId: a.stateService.GetFocusedCollectionId(),
		SelectedRequestId:    a.stateService.GetFocusedRequestId(),
	}
}

func (a *AppService) UpdateFocusedRequest(patchRequest entity.PatchRequest) {
	rId := a.stateService.GetFocusedRequestId()
	request := a.requestsService.GetRequest(rId)
	request.ApplyPatch(patchRequest)
	a.requestsService.UpdateRequest(request)
}

func (a *AppService) SendHttpRequestById(id string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := a.requestsService.GetRequest(id)
	return http.SendRequest(ctx, req.AsHttpRequest())
}

func (a *AppService) ChangeFocusedCollection(focusedCollectionId string) entity.BasicFocusData {
	a.stateService.SetFocusedCollection(focusedCollectionId)
	return a.FetchBasicFocusData()
}

func (a *AppService) ChangeFocusedRequest(selectedRequest entity.RequestBasic) {
	cId := a.stateService.GetFocusedCollectionId()
	a.stateService.SetFocusedRequest(cId, selectedRequest.Id)
}

func (a *AppService) AddRequest(position int) {
	cId := a.stateService.GetFocusedCollectionId()
	a.requestsService.CreateRequest(cId, position)
}

func (a *AppService) RemoveRequest(requestId string, position int) {
	cId := a.stateService.GetFocusedCollectionId()
	a.requestsService.DeleteRequest(cId, requestId, position)
}

func (a *AppService) FetchBasicFocusData() entity.BasicFocusData {
	collections := a.collectionsService.GetCollections()
	cId := a.stateService.GetFocusedCollectionId()

	requests := a.requestsService.GetRequestsBasic(cId)

	return entity.BasicFocusData{
		Collections:          collections,
		RequestsBasic:        requests,
		SelectedCollectionId: a.stateService.GetFocusedCollectionId(),
		SelectedRequestId:    a.stateService.GetFocusedRequestId(),
	}
}
