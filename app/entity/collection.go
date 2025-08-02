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

type UpdateCollection struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	RowPosition *int    `json:"row_position,omitempty"`
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

func (c *Collection) Update(update UpdateCollection) {
	now := time.Now()

	if update.Name != nil {
		c.Name = *update.Name
	}
	if update.Description != nil {
		c.Description = *update.Description
	}
	if update.RowPosition != nil {
		c.RowPosition = *update.RowPosition
	}

	c.ModifiedAt = &now
}
