package memmock

import (
	"errors"
	"slices"
	"snap-rq/app"
	"snap-rq/app/database/memmock/random"
	"time"
)

type MemMockCollectionService struct {
	collections []app.Collection
}

func NewCollectionService() *MemMockCollectionService {
	return &MemMockCollectionService{
		collections: random.Collection(20),
	}
}

func (s *MemMockCollectionService) GetCollections() ([]app.Collection, error) {
	return s.collections, nil
}

func (s *MemMockCollectionService) CreateCollection(c *app.Collection) error {
	if c.Id == "" {
		newCol := app.NewCollection(c.Name, c.Description)
		*c = newCol
	}
	c.CreatedAt = time.Now()
	s.collections = append(s.collections, *c)
	return nil
}

func (s *MemMockCollectionService) DeleteCollection(id string) error {
	for i, col := range s.collections {
		if col.Id == id {
			s.collections = slices.Delete(s.collections, i, i+1)
			return nil
		}
	}
	return errors.New("collection not found")
}

func (s *MemMockCollectionService) GetCollection(id string) (app.Collection, error) {
	for _, col := range s.collections {
		if col.Id == id {
			return col, nil
		}
	}
	return app.Collection{}, errors.New("collection not found")
}

func (s *MemMockCollectionService) UpdateCollection(updated app.Collection) (app.Collection, error) {
	for i, col := range s.collections {
		if col.Id == updated.Id {
			now := time.Now()
			updated.ModifiedAt = &now
			s.collections[i] = updated
			return updated, nil
		}
	}
	return app.Collection{}, errors.New("collection not found")
}
