package actions

import "github.com/tcolar/goed/core"

func (a *ar) CmdbarBackspace() {
	d(cmdbarBackspace{})
}

func (a *ar) CmdbarClear() {
	d(cmdbarClear{})
}

func (a *ar) CmdbarDelete() {
	d(cmdbarDelete{})
}

func (a *ar) CmdbarInsert(s string) {
	d(cmdbarInsert{s: s})
}

func (a *ar) CmdbarNewLine() {
	d(cmdbarNewLine{})
}

func (a *ar) CmdbarCursorMvmt(m core.CursorMvmt) {
	d(cmdbarCursorMvmt{m: m})
}

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
type cmdbarBackspace struct{}

func (a cmdbarBackspace) Run() {
	core.Ed.Commandbar().Backspace()
}

type cmdbarClear struct{}

func (a cmdbarClear) Run() {
	core.Ed.Commandbar().Clear()
}

type cmdbarDelete struct{}

func (a cmdbarDelete) Run() {
	core.Ed.Commandbar().Delete()
}

type cmdbarInsert struct {
	s string
}

func (a cmdbarInsert) Run() {
	core.Ed.Commandbar().Insert(a.s)
}

type cmdbarCursorMvmt struct {
	m core.CursorMvmt
}

func (a cmdbarCursorMvmt) Run() {
	core.Ed.Commandbar().CursorMvmt(a.m)
}

type cmdbarNewLine struct{}

func (a cmdbarNewLine) Run() {
	core.Ed.Commandbar().NewLine()
}

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
