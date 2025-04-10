package data

type Store interface {
	StoreCollection(collection Collection) error
	StoreRequest(request Request) error
	LoadAllCollections() ([]Collection, error)
	LoadAllRequests() ([]Request, error)
	GetCollectionsForRequest(requestID int) ([]Collection, error)
}
