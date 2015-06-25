package core

import "fmt"

// Selection represents some selected text in a view
type Selection struct {
	LineFrom, ColFrom int // selection start point
	LineTo, ColTo     int // selection end point (colto=-1 means whole lines)
}

func NewSelection(l1, c1, l2, c2 int) *Selection {
	s := &Selection{
		LineFrom: l1,
		ColFrom:  c1,
		LineTo:   l2,
		ColTo:    c2,
	}
	s.Normalize()
	return s
}

// String return the selection in the form "line1 col1 line2 col2"
func (s Selection) String() string {
	return fmt.Sprintf("%d %d %d %d", s.LineFrom, s.ColFrom, s.LineTo, s.ColTo)
}

// Normalize the slectin such as l1,c1 is "before" l2, c2
func (s *Selection) Normalize() {
	// Deal with "reversed" selection
	if s.LineFrom == s.LineTo && s.ColFrom > s.ColTo {
		s.ColFrom, s.ColTo = s.ColTo, s.ColFrom
	} else if s.LineFrom > s.LineTo {
		s.LineFrom, s.LineTo = s.LineTo, s.LineFrom
		s.ColFrom, s.ColTo = s.ColTo, s.ColFrom
	}
}
