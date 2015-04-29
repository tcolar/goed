package core

type Slice struct {
	text           [][]rune
	R1, C1, R2, C2 int //bounds
}

func (s *Slice) Text() *[][]rune {
	return &s.text
}

func NewSlice(r1, c1, r2, c2 int, text [][]rune) *Slice {
	return &Slice{
		R1:   r1,
		C1:   c1,
		R2:   r2,
		C2:   c2,
		text: text,
	}
}
