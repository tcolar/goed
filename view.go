// History: Oct 02 14 tcolar Creation

package main

type View struct {
	Y1, Y2 int
	Title  string
	Dirty  bool
}

func (e *Editor) RenderView(c Col, v View) {
	e.FB(e.Theme.Viewbar.Fg, e.Theme.Viewbar.Bg)
	e.Fill(e.Theme.Viewbar.Rune, c.X1+1, v.Y1+1, c.X2-c.X1, 1)
	e.FB(e.Theme.ViewbarText, e.Theme.Viewbar.Bg)
	e.Str(c.X1+2, v.Y1+1, v.Title)
	e.RenderScroll(c, v)
	e.RenderIsDirty(c, v)
}

func (e *Editor) RenderScroll(c Col, v View) {
	e.FB(e.Theme.Scrollbar.Fg, e.Theme.Scrollbar.Bg)
	e.Fill(e.Theme.Scrollbar.Rune, c.X1, v.Y1+2, 1, v.Y2)
}

func (e *Editor) RenderIsDirty(c Col, v View) {
	style := e.Theme.FileClean
	if v.Dirty {
		style = e.Theme.FileDirty
	}
	e.FB(style.Fg, style.Bg)
	e.Char(c.X1, v.Y1+1, style.Rune)
}
