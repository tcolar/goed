package ui

import (
	"log"
	"path"
	"path/filepath"
	"time"

	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

const tabSize = 4

var _ core.Viewable = (*View)(nil)

// View represents an individual view pane(file) in the editor.
type View struct {
	Widget
	id      int64
	dirty   bool
	backend core.Backend
	workDir string
	// relative position of cursor in view (0 index)
	CursorX, CursorY int
	// absolute view offset (scrolled down/right) (0 index)
	offx, offy       int
	HeightRatio      float64
	selections       []core.Selection
	title            string
	lastCloseTs      time.Time   // Timestamp of previous view close request
	slice            *core.Slice // curSlice
	autoScrollX      int
	autoScrollY      int
	autoScrollSelect bool
	viewType         core.ViewType
	highlighter      core.Highlighter
}

func (e *Editor) NewView(loc string) *View {
	d, _ := filepath.Abs(".")
	if len(loc) > 0 {
		d = path.Dir(loc)
	}
	v := &View{
		id:          e.genViewId(),
		HeightRatio: 0.5,
		workDir:     d,
		slice:       core.NewSlice(0, 0, 0, 0, [][]rune{}),
		highlighter: &CodeHighlighter{},
	}
	e.views[v.id] = v
	v.backend, _ = backend.NewMemBackend(loc, v.Id())
	return v
}

// NewFileView creates a view for a given file
func (e *Editor) NewFileView(loc string) *View {
	v := e.NewView(loc)
	e.Open(loc, v.Id(), "", true)
	return v
}

func (e *Editor) genViewId() int64 {
	return time.Now().UnixNano()
}

func (v *View) Reset() {
	v.CursorX, v.CursorY, v.offx, v.offy = 0, 0, 0, 0
	v.ClearSelections()
}

func (v *View) Render() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Viewbar.Fg, t.Viewbar.Bg)
	e.TermFill(t.Viewbar.Rune, v.y1, v.x1+1, v.y1, v.x2)
	fg := t.ViewbarText
	if v.Id() == e.CurViewId() {
		fg = fg.WithAttr(core.Bold)
	}
	e.TermFB(fg, t.Viewbar.Bg)
	ti := v.Title()
	if v.x2-v.x1 > 4 && v.x1+len(ti) > v.x2-4 {
		ti = ti[:v.x2-v.x1-4]
	}
	e.TermStr(v.y1, v.x1+2, ti)
	v.renderClose()
	v.renderScroll()
	v.renderIsDirty()
	v.renderMargin()
	if v.backend != nil {
		v.renderText()
	}
}

func (v *View) renderMargin() {
	e := core.Ed
	t := e.Theme()
	margin := e.Config().LineWidthIndicator
	if v.offx < margin && v.offx+v.LastViewCol() >= margin {
		for i := 0; i <= v.LastViewLine(); i++ {
			e.TermFB(t.Margin.Fg, t.Margin.Bg)
			e.TermChar(v.y1+2+i, v.x1+2+margin-v.offx, t.Margin.Rune)
			e.TermFB(t.Fg, t.Bg)
		}
	}
}

func (v *View) renderScroll() {
	e := core.Ed
	t := e.Theme()
	viewLines := v.y2 - v.y1 - 1
	textLines := v.LineCount()
	topLine := v.slice.R1
	e.TermFB(t.Scrollbar.Fg, t.Scrollbar.Bg)
	e.TermFill(t.Scrollbar.Rune, v.y1+1, v.x1, v.y2, v.x1)
	if textLines < viewLines || viewLines <= 0 {
		return // no scrollbar needed
	}
	size := int(float64(viewLines) * (float64(viewLines) / float64(textLines)))
	loc := int(float64(viewLines) * (float64(topLine) / float64(textLines)))
	if size < 2 {
		size = 2 // minimum scrollbar handle size
	}
	if topLine > 0 && loc == 0 {
		loc = 1 // if we are past first page make sure the scroll bar shows a bit of scrolling
	} else if topLine > textLines-viewLines {
		loc = viewLines - size // if on last page, make sure scrolbar hugs bottom
	}
	e.TermFB(t.ScrollTab.Fg, t.ScrollTab.Bg)
	e.TermFill(t.ScrollTab.Rune, v.y1+1+loc, v.x1, v.y1+1+loc+size, v.x1)
}

func (v *View) renderIsDirty() {
	e := core.Ed
	t := e.Theme()
	style := t.FileClean
	if v.Dirty() {
		style = t.FileDirty
	}
	e.TermFB(style.Fg, style.Bg)
	e.TermChar(v.y1, v.x1, style.Rune)
}

func (v *View) renderClose() {
	e := core.Ed
	t := e.Theme()
	e.TermFB(t.Close.Fg, t.Close.Bg)
	e.TermChar(v.y1, v.x2-1, t.Close.Rune)
}

func (v *View) renderText() {
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
		e.TermChar(y-1, v.x1+1, t.MoreTextUp.Rune)
		e.TermFB(fg, bg)
	}
	// Note: using full lines
	v.slice = v.backend.Slice(v.offy, 0, v.offy+v.LastViewLine(), -1)
	if e.Config().SyntaxHighlighting {
		v.highlighter.UpdateHighlights(v)
	}
	for lnc, l := range *v.slice.Text() {
		x := v.x1 + 2
		if v.offx >= len(l) {
			y++
			continue
		}
		start := 0
		if v.offx > 0 {
			// More text to our left
			e.TermFB(t.MoreTextSide.Fg, t.MoreTextSide.Bg)
			e.TermChar(y, x-1, t.MoreTextSide.Rune)
			e.TermFB(fg, bg)
			// skip letters until we get to or past offx
			sz := 0
			for sz < v.offx {
				sz += v.runeSize(l[start])
				start++
			}
			// if we went "past" offx it means there where some
			// tabs leftover spaces that we need to render
			for i := v.offx; i != sz; i++ {
				e.TermChar(y, x, ' ')
				x++
			}
		}
		for colc, c := range l[start:] {
			sy := v.offy + y - 2 - v.y1
			sx := v.offx + x - 2 - v.x1
			sx = v.LineRunesTo(v.slice, sy, sx)
			selected, _ := v.Selected(sx, sy)
			if selected != inSelection {
				inSelection = selected
				if selected {
					fg, bg = t.FgSelect, t.BgSelect
				} else {
					fg, bg = t.Fg, t.Bg
				}
				e.TermFB(fg, bg)
			}
			if c == '\t' { // tab
				e.TermFB(t.TabChar.Fg, bg)
				e.TermStr(y, x, tab)
				e.TermFB(fg, bg)
			} else if c < 32 { // other unprintable control char
				e.TermFB(t.TabChar.Fg, bg)
				e.TermChar(y, x, 0x1A) // ASCII substitute char (invisible)
				e.TermFB(fg, bg)
			} else { // normal char
				if e.Config().SyntaxHighlighting && !inSelection {
					v.highlighter.ApplyHighlight(v, v.offy, lnc, start+colc)
				}
				e.TermChar(y, x, c)
			}
			x += v.runeSize(c)
			if x > v.x2-1 {
				// More text to our right
				e.TermFB(t.MoreTextSide.Fg, t.MoreTextSide.Bg)
				e.TermChar(y, x-1, t.MoreTextSide.Rune)
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
		e.TermChar(y, v.x1+1, t.MoreTextDown.Rune)
		e.TermFB(fg, t.Bg)
	}
}

// LastViewLines returns the last Line of this view (~ number of visible lines)
func (v *View) LastViewLine() int {
	return v.y2 - v.y1 - 3
}

// LastViewCol returns the last column of this view (~ number of visible columns)
func (v *View) LastViewCol() int {
	return v.x2 - v.x1 - 3
}

// Same as MoveCursor but with "rolling" to next/prev line if overflowed.
func (v *View) MoveCursorRoll(y, x int) {
	slice := v.slice
	curLine := v.CurLine()
	lastLine := v.LineCount() - 1

	if v.CurCol()+x < 0 {
		// wrap to after end of previous line
		y--
		v.MoveCursor(y, 0)
		v.MoveCursor(0, v.lineCols(slice, curLine+y)-v.CurCol())
		return
	} else if v.CurCol()+x > v.lineCols(slice, curLine+y) {
		if y == 0 && curLine+y < lastLine {
			// moved (right) passed eol, wrap to beginning of next line
			y++
			v.MoveCursor(y, 0)
			v.MoveCursor(0, -v.CurCol())
			return
		} else {
			// when movin up/down, don't go passed eol
			x = v.lineCols(slice, curLine+y) - v.CurCol()
		}
	}
	v.MoveCursor(y, x)
}

func (v *View) SyncSlice() {
	v.slice = v.backend.Slice(v.offy, 0, v.offy+v.LastViewLine(), -1)
}

// MoveCursor : Move the cursor from it's current position by the y, x offsets (**in runes**)
func (v *View) MoveCursor(y, x int) {
	ln, col := v.CurTextPos()
	v.SetCursorPos(ln+y, col+x)
}

// SetCursor : Set the cursor text position
// This makes all the checks to make sure it's in a valid location,
// as well as scrolling the view as needed.
func (v *View) SetCursorPos(y, x int) {
	lastLine := v.LineCount()
	ln := y
	if ln < 0 {
		ln = 0
	} else if ln >= lastLine {
		ln = lastLine - 1
	}

	// slice for the area we will be in after scrolling
	slice := v.slice
	if !slice.ContainsLine(ln) {
		slice = v.backend.Slice(ln, 0, ln, -1)
	}

	col := v.lineColsTo(slice, ln, x)
	// check for col overflow
	cols := v.lineCols(slice, ln)
	if col < 0 {
		col = 0
	} else if col > cols {
		col = cols // put at EOL
	}

	// scroll vertically if needed
	if ln < v.offy && ln >= 0 {
		v.offy = ln
	} else if ln > v.offy+v.LastViewLine() {
		v.offy = ln - v.LastViewLine()
		if v.offy < 0 {
			v.offy = 0
		} else if v.offy > lastLine {
			v.offy = lastLine
		}
	}

	// scroll horizontally if needed
	if col < v.offx && col >= 0 {
		v.offx = col
	} else if col >= v.offx+v.LastViewCol() {
		v.offx = col - v.LastViewCol()
		if v.offx < 0 {
			v.offx = 0
		}
	}

	v.CursorY = ln - v.offy
	v.CursorX = col - v.offx
	v.updateCursor(slice)
}

// Update the editor cursor to be this view current cursor
func (v *View) updateCursor(slice *core.Slice) {
	v.NormalizeCursor(slice)
	if !slice.ContainsLine(v.CursorY) {
		v.SyncSlice()
	}
	core.Ed.SetCursor(v.y1+2+v.CursorY, v.x1+2+v.CursorX)
}

func (v *View) NormalizeCursor(slice *core.Slice) {
	lastLine := v.LineCount()
	if v.CursorY < 0 {
		v.CursorY = 0
		v.CursorX = 0
		return
	}
	ln := v.offy + v.CursorY
	if ln > lastLine {
		v.CursorY = lastLine - v.offy
		v.CursorX = v.lineCols(slice, v.offy+v.CursorY)
		return
	}
	lc := v.lineCols(slice, ln)
	if v.CursorX < 0 {
		v.CursorX = 0
		return
	}
	if v.offx+v.CursorX > lc+1 {
		v.CursorX = lc - v.offx
		return
	}
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

// Return the current UI line (0 indexed)
func (v *View) CurLine() int {
	return v.CursorY + v.offy
}

// Return the current UI column (0 indexed)
func (v *View) CurCol() int {
	return v.CursorX + v.offx
}

// Return the postion of the cursor in the view's text.
func (v *View) CurTextPos() (ln int, col int) {
	return v.CurLine(), v.LineRunesTo(v.Slice(), v.CurLine(), v.CurCol())
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

func (v *View) Backend() core.Backend {
	return v.backend
}

func (v *View) Dirty() bool {
	if v.viewType == core.ViewTypeShell {
		return false
	}
	return v.dirty
}

func (v *View) SetWorkDir(dir string) {
	v.workDir = dir
	log.Printf("View %d : workdir: %s", v.id, v.workDir)
}

func (v *View) WorkDir() string {
	return v.workDir
}

func (v *View) SetTitle(title string) {
	v.title = title
}

func (v *View) SetDirty(dirty bool) {
	if v.Type() == core.ViewTypeStandard {
		v.dirty = dirty
	}
}

func (v *View) SetBackend(b core.Backend) {
	v.backend = b
}

func (v *View) Selections() *[]core.Selection {
	return &v.selections
}

func (v *View) Id() int64 {
	if v == nil {
		return -1
	}
	return v.id
}

func (v *View) Slice() *core.Slice {
	return v.slice
}

func (v *View) SetAutoScroll(y, x int, isSelect bool) {
	v.autoScrollX, v.autoScrollY = x, y
	v.autoScrollSelect = isSelect
}

func (v *View) SetViewType(t core.ViewType) {
	v.viewType = t
}

func (v *View) CursorMvmt(mvmt core.CursorMvmt) {
	ln, col := v.CurLine(), v.CurCol()
	switch mvmt {
	case core.CursorMvmtRight:
		v.MoveCursorRoll(0, 1)
	case core.CursorMvmtLeft:
		v.MoveCursorRoll(0, -1)
	case core.CursorMvmtUp:
		v.MoveCursor(-1, 0)
	case core.CursorMvmtDown:
		v.MoveCursor(1, 0)
	case core.CursorMvmtPgDown:
		dist := v.LastViewLine() + 1
		if v.LineCount()-ln < dist {
			dist = v.LineCount() - ln - 1
		}
		v.MoveCursor(dist, 0)
	case core.CursorMvmtPgUp:
		dist := v.LastViewLine() + 1
		if dist > ln {
			dist = ln
		}
		v.MoveCursor(-dist, 0)
	case core.CursorMvmtEnd:
		v.MoveCursor(0, v.lineCols(v.slice, ln)-col)
	case core.CursorMvmtHome:
		v.MoveCursor(0, -col)
	case core.CursorMvmtTop:
		v.MoveCursor(-v.CurLine(), -col)
	case core.CursorMvmtBottom:
		c := 0
		if v.LineCount() > 0 {
			slice := v.backend.Slice(v.LineCount()-1, 1, v.LineCount()-1, -1)
			c = v.lineCols(slice, v.LineCount()-1) + 1 - col
		}
		v.MoveCursor(v.LineCount()-1-v.CurLine(), c)
	case core.CursorMvmtScrollDown:
		v.MoveCursor(7, 0)
	case core.CursorMvmtScrollUp:
		v.MoveCursor(-7, 0)
	}
}

func (v *View) SetVtCols(cols int) {
	v.backend.SetVtCols(cols)
}

func (v *View) ScrollPos() (ln, col int) {
	return v.offy, v.offx
}

func (v *View) SetScrollPos(ln, col int) {
	v.offy, v.offx = ln, col
}

func (v *View) Text(ln1, col1, ln2, col2 int) [][]rune {
	return v.SelectionText(core.NewSelection(ln1, col1, ln2, col2))
}

func (v *View) Type() core.ViewType {
	return v.viewType
}

func (v *View) SetScrollPct(ypct int) {
	if ypct < 0 {
		ypct = 0
	}
	if ypct > 100 {
		ypct = 100
	}
	lc := v.LineCount() - v.LastViewLine()
	v.SetScrollPos(lc*ypct/100, 0)
}
