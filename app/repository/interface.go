package repository

import "snap-rq/app/entity"

type StateRepository interface {
	GetState() (entity.AppViewState, error)
	SetState(state entity.AppViewState) error
}

type RequestsRepository interface {
	GetRequest(id string) (entity.Request, error)
	CreateRequest(r entity.Request) error
	DeleteRequest(id string) error
	UpdateRequest(updated entity.Request) (entity.Request, error)
	GetRequestsBasic(collectionId string) ([]entity.RequestBasic, error)
	ShiftRequests(collectionId string, startingPosition int, direction string) error
}

type CollectionsRepository interface {
	GetCollections() ([]entity.Collection, error)
	CreateCollection(c *entity.Collection) error
	DeleteCollection(id string) error
	GetCollection(id string) (entity.Collection, error)
	UpdateCollection(updated entity.Collection) (entity.Collection, error)
	ShiftCollections(position int, direction string) error
}


const (
	SHIFT_UP   = "UP"
	SHIFT_DOWN = "DOWN"
)
