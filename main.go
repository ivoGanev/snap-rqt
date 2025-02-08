package main

import (
	"context"
	"snap-rq/internal/http"
	"time"

	"github.com/rivo/tview"
)

func main() {
	var app = tview.NewApplication()

	requestsList := tview.NewList().ShowSecondaryText(false)
	responseView := getResponseView(app)

	// load mocks
	for _, val := range http.GetMockRequestsAsNodes() {
		requestsList.AddItem(val.Name, "", 0, nil)
	}

	// the rest
	requestsList.SetBorder(true)
	requestsList.SetTitle("Requests")

	requestsList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		request := http.GetMockRequestsAsNodes()[index].Data
		var response = ""
		response = http.SendRequest(ctx, *request)
		responseView.SetText(response)
	})

	var flex = tview.NewFlex()
	flex.AddItem(requestsList, 0, 1, true).
		AddItem(responseView, 0, 1, false)

	if err := app.SetFocus(flex).
		SetRoot(flex, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func getResponseView(app *tview.Application) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	textView.SetText("No response data")
	textView.SetBorder(true)
	textView.SetTitle("Response")

	return textView
}
