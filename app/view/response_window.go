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
	      *tview.TextView
	app        *tview.Application
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewResponseWindow(app *tview.Application) *ResponseWindow {
	responseView := ResponseWindow{
		TextView: tview.NewTextView(),
		app:  app,
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

	logger.Println("Awaiting response")
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
	r.app.QueueUpdateDraw(func() {
		r.SetText(response.Body)
	})
}

func (r *ResponseWindow) stopPreviousAnimation() {
	if r.cancelFunc != nil {
		r.cancelFunc()
	}
}
