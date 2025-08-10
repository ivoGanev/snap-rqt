package view

import (
	"snap-rq/app/entity"
	"snap-rq/app/input"

	"github.com/rivo/tview"
)

const (
	EDITOR_VIEW_FOCUS_HEADERS = 1
	EDITOR_VIEW_FOCUS_BODY    = 2
)

type EditorViewListener interface {
	OnEditorModeChanged(focusedComponent int)
	OnEditorEditTextArea(editorMode int, edit string)
	OnEditorTextAreaSelected()
	OnEditorInputDone()
}

func (r *EditorView) SetListener(l EditorViewListener) {
	r.listener = l
}

type EditorView struct {
	*tview.Flex
	HeadersButton *tview.Button
	BodyButton    *tview.Button
	TextArea      *tview.TextArea
	logicalFocus  int // this is to track what the user intends to edit, not the focused UI element itself
	listener      EditorViewListener
	inputHandler  *input.Handler
}

func NewEditorView(inputHandler *input.Handler) *EditorView {
	editorView := EditorView{
		Flex:          tview.NewFlex(),
		HeadersButton: tview.NewButton("(h) Headers"),
		BodyButton:    tview.NewButton("(b) Body"),
		logicalFocus:  EDITOR_VIEW_FOCUS_HEADERS,
		inputHandler:  inputHandler,
	}
	return &editorView
}

func (r *EditorView) Init(app *AppView) {
	r.SetBorder(true)
	r.SetTitle("Edit Request")
	r.SetDirection(tview.FlexRow)

	r.TextArea = tview.NewTextArea()
	r.TextArea.SetBorder(true)
	r.TextArea.SetChangedFunc(func() {
		r.listener.OnEditorEditTextArea(r.logicalFocus, r.TextArea.GetText())
	})

	// Set input capture
	r.inputHandler.RegisterInputElement(r.TextArea)
	r.inputHandler.SetInputCapture(r.Flex, input.SourceRequestEditor, func(action input.Action) {
		switch action {
		case input.ActionRequestEditorSwitchToBody:
			r.logicalFocus = EDITOR_VIEW_FOCUS_BODY
			r.listener.OnEditorModeChanged(EDITOR_VIEW_FOCUS_BODY)
		case input.ActionRequestEditorSwitchToHeaders:
			r.logicalFocus = EDITOR_VIEW_FOCUS_HEADERS
			r.listener.OnEditorModeChanged(EDITOR_VIEW_FOCUS_HEADERS)
		case input.ActionRequestEditorEdit:
			r.listener.OnEditorTextAreaSelected()
		case input.ActionRequestEditorDone:
			r.listener.OnEditorInputDone()
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
	return r.logicalFocus
}

func (r *EditorView) HasFocusOnInput() bool {
	return r.TextArea.HasFocus()
}

func (r *EditorView) SetTextArea(request entity.Request) {
	if r.logicalFocus == EDITOR_VIEW_FOCUS_HEADERS {
		r.TextArea.SetText(request.Headers, false)
	} else {
		r.TextArea.SetText(request.Body, false)
	}
}
