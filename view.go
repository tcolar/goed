package main

import (
	"path/filepath"
	"time"
)

const tabSize = 4

var id int = 0

type View struct {
	Widget
	Id               int
	Dirty            bool
	backend          Backend
	WorkDir          string
	CursorX, CursorY int // realtive position of cursor in view (0 index)
	offx, offy       int // absolute view offset (scrolled down/right) (0 index)
	HeightRatio      float64
	Selections       []Selection
	title            string
	lastCloseTs      time.Time // Timestamp of previous view close request
	slice            *Slice    // curSlice
}

func (e *Editor) NewView() *View {
	id++
	d, _ := filepath.Abs(".")
	v := &View{
		Id:          id,
		HeightRatio: 0.5,
		WorkDir:     d,
		slice:       &Slice{text: [][]rune{}},
	}
	v.backend, _ = e.NewFileBackend("", v.Id)
	return v
}

func (e *Editor) NewFileView(path string) *View {
	v := e.NewView()
	e.Open(path, v, "")
	return v
}

func (v *View) Reset() {
	v.CursorX, v.CursorY, v.offx, v.offy = 0, 0, 0, 0
	v.Selections = []Selection{}
}

func (v *View) Render() {
	Ed.TermFB(Ed.Theme.Viewbar.Fg, Ed.Theme.Viewbar.Bg)
	Ed.TermFill(Ed.Theme.Viewbar.Rune, v.x1+1, v.y1, v.x2, v.y1)
	fg := Ed.Theme.ViewbarText
	if v.Id == Ed.CurView.Id {
		fg = fg.WithAttr(Bold)
	}
	Ed.TermFB(fg, Ed.Theme.Viewbar.Bg)
	t := v.Title()
	if v.x1+len(t) > v.x2-4 {
		t = t[:v.x2-v.x1-4]
	}
	Ed.TermStr(v.x1+2, v.y1, t)
	v.RenderClose()
	v.RenderScroll()
	v.RenderIsDirty()
	v.RenderMargin()
	if v.backend != nil {
		v.RenderText()
	}
}

func (v *View) RenderMargin() {
	if v.offx < 80 && v.offx+v.LastViewCol() >= 80 {
		for i := 0; i <= v.LastViewLine(); i++ {
			Ed.TermFB(Ed.Theme.Margin.Fg, Ed.Theme.Margin.Bg)
			Ed.TermChar(v.x1+2+80-v.offx, v.y1+2+i, Ed.Theme.Margin.Rune)
			Ed.TermFB(Ed.Theme.Fg, Ed.Theme.Bg)
		}
	}
}

func (v *View) RenderScroll() {
	Ed.TermFB(Ed.Theme.Scrollbar.Fg, Ed.Theme.Scrollbar.Bg)
	Ed.TermFill(Ed.Theme.Scrollbar.Rune, v.x1, v.y1+1, v.x1, v.y2)
}

func (v *View) RenderIsDirty() {
	style := Ed.Theme.FileClean
	if v.Dirty {
		style = Ed.Theme.FileDirty
	}
	Ed.TermFB(style.Fg, style.Bg)
	Ed.TermChar(v.x1, v.y1, style.Rune)
}

func (v *View) RenderClose() {
	Ed.TermFB(Ed.Theme.Close.Fg, Ed.Theme.Close.Bg)
	Ed.TermChar(v.x2-1, v.y1, Ed.Theme.Close.Rune)
}

func (v *View) RenderText() {
	y := v.y1 + 2
	fg := Ed.Theme.Fg
	bg := Ed.Theme.Bg
	Ed.TermFB(fg, bg)
	inSelection := false
	tab := string(Ed.Theme.TabChar.Rune)
	for j := 1; j < tabSize; j++ {
		tab += " "
	}
	if v.offy > 0 {
		// More text above
		Ed.TermFB(Ed.Theme.MoreTextUp.Fg, Ed.Theme.MoreTextUp.Bg)
		Ed.TermChar(v.x1+1, y-1, Ed.Theme.MoreTextUp.Rune)
		Ed.TermFB(fg, bg)
	}
	// Note: using full lines
	v.slice = v.backend.Slice(v.offy+1, 1, v.offy+v.LastViewLine()+1, -1)
	for _, l := range v.slice.text {
		x := v.x1 + 2
		if v.offx >= len(l) {
			y++
			continue
		}
		start := 0
		if v.offx > 0 {
			// More text to our left
			Ed.TermFB(Ed.Theme.MoreTextSide.Fg, Ed.Theme.MoreTextSide.Bg)
			Ed.TermChar(x-1, y, Ed.Theme.MoreTextSide.Rune)
			Ed.TermFB(fg, bg)
			// skip letters until we get to or past offx
			sz := 0
			for sz < v.offx {
				sz += v.runeSize(l[start])
				start++
			}
			// if we went "past" offx it means there where some
			// tabs leftover spaces taht we need to render
			for i := v.offx; i != sz; i++ {
				Ed.TermChar(x, y, ' ')
				x++
			}
		}
		for _, c := range l[start:] {
			sx, sy := v.CursorTextPos(v.slice, v.offx+x-2-v.x1, v.offy+y-2-v.y1)
			selected, _ := v.Selected(sx+1, sy+1)
			if selected != inSelection {
				inSelection = selected
				if selected {
					fg, bg = Ed.Theme.FgSelect, Ed.Theme.BgSelect
				} else {
					fg, bg = Ed.Theme.Fg, Ed.Theme.Bg
				}
				Ed.TermFB(fg, bg)
			}
			if c == '\t' {
				Ed.TermFB(Ed.Theme.TabChar.Fg, bg)
				Ed.TermStr(x, y, tab)
				x += tabSize - 1
				Ed.TermFB(fg, bg)
			} else {
				Ed.TermChar(x, y, c)
			}
			x++
			if x > v.x2-1 {
				// More text to our right
				Ed.TermFB(Ed.Theme.MoreTextSide.Fg, Ed.Theme.MoreTextSide.Bg)
				Ed.TermChar(x-1, y, Ed.Theme.MoreTextSide.Rune)
				Ed.TermFB(fg, bg)
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
		Ed.TermFB(Ed.Theme.MoreTextDown.Fg, Ed.Theme.MoreTextDown.Bg)
		Ed.TermChar(v.x1+1, y, Ed.Theme.MoreTextDown.Rune)
		Ed.TermFB(fg, bg)
	}
}

func (v *View) LastViewLine() int {
	return v.y2 - v.y1 - 3
}

func (v *View) LastViewCol() int {
	return v.x2 - v.x1 - 3
}

// Same as MoveCursor but with "rolling" to next/prev line if overflowed.
func (v *View) MoveCursorRoll(x, y int) {
	slice := v.slice
	curCol := v.CurCol()
	curLine := v.CurLine()
	lastLine := v.LineCount() - 1
	ln := v.lineCols(slice, curLine+y)

	if curCol+x < 0 {
		// wrap to after end of previous line
		y--
		x = v.lineCols(slice, curLine+y) - curCol
	} else if curCol+x > ln {
		ln = v.lineCols(slice, curLine+y)
		if y == 0 && curLine+y < lastLine {
			// moved (right) passed eol, wrap to beginning of next line
			x = -curCol
			y++
		} else {
			// when movin up/down, don't go passed eol
			x = ln - curCol
		}
	}
	v.MoveCursor(x, y)
}

// MoveCursor : Move the cursor from it's current position by the x,y offsets (**in runes**)
// This makes all the checks to make sure it's in a valid location
// as well as scrolling the view as needed.
func (v *View) MoveCursor(x, y int) {

	slice := v.slice

	curCol := v.CurCol()
	curLine := v.CurLine()
	lastLine := v.LineCount() - 1

	// check for overflows
	if curLine+y < 0 {
		y = -curLine
	} else if curLine+y > lastLine {
		y = lastLine - curLine
	}
	if curCol+x < 0 {
		x = -curCol
	}
	ln := v.lineCols(slice, curLine+y)
	if curCol+x > ln {
		x = ln - curCol // put at EOL
	}

	v.CursorX += x
	v.CursorY += y

	// Special handling for tabs
	c, textX, textY := v.CurChar()
	if c != nil && *c == '\t' {
		from := v.CursorX
		// align cursor with beginning of tab
		v.CursorX = v.lineColsTo(slice, textY, textX) - v.offx
		x -= v.CursorX - from
	}

	// No scrolling needed
	if curCol+x >= v.offx && curCol+x <= v.offx+v.LastViewCol() &&
		curLine+y >= v.offy && curLine+y <= v.offy+v.LastViewLine() {
		v.setCursor(v.x1+2+v.CursorX, v.y1+2+v.CursorY)
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

	v.setCursor(tox, toy)
}

func (v *View) Title() string {
	if len(v.title) != 0 {
		return v.title
	}
	if v.backend == nil || len(v.backend.SrcLoc()) == 0 {
		v.title = "~~ NEW ~~"
		return v.title
	}
	v.title = filepath.Base(v.backend.SrcLoc())
	return v.title
}

// Return the current line (0 indexed)
func (v *View) CurLine() int {
	return v.CursorY + v.offy
}

// Return the current column (0 indexed)
func (v *View) CurCol() int {
	return v.CursorX + v.offx
}

// canClose checks if the view can be closed
// that is true if the view is not dirty
// otherwise, if dirty, returns true if we get 2 lose request in a short timespan
func (v *View) canClose() bool {
	if !v.Dirty {
		return true
	}
	if v.lastCloseTs.IsZero() || time.Now().Sub(v.lastCloseTs) > 10*time.Second {
		v.lastCloseTs = time.Now()
		return false
	}
	// 2 "quick" close request in a row
	return true
}

func (v *View) setCursor(x, y int) {
	Ed.term.SetCursor(x, y)
}
