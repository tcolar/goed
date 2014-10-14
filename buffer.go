package main

import (
	"bytes"
	"io/ioutil"
)

type Buffer struct {
	text [][]rune
}

// For now just load the whole thing in memory, might chnage this later
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
	}
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
	v.MoveCursor(1, 0)
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
		ln++
		if r == '\t' {
			ln += tabSize - 1
		}
	}
	return ln
}

// lineRunesTo returns the number of raw runes to the given line column
func (v View) lineRunesTo(lnIndex, column int) int {
	runes := 0
	ln := v.Line(lnIndex)
	for i := 0; i < column; i++ {
		runes++
		if ln[i] == '\t' {
			i += tabSize - 1
		}
	}
	return runes
}
