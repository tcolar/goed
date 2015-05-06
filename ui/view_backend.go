package ui

import (
	"bytes"

	"github.com/tcolar/goed/core"
)

func (v *View) Save() {
	e := core.Ed
	err := v.backend.Save(v.backend.SrcLoc())
	if err != nil {
		e.SetStatusErr("Saving Failed " + err.Error())
		return
	}
	v.SetDirty(false)
	e.SetStatus("Saved " + v.backend.SrcLoc())
}

func (v *View) InsertCur(s string) {
	_, x, y := v.CurChar()
	if len(v.selections) > 0 {
		s := v.selections[0]
		v.MoveCursorRoll(s.ColFrom-x-1, s.LineFrom-y-1)
		v.SelectionDelete(&s)
		v.ClearSelections()
	}
	_, x, y = v.CurChar()
	v.Insert(y, x, s)
}

// Insert inserts text at the given text location
func (v *View) Insert(row, col int, s string) {
	e := core.Ed
	// backend is 1-based indexed
	err := v.backend.Insert(row+1, col+1, s)
	if err != nil {
		e.SetStatusErr("Insert Failed " + err.Error())
		return
	}
	// move the cursor to after insertion
	b := []byte(s)
	offy := bytes.Count(b, core.LineSep)
	idx := bytes.LastIndex(b, core.LineSep)
	if idx < 0 {
		idx = 0
	}
	offx := v.strSize(string(b[idx:]))
	v.Render()
	v.MoveCursor(offx, offy)
}

func (v *View) InsertNewLineCur() {
	v.InsertCur("\n")
}

// InsertNewLine inserts a "newline"(Enter key) in the buffer
func (v *View) InsertNewLine(row, col int) {
	v.Insert(row, col, "\n")
}

func (v *View) Reload() {
	err := v.backend.Reload()
	if err != nil {
		core.Ed.SetStatusErr(err.Error())
	}
}

// Delete removes characters at the given text location
func (v *View) Delete(row1, col1, row2, col2 int) {
	err := v.backend.Remove(row1+1, col1+1, row2+1, col2+1)
	if err != nil {
		core.Ed.SetStatusErr("Delete Failed " + err.Error())
		return
	}
	v.Render()
}

// DeleteCur removes a selection or the curent character
func (v *View) DeleteCur() {
	c, x, y := v.CurChar()
	if len(v.selections) > 0 {
		s := v.selections[0]
		v.MoveCursorRoll(s.ColFrom-x-1, s.LineFrom-y-1)
		v.SelectionDelete(&s)
		v.ClearSelections()
		return
	}
	if c != nil {
		v.Delete(y, x, y, x)
	}
}

// Backspace removes a selection or character before the current location
func (v *View) Backspace() {
	if v.CursorY == 0 && v.CursorX == 0 {
		return
	}
	if len(v.selections) == 0 {
		v.MoveCursorRoll(-1, 0)
	}
	v.DeleteCur()
}

// LineCount return the number of lines in the  buffer
// if the last line is a blank line, do not count it
func (v *View) LineCount() int {
	return v.backend.LineCount()
}

// Line return the line at the given index
func (v *View) Line(s *core.Slice, lnIndex int) []rune {
	// backend is 1-based indexed
	index := lnIndex + 1 - s.R1
	if index < 0 || index >= len(*s.Text()) {
		return []rune{}
	}
	return (*s.Text())[index]
}

// LineLen returns the length of a line (raw runes length)
func (v *View) LineLen(s *core.Slice, lnIndex int) int {
	return len(v.Line(s, lnIndex))
}

// LineCol returns the number of columns used for the given lines
// ie: a tab uses multiple columns
func (v *View) lineCols(s *core.Slice, lnIndex int) int {
	return v.lineColsTo(s, lnIndex, v.LineLen(s, lnIndex))
}

// LineColsTo returns the number of columns up to the given line index
// ie: a tab uses multiple columns
func (v *View) lineColsTo(s *core.Slice, lnIndex, to int) int {
	line := v.Line(s, lnIndex)
	if lnIndex > v.LineCount() || to > len(line) {
		return 0
	}
	ln := 0
	for _, r := range line[:to] {
		ln += v.runeSize(r)
	}
	return ln
}

// lineRunesTo returns the number of raw runes to the given line column
func (v View) lineRunesTo(s *core.Slice, lnIndex, column int) int {
	runes := 0
	if len(*s.Text()) == 0 || lnIndex >= v.LineCount() || lnIndex < 0 {
		return 0
	}
	ln := v.Line(s, lnIndex)
	for i := 0; i <= column && runes < len(ln); {
		i += v.runeSize(ln[runes])
		if i <= column {
			runes++
		}
	}
	return runes
}

// CursorTextPos returns the position in the text buffer for a cursor location
func (v *View) CursorTextPos(s *core.Slice, cursorX, cursorY int) (int, int) {
	l := cursorY
	return v.lineRunesTo(s, l, cursorX), l
}

// CursorChar returns the rune at the given cursor location
// Also returns the position of the char in the text buffer
func (v *View) CursorChar(s *core.Slice, cursorX, cursorY int) (r *rune, textX, textY int) {
	// backend is 1-based indexed
	x, y := v.CursorTextPos(s, cursorX, cursorY)
	ln := v.Line(s, y)
	if len(ln) <= x { // EOL
		nl := '\n'
		return &nl, x, y
	} else if len(ln) <= x {
		return nil, x, y
	}
	return &ln[x], x, y
}

// CurChar returns the rune at the current cursor location
func (v *View) CurChar() (r *rune, textX, textY int) {
	return v.CursorChar(v.slice, v.CurCol(), v.CurLine())
}

// The runeSize (on screen)
// tabs are a special case
func (v *View) runeSize(r rune) int {
	if r == '\t' {
		return tabSize
	}
	return 1
}

// The string size (on screen)
// tabs are a special case
func (v *View) strSize(s string) int {
	ln := 0
	for _, r := range s {
		ln += v.runeSize(r)
	}
	return ln
}
