package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestView(t *testing.T) {
	var err error
	v := Ed.NewView()
	v.SetBounds(5, 5, 40, 30)

	err = Ed.Open("test_data/file1.txt", v, "")
	assert.Nil(t, err, "open")

	assertCursor(t, v, 0, 0, 0, 0, "mc")
	assert.True(t, strings.HasSuffix(v.backend.SrcLoc(), "test_data/file1.txt"), "srcloc")
	assert.True(t, strings.HasSuffix(v.WorkDir, "test_data"), "workdir")
	assert.Equal(t, v.backend.BufferLoc(), Ed.BufferFile(v.Id), "bufferloc")
	assert.False(t, v.Dirty, "dirty")
	assert.Equal(t, v.Title(), "file1.txt")
	assert.Equal(t, v.LineCount(), 12, "lineCount")
	assert.Equal(t, v.lineCols(0), 10, "lineCols")
	assert.Equal(t, v.LastViewLine(), 30-5-3, "lastViewLine")
	assert.Equal(t, v.LastViewCol(), 40-5-3, "lastViewCol")
	v.MoveCursor(0, 0)
	assertCursor(t, v, 0, 0, 0, 0, "mc1")
	v.MoveCursor(5, 0)
	assertCursor(t, v, 5, 0, 0, 0, "mc2")
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
	v.MoveCursor(-10, -10) // should do nothing
	assertCursor(t, v, 0, 1, 0, 0, "mc7")
	v.MoveCursor(100, 100) // should do nothing
	assertCursor(t, v, 0, 1, 0, 0, "mc8")
}

func TestViewSelections(t *testing.T) {
	// TODO: test selection stuff / copy/paste
}

func assertCursor(t *testing.T, v *View, x, y, offsetX, offsetY int, msg string) {
	assert.Equal(t, v.CursorX, x, msg+" CursorX")
	assert.Equal(t, v.CursorY, y, msg+" CursorY")
	assert.Equal(t, v.offx, offsetX, msg+" offsetX")
	assert.Equal(t, v.offy, offsetY, msg+" offsetY")
}
