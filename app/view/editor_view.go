package view

import (
	"snap-rq/app/entity"
	"snap-rq/app/input"

	"github.com/rivo/tview"
)

const (
	EDITOR_VIEW_MODE_HEADERS = 1
	EDITOR_VIEW_MODE_BODY    = 2
)

type EditorViewListener interface {
	OnEditorModeChanged()
	OnEditorEditTextArea(editorMode int, edit string)
}

func (r *EditorView) SetListener(l EditorViewListener) {
	r.listener = l
}

type EditorView struct {
	*tview.Flex
	app           *tview.Application
	headersButton *tview.Button
	bodyButton    *tview.Button
	textArea      *tview.TextArea
	currentMode   int
	listener      EditorViewListener
	inputHandler  *input.Handler
}

func NewEditorView(app *tview.Application, inputHandler *input.Handler) *EditorView {
	editorView := EditorView{
		Flex:          tview.NewFlex(),
		app:           app,
		headersButton: tview.NewButton("(h) Headers"),
		bodyButton:    tview.NewButton("(b) Body"),
		currentMode:   EDITOR_VIEW_MODE_HEADERS,
		inputHandler:  inputHandler,
	}
	return &editorView
}

func (r *EditorView) Init() {
	r.SetBorder(true)
	r.SetTitle("Edit Request")
	r.SetDirection(tview.FlexRow)

	r.textArea = tview.NewTextArea()
	r.textArea.SetBorder(true)
	r.textArea.SetChangedFunc(func() {
		r.listener.OnEditorEditTextArea(r.currentMode, r.textArea.GetText())
	})

	// Update buttons based on selected mode
	r.updateButtonLabels()
	
	// Set input capture
	r.inputHandler.SetInputCapture(r.Box, input.SourceRequestEditor, func(action input.Action) {
		switch action {
		case input.ActionSwitchToBody:
			r.currentMode = EDITOR_VIEW_MODE_BODY
			r.updateButtonLabels()
			r.app.SetFocus(r.textArea)
			r.listener.OnEditorModeChanged()
		case input.ActionSwitchToHeaders:
			r.currentMode = EDITOR_VIEW_MODE_HEADERS
			r.updateButtonLabels()
			r.app.SetFocus(r.textArea)
			r.listener.OnEditorModeChanged()
		}
	})

	top := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(r.headersButton, 0, 1, false).
		AddItem(r.bodyButton, 0, 1, false)

	r.AddItem(top, 3, 0, false)
	r.AddItem(r.textArea, 0, 1, true)
}

func (r *EditorView) SetTextArea(request entity.Request) {
	if r.currentMode == EDITOR_VIEW_MODE_HEADERS {
		r.textArea.SetText(request.Headers, false)
	} else {
		r.textArea.SetText(request.Body, false)
	}
}

func (r *EditorView) updateButtonLabels() {
	if r.currentMode == EDITOR_VIEW_MODE_HEADERS {
		r.headersButton.SetLabel("[::b][*] (h) Headers[::-]")
		r.bodyButton.SetLabel("   (b) Body")
	} else {
		r.headersButton.SetLabel("   (h) Headers")
		r.bodyButton.SetLabel("[::b][*] (b) Body[::-]")
	}
}
