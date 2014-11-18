// History: Oct 02 14 tcolar Creation

package main

import (
	"path/filepath"

	"github.com/tcolar/termbox-go"
)

const tabSize = 4

var id int = 0

type View struct {
	Widget
	Id               int
	Dirty            bool
	Buffer           *Buffer
	CursorX, CursorY int
	offx, offy       int
	HeightRatio      float64
	Selections       []Selection
}

func (e *Editor) NewView() *View {
	id++
	return &View{
		Id:          id,
		Buffer:      &Buffer{},
		HeightRatio: 0.5,
	}
}

func (e *Editor) NewFileView(path string) *View {
	v := e.NewView()
	Ed.OpenFile(path, v)
	return v
}

func (v *View) Render() {
	Ed.FB(Ed.Theme.Viewbar.Fg, Ed.Theme.Viewbar.Bg)
	Ed.Fill(Ed.Theme.Viewbar.Rune, v.x1+1, v.y1, v.x2, v.y1)
	fg := Ed.Theme.ViewbarText
	if v.Id == Ed.CurView.Id {
		fg = fg.WithAttr(Bold)
	}
	Ed.FB(fg, Ed.Theme.Viewbar.Bg)
	t := v.Title()
	if v.x1+len(t) > v.x2-2 {
		t = t[:v.x2-v.x1-2]
	}
	Ed.Str(v.x1+2, v.y1, t)
	v.RenderScroll()
	v.RenderIsDirty()
	v.RenderMargin()
	if v.Buffer != nil {
		v.RenderText()
	}
}

func (v *View) RenderMargin() {
	if v.offx < 80 && v.offx+v.LastViewCol() >= 80 {
		for i := 0; i <= v.LastViewLine(); i++ {
			Ed.FB(Ed.Theme.Margin.Fg, Ed.Theme.Margin.Bg)
			Ed.Char(v.x1+2+80-v.offx, v.y1+2+i, Ed.Theme.Margin.Rune)
			Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)
		}
	}
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

// Gives us the lines to show in the view as text (\t as spaces)
func (v *View) viewLines() [][]rune {
	lines := [][]rune{}
	for i := v.offy; i < v.LineCount() && i <= v.offy+v.LastViewLine(); i++ {
		lines = append(lines, v.viewLine(i))
	}
	return lines
}

// single line to display (return only the relevant parts for display)
func (v *View) viewLine(index int) []rune {
	line := []rune{}
	if index >= v.LineCount() {
		return line
	}
	ln := v.Line(index)
	x := 0
	// Get what we can "See" in the viewport, plus an extra to the right so
	// we can know if we have it the end of the line or not.
	for i := 0; i < len(ln) && len(line) < v.LastViewCol()+1; i++ {
		c := ln[i]
		if x >= v.offx {
			// we have skipped all that is to the left of the viewport
			if len(line) == 0 {
				// Special case for part of a tab showing at the beginning of the line
				for j := v.offx; j < x; j++ {
					line = append(line, ' ')
				}
			}
			line = append(line, c)
			if c == '\t' {
				for j := 1; j < tabSize; j++ {
					x += tabSize - 1
					line = append(line, ' ')
				}
			}
		}
		x++
		if c == '\t' {
			x += tabSize - 1
		}
	}
	return line
}

func (v *View) RenderText() {
	y := v.y1 + 2
	fg := Ed.Theme.Fg
	bg := Ed.Theme.Bg
	Ed.FB(fg, bg)
	inSelection := false
	if v.offy > 0 {
		// More text above
		Ed.FB(Ed.Theme.MoreTextUp.Fg, Ed.Theme.MoreTextUp.Bg)
		Ed.Char(v.x1+1, y-1, Ed.Theme.MoreTextUp.Rune)
		Ed.FB(fg, bg)
	}
	for _, l := range v.viewLines() {
		x := v.x1 + 2
		if v.offx > 0 {
			// More text to our left
			Ed.FB(Ed.Theme.MoreTextSide.Fg, Ed.Theme.MoreTextSide.Bg)
			Ed.Char(x-1, y, Ed.Theme.MoreTextSide.Rune)
			Ed.FB(fg, bg)
		}
		for _, c := range l {
			selected, _ := v.Selected(v.CursorTextPos(v.offx+x-2-v.x1, v.offy+y-2-v.y1))
			if selected != inSelection {
				inSelection = selected
				if selected {
					fg, bg = Ed.Theme.FgSelect, Ed.Theme.BgSelect
				} else {
					fg, bg = Ed.Theme.Fg, Ed.Theme.Bg
				}
				Ed.FB(fg, bg)
			}
			if c == '\t' {
				Ed.FB(Ed.Theme.TabChar.Fg, bg)
				Ed.Char(x, y, Ed.Theme.TabChar.Rune)
				Ed.FB(fg, bg)
			} else {
				Ed.Char(x, y, c)
			}
			x++
			if x > v.x2-1 {
				// More text to our right
				Ed.FB(Ed.Theme.MoreTextSide.Fg, Ed.Theme.MoreTextSide.Bg)
				Ed.Char(x-1, y, Ed.Theme.MoreTextSide.Rune)
				Ed.FB(fg, bg)
				break
			}
		}
		y++
		if y > v.y2-1 {
			break
		}
	}
	if v.offy+v.LastViewLine() < v.LineCount()-1 {
		// More text below
		Ed.FB(Ed.Theme.MoreTextDown.Fg, Ed.Theme.MoreTextDown.Bg)
		Ed.Char(v.x1+1, y, Ed.Theme.MoreTextDown.Rune)
		Ed.FB(fg, bg)
	}
}

func (v *View) LastViewLine() int {
	return v.y2 - v.y1 - 3
}

func (v *View) LastViewCol() int {
	return v.x2 - v.x1 - 3
}

// MoveCursor : Move the cursor from it's current position by the x,y offsets (**in runes**)
// This makes all the checks to make sure it's in a valid location
// also takes care to wrapping to previous/next line as needed
// as well as scrolling the view as needed.
func (v *View) MoveCursor(x, y int) {
	curCol := v.CurCol()
	curLine := v.CurLine()
	lastLine := v.LineCount() - 1

	if curLine+y < 0 || // before first line
		curLine+y > lastLine || // after last line
		(curLine <= 0 && curCol+x < 0) || // before beginning of file
		(curLine >= lastLine && curCol+x > v.lineCols(lastLine)) { // after eof
		return
	}

	ln := v.lineCols(curLine + y)

	if curCol+x < 0 {
		// wrap to after end of previous line
		y--
		x = v.lineCols(curLine+y) - curCol
	} else if curCol+x > ln {
		ln := v.lineCols(curLine + y)
		if y == 0 {
			// moved (right) passed eol, wrap to beginning of next line
			x = -curCol
			y++
		} else {
			// when movin up/down, don't go passed eol
			x = ln - curCol
		}
	}

	v.CursorX += x
	v.CursorY += y
	ln = v.LineLen(curLine + y)

	// Special handling for tabs
	c, textX, textY := v.CurChar()
	if c != nil && *c == '\t' {
		from := v.CursorX
		// align cursor with beginning of tab
		v.CursorX = v.lineColsTo(textY, textX) - v.offx
		x -= v.CursorX - from
	}

	// No scrolling needed
	if curCol+x >= v.offx && curCol+x <= v.offx+v.LastViewCol() &&
		curLine+y >= v.offy && curLine+y <= v.offy+v.LastViewLine() {
		termbox.SetCursor(v.x1+2+v.CursorX, v.y1+2+v.CursorY)
		return
	}

	// scrolling needed
	if curCol+x < v.offx {
		v.offx = curCol + x
		v.CursorX = 0
	} else if curCol+x >= v.offx+v.LastViewCol() {
		v.offx = curCol + x - v.LastViewCol()
		v.CursorX = v.LastViewCol()
	}
	if curLine+y < v.offy && curLine+y >= 0 {
		v.offy = curLine + y
		v.CursorY = 0
	} else if curLine+y > v.offy+v.LastViewLine() {
		v.offy = curLine + y - v.LastViewLine()
		v.CursorY = v.LastViewLine()
	}

	tox := v.x1 + 2 + v.CursorX
	toy := v.y1 + 2 + v.CursorY

	termbox.SetCursor(tox, toy)
}

func (v *View) Title() string {
	if len(v.Buffer.file) == 0 {
		return "~~ NEW ~~"
	}
	return filepath.Base(v.Buffer.file)
}

// Return the current line (zero indexed)
func (v *View) CurLine() int {
	return v.CursorY + v.offy
}

// Return the current column (zero indexed)
func (v *View) CurCol() int {
	return v.CursorX + v.offx
}
