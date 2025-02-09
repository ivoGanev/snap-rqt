package main

import (
	"context"
	"fmt"
	"snap-rq/internal/data"
	"snap-rq/internal/http"
	"snap-rq/internal/mocks"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppState struct {
	SelectedRequestCell int
	SelectedNode        data.Node[http.Request]
	Mocks               []data.Node[http.Request]
}

var state = AppState{
	Mocks: mocks.GenerateMockRequests(1000),
}

type AppComponents struct {
	App                  *tview.Application
	Debugger             *tview.TextArea
	Requests             *tview.Table
	MethodSelectionModal *tview.Grid
	Pages                *tview.Pages
}

var c = AppComponents{
	App:                  tview.NewApplication(),
	Debugger:             tview.NewTextArea(),
	Requests:             tview.NewTable(),
	MethodSelectionModal: tview.NewGrid(),
	Pages:                tview.NewPages(),
}

const (
	PAGE_REQUEST_METHOD_PICKER_MODAL = "request-method-picker"
	PAGE_LANDING_VIEW                = "landing-view"
)

func main() {
	c.App.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		c.Debugger.SetText(state.SelectedNode.String(), false)
		return false // Allow normal drawing to continue
	})

	responseView := getResponseView(c)

	requests := c.Requests

	requests.SetBorder(true)
	requests.SetTitle("Requests")

	size := len(state.Mocks)
	count := 0

	for count < size {
		request := state.Mocks[count]

		method := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(request.Data.Method), string(request.Data.Method))

		methodCell := tview.NewTableCell(method)
		methodCell.SetReference(request)

		nameCell := tview.NewTableCell(string(request.Name))
		nameCell.SetReference(request)
		requests.SetCell(count, 0, methodCell)
		requests.SetCell(count, 1, nameCell)
		count++
	}

	requests.SetSelectable(true, true)
	requests.Select(0, 1)
	state.SelectedNode = state.Mocks[0]
	state.SelectedRequestCell = 0

	requests.SetSelectionChangedFunc(func(row int, column int) {
		state.SelectedNode = state.Mocks[row]
		state.SelectedRequestCell = row
	})

	requests.SetSelectedFunc(func(row int, column int) {
		ref := requests.GetCell(row, column).GetReference()
		request, ok := ref.(data.Node[http.Request])
		if !ok {
			panic("Failed to cast reference to *http.Request")
		} else {
			if column == 0 {
				c.Pages.ShowPage(PAGE_REQUEST_METHOD_PICKER_MODAL)
			} else {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var response = ""
				response = http.SendRequest(ctx, *request.Data)
				responseView.SetText(response)
			}
		}
	})

	var mainWindows = tview.NewFlex()

	mainWindows.AddItem(requests, 0, 1, true).
		AddItem(responseView, 0, 1, false)

	var rootFlex = tview.NewFlex()

	rootFlex.SetDirection(tview.FlexRow).
		AddItem(mainWindows, 0, 10, true).
		AddItem(c.Debugger, 0, 1, false)

	c.Pages.
		AddPage(PAGE_LANDING_VIEW, rootFlex, true, true).
		AddPage(PAGE_REQUEST_METHOD_PICKER_MODAL, getMethodSelectionView(), true, false)

	if err := c.App.SetFocus(c.Pages).
		SetRoot(c.Pages, true).
		EnableMouse(true).
		Run(); err != nil {
		panic(err)
	}
}

func getMethodSelectionView() tview.Primitive {
	methodList := tview.NewList()
	methodList.ShowSecondaryText(false).
		SetBorder(true)

	for _, method := range http.RequestMethods {
		methodList.AddItem(string(method), "", 0, nil).
			SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				for _, m := range state.Mocks {
					if m.Id == state.SelectedNode.Id { // this needs to be optimised with a hashmap
						m.Data.Method = http.RequestMethod(mainText)

						method := fmt.Sprintf("%s %s [white]", http.GetTcellColorForRequest(m.Data.Method), mainText)
						c.Requests.GetCell(state.SelectedRequestCell, 0).SetText(method)
						c.Pages.HidePage(PAGE_REQUEST_METHOD_PICKER_MODAL)
						c.App.SetFocus(c.Requests)
					}
				}
			})
	}

	modal := tview.NewGrid().
		SetColumns(0, 10, 0).
		SetRows(0, 7, 0).
		AddItem(methodList, 1, 1, 1, 1, 0, 0, true)

	return modal
}

func getResponseView(c AppComponents) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			c.App.Draw()
		})

	textView.SetText("No response data")
	textView.SetBorder(true)
	textView.SetTitle("Response")

	return textView
}
