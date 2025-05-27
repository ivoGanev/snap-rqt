package view

import (
	"github.com/rivo/tview"
)

type UrlInputListener interface {
	OnUrlInputTextChanged(text string)
}

func (r *UrlInputView) SetListener(listener UrlInputListener) {
	r.listener = listener
}

// URL Input view displays the URL for a given HTTP request
type UrlInputView struct {
	*tview.InputField
	listener UrlInputListener
}

func NewUrlInput() *UrlInputView {
	return &UrlInputView{
		InputField:    tview.NewInputField(),
	}
}

func (r *UrlInputView) Init() {
	r.SetTitle("Url")
	r.SetBorder(true)
	r.SetChangedFunc(func(text string) {
		r.listener.OnUrlInputTextChanged(text)
	})
}

func (r *UrlInputView) SetUrlText(text string) {
	r.SetText(text)
}
