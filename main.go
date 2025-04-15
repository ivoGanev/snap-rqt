package main

import (
	"snap-rq/app"
	memmock "snap-rq/app/database/memmock"
)

func main() {
	collectionService := memmock.NewCollectionService()
	requestsService := memmock.NewRequestsService(*collectionService)

	services := app.Services{
		RequestsService:   requestsService,
		CollectionService: collectionService,
	}

	app := app.NewApp(&services)
	app.Init()
}
