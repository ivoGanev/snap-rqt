package model

import (
	"slices"
	"snap-rq/internal/data"
)

// Event handlers
type CollectionsListener interface {
	OnCollectionsModelChanged(collections *[]data.Collection, operation CrudOp)
}

func (r *Collections) AddListener(l CollectionsListener) {
	r.listeners = append(r.listeners, l)
}

func (r *Collections) RemoveListener(l CollectionsListener) {
	for i, lis := range r.listeners {
		if lis == l {
			r.listeners = slices.Delete(r.listeners, i, i+1)
			return
		}
	}
}

// Initialisation
type Collections struct {
	data                        map[string]data.Collection
	listeners                   []CollectionsListener
	currentlySelectedCollection string
}

func NewCollectionsModel() *Collections {
	return &Collections{
		data: make(map[string]data.Collection),
	}
}

// CRUD operations
func (r *Collections) SetCollections(data *[]data.Collection) {
	for _, value := range *data {
		r.data[value.Id] = value
	}
	for _, l := range r.listeners {
		l.OnCollectionsModelChanged(data, UPDATE)
	}
}
