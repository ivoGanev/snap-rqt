package memmock
// var collectionNames = []string{
// 	"User APIs",
// 	"Order Services",
// 	"Product Endpoints",
// 	"Authentication Flows",
// 	"Billing & Payments",
// 	"Admin Tools",
// 	"Notification Services",
// 	"Analytics APIs",
// 	"Third-Party Integrations",
// 	"Debugging & Testing APIs",
// }

// func GenerateCollectionMocks(collectionCount int, requestGenerator func() *[]Request) *[]Collection {
// 	var collections []Collection

// 	for range collectionCount {
// 		collectionName := collectionNames[rand.Intn(len(collectionNames))]
// 		collectionDescription := "A collection of API requests for " + collectionName
// 		collectionNode := NewNode(collectionName, collectionDescription, requestGenerator())
// 		collection := Collection{ Node: collectionNode }
// 		collections = append(collections, collection)
// 	}

// 	return &collections
// }
