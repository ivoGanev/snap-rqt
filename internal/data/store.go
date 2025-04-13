package data

type Store interface {
	StoreCollection(collection *Collection) error
	StoreRequest(request *Request) error
	CheckExistsRequest(request *Request) bool
	CheckExistsCollection(collection *Collection) bool
	GetCollections() ([]Collection, error)
	GetCollectionsSimple() ([]CollectionSimple, error)
	GetCollectionSimple(collectionId string) (CollectionSimple, error)
	LoadAllRequests() ([]Request, error)
	LoadSessionData() (*UserSession, error)
	GetCollection(id string) (Collection, error)
	GetRequest(id string) (Request, error)
	UpdateRequest(request Request) (Request, error)
}
