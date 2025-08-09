package input

import (
	"slices"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ActionListener func(Action)

type Mode int

const (
	ModeNormal Mode = iota
	ModeTextInput
)

type ActionBindingSettings struct {
	Action      Action
	RunsInModes []Mode
}

type Handler struct {
	bindings  map[Source]map[Binding]ActionBindingSettings
	listeners []ActionListener
	mode      Mode
}

func NewHandler() *Handler {
	keyBindings := map[Source]map[Binding]ActionBindingSettings{
		// App focus key bindings
		SourceApp: {
			NewCodeBinding(tcell.KeyTAB): ActionBindingSettings{
				Action:      ActionSwapFocus,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('q'): ActionBindingSettings{
				Action:      ActionFocusCollections,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('w'): ActionBindingSettings{
				Action:      ActionFocusRequests,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('e'): ActionBindingSettings{
				Action:      ActionToggleViewMode,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEsc): ActionBindingSettings{
				Action:      ActionQuit,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEscape): ActionBindingSettings{
				Action:      ActionQuit,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyESC): ActionBindingSettings{
				Action:      ActionQuit,
				RunsInModes: []Mode{ModeNormal},
			},
		},
		// Collection list key bindings
		SourceCollectionsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:      ActionAddCollection,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('n'): ActionBindingSettings{
				Action:      ActionEditCollectionName,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:      ActionRemoveCollection,
				RunsInModes: []Mode{ModeNormal},
			},
		},
		// Requests lists key bindings
		SourceRequestsList: {
			NewRuneBinding('a'): ActionBindingSettings{
				Action:      ActionAddRequest,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('d'): ActionBindingSettings{
				Action:      ActionDuplicateRequest,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('n'): ActionBindingSettings{
				Action:      ActionEditRequestName,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyDelete): ActionBindingSettings{
				Action:      ActionRemoveRequest,
				RunsInModes: []Mode{ModeNormal},
			},
		},
		// Modal input editor (e.g. edit names of collection/request)
		SourceModalEditor: {
			NewCodeBinding(tcell.KeyEnter): ActionBindingSettings{
				Action:      ActionModalSave,
				RunsInModes: []Mode{ModeNormal},
			},
			NewCodeBinding(tcell.KeyEsc): ActionBindingSettings{
				Action:      ActionModalCancel,
				RunsInModes: []Mode{ModeNormal},
			},
		},
		// Request editor view keys (where we edit the request body/headers etc.
		SourceRequestEditor: {
			NewRuneBinding('b'): ActionBindingSettings{
				Action:      ActionSwitchToBody,
				RunsInModes: []Mode{ModeNormal},
			},
			NewRuneBinding('h'): ActionBindingSettings{
				Action:      ActionSwitchToHeaders,
				RunsInModes: []Mode{ModeNormal},
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
}

func (h *Handler) AddListener(listener ActionListener) {
	h.listeners = append(h.listeners, listener)
}

func (h *Handler) emit(action Action) {
	for _, l := range h.listeners {
		l(action)
	}
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
    default:
        return
    }

    setInputCaptureFunc(func(event *tcell.EventKey) *tcell.EventKey {
        var binding Binding
        if event.Key() == tcell.KeyRune {
            binding = NewRuneBinding(event.Rune())
        } else {
            binding = NewCodeBinding(event.Key())
        }

        if actionSetting, ok := h.bindings[source][binding]; ok {
            if slices.Contains(actionSetting.RunsInModes, h.mode) {
                if listener != nil {
                    listener(actionSetting.Action)
                }
                return nil
            }
        }
        return event
    })
}
