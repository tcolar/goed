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
		Ed.SetStatusErr("Saving Failed %s " + err.Error())
		return
	}
	Ed.SetStatus("Saved " + v.backend.SrcLoc())
}
func (v *View) Insert(s string) {
	// backend is 1-based indexed
	err := v.backend.Insert(v.CurLine()+1, v.CurCol()+1, s)
	if err != nil {
		Ed.SetStatusErr("Saving Failed %s " + err.Error())
		return
	}
	v.MoveCursor(v.strSize(s), 0)
}

// InsertNewLine inserts a "newline"(Enter key) in the buffer
func (v *View) InsertNewLine() {
	v.Insert("\n")
}

// Delete removes a character at the current location
func (v *View) Delete(s string) {
	// backend is 1-based indexed
	err := v.backend.Remove(v.CurLine()+1, v.CurCol()+1, s)
	if err != nil {
		Ed.SetStatusErr("Delete Failed %s " + err.Error())
		return
	}
}

// Backspace removes a character before the current location
func (v *View) Backspace() {
	c, _, _ := v.CurChar()
	if c == nil {
		return
	}
	err := v.backend.Remove(v.CurLine(), v.CurCol(), string(*c))
	if err != nil {
		Ed.SetStatusErr("Delete Failed %s " + err.Error())
		return
	}
	v.MoveCursor(-v.runeSize(*c), 0)
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
	s := v.backend.Slice(cursorY+1, cursorX+1, cursorY+1, cursorX+1)
	if len(s) == 0 || len(s[0]) == 0 {
		return nil, 0, 0
	}
	x, y := v.CursorTextPos(cursorX, cursorY)
	return &s[0][0], x, y
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
