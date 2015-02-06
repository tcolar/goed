package main

import (
	"fmt"

	"github.com/tcolar/termbox-go"
)

func (e *Editor) Size() (int, int) {
	return termbox.Size()
}

func (e *Editor) FB(fg, bg Style) {
	e.Fg = fg
	e.Bg = bg
}

func (e *Editor) Char(x, y int, c rune) {
	termbox.SetCell(x, y, c, termbox.Attribute(e.Fg.uint16), termbox.Attribute(e.Bg.uint16))
}

func (e *Editor) Str(x, y int, s string) {
	for _, c := range s {
		e.Char(x, y, c)
		x++
	}
}

func (e *Editor) Strv(x, y int, s string) {
	for _, c := range s {
		e.Char(x, y, c)
		y++
	}
}

func (e *Editor) Fill(c rune, x1, y1, x2, y2 int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			e.Char(x, y, c)
		}
	}
}

// Print colors to terminal to try it.
func testTerm() {
	fmt.Printf("Standard Colors (16):\n Plain      : ")
	for i := 0; i != 16; i++ {
		fmt.Printf("\033[3%dm%02X ", i, i)
	}
	fmt.Printf("\n Bold       : ")
	for i := 0; i != 16; i++ {
		fmt.Printf("\033[1;3%dm%02X ", i, i)
	}
	fmt.Printf("\033[0m\n Underlined : ")
	for i := 0; i != 16; i++ {
		fmt.Printf("\033[4;3%dm%02X ", i, i)
	}
	fmt.Println("\033[0m\n\nExtended Colors (256):")
	for i := 0; i != 256; i++ {
		fmt.Printf("\033[0;38;5;%dm%02X ", i, i)
	}
	fmt.Println("\n\nAscii Chars: a A 6 ¼ Ø \nUnicode chars: \u0E5B  ಠﭛಠ")
}

func detectColors() int {
	// TBD
	return 256
}
