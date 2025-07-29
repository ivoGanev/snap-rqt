package service

import (
	"context"
	"snap-rq/app/entity"
	"snap-rq/app/http"
	logger "snap-rq/app/log"
	"snap-rq/app/repository/sqlite"
	"time"
)

const APP_SERVICE_LOG_TAG = "[App Service]"

type AppService struct {
	stateService       *StateService
	collectionsService *CollectionsService
	requestsService    *RequestsService
	cancelFunc         context.CancelFunc
}

func NewAppService() *AppService {
	// collectionsRepository := memmock.NewCollectionRepository()
	// requestsRepository := memmock.NewRequestsRepository(collectionsRepository)
	// stateRepository := memmock.NewStateRepository(collectionsRepository, requestsRepository)

	db, err := sqlite.NewDb("requests.db")
	if err != nil {
		logger.Println(STATE_SERVICE_LOG_TAG, "Failed to initialise SQLite DB:", err)
	}
	collectionsRepository := sqlite.NewCollectionRepository(db)
	requestsRepository := sqlite.NewRequestsRepository(db)
	stateRepository := sqlite.NewStateRepository(db, collectionsRepository, requestsRepository)

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
		Collections:        collections,
		RequestsBasic:      requests,
		SelectedCollection: a.GetFocusedCollection(),
		SelectedRequest:    a.GetFocusedRequest(),
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
		Collections:        collections,
		RequestsBasic:      requests,
		SelectedCollection: a.GetFocusedCollection(),
		SelectedRequest:    a.GetFocusedRequest(),
	}
}

func (a *AppService) GetFocusedRequest() entity.Request {
	rId := a.stateService.GetFocusedRequestId()
	return a.requestsService.GetRequest(rId)
}

func (a *AppService) GetFocusedCollection() entity.Collection {
	cId := a.stateService.GetFocusedCollectionId()
	col, err := a.collectionsService.GetCollection(cId)
	if err != nil {
		// self-heal if we didn't find a collection and fallback to the first one available
		cols := a.collectionsService.GetCollections()
		col = cols[0]
		logger.Error(APP_SERVICE_LOG_TAG, "Collection Self-Heal: the focused collection was not found")
	}
	return col
}


func (a *AppService) CreateCollection(position int) {
	a.collectionsService.CreateCollection(position)
}

func (a *AppService) DeleteCollection(cId string, position int) {
	a.collectionsService.DeleteCollection(cId, position)
}
