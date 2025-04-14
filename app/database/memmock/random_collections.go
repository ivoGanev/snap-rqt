package memmock

import (
	"math/rand"
	"snap-rq/app"
)

var collectionNames = []string{
	"User APIs",
	"Order Services",
	"Product Endpoints",
	"Authentication Flows",
	"Billing & Payments",
	"Admin Tools",
	"Notification Services",
	"Analytics APIs",
	"Third-Party Integrations",
	"Debugging & Testing APIs",
}

func GenerateCollectionMocks(collectionCount int) *[]app.Collection {
	var collections []app.Collection

	for range collectionCount {
		collectionName := collectionNames[rand.Intn(len(collectionNames))]
		collectionDescription := "A collection of API requests for " + collectionName
		collection := app.NewCollection(collectionName, collectionDescription)
		collections = append(collections, collection)
	}

	return &collections
}
