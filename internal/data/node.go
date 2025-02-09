package data

import (
	"fmt"
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

func (n Node[T]) String() string {
	modified := "Never"
	if n.ModifiedAt != nil {
		modified = n.ModifiedAt.Format(time.RFC3339)
	}

	dataStr := "<nil>"
	if n.Data != nil {
		dataStr = fmt.Sprintf("%v", *n.Data) // Print the value inside the pointer
	}

	return fmt.Sprintf(
		"Node[%T]{\n  Id: %s\n  Name: %s\n  Description: %s\n  CreatedAt: %s\n  ModifiedAt: %s\n  Data: %s\n}",
		n.Data,
		n.Id,
		n.Name,
		n.Description,
		n.CreatedAt.Format(time.RFC3339),
		modified,
		dataStr,
	)
}