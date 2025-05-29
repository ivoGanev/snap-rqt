package memmock

import (
	"errors"
	"slices"
	"snap-rq/app/entity"
	"snap-rq/app/repository/memmock/random"
	"time"
)

type MemMockCollectionRepository struct {
	collections []entity.Collection
}

func NewCollectionRepository() *MemMockCollectionRepository {
	return &MemMockCollectionRepository{
		collections: random.Collection(20),
	}
}

func (s *MemMockCollectionRepository) GetCollections() ([]entity.Collection, error) {
	return s.collections, nil
}

func (s *MemMockCollectionRepository) CreateCollection(c *entity.Collection) error {
	if c.Id == "" {
		newCol := entity.NewCollection(c.Name, c.Description)
		*c = newCol
	}
	c.CreatedAt = time.Now()
	s.collections = append(s.collections, *c)
	return nil
}

func (s *MemMockCollectionRepository) DeleteCollection(id string) error {
	for i, col := range s.collections {
		if col.Id == id {
			s.collections = slices.Delete(s.collections, i, i+1)
			return nil
		}
	}
	return errors.New("collection not found")
}

func (s *MemMockCollectionRepository) GetCollection(id string) (entity.Collection, error) {
	for _, col := range s.collections {
		if col.Id == id {
			return col, nil
		}
	}
	return entity.Collection{}, errors.New("collection not found")
}

func (s *MemMockCollectionRepository) UpdateCollection(updated entity.Collection) (entity.Collection, error) {
	for i, col := range s.collections {
		if col.Id == updated.Id {
			now := time.Now()
			updated.ModifiedAt = &now
			s.collections[i] = updated
			return updated, nil
		}
	}
	return entity.Collection{}, errors.New("collection not found")
}
