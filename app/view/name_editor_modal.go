package view

import (
	"snap-rq/app/input"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	EDITOR_MODAL_COMPONENT_REQUESTS   = 0
	EDITOR_MODAL_COMPONENT_COLLETIONS = 1
)

type EditorModalListener interface {
	OnEditorModalSave(text string, component int)
	OnEditorModalCancel()
}

// NameEditorModal is a modal dialog with a single input field for editing a name.
type NameEditorModal struct {
	*tview.Flex
	Input     *tview.InputField
	Save      *tview.Button
	Cancel    *tview.Button
	Form      *tview.Flex
	listener  EditorModalListener
	component int
}

func (n *NameEditorModal) SetListener(listener EditorModalListener) {
	n.listener = listener
}

func NewNameEditorModal(inputHandler *input.Handler) *NameEditorModal {
	inputBox := tview.NewInputField().
		SetLabelColor(tcell.ColorYellow).
		SetFieldBackgroundColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorBlack)
	inputBox.SetBackgroundColor(tcell.ColorBlue)

	saveBtn := tview.NewButton("(Enter) Save")
	saveBtn.SetBorder(true)

	cancelBtn := tview.NewButton("(Esc) Cancel")
	cancelBtn.SetBorder(true)

	// Buttons row
	buttons := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(saveBtn, 0, 1, false).
		AddItem(cancelBtn, 0, 1, false)
	buttons.SetBackgroundColor(tcell.ColorBlue)

	// Main modal body with border and title
	form := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(inputBox, 0, 1, true).
		AddItem(buttons, 0, 1, false)

	form.SetBorder(true).SetTitle("Edit Name")
	form.SetBackgroundColor(tcell.ColorBlue)

	// Center the modal on screen
	vCenter := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(form, 120, 1, false).
		AddItem(nil, 0, 1, false)
	hCenter := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(vCenter, 7, 1, false).
		AddItem(nil, 0, 1, false)

	modal := &NameEditorModal{
		Flex:      hCenter,
		Input:     inputBox,
		Save:      saveBtn,
		Cancel:    cancelBtn,
		Form:      form,
		component: -1,
	}

	// Hook key events to listener
	inputHandler.SetInputCapture(hCenter, input.SourceModalEditor, func(action input.Action) {
		switch action {
		case input.ActionModalSave:
			modal.listener.OnEditorModalSave(modal.Input.GetText(), modal.component)
		case input.ActionModalCancel:
			modal.listener.OnEditorModalCancel()
		}
	})


	// Hook button presses to listener
	saveBtn.SetSelectedFunc(func() {
		if modal.listener != nil {
			modal.listener.OnEditorModalSave(modal.Input.GetText(), modal.component)
		}
	})
	cancelBtn.SetSelectedFunc(func() {
		if modal.listener != nil {
			modal.listener.OnEditorModalCancel()
		}
	})

	return modal
}

func (n *NameEditorModal) Edit(component int) {
	n.component = component
	var title string
	switch component {
	case EDITOR_MODAL_COMPONENT_REQUESTS:
		title = "Edit Request Name"
	case EDITOR_MODAL_COMPONENT_COLLETIONS:
		title = "Edit Collection Name"
	default:
		title = "Err ???"
	}

	n.Form.SetTitle(title)
}
