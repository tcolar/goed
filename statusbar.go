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
	pos := fmt.Sprintf("%d:%d [%d]", v.CurLine()+1, v.CurCol()+1, v.LineCount())
	Ed.Str(s.x2-len(pos)-1, s.y1, pos)
}
