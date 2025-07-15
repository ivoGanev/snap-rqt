package view

import "github.com/rivo/tview"

type EditorView struct {
	*tview.Flex
}

func NewEditorView() *EditorView {
	editorView := EditorView{
		Flex: tview.NewFlex(),
	}
	return &editorView
}

func (r EditorView) Init() {
	r.SetBorder(true)
	r.SetTitle("Edit Request")
}
