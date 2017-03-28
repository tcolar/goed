package ui

import "github.com/tcolar/goed/ui/style"

// TermFB sets the "active" forground and backgrounds colors.
func (e *Editor) TermFB(fg, bg style.Style) {
	e.Fg = fg
	e.Bg = bg
}

func (e *Editor) TermChar(y, x int, c rune) {
	e.term.Char(y, x, c, e.Fg, e.Bg)
}

// TermStr draws an horizonttal string to the terminal
func (e *Editor) TermStr(y, x int, s string) {
	for _, c := range s {
		e.term.Char(y, x, c, e.Fg, e.Bg)
		x++
	}
}

// TermStrv draws a vertical string to the terminal
func (e *Editor) TermStrv(y, x int, s string) {
	for _, c := range s {
		e.term.Char(y, x, c, e.Fg, e.Bg)
		y++
	}
}

// TermFill fills an area of the terminal
func (e *Editor) TermFill(c rune, y1, x1, y2, x2 int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			e.term.Char(y, x, c, e.Fg, e.Bg)
		}
	}
}
