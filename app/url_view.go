package app

import (
	"github.com/rivo/tview"
)

// URL Input view displays the URL for a given HTTP request
type UrlInputView struct {
	*tview.InputField
	controller UrlInputController
}

func NewUrlInputView(controller UrlInputController) UrlInputView {
	return UrlInputView{
		InputField: tview.NewInputField(),
		controller: controller,
	}
}

func (r *UrlInputView) Init() {
	r.SetTitle("Url")
	r.SetBorder(true)
	r.SetChangedFunc(func(text string) {
		r.controller.HandleUrlTextChanged(text)
	})
}

func (r *UrlInputView) SetUrlText(text string) {
	r.SetText(text)
}
