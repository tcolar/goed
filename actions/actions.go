/*
Set of actions that can be dispatched.
Actions are dispatched, and processed one at a time by the action bus for
concurency safety.
*/
package actions

import "github.com/tcolar/goed/core"

func d(action core.Action) {
	core.Bus.Dispatch(action)
}

func CmdbarEnable(on bool) {
	d(cmdbarEnable{on: on})
}

func CmdbarToggle() {
	d(cmdbarToggle{})
}

// ########  Impl ......

type cmdbarEnable struct {
	on bool
}

func (a cmdbarEnable) Run() error {
	core.Ed.SetCmdOn(a.on)
	return nil
}

type cmdbarToggle struct{}

func (a cmdbarToggle) Run() error {
	core.Ed.CmdbarToggle()
	return nil
}
