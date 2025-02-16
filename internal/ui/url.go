package ui

import "github.com/rivo/tview"

type Url struct {
	*tview.InputField
}

func NewUrl() *Url {
	return &Url{
		InputField: tview.NewInputField(),
	}
}

func (url *Url) Init() {

}