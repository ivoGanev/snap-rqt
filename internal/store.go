package internal

type Store interface {
	StoreCollection(collection *Collection) error
	StoreRequest(request *Request) error
	CheckExistsRequest(request *Request) bool
	CheckExistsCollection(collection *Collection) bool
	LoadAllCollections() ([]Collection, error)
	LoadAllRequests() ([]Request, error)
}
