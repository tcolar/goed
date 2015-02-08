package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/atotto/clipboard"
)

type Selection struct {
	LineFrom, ColFrom int // selection start point
	LineTo, ColTo     int // selection end point (colto=-1 means whole lines)
}

func (s Selection) String() string {
	return fmt.Sprintf("(%d,%d)-(%d,%d)", s.LineFrom, s.ColFrom, s.LineTo, s.ColTo)
}

// Text returns the text contained in the selection of the given view
func (s Selection) Text(v *View) [][]rune {
	runes := [][]rune{}
	cf := s.ColFrom
	ct := s.ColTo + 1
	lt := s.LineTo
	lf := s.LineFrom
	if v.LineCount() < s.LineFrom {
		return runes
	}
	if lt > v.LineCount() {
		lt = v.LineCount()
	}
	line := v.Line(s.LineFrom)
	line2 := v.Line(s.LineTo)
	if cf > len(line) {
		cf = len(line)
	}
	if ct > len(line2) {
		ct = len(line2)
	}
	if s.LineFrom == s.LineTo {
		runes = append(runes, line[cf:ct])
		return runes
	}
	runes = append(runes, line[cf:])
	for i := lf + 1; i < lt; i++ {
		runes = append(runes, v.Line(i))
	}
	runes = append(runes, line2[:ct])
	return runes
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

func (v *View) Copy(s Selection) {
	t := s.Text(v)
	Ed.SetStatus(fmt.Sprintf("Copied %d lines to clipboard.", len(t)))
	clipboard.WriteAll(Ed.RunesToString(t))
}

func (v *View) Paste() {
	text, err := clipboard.ReadAll()
	if err != nil {
		Ed.SetStatusErr(err.Error())
		return
	}
	v.Insert(text)
}

var locationRegexp = regexp.MustCompile(`([^"\s(){}[\]<>,?|+=&^%#@!;':]+)(:\d+)?(:\d+)?`)

// Try to select a "location" from the given position
// a location is a path with possibly a line number and maybe a column number as well
func (v *View) PathSelection(line, col int) *Selection {
	ln := string(v.Line(line))
	c := v.lineRunesTo(line, col)
	matches := locationRegexp.FindAllStringIndex(string(ln), -1)
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
	return &Selection{
		LineFrom: line,
		ColFrom:  best[0],
		LineTo:   line,
		ColTo:    best[1] - 1,
	}
}

// Parses a selection into a location (file, line, col)
func (v *View) selToLoc(selection Selection) (loc string, line, col int) {
	sub := locationRegexp.FindAllStringSubmatch(Ed.RunesToString(selection.Text(v)), 1)
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
