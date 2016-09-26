package core

// Slice represents a "matrix" of text (runes)
// coordinates are of a rectangle (unlike a selection whihc is ptA to ptB)
type Slice struct {
	text           [][]rune
	R1, C1, R2, C2 int //bounds
}

func (s *Slice) Text() *[][]rune {
	return &s.text
}

func NewSlice(r1, c1, r2, c2 int, text [][]rune) *Slice {
	s := &Slice{
		R1:   r1,
		C1:   c1,
		R2:   r2,
		C2:   c2,
		text: text,
	}
	s.Normalize()
	return s
}

func (s *Slice) Normalize() {
	if s.R2 != -1 && s.R1 > s.R2 {
		s.R1, s.R2 = s.R2, s.R1
	}
	if s.C2 != -1 && s.C1 > s.C2 {
		s.C1, s.C2 = s.C2, s.C1
	}
}

func (s *Slice) ContainsLine(lnIndex int) bool {
	if s.R1 == 0 && s.R2 == 0 && s.C1 == 0 && s.C2 == 0 {
		return false
	}
	return lnIndex >= s.R1 && lnIndex <= s.R2
}
