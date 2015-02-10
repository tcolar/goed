package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

type MemBackend struct {
	text [][]rune
	file string
	view *View
}

func (e *Editor) NewMemBackend(path string, viewId int) *MemBackend {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	runes := Ed.StringToRunes(data)
	return &MemBackend{
		text: runes,
		file: path,
	}
}

func (b *MemBackend) Save(loc string) error {
	if len(loc) == 0 {
		return fmt.Errorf("Save where ? Use save [path]")
	}
	f, err := os.Create(loc)
	if err != nil {
		return fmt.Errorf("Saving Failed ! %v", loc)
	}
	defer f.Close()
	buf := make([]byte, 4)
	for i, l := range b.text {
		for _, c := range l {
			n := utf8.EncodeRune(buf, c)
			_, err := f.Write(buf[0:n])
			if err != nil {
				return fmt.Errorf("Saved Failed failed %v", err.Error())
			}
		}
		if i != b.LineCount() || b.view.LineLen(i) != 0 {
			f.WriteString("\n")
		}
	}
	b.file = loc
	Ed.SetStatus("Saved " + b.file)
	return nil
}

func (b *MemBackend) SrcLoc() string {
	return b.file
}

func (b *MemBackend) BufferLoc() string {
	return "_MEM_" // TODO : BufferLoc for in-memory ??
}

func (b *MemBackend) Insert(row, col int, text string) error {
	runes := Ed.StringToRunes([]byte(text))
	if len(runes) == 0 {
		return nil
	}
	var tail []rune
	last := len(runes) - 1
	// Create a "hole" for the new lines to be inserted
	if len(runes) > 2 {
		for i := 2; i < len(runes); i++ {
			b.text = append(b.text, []rune{})
		}
		copy(b.text[row+len(runes)-2:], b.text[row+1:])
	}
	for i, ln := range runes {
		if i == 0 {
			tail = b.text[row][col:]
			b.text[row+i] = append(b.text[row+i], ln...)
		}
		if i == last {
			b.text[row+i] = append(ln, tail...)
		}
		if i > 0 && i < last {
			b.text[row+i] = runes[row+i]
		}
	}

	return nil
}

func (b *MemBackend) Remove(row, col int, text string) error {
	return nil
}

func (b *MemBackend) Slice(row, col, row2, col2 int) [][]rune {
	if row < 1 || col < 1 {
		return [][]rune{}
	}
	if row2 != -1 && row > row2 {
		row, row2 = row2, row
	}
	if col2 != -1 && col > col2 {
		col, col2 = col2, col
	}
	runes := [][]rune{}
	r := row
	for ; row2 == -1 || r <= row2; r++ {
		if col2 == -1 {
			runes = append(runes, b.text[row])
		} else {
			c, c2, l := col, col2, len(b.text[row])
			if c > l {
				c = l
			}
			if c2 > l {
				c2 = l
			}
			runes = append(runes, b.text[row][c:c2])
		}
	}
	return runes
}

func (b *MemBackend) LineCount() int {
	return len(b.text)
}

func (b *MemBackend) Close() error {
	return nil // Noop
}

/*
// Inserts a rune at the cursor location
func (v *View) Insert(c rune) {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	if l >= len(v.Buffer.text) {
		v.Buffer.text = append(v.Buffer.text, []rune{})
	}
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
	if l >= len(v.Buffer.text) {
		v.Buffer.text = append(v.Buffer.text, []rune{})
		return
	}
	line := v.Buffer.text[l]
	indent := []rune{}
	for _, r := range line {
		if unicode.IsSpace(r) {
			indent = append(indent, r)
		} else {
			break
		}
	}
	b := append(v.Buffer.text, []rune{}) // Extend buffer size with a new blank line
	copy(b[l+1:], b[l:])                 // Move buffer tail by one to create a "hole" (blank line)
	b[l] = line[:i]                      // truncate current line up to cursor
	b[l+1] = append(indent, line[i:]...) // make rest of current line it's own line
	v.Buffer.text = b
	v.MoveCursor(v.lineColsTo(l+1, len(indent))-v.CurCol(), 1)
}

// TODO: This is not most efficient
func (v *View) InsertLines(lines [][]rune) {
	for i, l := range lines {
		for _, r := range l {
			v.Insert(r)
		}
		if i < len(lines)-1 {
			v.InsertNewLine()
		}
	}
}

// Delete removes a character at the current location
func (v *View) Delete() {
	l := v.CurLine()
	i := v.lineRunesTo(l, v.CurCol())
	if l >= len(v.Buffer.text) {
		return
	}
	line := v.Buffer.text[l]
	if i < len(line) {
		// remove a char
		v.Buffer.text[l] = append(line[:i], line[i+1:]...)
	} else if l+1 < v.LineCount() {
		// at end of line, pull the next line up to end of current line
		v.Buffer.text[l] = append(line, v.Buffer.text[l+1]...)
		v.Buffer.text = append(v.Buffer.text[:l+1], v.Buffer.text[l+2:]...)
	}
}*/
