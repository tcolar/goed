package actions

import (
	"log"

	"github.com/kr/pretty"
	"github.com/tcolar/goed/core"
)

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

func (e actionBus) Dispatch(action core.Action) {
	e.actionChan <- action
}

func (e actionBus) Start() {
	for {
		select {
		case action := <-e.actionChan:
			if core.Trace {
				pretty.Logln(action)
			}
			err := action.Run()
			if err != nil {
				core.Ed.SetStatusErr(err.Error())
				log.Println(err.Error())
			}
		case <-e.quitc:
			break
		}
	}
}

func (e actionBus) Shutdown() {
	e.quitc <- struct{}{}
}
