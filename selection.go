package main

type Selection struct {
	LineFrom, ColFrom int // selection start point
	LineTo, ColTo     int // selection end point
}

// Selected returns whether the text at line, col is current selected
// alos returns the matching selection, if any.
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
