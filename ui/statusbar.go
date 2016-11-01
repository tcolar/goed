package ui

import (
	"fmt"

	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui/widgets"
)

// Statusbar widget
type Statusbar struct {
	widgets.BaseWidget
	msg   string
	isErr bool
}

func (s *Statusbar) Render() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Statusbar.Fg, t.Statusbar.Bg)
	y1, x1, y2, x2 := s.Bounds()
	e.TermFill(t.Statusbar.Rune, y1, x1, y2, x2)
	if s.isErr {
		e.TermFB(t.StatusbarTextErr, t.Statusbar.Bg)
	} else {
		e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	}
	e.TermStr(y1, x1, s.msg)
	s.RenderPos()
}

func (s *Statusbar) RenderPos() {
	e := core.Ed
	t := e.Theme()
	y1, _, _, x2 := s.Bounds()
	e.TermFB(t.StatusbarText, t.Statusbar.Bg)
	vid := e.CurViewId()
	if vid < 0 {
		return
	}
	v := e.ViewById(vid)
	if v == nil || v.Backend() == nil {
		return
	}
	ln, col := v.CurLine(), v.LineRunesTo(v.Slice(), v.CurLine(), v.CurCol())
	pos := fmt.Sprintf(" %d:%d [%d]", ln+1, col+1, v.LineCount())
	e.TermStr(y1, x2-len(pos), pos)
}
