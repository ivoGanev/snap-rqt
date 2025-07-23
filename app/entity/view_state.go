package entity

import (
	"fmt"
	"strings"
)

type AppViewState struct {
	FocusedView         string
	FocusedRequestIds   map[string]string // Mapping: []CollectionId -> RequestId
	FocusedCollectionId string
}

func (s AppViewState) String() string {
	var builder strings.Builder

	builder.WriteString("FocusedView: " + s.FocusedView + "\n")
	builder.WriteString("FocusedCollectionId: " + s.FocusedCollectionId + "\n")
	// builder.WriteString("FocusedRequestIds:\n")

	for collectionID, requestID := range s.FocusedRequestIds {
		builder.WriteString(fmt.Sprintf("  [%s] -> %s\n", collectionID, requestID))
	}

	return builder.String()
}
