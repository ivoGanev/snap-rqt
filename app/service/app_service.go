package service

import (
	"context"
	"errors"
	"snap-rq/app/constants"
	"snap-rq/app/entity"
	"snap-rq/app/http"
	logger "snap-rq/app/log"
	"snap-rq/app/repository"
	"snap-rq/app/repository/sqlite"
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

func (a *AppService) Start() {
	// purge any focus focus data that doesn't match with the actual stored data
	logger.Info(APP_SERVICE_LOG_TAG, "Verifying app view state integrity")
	_, err := a.repoCollections.GetCollections()
	tryHandleGenericError(err)

	appState, err := a.repoState.GetState()
	tryHandleGenericError(err)

	for collectionId := range appState.FocusedRequestIds {
		_, err := a.repoCollections.GetCollection(collectionId)
		if errors.Is(err, sqlite.ErrCollectionNotFound) {
			delete(appState.FocusedRequestIds, collectionId)
			logger.Warning(APP_SERVICE_LOG_TAG, "Verifying app view state integrity status: deleting non-existing focus collection id", collectionId)
			if appState.FocusedCollectionId == collectionId {
				appState.FocusedCollectionId = ""
			}
		}
	}
	a.repoState.SetState(appState)
}

func NewAppService() *AppService {
	db, err := sqlite.NewDb(DB_FILENAME)
	if err != nil {
		logger.Error(APP_SERVICE_LOG_TAG, "Failed to initialise SQLite DB:", err)
		panic(err)
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

func (a *AppService) GetBasicFocusData() entity.BasicFocusData {
	collections, err := a.repoCollections.GetCollections()
	tryHandleGenericError(err)

	// we cannot leave the app without a collection
	if len(collections) == 0 {
		collection := entity.NewCollection("Default Collection", "default collection for requests", 0)
		a.repoCollections.CreateCollection(&collection)

		// refresh the collections data
		collections, err = a.repoCollections.GetCollections()
		tryHandleGenericError(err)
	}

	cId := a.getFocusedCollectionId()

	requests, err := a.repoRequests.GetRequestsBasic(cId)
	tryHandleGenericError(err)

	return entity.BasicFocusData{
		Collections:        collections,
		RequestsBasic:      requests,
		SelectedCollection: a.getFocusedCollection(),
		SelectedRequest:    a.GetFocusedRequest(),
	}
}

func (a *AppService) UpdateFocusedRequest(update entity.UpdateRequest) {
	rId := a.getFocusedRequestId()
	request, err := a.repoRequests.GetRequest(rId)
	tryHandleGenericError(err)

	request.Update(update)
	a.repoRequests.UpdateRequest(request)
}


func (a *AppService) UpdateFocusedCollection(update entity.UpdateCollection) {
	cId := a.getFocusedCollectionId()
	collection, err := a.repoCollections.GetCollection(cId)
	tryHandleGenericError(err)

	collection.Update(update)
	a.repoCollections.UpdateCollection(collection)
}

func (a *AppService) SendHttpRequest(id string, onHttpResponse func(entity.HttpResult)) {
	ctx, cancel := context.WithCancel(context.Background())
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

	return a.GetBasicFocusData()
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

func (a *AppService) DeleteRequest(requestId string, position int) {
	cId := a.getFocusedCollectionId()
	a.repoRequests.ShiftRequests(cId, position, repository.SHIFT_DOWN)
	a.repoRequests.DeleteRequest(requestId)
	logger.Debug(APP_SERVICE_LOG_TAG, "Remove request from repository with id", requestId, "belonging to collection id", requestId)
}

func (a *AppService) GetFocusedRequest() entity.Request {
	rId := a.getFocusedRequestId()
	logger.Info(APP_SERVICE_LOG_TAG, "Fetched focused request id", rId)

	request, err := a.repoRequests.GetRequest(rId)
	tryHandleGenericError(err)
	return request
}

func (a *AppService) AddCollection(position int) {
	position += 1
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

	a.deleteFocusedCollection(cId)
}

// State

func (a *AppService) deleteFocusedCollection(cId string) {
	state, err := a.repoState.GetState()
	tryHandleGenericError(err)
	delete(state.FocusedRequestIds, cId)
	a.repoState.SetState(state)
}

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

func (a *AppService) getFocusedCollection() entity.Collection {
	cId := a.getFocusedCollectionId()
	logger.Info(APP_SERVICE_LOG_TAG, "Fetched focused collection id", cId)

	col, err := a.repoCollections.GetCollection(cId)
	if err != nil {
		logger.Error(APP_SERVICE_LOG_TAG, "Collection Self-Heal: the focused collection was not found")

		// self-heal if we didn't find a collection and fallback to the first one available
		cols, err := a.repoCollections.GetCollections()
		tryHandleGenericError(err)

		col = cols[0]
	}
	return col
}

// Generic error handler
func tryHandleGenericError(err error) {
	if err != nil {
		logger.Error(err)
		// panic(err) // remove this on production build
	}
}
