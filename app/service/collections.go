package service

import (
	"snap-rq/app/entity"
	"snap-rq/app/repository"
)

type CollectionsService struct {
	collections repository.CollectionsRepository
}

func NewCollectionService(collections repository.CollectionsRepository) *CollectionsService {
	return &CollectionsService{collections}
}

func (c *CollectionsService) GetCollections() []entity.Collection {
	collections, err := c.collections.GetCollections()
	if err != nil {
		panic(err)
	}
	return collections
}