// History: Oct 02 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

const (
	Plain uint16 = iota + (1 << 8)
	Bold
	Underlined
)

type Col struct {
	X1, X2 int
	Views  []View
}

func (e *Editor) Render() {
	e.FB(e.Theme.Fg, e.Theme.Bg)
	termbox.Clear(termbox.Attribute(e.Bg.uint16), termbox.Attribute(e.Bg.uint16))

	for _, c := range e.Cols {
		e.RenderCol(c)
	}

	e.RenderMenu()
	e.RenderStatus()
}

func (e *Editor) RenderMenu() {
	w, _ := e.Size()
	e.FB(e.Theme.Menubar.Fg, e.Theme.Menubar.Bg)
	e.Fill(e.Theme.Menubar.Rune, 0, 0, w, 1)
	e.FB(e.Theme.MenubarText, e.Theme.Menubar.Bg)
	e.Str(0, 0, "GoEd 0.0.1")
}

func (e *Editor) RenderStatus() {
	w, h := e.Size()
	e.FB(e.Theme.Statusbar.Fg, e.Theme.Statusbar.Bg)
	e.Fill(e.Theme.Statusbar.Rune, 0, h-1, w, 1)
	e.FB(e.Theme.StatusbarText, e.Theme.Statusbar.Bg)
	e.Str(0, h-1, "All is good !")
	e.RenderPos()
}

func (e *Editor) RenderPos() {
	w, h := e.Size()
	e.FB(e.Theme.StatusbarText, e.Theme.Statusbar.Bg)
	pos := "123:59"
	e.Str(w-len(pos)-1, h-1, pos)
}

func (e *Editor) RenderCol(c Col) {
	for _, v := range c.Views {
		e.RenderView(c, v)
	}
}
