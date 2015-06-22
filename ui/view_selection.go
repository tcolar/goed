package ui

import (
	"fmt"
	"regexp"
	"strconv"

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
		text = append(text, *v.backend.Slice(l, 1, l, -1).Text()...)
	}
	// last line
	text = append(text, *v.backend.Slice(lt, 1, lt, ct).Text()...)
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

func (v *View) SelectionCopy(s *core.Selection) {
	t := v.SelectionText(s)
	core.Ed.SetStatus(fmt.Sprintf("Copied %d lines to clipboard.", len(t)))
	clipboard.WriteAll(core.RunesToString(t))
}

func (v *View) SelectionDelete(s *core.Selection) {
	v.Delete(s.LineFrom-1, s.ColFrom-1, s.LineTo-1, s.ColTo-1)
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
	_, x, y := v.CurChar()
	v.Insert(y, x, text)
}

var locationRegexp = regexp.MustCompile(`([^"\s(){}[\]<>,?|+=&^%#@!;':]+)(:\d+)?(:\d+)?`)

// Try to select a "location" from the given position
// a location is a path with possibly a line number and maybe a column number as well
func (v *View) PathSelection(line, col int) *core.Selection {
	l := v.Line(v.slice, line-1)
	ln := string(l)
	slice := core.NewSlice(1, 1, 1, len(l)+1, [][]rune{l})
	c := v.lineRunesTo(slice, 0, col)
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
	// TODO: if a path like a go import, try to find that path up from curdir ?
	return core.NewSelection(line, best[0]+1, line, best[1])
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

// Expand a selection toward a new position
func (v *View) ExpandSelection(prevl, prevc, l, c int) {
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
