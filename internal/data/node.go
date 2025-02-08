package data

import (
	"time"

	"github.com/google/uuid"
)

type Node[T any] struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
	ModifiedAt  *time.Time
	Data        *T
}

func NewNode[T any](name, description string, data *T) Node[T] {
	now := time.Now()
	return Node[T]{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		ModifiedAt:  nil,
		Data:        data,
	}
}

func (n *Node[T]) UpdateModifiedAt() {
	now := time.Now()
	n.ModifiedAt = &now
}