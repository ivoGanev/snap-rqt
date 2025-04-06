package model

import (
	"slices"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
)

type RequestsListener interface {
	OnRequestsModelChanged(requests *[]data.Node[http.Request], operation CrudOp)
}

type Requests struct {
	data      map[string]data.Node[http.Request]
	listeners []RequestsListener
}

func NewRequestsModel() *Requests {
	return &Requests{
		data: make(map[string]data.Node[http.Request]),
	}
}

func (r *Requests) SetAllData(data *[]data.Node[http.Request]) {
	for _, value := range *data {
		r.data[value.Id] = value
	}
	for _, l := range r.listeners {
		l.OnRequestsModelChanged(data, UPDATE)
	}
}

func (r *Requests) SetData(requestId string, replace *data.Node[http.Request]) {
	if _, exists := r.data[requestId]; exists {
		r.data[requestId] = *replace
	}

	update := &[]data.Node[http.Request]{*replace}
	for _, l := range r.listeners {
		l.OnRequestsModelChanged(update, UPDATE)
	}
}

func (r *Requests) AddListener(l RequestsListener) {
	r.listeners = append(r.listeners, l)
}

func (r *Requests) RemoveListener(l RequestsListener) {
	for i, lis := range r.listeners {
		if lis == l {
			r.listeners = slices.Delete(r.listeners, i, i+1)
			return
		}
	}
}
