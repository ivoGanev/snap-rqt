package entity

type AppViewState struct {
	FocusedView         string
	FocusedRequestIds   map[string]string // Mapping: []CollectionId -> RequestId
	FocusedCollectionId string
}
