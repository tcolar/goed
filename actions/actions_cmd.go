package actions

import "github.com/tcolar/goed/core"

// Activate / desactivate the command bar
func (a *ar) CmdbarEnable(on bool) {
	d(cmdbarEnable{on: on})
}

// Toggle the command bar active or not
func (a *ar) CmdbarToggle() {
	d(cmdbarToggle{})
}

// ########  Impl ......

type cmdbarEnable struct {
	on    bool
	_help string
}

func (a cmdbarEnable) Run() {
	core.Ed.SetCmdOn(a.on)
}

type cmdbarToggle struct{}

func (a cmdbarToggle) Run() {
	core.Ed.CmdbarToggle()
}
