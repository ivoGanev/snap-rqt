package internal

import "slices"

type EditorModelEventListener interface {
	OnCollectionsModelChanged(collections *[]Collection, operation CrudOp)
	OnRequestsModelChanged(requests *[]Request, operation CrudOp)
}

func (e *EditorModel) AddListener(l EditorModelEventListener) {
	e.eventListeners = append(e.eventListeners, l)
}

func (e *EditorModel) RemoveListener(l EditorModelEventListener) {
	for i, lis := range e.eventListeners {
		if lis == l {
			e.eventListeners = slices.Delete(e.eventListeners, i, i+1)
			return
		}
	}
}

type EditorModel struct {
	store          Store
	collections    []Collection
	requests       []Request
	eventListeners []EditorModelEventListener
	colReqIndex    map[*Collection]map[*Request]bool
}

func NewEditorModel(store Store) EditorModel {
	return EditorModel{
		store: store,
	}
}

func (e *EditorModel) Init() {
	collections, err := e.store.LoadAllCollections()
	if err != nil {
		panic(err)
	}
	e.collections = collections

	requests, err := e.store.LoadAllRequests()
	if err != nil {
		panic(err)
	}
	e.requests = requests

	// build the index

}

func (e *EditorModel) AddCollection(collection *Collection) {
	if !e.store.CheckExistsCollection(collection) {
		e.store.StoreCollection(collection)
	}
}

func (e *EditorModel) UpdateCollection(collection *Collection) {

}

func (e *EditorModel) DeleteCollection(collection *Collection) {
	// delete all request pointers, and if it's the last collection, delete the request values
}

func (e *EditorModel) AddRequest(collection *Collection, request *Request) {
	if !e.store.CheckExistsCollection(collection) {
		e.store.StoreCollection(collection)
	}

	// add the request to persistent store if not there
	if !e.store.CheckExistsRequest(request) {
		e.store.StoreRequest(request)
	}

	// After all persistence, add the bindings
	e.colReqIndex[collection][request] = true

	for _, l := range e.eventListeners {
		l.OnRequestsModelChanged(&[]Request{*request}, CREATE)
	}
}

func (e *EditorModel) UpdateRequest(request *Request) {

}

func (e *EditorModel) RemoveRequest(collection *Collection, request *Request) {
	// if this is the only collection left for this request, then delete the request value from memory (not only the pointer)
}
