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
	e.TermFill(t.Statusbar.Rune, s.y1, s.x1, s.y2, s.x2)
	if s.isErr {
		e.TermFB(t.StatusbarTextErr, t.Statusbar.Bg)
	} else {
		e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	}
	e.TermStr(s.y1, s.x1, s.msg)
	s.RenderPos()
}

func (s *Statusbar) RenderPos() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	vid := e.CurViewId()
	if vid < 0 {
		return
	}
	v := e.ViewById(vid).(*View)
	if v == nil || v.Backend() == nil {
		return
	}
	col, ln := v.LineRunesTo(v.Slice(), v.CurLine(), v.CurCol()), v.CurLine()
	pos := fmt.Sprintf(" %d:%d [%d]", ln+1, col+1, v.LineCount())
	e.TermStr(s.y1, s.x2-len(pos), pos)
}
