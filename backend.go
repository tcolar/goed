package main

import "io"

// TODO: flush this out + File based impl
// TODO: use bufio
type Backend interface {
	SrcLoc() string    // "original" source
	BufferLoc() string // buffer location

	Insert(row, col int, text string) error
	Remove(row, col int, text string) error

	LineCount() int

	Save(loc string) error

	// Get a region ("rectangle") as a runes matrix
	Slice(row, col, line2, row2 int) [][]rune

	Close() error

	//Sync() error         // sync from source ?
	//IsStale() bool       // whether the source as changed under us (fsnotify)
	//IsBufferStale() bool // whether the buffer has changed under us

	//SourceMd5 or ts?
	//BufferMd5 or ts ?

	//Reset() // TODO: refresh from content (rerun command or refresh src file)
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
func (v *View) Insert(s string) {
	// backend is 1-based indexed
	_, x, y := v.CurChar()
	err := v.backend.Insert(y+1, x+1, s)
	if err != nil {
		Ed.SetStatusErr("Insert Failed " + err.Error())
		return
	}
	v.MoveCursor(v.strSize(s), 0)
}

// InsertNewLine inserts a "newline"(Enter key) in the buffer
func (v *View) InsertNewLine() {
	v.Insert("\n")
}

// Delete removes characters at the current location
func (v *View) Delete(s string) {
	_, x, y := v.CurChar()
	Ed.SetStatus("Removing " + s)
	// backend is 1-based indexed
	err := v.backend.Remove(y+1, x+1, s)
	if err != nil {
		Ed.SetStatusErr("Delete Failed " + err.Error())
		return
	}
}

// Backspace removes a character before the current location
func (v *View) Backspace() {
	if v.CursorY == 0 && v.CursorX == 0 {
		return
	}
	v.MoveCursor(-1, 0)
	c, _, _ := v.CurChar()
	if c == nil {
		return
	}
	v.Delete(string(*c))
}

// LineCount return the number of lines in the buffer
func (v *View) LineCount() int {
	return v.backend.LineCount()
}

// Line return the line at the given index
func (v *View) Line(lnIndex int) []rune {
	// backend is 1-based indexed
	ln := v.backend.Slice(lnIndex+1, 1, lnIndex+1, -1)
	if len(ln) < 1 {
		return []rune{}
	}
	return ln[0]
}

// LineLen returns the length of a line (raw runes length)
func (v *View) LineLen(lnIndex int) int {
	return len(v.Line(lnIndex))
}

// LineCol returns the number of columns used for the given lines
// ie: a tab uses multiple columns
func (v *View) lineCols(lnIndex int) int {
	return v.lineColsTo(lnIndex, v.LineLen(lnIndex))
}

// LineColsTo returns the number of columns up to the given line index
// ie: a tab uses multiple columns
func (v *View) lineColsTo(lnIndex, to int) int {
	if v.LineCount() <= lnIndex || v.LineLen(lnIndex) < to {
		return 0
	}
	ln := 0
	for _, r := range v.Line(lnIndex)[:to] {
		ln += v.runeSize(r)
	}
	return ln
}

// lineRunesTo returns the number of raw runes to the given line column
func (v View) lineRunesTo(lnIndex, column int) int {
	runes := 0
	if lnIndex >= v.LineCount() || lnIndex < 0 {
		return 0
	}
	ln := v.Line(lnIndex)
	for i := 0; i <= column && runes < len(ln); {
		i += v.runeSize(ln[runes])
		if i <= column {
			runes++
		}
	}
	return runes
}

// CursorTextPos returns the position in the text buffer for a cursor location
func (v *View) CursorTextPos(cursorX, cursorY int) (int, int) {
	l := cursorY
	return v.lineRunesTo(l, cursorX), l
}

// CursorChar returns the rune at the given cursor location
// Also returns the position of the char in the text buffer
func (v *View) CursorChar(cursorX, cursorY int) (r *rune, textX, textY int) {
	// backend is 1-based indexed
	x, y := v.CursorTextPos(cursorX, cursorY)
	ln := v.Line(y)
	if len(ln) == x { // EOL
		nl := '\n'
		return &nl, x, y
	} else if len(ln) <= x {
		return nil, x, y
	}
	return &ln[x], x, y
}

// CurChar returns the rune at the current cursor location
func (v *View) CurChar() (r *rune, textX, textY int) {
	return v.CursorChar(v.CurCol(), v.CurLine())
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
func (v *View) strSize(s string) int {
	ln := 0
	for _, r := range s {
		ln += v.runeSize(r)
	}
	return ln
}
