package random

import (
	"math/rand"
	"snap-rq/app/entity"
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

func Collection(collectionCount int) []entity.Collection {
	var collections []entity.Collection

	for range collectionCount {
		collectionName := collectionNames[rand.Intn(len(collectionNames))]
		collectionDescription := "A collection of API requests for " + collectionName
		collection := entity.NewCollection(collectionName, collectionDescription, 0)
		collections = append(collections, collection)
	}

	return collections
}
