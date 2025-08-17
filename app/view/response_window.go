package view

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"snap-rq/app/entity"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseWindow struct {
	*tview.TextView
	app        *tview.Application
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewResponseWindow(app *tview.Application) *ResponseWindow {
	responseView := ResponseWindow{
		TextView: tview.NewTextView(),
		app:      app,
	}

	return &responseView
}

func (r *ResponseWindow) Init() {
	r.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'p' {
			return tcell.NewEventKey(tcell.KeyCtrlV, 'v', tcell.ModNone)
		}
		if event.Rune() == 'c' {
			return tcell.NewEventKey(tcell.KeyCtrlQ, 'c', tcell.ModNone)
		}
		return event
	})

	r.SetText("No response data")
	r.SetBorder(true)
	r.SetTitle("Response")
}

func (r *ResponseWindow) AwaitResponse() {
	r.stopPreviousAnimation()

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())

	frames := []string{".", "..", "..."}
	current := 0
	r.SetText(fmt.Sprintf("Requesting data%s", frames[current]))
	current = (current + 1) % len(frames)

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
				r.app.QueueUpdateDraw(func() {
					r.SetText(fmt.Sprintf("Requesting data%s", frames[current]))
					current = (current + 1) % len(frames)
				})
			}
		}
	}()
}

func (r *ResponseWindow) SetError(err error) {
	r.stopPreviousAnimation()
	r.app.QueueUpdateDraw(func() {
		r.SetText(err.Error())
	})
}

func (r *ResponseWindow) SetHttpResponse(response entity.HttpResponse) {
	r.stopPreviousAnimation()

	contentTypes, ok := response.Header["Content-Type"]
	body := response.Body

	var output string

	// Pretty-print JSON
	if ok && len(contentTypes) > 0 && strings.HasPrefix(contentTypes[0], "application/json") {
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, []byte(body), "", "  "); err != nil {
			output = fmt.Sprintf("invalid JSON: %v\n\n%s", err, body)
		} else {
			output = pretty.String()
		}
	} else {
		output = body
	}

	r.app.QueueUpdateDraw(func() {
		r.SetText(output)
	})
}

func (r *ResponseWindow) stopPreviousAnimation() {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
}
