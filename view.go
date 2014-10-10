// History: Oct 02 14 tcolar Creation

package main

type View struct {
	Widget
	Id               int
	Title            string
	Dirty            bool
	Buffer           [][]rune
	CursorX, CursorY int
}

func (v *View) Render() {
	Ed.FB(Ed.Theme.Viewbar.Fg, Ed.Theme.Viewbar.Bg)
	Ed.Fill(Ed.Theme.Viewbar.Rune, v.x1+1, v.y1, v.x2, v.y1)
	fg := Ed.Theme.ViewbarText
	if v.Id == Ed.CurView.Id {
		fg = fg.WithAttr(Bold)
	}
	Ed.FB(fg, Ed.Theme.Viewbar.Bg)
	Ed.Str(v.x1+2, v.y1, v.Title)
	v.RenderScroll()
	v.RenderIsDirty()
	v.RenderText()
}

func (v *View) RenderScroll() {
	Ed.FB(Ed.Theme.Scrollbar.Fg, Ed.Theme.Scrollbar.Bg)
	Ed.Fill(Ed.Theme.Scrollbar.Rune, v.x1, v.y1+1, v.x1, v.y2)
}

func (v *View) RenderIsDirty() {
	style := Ed.Theme.FileClean
	if v.Dirty {
		style = Ed.Theme.FileDirty
	}
	Ed.FB(style.Fg, style.Bg)
	Ed.Char(v.x1, v.y1, style.Rune)
}

func (v *View) RenderText() {
	y := v.y1 + 2
	Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)
	for _, line := range v.Buffer {
		x := v.x1 + 2
		for _, rune := range line {
			if rune == '\t' {
				Ed.Char(x, y, ' ') // use a special background for leading/trailing spaces & tabs ?
				x++
			}
			Ed.Char(x, y, rune)
			x++
			if x >= v.x2-1 {
				Ed.FB(Ed.Theme.MoreText.Fg, Ed.Theme.MoreText.Bg)
				Ed.Char(x, y, Ed.Theme.MoreText.Rune)
				Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)
				break
			}
		}
		y++
		if y > v.y2-1 {
			break
		}
	}
}
