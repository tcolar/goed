// History: Oct 02 14 tcolar Creation

package main

import (
	"fmt"

	"github.com/tcolar/termbox-go"
)

const (
	Plain uint16 = iota + (1 << 8)
	Bold
	Underlined
)

func (e *Editor) WidgetAt(x, y int) Renderer {
	_, h := e.Size()
	if y == 1 {
		return e.Menubar
	}
	if y == h-1 {
		return e.Statusbar
	}
	for _, v := range e.Views {
		if x >= v.x1 && x <= v.x2 && y >= v.y1 && y <= v.y2 {
			return &v
		}
	}
	return nil
}

func (e *Editor) Render() {
	e.FB(e.Theme.Fg, e.Theme.Bg)
	termbox.Clear(termbox.Attribute(e.Bg.uint16), termbox.Attribute(e.Bg.uint16))

	for _, v := range e.Views {
		v.Render()
	}

	e.Menubar.Render()
	e.Statusbar.Render()

	termbox.Flush()
}

type Renderer interface {
	Bounds() (int, int, int, int)
	Render()
	SetBounds(x1, y1, x2, y2 int)
	Event(*termbox.Event)
}

// Widget implements the base of UI widgets
type Widget struct {
	x1, x2, y1, y2 int
}

func (w *Widget) Bounds() (int, int, int, int) {
	return w.x1, w.y1, w.x2, w.y2
}

func (w *Widget) SetBounds(x1, y1, x2, y2 int) {
	w.x1 = x1
	w.x2 = x2
	w.y1 = y1
	w.y2 = y2
}

// Menubar widget
type Menubar struct {
	Widget
}

func (m *Menubar) Render() {
	Ed.FB(Ed.Theme.Menubar.Fg, Ed.Theme.Menubar.Bg)
	Ed.Fill(Ed.Theme.Menubar.Rune, m.x1, m.y1, m.x2, m.y2)
	Ed.FB(Ed.Theme.MenubarText, Ed.Theme.Menubar.Bg)
	Ed.Str(m.x1, m.y1, "save saveall | cut copy paste | look | new del | newcol delcol | exit")

}

// Statusbar widget
type Statusbar struct {
	Widget
}

func (s *Statusbar) Render() {
	Ed.FB(Ed.Theme.Statusbar.Fg, Ed.Theme.Statusbar.Bg)
	Ed.Fill(Ed.Theme.Statusbar.Rune, s.x1, s.y1, s.x2, s.y2)
	Ed.FB(Ed.Theme.StatusbarText, Ed.Theme.Statusbar.Bg)
	Ed.Str(s.x1, s.y1, "All is good !")
	s.RenderPos()
}

func (s *Statusbar) RenderPos() {
	Ed.FB(Ed.Theme.StatusbarText, Ed.Theme.Statusbar.Bg)
	pos := fmt.Sprintf("%d:%d", Ed.CurView.CursorX, Ed.CurView.CursorY)
	Ed.Str(s.x2-len(pos)-1, s.y1, pos)
}
