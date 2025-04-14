package main

import (
	"snap-rq/app"
	memmock "snap-rq/app/database/memmock"
)

func main() {
	requestsService := memmock.NewRequestsService()
	collectionService := memmock.NewCollectionService()

	services := app.Services{
		RequestsService:   requestsService,
		CollectionService: collectionService,
	}

	app := app.NewApp(&services)
	app.Init()
}
