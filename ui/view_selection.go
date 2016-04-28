package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tcolar/goed/core"
)

func (v *View) ClearSelections() {
	v.selections = []core.Selection{}
}

// Text returns the text contained in the selection of the given view
// Note: **NOT** a rectangle but from pt1 to pt2
func (v *View) SelectionText(s *core.Selection) [][]rune {
	cf := s.ColFrom
	ct := s.ColTo
	lt := s.LineTo
	lf := s.LineFrom
	if lt == -1 {
		lt = v.LineCount() - 1
	}
	if lf == lt {
		return *v.backend.Slice(lf, cf, lt, ct).Text()
	}
	// first line
	text := *v.backend.Slice(lf, cf, lf, -1).Text()
	for l := lf + 1; l < lt; l++ {
		// middle
		text = append(text, *v.backend.Slice(l, 0, l, -1).Text()...)
	}
	// last line
	text = append(text, *v.backend.Slice(lt, 0, lt, ct).Text()...)
	return text
}

// Selected returns whether the text at line, col is current selected
// also returns the matching selection, if any.
func (v *View) Selected(col, line int) (bool, *core.Selection) {
	for _, s := range v.selections {
		if line < s.LineFrom || line > s.LineTo {
			continue
		} else if line > s.LineFrom && line < s.LineTo {
			return true, &s
		} else if s.LineFrom == s.LineTo {
			return col >= s.ColFrom && col <= s.ColTo, &s
		} else if line == s.LineFrom && col >= s.ColFrom {
			return true, &s
		} else if line == s.LineTo && col <= s.ColTo {
			return true, &s
		}
	}
	return false, nil
}

func (v *View) Cut() {
	if len(v.selections) == 0 {
		v.SelectLine(v.CurLine())
	}
	v.SelectionCopy(&v.selections[0])
	v.SelectionDelete(&v.selections[0])
}

func (v *View) Copy() {
	if len(v.selections) == 0 {
		v.SelectLine(v.CurLine())
	}
	v.SelectionCopy(&v.selections[0])
}

func (v *View) SelectionCopy(s *core.Selection) {
	t := v.SelectionText(s)
	text := core.RunesToString(t)
	if len(text) == 0 {
		return
	}
	// when copying full lines, add a "\n" at the end of the copy
	if s.ColTo == -1 {
		text += "\n"
	}
	core.Ed.SetStatus(fmt.Sprintf("Copied %d lines to clipboard.", len(t)))
	core.ClipboardWrite(text)
}

func (v *View) SelectionDelete(s *core.Selection) {
	colTo := s.ColTo
	if colTo == -1 {
		colTo = v.LineLen(v.slice, s.LineTo)
	}
	v.Delete(s.LineFrom, s.ColFrom, s.LineTo, colTo, true)
}

func (v *View) Paste() {
	text, err := core.ClipboardRead()
	if err != nil {
		core.Ed.SetStatusErr(err.Error())
		return
	}
	if len(v.selections) > 0 {
		v.DeleteCur()
	}
	ln, col := v.CurTextPos()
	v.Insert(ln, col, text, true)
}

var locationRegexp = regexp.MustCompile(`([^"\s(){}[\]<>,?|+=&^%#@!;':\x1B]+)(:\d+)?(:\d+)?`)

// Try to select a "location" from the given position
// a location is a path with possibly a line number and maybe a column number as well
func (v *View) ExpandSelectionPath(line, col int) *core.Selection {
	l := v.Line(v.slice, line)
	ln := string(l)
	// Note: Indexes taken and returned by FindAllStringIndex are in BYTES, not runes
	// not very inutitive to say the least
	c := core.RunesLen(l[:col])
	matches := locationRegexp.FindAllStringIndex(ln, -1)
	var best []int
	// Find the "narrowest" match around the cursor
	for _, s := range matches {
		if s[0] <= c && s[1] >= c {
			if best == nil || s[1]-s[0] < best[1]-best[0] {
				best = s
			}
		}
	}
	if best == nil {
		return nil
	}
	// convert byte indexes back to rune count
	r1 := utf8.RuneCountInString(ln[:best[0]])
	r2 := utf8.RuneCountInString(ln[:best[1]])
	return core.NewSelection(line, r1, line, r2-1)
}

// Try to select the longest "word" from current position.
func (v *View) ExpandSelectionWord(line, col int) *core.Selection {
	l := v.Line(v.slice, line)
	c := v.LineRunesTo(v.slice, line, col)
	if c < 0 || c >= len(l) {
		return nil
	}
	c1, c2 := c, c
	for ; c1 >= 0 && isWordRune(l[c1]); c1-- {
	}
	c1++
	for ; c2 < len(l) && isWordRune(l[c2]); c2++ {
	}
	c2--
	if c1 >= c2 {
		return nil
	}
	return core.NewSelection(line, c1, line, c2)
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func (v *View) SelectAll() {
	lastLn := v.LineCount() - 1
	slice := v.Backend().Slice(lastLn, 0, lastLn, -1)
	s := core.NewSelection(0, 0, lastLn, v.LineLen(slice, lastLn)-1)
	v.selections = []core.Selection{
		*s,
	}
}

// Select the whole given line
func (v *View) SelectLine(line int) {
	s := core.NewSelection(line, 0, line, -1)
	v.selections = []core.Selection{
		*s,
	}
}

// Select a word at the given location (if any)
func (v *View) SelectWord(line, col int) {
	s := v.ExpandSelectionWord(line, col)
	if s != nil {
		v.selections = []core.Selection{
			*s,
		}
	}
}

// Parses a selection into a location (file, line, col)
func (v *View) SelectionToLoc(sel *core.Selection) (loc string, line, col int) {
	sub := locationRegexp.FindAllStringSubmatch(core.RunesToString(v.SelectionText(sel)), 1)
	if len(sub) == 0 {
		return
	}
	s := sub[0]
	if len(s) >= 1 {
		loc = s[1]
	}
	if len(s[2]) > 0 {
		line, _ = strconv.Atoi(s[2][1:])
	}
	if len(s[3]) > 0 {
		col, _ = strconv.Atoi(s[3][1:])
	}
	return loc, line, col
}

// Stretch a selection toward a new position
func (v *View) StretchSelection(prevl, prevc, l, c int) {
	if len(v.selections) != 0 {
		s := &v.selections[0]
		s.LineFrom, s.LineTo, s.ColFrom, s.ColTo = prevl, l, prevc, c
		s.Normalize()
	} else {
		s := core.NewSelection(prevl, prevc, l, c)
		s.Normalize()
		v.selections = []core.Selection{*s}
	}
}

// Open what's selected or under the cursor
// if newView is true then open in a new view, otherwise
// replace content of v
func (v *View) OpenSelection(newView bool) {
	ed := core.Ed.(*Editor)
	newView = newView || v.Dirty()
	if len(v.selections) == 0 {
		ln, col := v.CurTextPos()
		selection := v.ExpandSelectionPath(ln, col)
		if selection == nil {
			ed.SetStatusErr("Could not expand location from cursor location.")
			return
		}
		v.selections = []core.Selection{*selection}
	}
	loc, line, col := v.SelectionToLoc(&v.selections[0])
	line-- // we use 0 indexes in views
	col--
	isDir := false
	loc, isDir = core.LookupLocation(v.WorkDir(), strings.TrimSpace(loc))
	vid := ed.ViewByLoc(loc)
	if vid >= 0 {
		vv := ed.ViewById(vid)
		if vv != nil {
			// Already open
			vv.SetCursorPos(line, col)
			ed.ViewActivate(vv.Id())
			return
		}
	}
	v2 := ed.NewView(loc)
	if _, err := ed.Open(loc, v2.Id(), v.WorkDir(), false); err != nil {
		ed.SetStatusErr(err.Error())
		return
	}
	if newView {
		if isDir {
			ed.AddDirViewSmart(v2)
		} else {
			ed.InsertViewSmart(v2)
		}
	} else {
		ed.ReplaceView(v, v2)
	}
	v2.SetCursorPos(line, col)
	ed.ViewActivate(v2.Id())
}
