package main

import (
	"context"
	"fmt"
	"snap-rq/internal/http"
	"time"
	"github.com/rivo/tview"
)

func main() {
	var app = tview.NewApplication()

	responseView := getResponseView(app)

	table := tview.NewTable()

	table.SetBorder(true)
	table.SetTitle("Requests")

	mocks := http.GenerateMockRequestsAsNodes(1000)
	size := len(mocks)
	count := 0

	for count < size {
		request := mocks[count]

		method := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(request.Data.Method), string(request.Data.Method))
		table.SetCell(count, 0, tview.NewTableCell(method))
		table.SetCell(count, 1, tview.NewTableCell(string(request.Name)))
		count++
	}

	table.SetSelectable(true, true)
	table.Select(0, 1)

	table.SetSelectedFunc(func(row int, column int) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		request := mocks[row].Data
		var response = ""
		response = http.SendRequest(ctx, *request)
		responseView.SetText(response)
	})

	var flex = tview.NewFlex()
	flex.AddItem(table, 0, 1, true).
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
