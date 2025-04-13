package model

import (
	"slices"
	"snap-rq/internal/data"
)

type CollectionsModelEventListener interface {
	OnCollectionsModelChanged(collections *[]data.Collection)
}


type CollectionsModel struct {
	collections    map[string][]data.Collection
	eventListeners []CollectionsModelEventListener
	store          data.Store
}

func (m *CollectionsModel) Load() {
	panic("unimplemented")
}

func NewCollectionModel(store data.Store) *CollectionsModel {
	return &CollectionsModel{
		store: store,
	}
}

func (m *CollectionsModel) AddListener(l CollectionsModelEventListener) {
	m.eventListeners = append(m.eventListeners, l)
}

func (m *CollectionsModel) RemoveListener(l CollectionsModelEventListener) {
	for i, lis := range m.eventListeners {
		if lis == l {
			m.eventListeners = slices.Delete(m.eventListeners, i, i+1)
			return
		}
	}
}

func (m *CollectionsModel) SetCollections(collections []data.Collection) {
	m.collections = collections
	for _, listener := range m.eventListeners {
		listener.OnCollectionsModelChanged(&m.collections)
	}
}

func (m *CollectionsModel) GetCollections() *[]data.Collection {
	return &m.collections
}

func (m *CollectionsModel) GetCollection(string collectionId) *[]data.Collection {
	return &m.collections
}
