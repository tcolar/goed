package main

import "fmt"

// Statusbar widget
type Statusbar struct {
	Widget
	msg   string
	isErr bool
}

func (s *Statusbar) Render() {
	Ed.FB(Ed.Theme.Statusbar.Fg, Ed.Theme.Statusbar.Bg)
	Ed.Fill(Ed.Theme.Statusbar.Rune, s.x1, s.y1, s.x2, s.y2)
	if s.isErr {
		Ed.FB(Ed.Theme.StatusbarTextErr, Ed.Theme.Statusbar.Bg)
	} else {
		Ed.FB(Ed.Theme.StatusbarText, Ed.Theme.Statusbar.Bg)
	}
	Ed.Str(s.x1, s.y1, s.msg)
	s.RenderPos()
}

func (s *Statusbar) RenderPos() {
	Ed.FB(Ed.Theme.StatusbarText, Ed.Theme.Statusbar.Bg)
	v := Ed.CurView
	if v == nil || v.backend == nil {
		return
	}
	col, ln := v.CursorTextPos(v.slice, v.CurCol(), v.CurLine())
	pos := fmt.Sprintf(" %d:%d [%d]", ln+1, col+1, v.LineCount())
	Ed.Str(s.x2-len(pos), s.y1, pos)
}
