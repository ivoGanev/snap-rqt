package input

type Source string

const (
	SourceRequestsList       Source = "requests"
	SourceCollectionsList    Source = "collections"
	SourceApp                Source = "app"
	SourceModalEditor        Source = "viewModalEditor"
	SourceRequestEditor      Source = "requestEditor"
	SourceRequestURLInputBox Source = "requestUrlInput"
)
