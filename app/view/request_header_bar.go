package view

import (
	"snap-rq/app/constants"
	"snap-rq/app/input"
	"snap-rq/app/style"

	"github.com/rivo/tview"
)

type RequestHeaderBarListener interface {
	OnUrlInputTextChanged(text string)
	OnUrlApply()
	OnMethodSelection(method string)
	OnUrlInputLoseFocus()
}

func (r *RequestHeaderBar) SetListener(listener RequestHeaderBarListener) {
	r.listener = listener
}

type RequestHeaderBar struct {
	*tview.Flex
	listener               RequestHeaderBarListener
	styles                 *style.StyleProvider
	RequestMethodDD        *tview.DropDown
	UrlInput               *tview.InputField
	inputHandler           *input.Handler
	suppressMethodCallback bool
}

func NewRequestHeaderBar(styles style.StyleProvider, inputHandler *input.Handler) *RequestHeaderBar {
	return &RequestHeaderBar{
		Flex:            tview.NewFlex(),
		UrlInput:        tview.NewInputField(),
		RequestMethodDD: tview.NewDropDown(),
		inputHandler:    inputHandler,
	}
}

func (r *RequestHeaderBar) Init() {
	r.SetDirection(tview.FlexColumn)

	r.RequestMethodDD.SetTitle("Method")
	r.RequestMethodDD.SetBorder(true)

	r.UrlInput.SetTitle("URL")
	r.UrlInput.SetBorder(true)

	r.AddItem(r.RequestMethodDD, 0, 1, false).
		AddItem(r.UrlInput, 0, 8, false)

	r.UrlInput.SetChangedFunc(func(text string) {
		r.listener.OnUrlInputTextChanged(text)
	})

	r.RequestMethodDD.SetOptions(constants.RequestMethodStrings(), func(text string, index int) {
		if r.suppressMethodCallback {
			return
		}
		r.listener.OnMethodSelection(text)
	})

	// Configure input handler
	r.inputHandler.RegisterInputElement(r.UrlInput)
	r.inputHandler.SetInputCapture(r.UrlInput, input.SourceRequestURLInputBox, func(action input.Action) {
		switch action {
		case input.ActionHeaderBarUrlApply:
			r.listener.OnUrlApply()
		case input.ActionLoseFocus:
			r.listener.OnUrlInputLoseFocus()
		}
	})
}

func (r *RequestHeaderBar) SetUrlText(text string) {
	r.UrlInput.SetText(text)
}

// with silent mode, selecting a method will not trigger its callback.
// This is useful when we don't want to wake up other UI element callbacks and introduce double calls
func (r *RequestHeaderBar) SetRequestMethod(method string, silent bool) {
	if silent {
		r.suppressMethodCallback = true
	}
	methodIndex := constants.RequestMethodIndex(constants.RequestMethod(method))
	r.RequestMethodDD.SetCurrentOption(methodIndex)
	if silent {
		r.suppressMethodCallback = false
	}
}
