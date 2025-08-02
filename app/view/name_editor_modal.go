package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// NameEditorModal is a modal dialog with a single input field for editing a name.
type NameEditorModal struct {
	tview.Primitive
	Input  *tview.InputField
	Save   *tview.Button
	Cancel *tview.Button
}

func NewNameEditorModal() *NameEditorModal {
	input := tview.NewInputField().
		SetLabelColor(tcell.ColorYellow).
		SetFieldBackgroundColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorBlack)
	input.SetBackgroundColor(tcell.ColorBlue)

	saveBtn := tview.NewButton("Save")
	saveBtn.SetBorder(true)

	cancelBtn := tview.NewButton("Cancel")
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
		AddItem(input, 0, 1, true).
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


	return &NameEditorModal{
		Primitive: hCenter,
		Input:     input,
		Save:      saveBtn,
		Cancel:    cancelBtn,
	}
}
