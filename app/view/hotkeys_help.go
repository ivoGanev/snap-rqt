package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Navigational help for the user. Shows the user the hot keys etc.
type HotkeysHelp struct {
	*tview.Table
}

func NewHotkeysHelp() *HotkeysHelp {
	return &HotkeysHelp{
		Table: tview.NewTable(),
	}
}

func (n *HotkeysHelp) Init() {
	shortcuts := [][]string{
		{"(e)", "Edit selected item"},
		{"(a)", "Add selected item"},
		{"(d)", "Duplicate selected item"},
		{"(Del)", "Delete selected item"},
		{"(Tab)", "Switch windows"},
		{"(Esc)", "Quit"},
	}
	mid := (len(shortcuts) + 1) / 2 // ceiling divide for uneven counts
	for row := range mid {
		left := shortcuts[row]
		n.SetCell(row, 0,
			tview.NewTableCell(left[0]).
				SetTextColor(tcell.ColorYellow).
				SetAlign(tview.AlignLeft))
		n.SetCell(row, 1,
			tview.NewTableCell(left[1]).
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignLeft))

		if row+mid < len(shortcuts) {
			right := shortcuts[row+mid]
			n.SetCell(row, 3,
				tview.NewTableCell(right[0]).
					SetTextColor(tcell.ColorYellow).
					SetAlign(tview.AlignLeft))
			n.SetCell(row, 4,
				tview.NewTableCell(right[1]).
					SetTextColor(tcell.ColorWhite).
					SetAlign(tview.AlignLeft))
		}
	}
}
