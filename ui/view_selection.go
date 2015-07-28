package ui

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	"github.com/atotto/clipboard"
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

func (v *View) Copy() {
	if len(v.selections) == 0 {
		v.SelectLine(v.CurLine())
	}
	v.SelectionCopy(&v.selections[0])
}

func (v *View) Delete() {
	if len(v.selections) == 0 {
		v.SelectLine(v.CurLine())
	}
	v.SelectionDelete(&v.selections[0])
}

func (v *View) SelectionCopy(s *core.Selection) {
	t := v.SelectionText(s)
	core.Ed.SetStatus(fmt.Sprintf("Copied %d lines to clipboard.", len(t)))
	clipboard.WriteAll(core.RunesToString(t))
}

func (v *View) SelectionDelete(s *core.Selection) {
	v.delete(s.LineFrom, s.ColFrom, s.LineTo, s.ColTo)
}

func (v *View) Paste() {
	text, err := clipboard.ReadAll()
	if err != nil {
		core.Ed.SetStatusErr(err.Error())
		return
	}
	if len(v.selections) > 0 {
		v.DeleteCur()
	}
	_, y, x := v.CurChar()
	v.Insert(y, x, text)
}

var locationRegexp = regexp.MustCompile(`([^"\s(){}[\]<>,?|+=&^%#@!;':]+)(:\d+)?(:\d+)?`)

// Try to select a "location" from the given position
// a location is a path with possibly a line number and maybe a column number as well
func (v *View) ExpandSelectionPath(line, col int) *core.Selection {
	l := v.Line(v.slice, line)
	ln := string(l)
	slice := core.NewSlice(0, 0, 0, len(l), [][]rune{l})
	c := v.LineRunesTo(slice, 0, col)
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
	return core.NewSelection(line, best[0], line, best[1]-1)
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

// Select the whole given line
func (v *View) SelectLine(line int) {
	s := core.NewSelection(line, 0, line, v.LineLen(v.slice, line))
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
	if len(v.selections) == 0 {
		s := *core.NewSelection(prevl, prevc, l, c)
		v.selections = []core.Selection{
			s,
		}
	} else {
		s := v.selections[0]
		if s.LineTo == prevl && s.ColTo == prevc {
			s.LineTo, s.ColTo = l, c
		} else {
			s.LineFrom, s.ColFrom = l, c
		}
		s.Normalize()
		v.selections[0] = s
	}
}

// Open what's selected or under the cursor
// if newView is true then open in a new view, otherwise
// replace content of v
func (v *View) OpenSelection(newView bool) {
	ed := core.Ed.(*Editor)
	newView = newView || v.Dirty()
	if len(v.selections) == 0 {
		selection := v.ExpandSelectionPath(v.CurLine(), v.CurCol())
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
	loc, isDir = core.LookupLocation(v.WorkDir(), loc)
	vv := ed.ViewByLoc(loc)
	if vv != nil {
		// Already open
		ed.ActivateView(vv, line, col)
		return
	}
	v2 := ed.NewView(loc)
	if _, err := ed.Open(loc, v2, v.WorkDir(), false); err != nil {
		ed.SetStatusErr(err.Error())
		return
	}
	if newView {
		if isDir {
			ed.InsertView(v2, v, 0.5)
		} else {
			ed.InsertViewSmart(v2)
		}
	} else {
		ed.ReplaceView(v, v2)
	}
	ed.ActivateView(v2, line, col)
}
