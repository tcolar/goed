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

func CmdbarEnableAction(on bool) {
	d(cmdbarEnableAction{on: on})
}

func CmdbarToggleAction() {
	d(cmdbarToggleAction{})
}

// ########  Impl ......

type cmdbarToggleAction struct{}

func (a cmdbarToggleAction) Run() error {
	core.Ed.CmdbarToggle()
	return nil
}

type cmdbarEnableAction struct {
	on bool
}

func (a cmdbarEnableAction) Run() error {
	core.Ed.SetCmdOn(a.on)
	return nil
}
