package main

import (
	"strings"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	var err error
	v := Ed.NewView()
	v.SetBounds(0, 0, 40, 25)

	err = Ed.Open("test_data/file1.txt", v, "")
	assert.Nil(t, err, "open")

	v.slice = v.backend.Slice(v.offy+1, v.offx+1, v.offy+v.LastViewLine()+1, v.offx+v.LastViewCol()+1)

	assertCursor(t, v, 0, 0, 0, 0, "mc")
	assert.True(t, strings.HasSuffix(v.backend.SrcLoc(), "test_data/file1.txt"), "srcloc")
	assert.True(t, strings.HasSuffix(v.WorkDir, "test_data"), "workdir")
	assert.Equal(t, v.backend.BufferLoc(), Ed.BufferFile(v.Id), "bufferloc")
	assert.False(t, v.Dirty, "dirty")
	assert.Equal(t, v.Title(), "file1.txt")
	assert.Equal(t, v.LineCount(), 12, "lineCount")
	assert.Equal(t, v.lineCols(v.slice, 0), 10, "lineCols")
	assert.Equal(t, v.LastViewLine(), 25-3, "lastViewLine")
	assert.Equal(t, v.LastViewCol(), 40-3, "lastViewCol")
	assert.Equal(t, string(v.Line(v.slice, 0)), "1234567890", "line")
	assert.Equal(t, v.LineLen(v.slice, 3), 26, "lineLen")
	assert.Equal(t, v.lineColsTo(v.slice, 0, 4), 4, "lineColsTo1")
	assert.Equal(t, v.lineColsTo(v.slice, 9, 4), 10, "lineColsTo2") //\t\t a
	assert.Equal(t, v.lineRunesTo(v.slice, 0, 4), 4, "lineRunesTo1")
	assert.Equal(t, v.lineRunesTo(v.slice, 9, 10), 4, "lineRunesTo2")
	x, y := v.CursorTextPos(v.slice, 4, 0)
	assert.Equal(t, x, 4, "cursortextpos_x1")
	assert.Equal(t, y, 0, "cursortextpos_y1")
	x, y = v.CursorTextPos(v.slice, 10, 9)
	assert.Equal(t, x, 4, "cursortextpos_x2")
	assert.Equal(t, y, 9, "cursortextpos_y2")
	c, x, y := v.CursorChar(v.slice, 3, 3)
	assert.Equal(t, *c, 'D', "cursorchar_c")
	assert.Equal(t, x, 3, "cursorchar_x")
	assert.Equal(t, y, 3, "cursorchar_y")
	assert.Equal(t, v.runeSize('a'), 1, "runSize1")
	assert.Equal(t, v.runeSize('\t'), tabSize, "runSize2")
	assert.Equal(t, v.strSize("abc"), 3, "strSize1")
	assert.Equal(t, v.strSize("a\tb\tc"), 3+2*tabSize, "strSize2")
	v.MoveCursor(0, 0)
	assertCursor(t, v, 0, 0, 0, 0, "mc1")
	//cursortextpos
	v.MoveCursor(5, 0)
	assertCursor(t, v, 5, 0, 0, 0, "mc2")
	c, x, y = v.CurChar()
	assert.Equal(t, x, 5, "curchar_x")
	assert.Equal(t, y, 0, "curchar_y")
	assert.Equal(t, *c, '6', "curchar_c")
	v.MoveCursor(-3, 0)
	assertCursor(t, v, 2, 0, 0, 0, "mc3")
	v.MoveCursor(2, 3)
	assertCursor(t, v, 4, 3, 0, 0, "mc4")
	v.MoveCursor(2, -1)
	assert.Equal(t, v.CurCol(), 6, "curcol")
	assert.Equal(t, v.CurLine(), 2, "curline")
	assertCursor(t, v, 6, 2, 0, 0, "mc5")
	v.MoveCursor(2, -1)
	// Note: x=0 because line "1" is blank
	assertCursor(t, v, 0, 1, 0, 0, "mc6")
	v.MoveCursor(-10, -10)
	assertCursor(t, v, 0, 0, 0, 0, "mc7")
	v.MoveCursor(100, 100)
	assertCursor(t, v, 36, 11, 0, 0, "mc8")
	v.MoveCursor(-100, -100)
	assertCursor(t, v, 0, 0, 0, 0, "mc9")
	v.MoveCursorRoll(10, 0)
	assertCursor(t, v, 10, 0, 0, 0, "mc10")
	v.MoveCursorRoll(1, 0)
	assertCursor(t, v, 0, 1, 0, 0, "mc11")
	v.MoveCursorRoll(-2, 0)
	assertCursor(t, v, 10, 0, 0, 0, "mc11")
}

func assertCursor(t *testing.T, v *View, x, y, offsetX, offsetY int, msg string) {
	assert.Equal(t, v.CursorX, x, msg+" CursorX")
	assert.Equal(t, v.CursorY, y, msg+" CursorY")
	assert.Equal(t, v.offx, offsetX, msg+" offsetX")
	assert.Equal(t, v.offy, offsetY, msg+" offsetY")
}

func TestViewSelections(t *testing.T) {
	var err error
	v := Ed.NewView()
	v.SetBounds(5, 5, 140, 30)
	err = Ed.Open("test_data/file1.txt", v, "")
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy+1, v.offx+1, v.offy+v.LastViewLine()+1, v.offx+v.LastViewCol()+1)

	s := NewSelection(1, 1, 1, 1)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "1")
	s = NewSelection(3, 2, 4, 8)
	v.Selections = append(v.Selections, *s)
	assert.Equal(t, s.String(), "(3,2)-(4,8)", "string")
	text := s.Text(v)
	assert.Equal(t, len(text), 2, "text length")
	assert.Equal(t, len(text[0]), 7, "text[0] length")
	assert.Equal(t, len(text[1]), 7, "text[1] length")
	assert.Equal(t, string(text[0]), "bcdefgh", "text[0]")
	assert.Equal(t, string(text[1]), "BCDEFGH", "text[1]")
	b, sel := v.Selected(4, 1)
	assert.False(t, b, "4,1")
	assert.Nil(t, sel, "sel 4,1")
	b, sel = v.Selected(1, 3)
	assert.False(t, b, "1,3")
	assert.Nil(t, sel, "sel 1,3")
	b, sel = v.Selected(4, 5)
	assert.False(t, b, "4,5")
	assert.Nil(t, sel, "sel 4,5")
	b, sel = v.Selected(3, 3)
	assert.True(t, b, "3, 3")
	assert.Equal(t, sel.String(), s.String(), "sel 3,3")
	v.Copy(*s)
	cb, _ := clipboard.ReadAll()
	assert.Equal(t, cb, "bcdefgh\nBCDEFGH", "copy")
	s = v.PathSelection(1, 3)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "1234567890", "path2")
	s = v.PathSelection(11, 1)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "aaa", "ps1")
	s = v.PathSelection(11, 6)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "aaa.go", "ps2")
	s = v.PathSelection(11, 22)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "/tmp/aaa.go", "ps3")
	s = v.PathSelection(11, 27)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "aaa.go:23", "ps4")
	s = v.PathSelection(11, 39)
	assert.Equal(t, Ed.RunesToString(s.Text(v)), "/tmp/aaa.go:23:7", "ps5")
	loc, ln, col := v.selToLoc(*s)
	assert.Equal(t, loc, "/tmp/aaa.go", "loc")
	assert.Equal(t, ln, 23, "ln")
	assert.Equal(t, col, 7, "col")
}

func TestViewEdition(t *testing.T) {
	var err error
	v := Ed.NewView()
	v.SetBounds(5, 5, 140, 30)
	err = Ed.Open("test_data/empty.txt", v, "")
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy+1, v.offx+1, v.offy+v.LastViewLine()+1, v.offx+v.LastViewCol()+1)

	assert.Equal(t, v.LineCount(), 0, "lineCount")
	// TODO : insert
	// TODO : insertnewline
	// TODO : delete
	// TODO : backspace
}

// TODO: test scrolling etc...
func TestViewScrolling(t *testing.T) {
}

// TODO: test term mock
// TODO: save etc ....
