package input

import (
	"fmt"
	"slices"
	logger "snap-rq/app/log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ActionListener func(Action)

type Mode int

const (
	ModeNormal Mode = iota
	ModeTextInput
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "ModeNormal"
	case ModeTextInput:
		return "ModeTextInput"
	default:
		return "UnknownMode"
	}
}

type ActionBindingSettings struct {
	Action         Action
	AllowedInModes []Mode
}

type Handler struct {
	bindings   map[Source]map[Binding]ActionBindingSettings
	listeners  []ActionListener
	inputViews []tview.Primitive
}

// Warning: If the input capture is set on the SourceApp,
// its key bindings will take precedence and override bindings set on
// other views or primitives. This means app-level bindings can
// intercept keys before any view-level handlers get them. In these cases,
// think about if the key should be global to the app or move it to the local UI component instead.
func NewHandler() *Handler {
	keyBindings := map[Source]map[Binding]ActionBindingSettings{
		// App focus key bindings
		SourceApp: {
			NewCodeBinding(tcell.KeyTAB): ActionBindingSettings{
				Action:         ActionSwapFocus,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('q'): ActionBindingSettings{
				Action:         ActionFocusCollections,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('w'): ActionBindingSettings{
				Action:         ActionFocusRequests,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('e'): ActionBindingSettings{
				Action:         ActionSwapPuppetModes,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('u'): ActionBindingSettings{
				Action:         ActionHeaderBarEditUrl,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('m'): ActionBindingSettings{
				Action:         ActionHeaderBarSelectMethod,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Collection list key bindings
		SourceCollectionsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:         ActionAddCollection,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyF2): ActionBindingSettings{
				Action:         ActionEditCollectionName,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:         ActionRemoveCollection,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Requests lists key bindings
		SourceRequestsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:         ActionAddRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('d'): ActionBindingSettings{
				Action:         ActionDuplicateRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyF2): ActionBindingSettings{
				Action:         ActionEditRequestName,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:         ActionRemoveRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:         ActionSelectRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Modal input editor (e.g. edit names of collection/request)
		SourceModalEditor: {
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:         ActionModalSave,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
			NewCodeBinding(tcell.KeyESC): ActionBindingSettings{
				Action:         ActionModalCancel,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
		},
		// Request editor view keys (where we edit the request body/headers etc.
		SourceRequestEditor: {
			NewRuneBinding('b'): ActionBindingSettings{
				Action:         ActionRequestEditorSwitchToBody,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('h'): ActionBindingSettings{
				Action:         ActionRequestEditorSwitchToHeaders,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyESC): ActionBindingSettings{
				Action:         ActionRequestEditorDone,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:         ActionRequestEditorEdit,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDown): ActionBindingSettings{
				Action:         ActionRequestEditorEdit,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		SourceRequestURLInputBox: {
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:         ActionHeaderBarUrlApply,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
			NewCodeBinding(tcell.KeyESC): ActionBindingSettings{
				Action:         ActionLoseFocus,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
		},
	}

	handler := &Handler{
		bindings: keyBindings,
	}

	return handler
}

func (h *Handler) SetInputCapture(
	tviewObject any,
	source Source,
	listener func(action Action),
) {
	var setInputCaptureFunc func(func(event *tcell.EventKey) *tcell.EventKey)

	switch v := tviewObject.(type) {
	case *tview.Box:
		setInputCaptureFunc = func(f func(event *tcell.EventKey) *tcell.EventKey) {
			_ = v.SetInputCapture(f)
		}
	case *tview.Application:
		setInputCaptureFunc = func(f func(event *tcell.EventKey) *tcell.EventKey) {
			_ = v.SetInputCapture(f)
		}
	case *tview.Flex:
		setInputCaptureFunc = func(f func(event *tcell.EventKey) *tcell.EventKey) {
			_ = v.SetInputCapture(f)
		}
	case *tview.TextArea:
		setInputCaptureFunc = func(f func(event *tcell.EventKey) *tcell.EventKey) {
			_ = v.SetInputCapture(f)
		}
	case *tview.InputField:
		setInputCaptureFunc = func(f func(event *tcell.EventKey) *tcell.EventKey) {
			_ = v.SetInputCapture(f)
		}
	default:
		panic(fmt.Sprintf("unsupported input capture type: %T", tviewObject))
	}

	if listener != nil {
		h.listeners = append(h.listeners, listener)
	}

	setInputCaptureFunc(func(event *tcell.EventKey) *tcell.EventKey {
		var binding Binding
		logger.Info("Input key pressed", "key", event.Key(), "rune", event.Rune(),
			"mods", event.Modifiers(), "source", source)

		if event.Key() == tcell.KeyRune {
			binding = NewRuneBindingWithModifier(event.Rune(), event.Modifiers())
		} else {
			binding = NewCodeBindingWithModifier(event.Key(), event.Modifiers())
		}
		if actionSetting, ok := h.bindings[source][binding]; ok {
			// If the element allows input (e.g. InputField), accept also bindings allowed in ModeTextInput
			// Otherwise, only accept bindings allowed in ModeNormal
			allowed := true

			for _, elem := range h.inputViews {
				if elem.HasFocus() && !slices.Contains(actionSetting.AllowedInModes, ModeTextInput) {
					allowed = false
				}
			}

			if allowed {
				for _, l := range h.listeners {
					l(actionSetting.Action)
				}
				return nil
			}
		}
		return event
	})

}

func (h *Handler) RegisterInputElement(tviewObject tview.Primitive) {
	h.inputViews = append(h.inputViews, tviewObject)
}
