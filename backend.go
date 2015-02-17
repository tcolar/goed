package main

import (
	"bytes"
	"io"
)

type Backend interface {
	SrcLoc() string    // "original" source
	BufferLoc() string // buffer location

	Insert(row, col int, text string) error
	Remove(row1, col1, row, col2 int) error

	LineCount() int

	Save(loc string) error

	// Get a region ("rectangle") as a runes matrix
	Slice(row, col, row2, col2 int) *Slice

	Close() error

	ViewId() int

	// Completely clears the buffer (empty)
	Wipe()

	//Sync() error         // sync from source ?
	//IsStale() bool       // whether the source as changed under us (fsnotify)
	//IsBufferStale() bool // whether the buffer has changed under us

	//SourceMd5 or ts?
	//BufferMd5 or ts?

	//Refresh() // TODO: refresh from content (rerun command or refresh src file)
}

type Slice struct {
	text           [][]rune
	r1, c1, r2, c2 int //bounds
}

func NewSlice(r1, c1, r2, c2 int, text [][]rune) *Slice {
	return &Slice{
		r1:   r1,
		c1:   c1,
		r2:   r2,
		c2:   c2,
		text: text,
	}
}

type Rwsc interface {
	io.Reader
	io.Writer
	io.ReaderAt
	io.WriterAt
	io.Seeker
	io.Closer
}

func (v *View) Save() {
	err := v.backend.Save(v.backend.SrcLoc())
	if err != nil {
		Ed.SetStatusErr("Saving Failed " + err.Error())
		return
	}
	v.Dirty = false
	Ed.SetStatus("Saved " + v.backend.SrcLoc())
}

// Insert inserts text at the given text location
func (v *View) Insert(row, col int, s string) {
	// backend is 1-based indexed
	err := v.backend.Insert(row+1, col+1, s)
	if err != nil {
		Ed.SetStatusErr("Insert Failed " + err.Error())
		return
	}
	// move the cursor to after insertion
	b := []byte(s)
	offy := bytes.Count(b, LineSep)
	idx := bytes.LastIndex(b, LineSep)
	if idx < 0 {
		idx = 0
	}
	offx := v.strSize(string(b[idx:]))
	v.Render()
	v.MoveCursor(offx, offy)
}

// InsertNewLine inserts a "newline"(Enter key) in the buffer
func (v *View) InsertNewLine(row, col int) {
	v.Insert(row, col, "\n")
}

// Delete removes characters at the given text location
func (v *View) Delete(row1, col1, row2, col2 int) {
	err := v.backend.Remove(row1+1, col1+1, row2+1, col2+1)
	if err != nil {
		Ed.SetStatusErr("Delete Failed " + err.Error())
		return
	}
	v.Render()
}

// Backspace removes a character before the current location
func (v *View) Backspace() {
	if v.CursorY == 0 && v.CursorX == 0 {
		return
	}
	v.MoveCursorRoll(-1, 0)
	_, x, y := v.CurChar()
	v.Delete(y, x, y, x)
}

// LineCount return the number of lines in the  buffer
// if the last line is a blank line, do not count it
func (v *View) LineCount() int {
	return v.backend.LineCount()
}

// Line return the line at the given index
func (v *View) Line(s *Slice, lnIndex int) []rune {
	// backend is 1-based indexed
	index := lnIndex + 1 - s.r1
	if index < 0 || index >= len(s.text) {
		return []rune{}
	}
	return s.text[index]
}

// LineLen returns the length of a line (raw runes length)
func (v *View) LineLen(s *Slice, lnIndex int) int {
	return len(v.Line(s, lnIndex))
}

// LineCol returns the number of columns used for the given lines
// ie: a tab uses multiple columns
func (v *View) lineCols(s *Slice, lnIndex int) int {
	return v.lineColsTo(s, lnIndex, v.LineLen(s, lnIndex))
}

// LineColsTo returns the number of columns up to the given line index
// ie: a tab uses multiple columns
func (v *View) lineColsTo(s *Slice, lnIndex, to int) int {
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
func (v View) lineRunesTo(s *Slice, lnIndex, column int) int {
	runes := 0
	if len(s.text) == 0 || lnIndex >= v.LineCount() || lnIndex < 0 {
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
func (v *View) CursorTextPos(s *Slice, cursorX, cursorY int) (int, int) {
	l := cursorY
	return v.lineRunesTo(s, l, cursorX), l
}

// CursorChar returns the rune at the given cursor location
// Also returns the position of the char in the text buffer
func (v *View) CursorChar(s *Slice, cursorX, cursorY int) (r *rune, textX, textY int) {
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
