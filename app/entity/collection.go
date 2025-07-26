package entity

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
	RowPosition int        `json:"row_position"` // User's logical position of the collection
}

func NewCollection(name string, description string, rowPosition int) Collection {
	return Collection{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		RowPosition: rowPosition,
	}
}
