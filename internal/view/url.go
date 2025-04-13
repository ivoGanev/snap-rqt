package view

import (
	"github.com/rivo/tview"
)

type UrlInputListener interface {
	OnUrlTextChanged(string)
}

func (r *UrlInput) SetUrlInputListener(l UrlInputListener) {
	r.eventListener = l
}

// URL Input view displays the URL for a given HTTP request
type UrlInput struct {
	*tview.InputField
	eventListener   UrlInputListener
}

func NewUrlInput() *UrlInput {
	return &UrlInput{
		InputField: tview.NewInputField(),
	}
}

func (r *UrlInput) Init() {
	r.SetTitle("Url")
	r.SetBorder(true)
	r.SetChangedFunc(func(text string) {
		if r.eventListener != nil {
			r.eventListener.OnUrlTextChanged(text)
		}
	})
}

func (r *UrlInput) SetUrlText(text string) {
	r.SetText(text)
}