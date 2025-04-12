package internal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type NavHelp struct {
	*tview.Table
}

func NewNavigationHelp() *NavHelp {
	return &NavHelp{
		Table: tview.NewTable(),
	}
}

func (n *NavHelp) Init() {
	shortcuts := [][]string{
		{"(e)", "Edit request"},
		{"(c)", "Select collection"},
		{"(q)", "Quit"},
		{"(s)", "Save"},
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
