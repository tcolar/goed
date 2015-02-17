package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/atotto/clipboard"
)

// Selection : 1 indexed
type Selection struct {
	LineFrom, ColFrom int // selection start point
	LineTo, ColTo     int // selection end point (colto=-1 means whole lines)
}

func NewSelection(l1, c1, l2, c2 int) *Selection {
	return &Selection{
		LineFrom: l1,
		ColFrom:  c1,
		LineTo:   l2,
		ColTo:    c2,
	}
}

func (v *View) ClearSelections() {
	v.Selections = []Selection{}
}

func (s Selection) String() string {
	return fmt.Sprintf("(%d,%d)-(%d,%d)", s.LineFrom, s.ColFrom, s.LineTo, s.ColTo)
}

// Text returns the text contained in the selection of the given view
// Note: **NOT** a rectangle but from pt1 to pt2
func (s Selection) Text(v *View) [][]rune {
	cf := s.ColFrom
	ct := s.ColTo
	lt := s.LineTo
	lf := s.LineFrom
	if lf == lt {
		return v.backend.Slice(lf, cf, lt, ct).text
	}
	// first line
	text := v.backend.Slice(lf, cf, lf, -1).text
	for l := lf + 1; l < lt; l++ {
		// middle
		text = append(text, v.backend.Slice(l, 1, l, -1).text...)
	}
	// last line
	text = append(text, v.backend.Slice(lt, 1, lt, ct).text...)
	return text
}

// Selected returns whether the text at line, col is current selected
// also returns the matching selection, if any.
func (v *View) Selected(col, line int) (bool, *Selection) {
	for _, s := range v.Selections {
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

func (s Selection) Copy(v *View) {
	t := s.Text(v)
	Ed.SetStatus(fmt.Sprintf("Copied %d lines to clipboard.", len(t)))
	clipboard.WriteAll(Ed.RunesToString(t))
}

func (s Selection) Delete(v *View) {
	v.Delete(s.LineFrom-1, s.ColFrom-1, s.LineTo-1, s.ColTo-1)
}

func (v *View) Paste() {
	text, err := clipboard.ReadAll()
	if err != nil {
		Ed.SetStatusErr(err.Error())
		return
	}
	_, x, y := v.CurChar()
	v.Insert(y, x, text)
}

var locationRegexp = regexp.MustCompile(`([^"\s(){}[\]<>,?|+=&^%#@!;':]+)(:\d+)?(:\d+)?`)

// Try to select a "location" from the given position
// a location is a path with possibly a line number and maybe a column number as well
func (v *View) PathSelection(line, col int) *Selection {
	l := v.Line(v.slice, line-1)
	ln := string(l)
	slice := NewSlice(1, 1, 1, len(l)+1, [][]rune{l})
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
	return NewSelection(line, best[0]+1, line, best[1])
}

// Parses a selection into a location (file, line, col)
func (sel Selection) ToLoc(v *View) (loc string, line, col int) {
	sub := locationRegexp.FindAllStringSubmatch(Ed.RunesToString(sel.Text(v)), 1)
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
