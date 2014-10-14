// History: Oct 02 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

const tabSize = 4

type View struct {
	Widget
	Id               int
	Title            string
	Dirty            bool
	Buffer           [][]rune
	CursorX, CursorY int
	offx, offy       int
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

// Gives us the lines to show in the view as text (\t as spaces)
func (v *View) viewLines() [][]rune {
	lines := [][]rune{}
	for i := v.offy; i < len(v.Buffer) && i <= v.offy+v.LastViewLine(); i++ {
		lines = append(lines, v.viewLine(i))
	}
	return lines
}

// single line to display (return only the relevant pasrt for display)
func (v *View) viewLine(index int) []rune {
	line := []rune{}
	if index >= len(v.Buffer) {
		return line
	}
	ln := v.Buffer[index]
	x := 0
	// Get what we can "See" in the viewport, plus an extra to the right so
	// we can know if we have it the end of the line or not.
	for i := 0; i < len(ln) && len(line) < v.LastViewCol()+1; i++ {
		if x >= v.offx {
			// we have skipped all that is to the left of the viewport
			if len(line) == 0 {
				// Special case for part of a tab showing at the beginning of the line
				for j := v.offx; j < x; j++ {
					line = append(line, ' ')
				}
			}
			if ln[i] == '\t' {
				for j := 1; j < tabSize; j++ {
					x += tabSize - 1
					line = append(line, ' ')
				}
			}
			line = append(line, ln[i])
		}
		x++
		if ln[i] == '\t' {
			x += tabSize - 1
		}
	}
	return line
}

func (v *View) RenderText() {
	y := v.y1 + 2
	Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)
	for _, l := range v.viewLines() {
		x := v.x1 + 2
		if v.offx > 0 {
			// More text to our left
			Ed.FB(Ed.Theme.MoreText.Fg, Ed.Theme.MoreText.Bg)
			Ed.Char(v.x1+1, y, Ed.Theme.MoreText.Rune)
			Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)
		}
		for _, c := range l {
			Ed.Char(x, y, c)
			x++
			if x > v.x2-1 {
				// More text to our right
				Ed.FB(Ed.Theme.MoreText.Fg, Ed.Theme.MoreText.Bg)
				Ed.Char(x-1, y, Ed.Theme.MoreText.Rune)
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

func (v *View) LastViewLine() int {
	return v.y2 - v.y1 - 3
}

func (v *View) LastViewCol() int {
	return v.x2 - v.x1 - 3
}

func (v *View) lineLn(lnIndex int) int {
	return v.lineLnTo(lnIndex, len(v.Buffer[lnIndex]))
}

func (v *View) lineLnTo(lnIndex, to int) int {
	if len(v.Buffer) <= lnIndex || len(v.Buffer[lnIndex]) < to {
		return 0
	}
	ln := len(v.Buffer[lnIndex][:to])
	for _, r := range v.Buffer[lnIndex][:to] {
		if r == '\t' {
			ln += tabSize - 1
		}
	}
	return ln
}

// MoveCursor : Move the cursor from it's current position by the x,y offsets
// This makes all the checks to make sure it's in a valid location
// also takes care to wrapping to previous/next line as needed
// as well as scrolling the view as needed.
func (v *View) MoveCursor(x, y int) {
	curCol := v.CursorX + v.offx
	curLine := v.CursorY + v.offy
	lastLine := len(v.Buffer) - 1

	if curLine+y < 0 || // before first line
		curLine+y > lastLine || // after last line
		(curLine <= 0 && curCol+x < 0) || // before beginning of file
		(curLine >= lastLine && curCol+x > len(v.Buffer[lastLine])) { // after eof
		return
	}

	ln := v.lineLn(curLine + y)

	if curCol+x < 0 {
		// wrap to after end of previous line
		y--
		x = v.lineLn(curLine+y) - curCol
	} else if curCol+x > ln {
		ln := v.lineLn(curLine + y)
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
	ln = v.lineLn(curLine + y)

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

	// TODO: Deal with cursor falling in the "middle" of a tab and adjust it
	tox := v.x1 + 2 + v.CursorX
	toy := v.y1 + 2 + v.CursorY

	termbox.SetCursor(tox, toy)
}

// Return the current line (zero indexed)
func (v *View) CurLine() int {
	return v.CursorX + v.offx
}

// Return the current column (zero indexed)
func (v *View) CurCol() int {
	return v.CursorY + v.offy
}
