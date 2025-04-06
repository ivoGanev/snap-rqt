package ui

import (
	"snap-rq/internal/data"
	"snap-rq/internal/http"

	"github.com/rivo/tview"
)

type Collections struct {
	*tview.Table
	app *App
	SelectedNode    *data.Node[http.Request]
}
