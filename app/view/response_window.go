package view

import (
	"context"
	"fmt"
	"snap-rq/app/entity"
	logger "snap-rq/app/log"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseWindow struct {
	view                *tview.TextView
	app                 *tview.Application
	ctx         context.Context
	cancelFunc  context.CancelFunc	
}

func NewResponseWindow(app *tview.Application) *ResponseWindow {
	responseView := ResponseWindow{
		view: tview.NewTextView(),
		app:  app,
	}

	return &responseView
}

func (r *ResponseWindow) Init() {
	r.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'p' {
			return tcell.NewEventKey(tcell.KeyCtrlV, 'v', tcell.ModNone)
		}
		if event.Rune() == 'c' {
			return tcell.NewEventKey(tcell.KeyCtrlQ, 'c', tcell.ModNone)
		}
		return event
	})

	r.view.SetText("No response data")
	r.view.SetBorder(true)
	r.view.SetTitle("Response")
}

func (r *ResponseWindow) AwaitResponse() {
	r.stopPreviousAnimation()

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())

	frames := []string{".", "..", "..."}
	current := 0
	r.view.SetText(fmt.Sprintf("Requesting data%s", frames[current]))
	current = (current + 1) % len(frames)

	logger.Println("Awaiting response")
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
				r.app.QueueUpdateDraw(func() {
					r.view.SetText(fmt.Sprintf("Requesting data%s", frames[current]))
					current = (current + 1) % len(frames)
				})
			}
		}
	}()
}

func (r *ResponseWindow) SetError(err error) {
	r.stopPreviousAnimation()
	r.app.QueueUpdateDraw(func() {
		r.view.SetText(err.Error())
	})
}

func (r *ResponseWindow) SetHttpResponse(response entity.HttpResponse) {
	r.stopPreviousAnimation()
	r.app.QueueUpdateDraw(func() {
		r.view.SetText(response.Body)
	})
}

func (r *ResponseWindow) stopPreviousAnimation() {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
}