package view

import "github.com/rivo/tview"

type StatusBar struct {
	*tview.TextView
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		TextView: tview.NewTextView(),
	}
}

func Init() {

}

func (s *StatusBar) SetStatusText(text string) {
	s.SetText(text)
}