package main

import (
	"fmt"

	"github.com/atotto/clipboard"
)

type Selection struct {
	LineFrom, ColFrom int // selection start point
	LineTo, ColTo     int // selection end point
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
	v.InsertLines(Ed.StringToRunes([]byte(text)))
}
