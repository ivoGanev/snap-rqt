package repository

import "snap-rq/app/entity"

type StateRepository interface {
	GetState() entity.State
	SetState(state entity.State)
}

type RequestsRepository interface {
	GetRequests() ([]entity.Request, error)
	GetRequest(id string) (entity.Request, error)
	CreateRequest(r entity.Request) error
	DeleteRequest(id string) (entity.Request, error)
	UpdateRequest(updated entity.Request) (entity.Request, error)
	SaveRequest(r *entity.Request) error
	GetRequestsBasic(collectionId string) ([]entity.RequestBasic, error)
}

type CollectionsRepository interface {
	GetCollections() ([]entity.Collection, error)
	CreateCollection(c *entity.Collection) error
	DeleteCollection(id string) error
	GetCollection(id string) (entity.Collection, error)
	UpdateCollection(updated entity.Collection) (entity.Collection, error)
}
