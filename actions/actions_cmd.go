package actions

import "github.com/tcolar/goed/core"

// Activate / desactivate the command bar
func (a *ar) CmdbarEnable(on bool) {
	d(cmdbarEnable{on: on})
}

// check if cmdbar is enabled
func (a *ar) CmdbarEnabled() bool {
	answer := make(chan bool, 1)
	d(cmdbarEnabled{answer: answer})
	return <-answer
}

// Toggle the command bar active or not
func (a *ar) CmdbarToggle() {
	d(cmdbarToggle{})
}

// ########  Impl ......

type cmdbarEnable struct {
	on bool
}

func (a cmdbarEnable) Run() {
	core.Ed.SetCmdOn(a.on)
}

type cmdbarEnabled struct {
	answer chan bool
}

func (a cmdbarEnabled) Run() {
	a.answer <- core.Ed.CmdOn()
}

type cmdbarToggle struct{}

func (a cmdbarToggle) Run() {
	core.Ed.CmdbarToggle()
}
