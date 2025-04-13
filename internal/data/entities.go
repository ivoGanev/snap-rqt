package data

import "snap-rq/internal/http"

type Collection struct {
	Node[[]Request]
}

type CollectionSimple struct {
	Node[[]RequestSimple]
}

type Request struct {
	Node[http.Request]
}

type UserSession struct {
	RequestId    string
	CollectionId string
}

type RequestSimple struct {
	Id         string
	Url        string
	Name       string
	MethodType http.RequestMethod
}

func NewRequestSimple(r Request) RequestSimple {
	return RequestSimple{
		Id:         r.Id,
		Url:        r.Data.Url,
		Name:       r.Name,
		MethodType: r.Data.Method,
	}
}
