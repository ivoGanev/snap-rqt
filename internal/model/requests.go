package model

import (
	"slices"
	"snap-rq/internal/data"
)

type RequestsModelEventListener interface {
	OnRequestsModelChanged(requests *[]data.Request)
}

type RequestsModel struct {
	request        data.Request
	eventListeners []RequestsModelEventListener
	store          data.Store
}


func (m *RequestsModel) UpdateRequestUrl(selectedRequest *data.Request, text string) {
	panic("unimplemented")
}

func NewRequestModel(store data.Store) *RequestsModel {
	return &RequestsModel{
		store:       store,
	}
}

func (m *RequestsModel) AddListener(l RequestsModelEventListener) {
	m.eventListeners = append(m.eventListeners, l)
}

func (m *RequestsModel) RemoveListener(l RequestsModelEventListener) {
	for i, lis := range m.eventListeners {
		if lis == l {
			m.eventListeners = slices.Delete(m.eventListeners, i, i+1)
			return
		}
	}
}

func (m *RequestsModel) SetRequest(request data.Request) {
	m.requests[request.Id] = request
	for _, listener := range m.eventListeners {
		listener.OnRequestsModelChanged(&m.requests)
	}
}

func (m *RequestsModel) GetRequests() *[]data.Request {
	return &m.requests
}
