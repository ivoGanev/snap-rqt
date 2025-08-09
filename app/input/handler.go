package input

import (
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
	Action      Action
	AllowedInModes []Mode
}

type Handler struct {
	bindings  map[Source]map[Binding]ActionBindingSettings
	listeners []ActionListener
	mode      Mode
}


// Warning: If the input capture is set on the App (tview.Application),
// its key bindings will take precedence and override bindings set on
// other views or primitives. This means app-level bindings can
// intercept keys before any view-level handlers get them. In these cases,
// think about if the key should be global to the app or move it to the local UI component insted
func NewHandler() *Handler {
	keyBindings := map[Source]map[Binding]ActionBindingSettings{
		// App focus key bindings
		SourceApp: {
			NewCodeBinding(tcell.KeyTAB): ActionBindingSettings{
				Action:      ActionSwapFocus,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('q'): ActionBindingSettings{
				Action:      ActionFocusCollections,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('w'): ActionBindingSettings{
				Action:      ActionFocusRequests,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('e'): ActionBindingSettings{
				Action:      ActionToggleViewMode,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Collection list key bindings
		SourceCollectionsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:      ActionAddCollection,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:      ActionEditCollectionName,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:      ActionRemoveCollection,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Requests lists key bindings
		SourceRequestsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:      ActionAddRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('d'): ActionBindingSettings{
				Action:      ActionDuplicateRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:      ActionEditRequestName,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:      ActionRemoveRequest,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
		// Modal input editor (e.g. edit names of collection/request)
		SourceModalEditor: {
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:      ActionModalSave,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
			NewCodeBinding(tcell.KeyESC): ActionBindingSettings{
				Action:      ActionModalCancel,
				AllowedInModes: []Mode{ModeNormal, ModeTextInput},
			},
		},
		// Request editor view keys (where we edit the request body/headers etc.
		SourceRequestEditor: {
			NewRuneBinding('b'): ActionBindingSettings{
				Action:      ActionSwitchToBody,
				AllowedInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('h'): ActionBindingSettings{
				Action:      ActionSwitchToHeaders,
				AllowedInModes: []Mode{ModeNormal},
			},
		},
	}

	handler := &Handler{
		mode:     ModeNormal,
		bindings: keyBindings,
	}

	return handler
}

func (h *Handler) SetMode(mode Mode) {
	h.mode = mode
	logger.Info("Input mode set to:", mode)
}

func (h *Handler) AddListener(listener ActionListener) {
	h.listeners = append(h.listeners, listener)
}

func (h *Handler) SetInputCapture(
	p any,
	source Source,
	listener func(action Action),
) {
	var setInputCaptureFunc func(func(event *tcell.EventKey) *tcell.EventKey)

	switch v := p.(type) {
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
	default:
		return
	}

	if listener != nil {
		h.AddListener(listener)
	}

	setInputCaptureFunc(func(event *tcell.EventKey) *tcell.EventKey {
		var binding Binding
		logger.Info("Input key pressed", "key", event.Key(), "rune", event.Rune(), "source", source)
		if event.Key() == tcell.KeyRune {
			binding = NewRuneBinding(event.Rune())
		} else {
			binding = NewCodeBinding(event.Key())
		}

		if actionSetting, ok := h.bindings[source][binding]; ok {
			if slices.Contains(actionSetting.AllowedInModes, h.mode) {
				for _, l := range h.listeners {
					l(actionSetting.Action)
				}
				return nil
			}
		}
		return event
	})
}
