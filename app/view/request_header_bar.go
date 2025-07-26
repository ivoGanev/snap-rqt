package view

import (
	"snap-rq/app/constants"
	"snap-rq/app/style"

	"github.com/rivo/tview"
)

type RequestHeaderBarListener interface {
	OnUrlInputTextChanged(text string)
	OnMethodSelection(method string)
}

func (r *RequestHeaderBar) SetListener(listener RequestHeaderBarListener) {
	r.listener = listener
}

type RequestHeaderBar struct {
	*tview.Flex
	listener      RequestHeaderBarListener
	styles        *style.StyleProvider
	requestMethod *tview.DropDown
	urlInput      *tview.InputField
}

func NewRequestHeaderBar(styles style.StyleProvider) *RequestHeaderBar {
	return &RequestHeaderBar{
		Flex:          tview.NewFlex(),
		urlInput:      tview.NewInputField(),
		requestMethod: tview.NewDropDown(),
	}
}

func (r *RequestHeaderBar) Init() {
	r.SetDirection(tview.FlexColumn)

	r.requestMethod.SetTitle("Method")
	r.requestMethod.SetBorder(true)

	r.urlInput.SetTitle("URL")
	r.urlInput.SetBorder(true)

	r.AddItem(r.requestMethod, 0, 1, false).
		AddItem(r.urlInput, 0, 8, false)

	r.urlInput.SetChangedFunc(func(text string) {
		r.listener.OnUrlInputTextChanged(text)
	})

	r.requestMethod.SetOptions(constants.RequestMethodStrings(), func(text string, index int) {
		r.listener.OnMethodSelection(text)
	})
}

func (r *RequestHeaderBar) SetUrlText(text string) {
	r.urlInput.SetText(text)
}

func (r *RequestHeaderBar) SetRequestMethod(method string) {
	methodIndex := constants.RequestMethodIndex(constants.RequestMethod(method))
	r.requestMethod.SetCurrentOption(methodIndex)
}
