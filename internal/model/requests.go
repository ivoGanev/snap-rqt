package model

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"
)

type RequestsListener interface {
	OnRequestAdded(string)
	OnRequestRemoved(string)
	OnRequestChanged(string)
	OnRequestsSet([]data.Node[http.Request])
}

type Requests struct {
	data      []data.Node[http.Request]
	listeners []RequestsListener
}

func NewRequests() *Requests {
	return &Requests{}
}

func (r *Requests) SetData(data []data.Node[http.Request]) {
	for _, l := range r.listeners {
		l.OnRequestsSet(data)
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

func (r *Requests) Data() *[]data.Node[http.Request] {
	return &r.data
}
