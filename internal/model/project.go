package model

import "snap-rq/internal/data"

// A project contains many collections.
// A single collection holds many requets.
const (
	CREATE = iota
	READ
	UPDATE
	DELETE
)

type Operation int

var Operations = []Operation{CREATE, READ, UPDATE, DELETE}

type CollectionsEventListener interface {
	OnCollectionChanged(collection *[]data.Collection, operation Operation)
}

type RequestEventListener interface {
	OnRequestChanged(request *data.Request, operation Operation)
}

type OnLoadSession interface {
	OnProjectModelLoaded()
}

type ProjectModel struct {
	store                         data.Store
	requestUpdateListeners        []RequestEventListener
	onProjectModelLoadedListeners []OnLoadSession
}

func NewProjectModel(store data.Store) *ProjectModel {
	return &ProjectModel{
		store: store,
	}
}


func (m *ProjectModel) UpdateRequestUrl(requestId string, text string) {
	request, _ := m.store.GetRequest(requestId)
	request.Data.Url = text
	m.store.UpdateRequest(request)

	for _, listener := range m.requestUpdateListeners {
		listener.OnRequestChanged(&request, UPDATE)
	}
}

func (m *ProjectModel) GetCollection(collectionId string) *data.Collection {
	collection, _ := m.store.GetCollection(collectionId)
	return &collection
}

func (m *ProjectModel) FilterRequestsByName(collectionId string) *data.Collection {
	collection, _ := m.store.GetCollection(collectionId)
	return &collection
}

func (m *ProjectModel) GetRequest(requestId string) *data.Request {
	request, _ := m.store.GetRequest(requestId)
	return &request
}


func (m *ProjectModel) GetRequestSummary(requestId string) *data.Request {
	request, _ := m.store.GetRequest(requestId)
	return &request
}
