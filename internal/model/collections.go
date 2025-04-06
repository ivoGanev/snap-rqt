package model

import (
	"snap-rq/internal/data"
)

type CollectionsListener interface {
	OnCollectionsModelChanged(collections *[]Collections, operation CrudOp)
}

type Collections struct {
	data      []data.Collection
	listeners []CollectionsListener
}

func NewCollectionsModel() *Collections {
	return &Collections{
		data: []data.Collection {},
	}
}