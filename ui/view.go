package ui

import (
	"path/filepath"
	"time"

	"github.com/tcolar/goed/core"
)

const tabSize = 4

var id int = 0

type View struct {
	Widget
	id               int
	dirty            bool
	backend          core.Backend
	workDir          string
	CursorX, CursorY int // realtive position of cursor in view (0 index)
	offx, offy       int // absolute view offset (scrolled down/right) (0 index)
	HeightRatio      float64
	selections       []core.Selection
	title            string
	lastCloseTs      time.Time   // Timestamp of previous view close request
	slice            *core.Slice // curSlice
}

func (v *View) Reset() {
	v.CursorX, v.CursorY, v.offx, v.offy = 0, 0, 0, 0
	v.ClearSelections()
}

func (v *View) Render() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Viewbar.Fg, t.Viewbar.Bg)
	e.TermFill(t.Viewbar.Rune, v.x1+1, v.y1, v.x2, v.y1)
	fg := t.ViewbarText
	if v.Id() == e.CurView().Id() {
		fg = fg.WithAttr(core.Bold)
	}
	e.TermFB(fg, t.Viewbar.Bg)
	ti := v.Title()
	if v.x2-v.x1 > 4 && v.x1+len(ti) > v.x2-4 {
		ti = ti[:v.x2-v.x1-4]
	}
	e.TermStr(v.x1+2, v.y1, ti)
	v.RenderClose()
	v.RenderScroll()
	v.RenderIsDirty()
	v.RenderMargin()
	if v.backend != nil {
		v.RenderText()
	}
}

func (v *View) RenderMargin() {
	e := core.Ed
	t := e.Theme()
	if v.offx < 80 && v.offx+v.LastViewCol() >= 80 {
		for i := 0; i <= v.LastViewLine(); i++ {
			e.TermFB(t.Margin.Fg, t.Margin.Bg)
			e.TermChar(v.x1+2+80-v.offx, v.y1+2+i, t.Margin.Rune)
			e.TermFB(t.Fg, t.Bg)
		}
	}
}

func (v *View) RenderScroll() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Scrollbar.Fg, t.Scrollbar.Bg)
	e.TermFill(t.Scrollbar.Rune, v.x1, v.y1+1, v.x1, v.y2)
}

func (v *View) RenderIsDirty() {
	e := core.Ed
	t := e.Theme()
	style := t.FileClean
	if v.Dirty() {
		style = t.FileDirty
	}
	e.TermFB(style.Fg, style.Bg)
	e.TermChar(v.x1, v.y1, style.Rune)
}

func (v *View) RenderClose() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Close.Fg, t.Close.Bg)
	e.TermChar(v.x2-1, v.y1, t.Close.Rune)
}

func (v *View) RenderText() {
	e := core.Ed
	t := e.Theme()
	y := v.y1 + 2
	fg := t.Fg
	bg := t.Bg
	e.TermFB(fg, bg)
	inSelection := false
	tab := string(t.TabChar.Rune)
	for j := 1; j < tabSize; j++ {
		tab += " "
	}
	if v.offy > 0 {
		// More text above
		e.TermFB(t.MoreTextUp.Fg, t.MoreTextUp.Bg)
		e.TermChar(v.x1+1, y-1, t.MoreTextUp.Rune)
		e.TermFB(fg, bg)
	}
	// Note: using full lines
	v.slice = v.backend.Slice(v.offy+1, 1, v.offy+v.LastViewLine()+1, -1)
	for _, l := range *v.slice.Text() {
		x := v.x1 + 2
		if v.offx >= len(l) {
			y++
			continue
		}
		start := 0
		if v.offx > 0 {
			// More text to our left
			e.TermFB(t.MoreTextSide.Fg, t.MoreTextSide.Bg)
			e.TermChar(x-1, y, t.MoreTextSide.Rune)
			e.TermFB(fg, bg)
			// skip letters until we get to or past offx
			sz := 0
			for sz < v.offx {
				sz += v.runeSize(l[start])
				start++
			}
			// if we went "past" offx it means there where some
			// tabs leftover spaces taht we need to render
			for i := v.offx; i != sz; i++ {
				e.TermChar(x, y, ' ')
				x++
			}
		}
		for _, c := range l[start:] {
			sx, sy := v.CursorTextPos(v.slice, v.offx+x-2-v.x1, v.offy+y-2-v.y1)
			selected, _ := v.Selected(sx+1, sy+1)
			if selected != inSelection {
				inSelection = selected
				if selected {
					fg, bg = t.FgSelect, t.BgSelect
				} else {
					fg, bg = t.Fg, t.Bg
				}
				e.TermFB(fg, bg)
			}
			if c == '\t' {
				e.TermFB(t.TabChar.Fg, bg)
				e.TermStr(x, y, tab)
				x += tabSize - 1
				e.TermFB(fg, bg)
			} else {
				e.TermChar(x, y, c)
			}
			x++
			if x > v.x2-1 {
				// More text to our right
				e.TermFB(t.MoreTextSide.Fg, t.MoreTextSide.Bg)
				e.TermChar(x-1, y, t.MoreTextSide.Rune)
				e.TermFB(fg, bg)
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
		e.TermFB(t.MoreTextDown.Fg, t.MoreTextDown.Bg)
		e.TermChar(v.x1+1, y, t.MoreTextDown.Rune)
		e.TermFB(fg, bg)
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
	if !v.Dirty() {
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
	core.Ed.SetCursor(x, y)
}

func (v *View) Backend() core.Backend {
	return v.backend
}

func (v *View) Dirty() bool {
	return v.dirty
}

func (v *View) SetWorkDir(dir string) {
	v.workDir = dir
}

func (v *View) WorkDir() string {
	return v.workDir
}

func (v *View) SetTitle(title string) {
	v.title = title
}

func (v *View) SetDirty(dirty bool) {
	v.dirty = dirty
}

func (v *View) SetBackend(b core.Backend) {
	v.backend = b
}

func (v *View) Selections() *[]core.Selection {
	return &v.selections
}

func (v *View) Id() int {
	return v.id
}

func (v *View) Slice() *core.Slice {
	return v.slice
}
