package ui

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

func TestView(t *testing.T) {
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(0, 0, 25, 40)

	_, err = Ed.Open("../test_data/file1.txt", v.Id(), "", false)
	assert.Nil(t, err, "open")

	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	assertCursor(t, v, 0, 0, 0, 0, "mc")
	assert.True(t, strings.HasSuffix(v.backend.SrcLoc(), "/test_data/file1.txt"), fmt.Sprintf("srcloc %s", v.backend.SrcLoc()))
	fmt.Println(v.workDir)
	assert.True(t, strings.HasSuffix(v.workDir, "/test_data"), fmt.Sprintf("workdir %s", v.workDir))
	assert.False(t, v.Dirty(), "dirty")
	assert.Equal(t, v.Title(), "file1.txt")
	assert.Equal(t, v.LineCount(), 12, "lineCount")
	assert.Equal(t, v.LastViewLine(), 25-3, "lastViewLine")
	assert.Equal(t, v.LastViewCol(), 40-3, "lastViewCol")
	assert.Equal(t, v.lineCols(v.slice, 0), 10, "lineCols")
	assert.Equal(t, string(v.Line(v.slice, 0)), "1234567890", "line")
	assert.Equal(t, v.LineLen(v.slice, 3), 26, "lineLen")
	assert.Equal(t, v.lineColsTo(v.slice, 0, 4), 4, "lineColsTo1")
	assert.Equal(t, v.lineColsTo(v.slice, 9, 4), 10, "lineColsTo2") //\t\t a
	assert.Equal(t, v.LineRunesTo(v.slice, 0, 4), 4, "lineRunesTo1")
	assert.Equal(t, v.LineRunesTo(v.slice, 9, 10), 4, "lineRunesTo2")
	x := v.LineRunesTo(v.slice, 0, 4)
	assert.Equal(t, x, 4, "cursortextpos_x1")
	x = v.LineRunesTo(v.slice, 9, 10)
	assert.Equal(t, x, 4, "cursortextpos_x2")
	c, y, x := v.CursorChar(v.slice, 3, 3)
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
	v.MoveCursor(0, 5)
	assertCursor(t, v, 0, 5, 0, 0, "mc2")
	c, y, x = v.CurChar()
	assert.Equal(t, x, 5, "curchar_x")
	assert.Equal(t, y, 0, "curchar_y")
	assert.Equal(t, *c, '6', "curchar_c")
	v.MoveCursor(0, -3)
	assertCursor(t, v, 0, 2, 0, 0, "mc3")
	v.MoveCursor(3, 2)
	assertCursor(t, v, 3, 4, 0, 0, "mc4")
	v.MoveCursor(-1, 2)
	assert.Equal(t, v.CurCol(), 6, "curcol")
	assert.Equal(t, v.CurLine(), 2, "curline")
	assertCursor(t, v, 2, 6, 0, 0, "mc5")
	v.MoveCursor(-1, 2)
	// Note: x=0 because line "1" is blank
	assertCursor(t, v, 1, 0, 0, 0, "mc6")
	v.MoveCursor(-10, -10)
	assertCursor(t, v, 0, 0, 0, 0, "mc7")
	v.MoveCursor(100, 100)
	assertCursor(t, v, 11, 36, 0, 0, "mc8")
	v.MoveCursor(-100, -100)
	assertCursor(t, v, 0, 0, 0, 0, "mc9")
	v.MoveCursorRoll(0, 10)
	assertCursor(t, v, 0, 10, 0, 0, "mc10")
	v.MoveCursorRoll(0, 1)
	assertCursor(t, v, 1, 0, 0, 0, "mc11")
	v.MoveCursorRoll(0, -2)
	assertCursor(t, v, 0, 10, 0, 0, "mc11")

}

func assertCursor(t *testing.T, v *View, y, x, offsetY, offsetX int, msg string) {
	assert.Equal(t, v.CursorX, x, msg+" CursorX")
	assert.Equal(t, v.CursorY, y, msg+" CursorY")
	assert.Equal(t, v.offx, offsetX, msg+" offsetX")
	assert.Equal(t, v.offy, offsetY, msg+" offsetY")
}

func TestViewSelections(t *testing.T) {
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/file1.txt", v.Id(), "", false)
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	s := core.NewSelection(0, 0, 0, 0)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "1")
	s = core.NewSelection(2, 1, 3, 7)
	v.selections = append(v.selections, *s)
	assert.Equal(t, s.String(), "2 1 3 7", "string")
	text := v.SelectionText(s)
	assert.Equal(t, len(text), 2, "text length")
	assert.Equal(t, len(text[0]), 25, "text[0] length")
	assert.Equal(t, len(text[1]), 8, "text[1] length")
	assert.Equal(t, string(text[0]), "bcdefghijklmnopqrstuvwxyz", "text[0]")
	assert.Equal(t, string(text[1]), "ABCDEFGH", "text[1]")
	b, sel := v.Selected(3, 0)
	assert.False(t, b, "3,0")
	assert.Nil(t, sel, "sel 3,0")
	b, sel = v.Selected(0, 2)
	assert.False(t, b, "0,2")
	assert.Nil(t, sel, "sel 0,2")
	b, sel = v.Selected(3, 4)
	assert.False(t, b, "3,4")
	assert.Nil(t, sel, "sel 3,4")
	b, sel = v.Selected(2, 2)
	assert.True(t, b, "2, 2")
	assert.Equal(t, sel.String(), s.String(), "sel 2,2")
	oldcb, _ := clipboard.ReadAll()
	v.SelectionCopy(s)
	cb, _ := clipboard.ReadAll()
	defer clipboard.WriteAll(oldcb) // restore the clipbaord after tests
	assert.Equal(t, cb, "bcdefghijklmnopqrstuvwxyz\nABCDEFGH", "copy")
	s = v.ExpandSelectionPath(0, 2)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "1234567890", "path2")
	s = v.ExpandSelectionPath(9, 2)
	assert.Equal(t, s.String(), "9 2 9 4", "abc")
	s = v.ExpandSelectionPath(10, 0)
	assert.Equal(t, s.String(), "10 0 10 2", "ps1 selection")
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa", "ps1")
	s = v.ExpandSelectionPath(10, 5)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa.go", "ps2")
	s = v.ExpandSelectionPath(10, 21)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "/tmp/aaa.go", "ps3")
	s = v.ExpandSelectionPath(10, 26)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa.go:23", "ps4")
	s = v.ExpandSelectionPath(10, 38)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "/tmp/aaa.go:23:7", "ps5")
	loc, ln, col := v.SelectionToLoc(s)
	assert.Equal(t, loc, "/tmp/aaa.go", "loc")
	assert.Equal(t, ln, 23, "ln")
	assert.Equal(t, col, 7, "col")
	s = v.ExpandSelectionWord(10, 0)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa", "ws1")
	s = v.ExpandSelectionWord(10, 1)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa", "ws2")
	s = v.ExpandSelectionWord(10, 2)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "aaa", "ws3")
	s = v.ExpandSelectionWord(10, 3)
	assert.Nil(t, s, "ws4")
	s = v.ExpandSelectionWord(10, 12)
	assert.Equal(t, core.RunesToString(v.SelectionText(s)), "tmp", "ws5")
}

func TestViewEdition(t *testing.T) {
	// Note: more tests done directly on backemds
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/empty.txt", v.Id(), "", false)
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	v.Insert(0, 0, "a", true)
	testChar(t, v, 0, 0, 'a')
	v.Insert(0, 0, "1", true)
	v.Render()
	testChar(t, v, 0, 0, '1')
	testChar(t, v, 0, 1, 'a')
	v.MoveCursor(0, 1)
	v.Backspace()
	testChar(t, v, 0, 0, '1')
	v.Delete(0, 0, 0, 0, true)
	v.Insert(0, 0, "b", true)
	testChar(t, v, 0, 0, 'b')
	v.Insert(0, 1, "cd", true)
	testChar(t, v, 0, 0, 'b')
	testChar(t, v, 0, 1, 'c')
	testChar(t, v, 0, 2, 'd')
	v.Delete(0, 0, 0, 1, true)
	testChar(t, v, 0, 0, 'd')
	v.Insert(0, 1, "e", true)
	testChar(t, v, 0, 0, 'd')
	testChar(t, v, 0, 1, 'e')
	v.InsertNewLine(0, 1)
	testChar(t, v, 0, 0, 'd')
	testChar(t, v, 1, 0, 'e')
	v.Delete(0, 1, 0, 1, true)
	testChar(t, v, 0, 0, 'd')
	testChar(t, v, 0, 1, 'e')
	v.Delete(0, 1, 0, 1, true)
	testChar(t, v, 0, 0, 'd')
	v.Delete(0, 1, 0, 1, true)
	testChar(t, v, 0, 0, 'd')
}

func TestDelete(t *testing.T) {
	// Test some edge cases
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/empty.txt", v.Id(), "", false)
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	v.Insert(0, 0, "ab\ncd\nef", true)
	s := core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "ab\ncd\nef")
	// backspace edge cases
	v.CursorX = 3
	v.CursorY = 0
	v.Backspace()
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abcd\nef")
	v.CursorX = 0
	v.CursorY = 0
	v.Backspace()
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abcd\nef")
	// delete edge cases
	v.CursorX = 5
	v.CursorY = 0
	v.DeleteCur()
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abcdef")
	v.CursorX = 7
	v.CursorY = 0
	v.DeleteCur()
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abcdef")
	v.Delete(0, 0, 0, 5, true)
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "")
	v.DeleteCur() // check no panic
	v.Backspace() // check no panic
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "")

	_, err = Ed.Open("../test_data/no_eol.txt", v.Id(), "", false)
	assert.Nil(t, err, "open")
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abc\n123\nz")
	v.MoveCursor(-100, -100)
	assert.Equal(t, v.CurLine(), 0)
	assert.Equal(t, v.CurCol(), 0)
	v.MoveCursor(2, 0)
	assert.Equal(t, v.CurLine(), 2)
	assert.Equal(t, v.CurCol(), 0)
	v.MoveCursorRoll(0, 1)
	assert.Equal(t, v.CursorY, 2)
	assert.Equal(t, v.CursorX, 1)
	v.InsertCur("Z")
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abc\n123\nzZ")
	v.MoveCursorRoll(0, -1)
	v.InsertCur("A")
	v.InsertCur("B")
	s = core.RunesToString(*v.backend.Slice(0, 0, 10, 10).Text())
	assert.Equal(t, s, "abc\n123\nzABZ")

}

func testChar(t *testing.T, v *View, y, x int, c rune) {
	Ed := core.Ed.(*Editor)
	s := v.backend.Slice(y, x, y, x)
	if (*s.Text())[0][0] != c {
		panic(c)
	}
	assert.Equal(t, (*s.Text())[0][0], c, "testchar slice "+string(c))
	c2, _, _ := v.CursorChar(v.slice, y, x)
	assert.Equal(t, *c2, c, "testchar cursorchar "+string(c))
	// Test mock term matches after rendering
	term := Ed.term.(*core.MockTerm)
	v.Render()
	tc := term.CharAt(y+v.y1+2, x+v.x1+2)
	assert.Equal(t, tc, c, fmt.Sprintf("term.charAt %s, y:%d ,x:%d"+string(c), y, x))
}

// TODO: test scrolling etc...{
func TestViewScrolling(t *testing.T) {
}

func TestUndo(t *testing.T) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	Ed.InsertViewSmart(v)
	v.SetBounds(0, 0, 100, 1000)
	v.slice = v.backend.Slice(0, 0, 100, 1000)
	s := core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "")
	v.InsertCur("abcd")
	v.InsertCur("123")
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd123")
	// insert in middle
	v.Insert(0, 2, "xyz", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abxyzcd123")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd123")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "")
	// Redo
	actions.Redo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd")
	// Insert several lines
	v.InsertCur("\n\t123\nXYX")
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t123\nXYX")
	// Line feed at eol -> indentation
	v.Insert(1, 4, "\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t123\n\t\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t123\nXYX")
	// \n at end of line
	v.Insert(1, 0, "I am hungry\nLet's go eat !\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\nI am hungry\nLet's go eat !\n\t123\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t123\nXYX")
	// Breaking line in half
	v.Insert(1, 2, "\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t1\n23\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Equal(t, s, "abcd\n\t123\nXYX")

	type state struct {
		x, y int
		txt  string
	}
	states := []state{}
	for i := 0; i != 15; i++ {
		size := 1 + rand.Int63()%int64(15)
		//ln := rand.Int63() % int64(v.LineCount())
		//col := rand.Int63() % (1 + int64(v.LineLen(v.slice, int(ln))))
		txt := core.RandString(int(size))
		//v.Insert(int(ln), int(col), txt, false)
		v.InsertCur(txt)
		s := state{
			txt: core.RunesToString(*v.Slice().Text()),
			x:   v.CurCol(),
			y:   v.CurLine(),
		}
		states = append(states, s)
	}

	for i := 0; i != 15; i++ {
		s := states[len(states)-i-1]
		assert.Equal(t, s.y, v.CurLine(), fmt.Sprintf("cl %d #v", i, s))
		assert.Equal(t, s.x, v.CurCol(), fmt.Sprintf("cc %d %#v", i, s))
		assert.Equal(t, s.txt, core.RunesToString(*v.Slice().Text()))
		actions.Undo(v.Id())
	}
}

// TODO: test term mock
// TODO: save etc ....
