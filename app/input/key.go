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
	Modifiers tcell.ModMask
}

func NewCodeBinding(k tcell.Key) Binding {
    return Binding{KeyType: KeyCode, Key: k}
}

func NewCodeBindingWithModifier(k tcell.Key, m tcell.ModMask) Binding {
    return Binding{KeyType: KeyCode, Key: k, Modifiers: m}
}

func NewRuneBinding(r rune) Binding {
    return Binding{KeyType: KeyRune, Rune: r}
}

func NewRuneBindingWithModifier(r rune, m tcell.ModMask) Binding {
    return Binding{KeyType: KeyRune, Rune: r, Modifiers: m}
}
