package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

// TODO: flush this out + File based impl
// TODO: use bufio
type Backend interface {
	Location() string  // "original" source
	BufferLoc() string // buffer location (ie: file or mem)
	Title() string
	Insert(c rune, row, col int) error
	//Remove(r1,c1,count int) error // remove n runes
	Save() error         // save to source
	Sync() error         // sync from source ?
	IsStale() bool       // whether the source as changed under us (fsnotify)
	IsBufferStale() bool // wether the buffer has changed under us
	LineCount() int

	// insert(lineId, index, []rune)
	// remove(lineId, from, len int)
	// insertLine(index, []rune)
	// removeLine(index)
}

type Buffer struct {
	text [][]rune
	file string
}

// For now just load the whole thing in memory, might change this later
func NewFileBuffer(path string) *Buffer {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(data, []byte("\n"))
	runes := [][]rune{}
	for i, l := range lines {
		// Ignore last line if empty
		if i != len(lines)-1 || len(l) != 0 {
			runes = append(runes, bytes.Runes(l))
		}
	}
	return &Buffer{
		text: runes,
		file: path,
	}
}

// TODO : this is an in memory "stupid" impl
// Later implement "direct to file" edition ?
func (v *View) Save() {
	f, err := os.Create("/tmp/a.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, 4)
	for i, l := range v.Buffer.text {
		for _, c := range l {
			n := utf8.EncodeRune(buf, c)
			_, err := f.Write(buf[0:n])
			if err != nil {
				panic(err)
			}
		}
		if i != v.LineCount() || v.LineLen(i) != 0 { // ??
			f.WriteString("\n")
		}
	}
	Ed.SetStatus("Saved " + v.Buffer.file)
}

// Inserts a rune at the cursor location
func (v *View) Insert(c rune) {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	line := v.Buffer.text[l]
	line = append(line, c)
	copy(line[i+1:], line[i:])
	line[i] = c
	v.Buffer.text[l] = line
	v.MoveCursor(v.runeSize(c), 0)
}

// InsertNewLine inserts a "newline"(Enter key) in the buffer
func (v *View) InsertNewLine() {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	line := v.Buffer.text[l]
	b := append(v.Buffer.text, []rune{}) // Extend buffer size with a new blank line
	copy(b[l+1:], b[l:])                 // Move buffer tail by one to create a "hole" (blank line)
	b[l] = line[:i]                      // truncate current line up to cursor
	b[l+1] = line[i:]                    // make rest of current line it's own line
	v.Buffer.text = b
	v.MoveCursor(1, 0)
}

// Delete removes a character at the current location
func (v *View) Delete() {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	line := v.Buffer.text[l]
	if i < len(line) {
		// remove a char
		v.Buffer.text[l] = append(line[:i], line[i+1:]...)
	} else if l+1 < v.LineCount() {
		// at end of line, pull the next line up to end of current line
		v.Buffer.text[l] = append(line, v.Buffer.text[l+1]...)
		v.Buffer.text = append(v.Buffer.text[:l+1], v.Buffer.text[l+2:]...)
	}
}

// Backspace removes a character before the current location
func (v *View) Backspace() {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	line := v.Buffer.text[l]
	if i > 0 {
		c := line[i-1]
		// remove a char
		v.Buffer.text[l] = append(line[:i-1], line[i:]...)
		v.MoveCursor(-v.runeSize(c), 0)
	} else if l > 0 {
		// at beginning of line, pull the current line to the end of prev line
		prevCols := v.lineColsTo(l-1, v.LineLen(l-1))
		v.Buffer.text[l-1] = append(v.Buffer.text[l-1], line...)
		v.Buffer.text = append(v.Buffer.text[:l], v.Buffer.text[l+1:]...)
		v.MoveCursor(prevCols, -1) // place cursor at end of prev line (before pull)
	}
}

// LineCount return the number of lines in the buffer
func (v *View) LineCount() int {
	return len(v.Buffer.text)
}

// Line return the line at the given index
func (v *View) Line(lnIndex int) []rune {
	return v.Buffer.text[lnIndex]
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
func (v *View) CursorChar(cursorX, cursorY int) (r *rune, textX int, textY int) {
	x, y := v.CursorTextPos(cursorX, cursorY)
	if y >= v.LineCount() || x >= len(v.Buffer.text[y]) {
		return nil, 0, 0
	}
	return &v.Buffer.text[y][x], x, y
}

// CurChar returns the rune at the current cursor location
func (v *View) CurChar() (r *rune, textX int, textY int) {
	return v.CursorChar(v.CurCol(), v.CurLine())
}

// The runeSize (on screen)
// tabs are a special case
func (v View) runeSize(r rune) int {
	if r == '\t' {
		return tabSize
	}
	return 1
}
