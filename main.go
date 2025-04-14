package main

import (
	"snap-rq/app"
	"snap-rq/app/database/memmock"
)

func main() {
	requestsService := memmock.NewRequestsService()

	services := app.Services{
		RequestsService: requestsService,
	}

	app := app.NewApp(&services)
	app.Init()
}
