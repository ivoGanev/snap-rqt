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
	cancelFunc         context.CancelFunc
}

func NewAppService() *AppService {
	collectionsRepository := memmock.NewCollectionRepository()
	requestsRepository := memmock.NewRequestsRepository(collectionsRepository)
	stateRepository := memmock.NewStateService(collectionsRepository, requestsRepository)

	stateService := NewStateService(stateRepository)
	collectionService := NewCollectionService(collectionsRepository)
	requestsService := NewRequestsService(requestsRepository)

	appService := &AppService{
		stateService:       stateService,
		collectionsService: collectionService,
		requestsService:    requestsService,
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

func (a *AppService) UpdateFocusedRequest(modRequest entity.ModRequest) {
	rId := a.stateService.GetFocusedRequestId()
	request := a.requestsService.GetRequest(rId)
	request.Mod(modRequest)
	a.requestsService.UpdateRequest(request)
}

func (a *AppService) SendHttpRequest(id string, onHttpResponse func(entity.HttpResult)) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	a.cancelFunc = cancel

	go func() {
		req := a.requestsService.GetRequest(id)
		httpResult := http.SendRequest(ctx, req.AsHttpRequest())

		if ctx.Err() == context.Canceled {
			return // User manually canceled – skip onHttpResponse
		}

		onHttpResponse(httpResult)
	}()
}

func (a *AppService) CancelSentHttpRequest() {
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
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
