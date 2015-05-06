package core

import "fmt"

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

// Return the selection in the form "row1 col1 row1 col2"
func (s Selection) String() string {
	return fmt.Sprintf("%d %d %d %d", s.LineFrom, s.ColFrom, s.LineTo, s.ColTo)
}
