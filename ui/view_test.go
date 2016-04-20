package ui

import (
	"math/rand"
	"strings"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (us *UiSuite) TestView(t *C) {
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(0, 0, 25, 40)

	_, err = Ed.Open("../test_data/file1.txt", v.Id(), "", false)
	assert.Nil(t, err)

	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	assertCursor(t, v, 0, 0, 0, 0)
	assert.True(t, strings.HasSuffix(v.backend.SrcLoc(), "/test_data/file1.txt"))
	assert.True(t, strings.HasSuffix(v.workDir, "/test_data"))
	assert.False(t, v.Dirty())
	assert.Eq(t, v.Title(), "file1.txt")
	assert.Eq(t, v.LineCount(), 12)
	assert.Eq(t, v.LastViewLine(), 25-3)
	assert.Eq(t, v.LastViewCol(), 40-3)
	assert.Eq(t, v.lineCols(v.slice, 0), 10)
	assert.Eq(t, string(v.Line(v.slice, 0)), "1234567890")
	assert.Eq(t, v.LineLen(v.slice, 3), 26)
	assert.Eq(t, v.lineColsTo(v.slice, 0, 4), 4)
	assert.Eq(t, v.lineColsTo(v.slice, 9, 4), 10) //\t\t a
	assert.Eq(t, v.LineRunesTo(v.slice, 0, 4), 4)
	assert.Eq(t, v.LineRunesTo(v.slice, 9, 10), 4)
	x := v.LineRunesTo(v.slice, 0, 4)
	assert.Eq(t, x, 4)
	x = v.LineRunesTo(v.slice, 9, 10)
	assert.Eq(t, x, 4)
	c, y, x := v.CursorChar(v.slice, 3, 3)
	assert.Eq(t, *c, 'D')
	assert.Eq(t, x, 3)
	assert.Eq(t, y, 3)
	assert.Eq(t, v.runeSize('a'), 1)
	assert.Eq(t, v.runeSize('\t'), tabSize)
	assert.Eq(t, v.strSize("abc"), 3)
	assert.Eq(t, v.strSize("a\tb\tc"), 3+2*tabSize)
	v.MoveCursor(0, 0)
	assertCursor(t, v, 0, 0, 0, 0)
	v.MoveCursor(0, 5)
	assertCursor(t, v, 0, 5, 0, 0)
	c, y, x = v.CurChar()
	assert.Eq(t, x, 5)
	assert.Eq(t, y, 0)
	assert.Eq(t, *c, '6')
	v.MoveCursor(0, -3)
	assertCursor(t, v, 0, 2, 0, 0)
	v.MoveCursor(3, 2)
	assertCursor(t, v, 3, 4, 0, 0)
	v.MoveCursor(-1, 2)
	assert.Eq(t, v.CurCol(), 6)
	assert.Eq(t, v.CurLine(), 2)
	assertCursor(t, v, 2, 6, 0, 0)
	v.MoveCursor(-1, 2)
	// Note: x=0 because line "1" is blank
	assertCursor(t, v, 1, 0, 0, 0)
	v.MoveCursor(-10, -10)
	assertCursor(t, v, 0, 0, 0, 0)
	v.MoveCursor(100, 100)
	assertCursor(t, v, 11, 36, 0, 0)
	v.MoveCursor(-100, -100)
	assertCursor(t, v, 0, 0, 0, 0)
	v.MoveCursorRoll(0, 10)
	assertCursor(t, v, 0, 10, 0, 0)
	v.MoveCursorRoll(0, 1)
	assertCursor(t, v, 1, 0, 0, 0)
	v.MoveCursorRoll(0, -2)
	assertCursor(t, v, 0, 10, 0, 0)
}

func assertCursor(t *C, v *View, y, x, offsetY, offsetX int) {
	assert.Eq(t, v.CursorX, x)
	assert.Eq(t, v.CursorY, y)
	assert.Eq(t, v.offx, offsetX)
	assert.Eq(t, v.offy, offsetY)
}

func (us *UiSuite) TestViewSelections(t *C) {
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/file1.txt", v.Id(), "", false)
	assert.Nil(t, err)
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	s := core.NewSelection(0, 0, 0, 0)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "1")
	s = core.NewSelection(2, 1, 3, 7)
	v.selections = append(v.selections, *s)
	assert.Eq(t, s.String(), "2 1 3 7")
	text := v.SelectionText(s)
	assert.Eq(t, len(text), 2)
	assert.Eq(t, len(text[0]), 25)
	assert.Eq(t, len(text[1]), 8)
	assert.Eq(t, string(text[0]), "bcdefghijklmnopqrstuvwxyz")
	assert.Eq(t, string(text[1]), "ABCDEFGH")
	b, sel := v.Selected(3, 0)
	assert.False(t, b)
	assert.Nil(t, sel)
	b, sel = v.Selected(0, 2)
	assert.False(t, b)
	assert.Nil(t, sel)
	b, sel = v.Selected(3, 4)
	assert.False(t, b)
	assert.Nil(t, sel)
	b, sel = v.Selected(2, 2)
	assert.True(t, b)
	assert.Eq(t, sel.String(), s.String())
	v.SelectionCopy(s)
	cb, _ := core.Ed.Clipboard.ReadAll()
	assert.Eq(t, cb, "bcdefghijklmnopqrstuvwxyz\nABCDEFGH")
	s = v.ExpandSelectionPath(0, 2)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "1234567890")
	s = v.ExpandSelectionPath(9, 2)
	assert.Eq(t, s.String(), "9 2 9 4")
	s = v.ExpandSelectionPath(10, 0)
	assert.Eq(t, s.String(), "10 0 10 2")
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa")
	s = v.ExpandSelectionPath(10, 5)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa.go")
	s = v.ExpandSelectionPath(10, 21)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "/tmp/aaa.go")
	s = v.ExpandSelectionPath(10, 26)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa.go:23")
	s = v.ExpandSelectionPath(10, 38)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "/tmp/aaa.go:23:7")
	loc, ln, col := v.SelectionToLoc(s)
	assert.Eq(t, loc, "/tmp/aaa.go")
	assert.Eq(t, ln, 23)
	assert.Eq(t, col, 7)
	s = v.ExpandSelectionWord(10, 0)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa")
	s = v.ExpandSelectionWord(10, 1)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa")
	s = v.ExpandSelectionWord(10, 2)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "aaa")
	s = v.ExpandSelectionWord(10, 3)
	assert.Nil(t, s)
	s = v.ExpandSelectionWord(10, 12)
	assert.Eq(t, core.RunesToString(v.SelectionText(s)), "tmp")
}

func (us *UiSuite) TestViewEdition(t *C) {
	// Note: more tests done directly on backemds
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/empty.txt", v.Id(), "", false)
	assert.Nil(t, err)
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

func (us *UiSuite) TestDelete(t *C) {
	// Test some edge cases
	var err error
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	v.SetBounds(5, 5, 30, 140)
	_, err = Ed.Open("../test_data/empty.txt", v.Id(), "", false)
	assert.Nil(t, err)
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())

	v.Insert(0, 0, "ab\ncd\nef", true)
	s := core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "ab\ncd\nef")
	// backspace edge cases
	v.CursorX = 0
	v.CursorY = 1
	v.Backspace()
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abcd\nef")
	v.CursorX = 0
	v.CursorY = 0
	v.Backspace()
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abcd\nef")
	// delete edge cases
	v.CursorX = 5
	v.CursorY = 0
	v.DeleteCur()
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abcdef")
	v.CursorX = 7
	v.CursorY = 0
	v.DeleteCur()
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abcdef")
	v.Delete(0, 0, 0, 5, true)
	v.SyncSlice()
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "")
	v.DeleteCur() // check no panic
	v.Backspace() // check no panic
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "")

	_, err = Ed.Open("../test_data/no_eol.txt", v.Id(), "", false)
	assert.Nil(t, err)
	v.slice = v.backend.Slice(v.offy, v.offx, v.offy+v.LastViewLine(), v.offx+v.LastViewCol())
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abc\n123\nz")
	v.MoveCursor(-100, -100)
	assert.Eq(t, v.CurLine(), 0)
	assert.Eq(t, v.CurCol(), 0)
	v.MoveCursor(2, 0)
	assert.Eq(t, v.CurLine(), 2)
	assert.Eq(t, v.CurCol(), 0)
	v.MoveCursorRoll(0, 1)
	assert.Eq(t, v.CursorY, 2)
	assert.Eq(t, v.CursorX, 1)
	v.InsertCur("Z")
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abc\n123\nzZ")
	v.MoveCursorRoll(0, -1)
	v.InsertCur("A")
	v.InsertCur("B")
	s = core.RunesToString(*v.slice.Text())
	assert.Eq(t, s, "abc\n123\nzABZ")
}

func testChar(t *C, v *View, y, x int, c rune) {
	Ed := core.Ed.(*Editor)
	s := v.backend.Slice(y, x, y, x)
	if (*s.Text())[0][0] != c {
		panic(c)
	}
	assert.Eq(t, (*s.Text())[0][0], c)
	c2, _, _ := v.CursorChar(v.slice, y, x)
	assert.Eq(t, *c2, c)
	// Test mock term matches after rendering
	term := Ed.term.(*core.MockTerm)
	v.Render()
	tc := term.CharAt(y+v.y1+2, x+v.x1+2)
	assert.Eq(t, tc, c)
}

// TODO: test scrolling etc...{
func (us *UiSuite) TestViewScrolling(t *C) {
}

func (us *UiSuite) TestUndo(t *C) {
	Ed := core.Ed.(*Editor)
	v := Ed.NewView("")
	Ed.InsertViewSmart(v)
	v.SetBounds(0, 0, 100, 1000)
	v.slice = v.backend.Slice(0, 0, 100, 1000)
	s := core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "")
	v.InsertCur("abcd")
	v.InsertCur("123")
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd123")
	// insert in middle
	v.Insert(0, 2, "xyz", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abxyzcd123")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd123")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "")
	// Redo
	actions.Redo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd")
	// Insert several lines
	v.InsertCur("\n\t123\nXYX")
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t123\nXYX")
	// Line feed at eol -> indentation
	v.Insert(1, 4, "\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t123\n\t\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t123\nXYX")
	// \n at end of line
	v.Insert(1, 0, "I am hungry\nLet's go eat !\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\nI am hungry\nLet's go eat !\n\t123\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t123\nXYX")
	// Breaking line in half
	v.Insert(1, 2, "\n", true)
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t1\n23\nXYX")
	actions.Undo(v.Id())
	s = core.RunesToString(*v.Slice().Text())
	assert.Eq(t, s, "abcd\n\t123\nXYX")

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
		assert.Eq(t, s.y, v.CurLine())
		assert.Eq(t, s.x, v.CurCol())
		assert.Eq(t, s.txt, core.RunesToString(*v.Slice().Text()))
		actions.Undo(v.Id())
	}
}

// TODO: test term mock
// TODO: save etc ....
