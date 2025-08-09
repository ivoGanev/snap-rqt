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
	OnEditorModeChanged(mode int)
	OnEditorEditTextArea(editorMode int, edit string)
	OnEditorTextAreaSelected()
}

func (r *EditorView) SetListener(l EditorViewListener) {
	r.listener = l
}

type EditorView struct {
	*tview.Flex
	HeadersButton *tview.Button
	BodyButton    *tview.Button
	TextArea      *tview.TextArea
	currentMode   int
	listener      EditorViewListener
	inputHandler  *input.Handler
}

func NewEditorView(inputHandler *input.Handler) *EditorView {
	editorView := EditorView{
		Flex:          tview.NewFlex(),
		HeadersButton: tview.NewButton("(h) Headers"),
		BodyButton:    tview.NewButton("(b) Body"),
		currentMode:   EDITOR_VIEW_MODE_HEADERS,
		inputHandler:  inputHandler,
	}
	return &editorView
}

func (r *EditorView) Init() {
	r.SetBorder(true)
	r.SetTitle("Edit Request")
	r.SetDirection(tview.FlexRow)

	r.TextArea = tview.NewTextArea()
	r.TextArea.SetBorder(true)
	r.TextArea.SetChangedFunc(func() {
		r.listener.OnEditorEditTextArea(r.currentMode, r.TextArea.GetText())
	})

	// Make sure key input mode is set correct
	r.TextArea.SetFocusFunc(func() {
		r.inputHandler.SetMode(input.ModeTextInput)
	})
	r.TextArea.SetBlurFunc(func() {
		r.inputHandler.SetMode(input.ModeNormal)
	})

	// Set input capture
	r.inputHandler.SetInputCapture(r.Box, input.SourceRequestEditor, func(action input.Action) {
		switch action {
		case input.ActionRequestEditorSwitchToBody:
			r.currentMode = EDITOR_VIEW_MODE_BODY
			r.listener.OnEditorModeChanged(EDITOR_VIEW_MODE_BODY)
		case input.ActionRequestEditorSwitchToHeaders:
			r.currentMode = EDITOR_VIEW_MODE_HEADERS
			r.listener.OnEditorModeChanged(EDITOR_VIEW_MODE_HEADERS)
		case input.ActionRequestEditorEdit:
			r.listener.OnEditorTextAreaSelected()
		}
	})

	top := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(r.HeadersButton, 0, 1, false).
		AddItem(r.BodyButton, 0, 1, false)

	r.AddItem(top, 3, 0, false)
	r.AddItem(r.TextArea, 0, 1, true)
}

func (r *EditorView) GetCurrentMode() int {
	return r.currentMode
}

func (r *EditorView) SetTextArea(request entity.Request) {
	if r.currentMode == EDITOR_VIEW_MODE_HEADERS {
		r.TextArea.SetText(request.Headers, false)
	} else {
		r.TextArea.SetText(request.Body, false)
	}
}
