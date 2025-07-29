package service

import (
	"snap-rq/app/entity"
	"snap-rq/app/repository"
)

type CollectionsService struct {
	repository repository.CollectionsRepository
}

func NewCollectionService(collections repository.CollectionsRepository) *CollectionsService {
	return &CollectionsService{collections}
}

func (c *CollectionsService) GetCollections() []entity.Collection {
	collections, err := c.repository.GetCollections()
	if err != nil {
		panic(err)
	}
	return collections
}

func (c *CollectionsService) GetCollection(cId string) (entity.Collection, error) {
	return c.repository.GetCollection(cId)
}

func (m *CollectionsService) CreateCollection(position int) {
	collection := entity.NewCollection("New Collection", "", position)

	if err := m.repository.ShiftCollections(position, repository.SHIFT_UP); err != nil {
		panic(err)
	}

	if err := m.repository.CreateCollection(&collection); err != nil {
		panic(err)
	}
}

func (m *CollectionsService) DeleteCollection(cId string, position int) {
	if err := m.repository.ShiftCollections(position, repository.SHIFT_DOWN); err != nil {
		panic(err)
	}

	if err := m.repository.DeleteCollection(cId); err != nil {
		panic(err)
	}
}

