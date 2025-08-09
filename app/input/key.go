package input

import "github.com/gdamore/tcell/v2"

type KeyType int

const (
	KeyRune KeyType = iota
	KeyCode
)

type Binding struct {
	KeyType KeyType
	Rune    rune      // valid if KeyType == KeyRune
	Key     tcell.Key // valid if KeyType == KeyCode
}

func NewRuneBinding(r rune) Binding {
	return Binding{KeyType: KeyRune, Rune: r}
}

func NewCodeBinding(k tcell.Key) Binding {
	return Binding{KeyType: KeyCode, Key: k}
}
