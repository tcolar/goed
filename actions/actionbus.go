package actions

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/tcolar/goed/core"
)

var _ core.ActionDispatcher = (*actionBus)(nil)
var latestRenderAction int64

type actionBus struct {
	actionChan chan (core.Action)
	quitc      chan (struct{})
}

func NewActionBus() core.ActionDispatcher {
	return actionBus{
		actionChan: make(chan (core.Action), 1000),
		quitc:      make(chan (struct{})),
	}
}

func (a actionBus) Dispatch(action core.Action) {
	a.actionChan <- action
}

func (a actionBus) Start() {
	pause := 10 * time.Millisecond
	// handle events
	for {
		select {
		case <-a.quitc:
			break

		case action := <-a.actionChan:
			switch a := action.(type) {
			case edRender:
				time.Sleep(pause)
				// If we have a bunch of render pending, only need to honor the most recent.
				if a.time == atomic.LoadInt64(&latestRenderAction) {
					core.Ed.Render()
				}
			default:
				if core.Trace {
					log.Printf("> %#v", action)
				}
				action.Run()
				if core.Trace {
					log.Printf("< %#v", action)
				}
			}
		}
	}
}

// Flush waits for all actions sent before it to have been processed
func (a actionBus) Flush() {
	c := make(chan (struct{}), 1)
	d(flushAction{c})
	<-c
}

func (a actionBus) Shutdown() {
	a.quitc <- struct{}{}
}

type flushAction struct {
	c chan (struct{})
}

func (a flushAction) Run() {
	a.c <- struct{}{}
}
