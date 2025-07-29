package service

import (
	"context"
	"snap-rq/app/constants"
	"snap-rq/app/entity"
	"snap-rq/app/http"
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
	"snap-rq/app/repository/sqlite"
	"time"
)

const (
	APP_SERVICE_LOG_TAG = "[App Service]"
	DB_FILENAME         = "requests.db"
)

type AppService struct {
	repoState       repository.StateRepository
	repoCollections repository.CollectionsRepository
	repoRequests    repository.RequestsRepository
	cancelFunc      context.CancelFunc
}

func NewAppService() *AppService {
	db, err := sqlite.NewDb(DB_FILENAME)
	if err != nil {
		logger.Error(APP_SERVICE_LOG_TAG, "Failed to initialise SQLite DB:", err)
	}
	collectionsRepository := sqlite.NewCollectionRepository(db)
	requestsRepository := sqlite.NewRequestsRepository(db)
	stateRepository := sqlite.NewStateRepository(db, collectionsRepository, requestsRepository)

	appService := &AppService{
		repoState:       stateRepository,
		repoCollections: collectionsRepository,
		repoRequests:    requestsRepository,
	}

	return appService
}

func (a *AppService) FetchLandingData() entity.BasicFocusData {
	collections, err := a.repoCollections.GetCollections()
	tryHandleGenericError(err)

	cId := collections[0].Id
	requests, err := a.repoRequests.GetRequestsBasic(cId)
	tryHandleGenericError(err)

	return entity.BasicFocusData{
		Collections:        collections,
		RequestsBasic:      requests,
		SelectedCollection: a.GetFocusedCollection(),
		SelectedRequest:    a.GetFocusedRequest(),
	}
}

func (a *AppService) UpdateFocusedRequest(modRequest entity.ModRequest) {
	rId := a.getFocusedRequestId()
	request, err := a.repoRequests.GetRequest(rId)
	tryHandleGenericError(err)

	request.Mod(modRequest)
	a.repoRequests.UpdateRequest(request)
}

func (a *AppService) SendHttpRequest(id string, onHttpResponse func(entity.HttpResult)) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	a.cancelFunc = cancel

	go func() {
		req, err := a.repoRequests.GetRequest(id)
		tryHandleGenericError(err)

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
	state, err := a.repoState.GetState()
	tryHandleGenericError(err)

	state.FocusedCollectionId = focusedCollectionId
	err = a.repoState.SetState(state)
	tryHandleGenericError(err)

	return a.FetchBasicFocusData()
}

func (a *AppService) ChangeFocusedRequest(selectedRequest entity.RequestBasic) {
	cId := a.getFocusedCollectionId()

	state, err := a.repoState.GetState()
	tryHandleGenericError(err)

	state.FocusedRequestIds[cId] = selectedRequest.Id
	logger.Debug(APP_SERVICE_LOG_TAG, "Setting focused request to", selectedRequest.Id)

	err = a.repoState.SetState(state)
	tryHandleGenericError(err)
}

func (a *AppService) AddRequest(position int) {
	cId := a.getFocusedCollectionId()
	request := entity.NewRequest(cId, "New Request", "", string(constants.GET), "", "", "", position)
	a.repoRequests.ShiftRequests(cId, position, repository.SHIFT_UP)
	a.repoRequests.CreateRequest(request)
	logger.Debug(APP_SERVICE_LOG_TAG, "Add new request to repository", request, "for collection id", cId)
}

func (a *AppService) RemoveRequest(requestId string, position int) {
	cId := a.getFocusedCollectionId()
	a.repoRequests.ShiftRequests(cId, position, repository.SHIFT_DOWN)
	a.repoRequests.DeleteRequest(requestId)
	logger.Debug(APP_SERVICE_LOG_TAG, "Remove request from repository with id", requestId, "belonging to collection id", requestId)
}

func (a *AppService) FetchBasicFocusData() entity.BasicFocusData {
	collections, err := a.repoCollections.GetCollections()
	tryHandleGenericError(err)

	cId := a.getFocusedCollectionId()

	requests, err := a.repoRequests.GetRequestsBasic(cId)
	tryHandleGenericError(err)

	return entity.BasicFocusData{
		Collections:        collections,
		RequestsBasic:      requests,
		SelectedCollection: a.GetFocusedCollection(),
		SelectedRequest:    a.GetFocusedRequest(),
	}
}

func (a *AppService) GetFocusedRequest() entity.Request {
	rId := a.getFocusedRequestId()
	logger.Info(APP_SERVICE_LOG_TAG, "Fetched focused request id", rId)

	request, err := a.repoRequests.GetRequest(rId)
	tryHandleGenericError(err)
	return request
}

func (a *AppService) GetFocusedCollection() entity.Collection {
	cId := a.getFocusedCollectionId()
	logger.Info(APP_SERVICE_LOG_TAG, "Fetched focused collection id", cId)

	col, err := a.repoCollections.GetCollection(cId)
	if err != nil {
		// self-heal if we didn't find a collection and fallback to the first one available
		cols, err := a.repoCollections.GetCollections()
		tryHandleGenericError(err)

		col = cols[0]
		logger.Error(APP_SERVICE_LOG_TAG, "Collection Self-Heal: the focused collection was not found")
	}
	return col
}

func (a *AppService) CreateCollection(position int) {
	collection := entity.NewCollection("New Collection", "", position)

	err := a.repoCollections.ShiftCollections(position, repository.SHIFT_UP)
	tryHandleGenericError(err)

	err = a.repoCollections.CreateCollection(&collection)
	tryHandleGenericError(err)
}

func (a *AppService) DeleteCollection(cId string, position int) {
	err := a.repoCollections.ShiftCollections(position, repository.SHIFT_DOWN)
	tryHandleGenericError(err)

	err = a.repoCollections.DeleteCollection(cId)
	tryHandleGenericError(err)
}

// State Getters

func (a *AppService) getFocusedRequestByCollection(cId string) string {
	state, err := a.repoState.GetState()
	tryHandleGenericError(err)
	return state.FocusedRequestIds[cId]
}

func (a *AppService) getFocusedCollectionId() string {
	state, err := a.repoState.GetState()
	tryHandleGenericError(err)
	return state.FocusedCollectionId
}

func (a *AppService) getFocusedRequestId() string {
	collectionID := a.getFocusedCollectionId()
	requestID := a.getFocusedRequestByCollection(collectionID)
	return requestID
}

// Generic error handler
func tryHandleGenericError(err error) {
	if err != nil {
		logger.Error(err)
		// panic(err) // remove this on production build
	}
}
