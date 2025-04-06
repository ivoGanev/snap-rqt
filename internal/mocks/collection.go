package mocks

import (
	"math/rand"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
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

func GenerateCollectionMocks(collectionCount, requestsPerCollection int) *[]data.Node[map[string]data.Node[http.Request]] {
	var collections []data.Node[map[string]data.Node[http.Request]]

	for range collectionCount {
		collectionName := collectionNames[rand.Intn(len(collectionNames))]
		collectionDescription := "A collection of API requests for " + collectionName

		requestNodes := GenerateMockRequests(requestsPerCollection)
		requestMap := make(map[string]data.Node[http.Request])

		for _, reqNode := range *requestNodes {
			requestMap[reqNode.Id] = reqNode
		}

		collection := data.NewNode(collectionName, collectionDescription, &requestMap)
		collections = append(collections, collection)
	}

	return &collections
}
