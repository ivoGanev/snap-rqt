package data

import (
	"snap-rq/internal/http"
)

type Collection struct {
	Node[[]Request]
}

type Request struct {
	Node[http.Request]
}

// When deleting or adding a request or a collection, we will also need to update this entity.
// Updating a request (such as copying it in another collection) means we also need to add to this entity
type RequestCollections struct {
	RequestId string
	CollectionId string
}