package model

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"
)

type RequestsListener interface {
	OnRequestsModelChanged(requests []data.Node[http.Request], operation CrudOp, multiplicity Multiplicity)
}

type Requests struct {
	data      map[string]data.Node[http.Request]
	listeners []RequestsListener
}

func NewRequests() *Requests {
	return &Requests{
		data: make(map[string]data.Node[http.Request]),
	}
}

func (r *Requests) SetData(data []data.Node[http.Request]) {
	for _, value := range data {
		r.data[value.Id] = value
	}
	for _, l := range r.listeners {
		l.OnRequestsModelChanged(data, UPDATE, MANY)
	}
}

func (r *Requests) AddListener(l RequestsListener) {
	r.listeners = append(r.listeners, l)
}

func (r *Requests) RemoveListener(l RequestsListener) {
	for i, lis := range r.listeners {
		if lis == l {
			r.listeners = append(r.listeners[:i], r.listeners[i+1:]...)
			return
		}
	}
}
