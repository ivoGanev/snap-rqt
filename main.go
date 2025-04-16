package main

import (
	"snap-rq/app"
	memmock "snap-rq/app/database/memmock"
)

func main() {
	collectionService := memmock.NewCollectionService()
	requestsService := memmock.NewRequestsService(*collectionService)
	UIStateService := memmock.NewStateService(collectionService, requestsService)

	services := app.Services{
		RequestsService:   requestsService,
		CollectionService: collectionService,
		StateService:      UIStateService,
	}

	app := app.NewApp(&services)
	app.Init()
}
