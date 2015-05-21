package ui

import (
	"fmt"

	"github.com/tcolar/goed/core"
)

// Statusbar widget
type Statusbar struct {
	Widget
	msg   string
	isErr bool
}

func (s *Statusbar) Render() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Statusbar.Fg, t.Statusbar.Bg)
	e.TermFill(t.Statusbar.Rune, s.x1, s.y1, s.x2, s.y2)
	if s.isErr {
		e.TermFB(t.StatusbarTextErr, t.Statusbar.Bg)
	} else {
		e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	}
	e.TermStr(s.x1, s.y1, s.msg)
	s.RenderPos()
}

func (s *Statusbar) RenderPos() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	v := e.CurView()
	if v == nil {
		return
	}
	col, ln := v.CursorTextPos(v.Slice(), v.CurCol(), v.CurLine())
	pos := fmt.Sprintf(" %d:%d [%d]", ln+1, col+1, v.LineCount())
	e.TermStr(s.x2-len(pos), s.y1, pos)
}
