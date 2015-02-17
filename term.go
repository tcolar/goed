package main

import (
	"fmt"

	"github.com/tcolar/termbox-go"
)

// Terminal interface
type Term interface {
	Close()
	Clear(fg, bg uint16)
	Char(x, y int, c rune)
	Flush()
	Init() error
	SetExtendedColors(bool)
	SetCursor(x, y int)
	Size() (int, int)
	SetMouseMode(termbox.MouseMode)
	SetInputMode(termbox.InputMode)
}

// ==================== Termbox impl ===========================

// Real Terinal implementation using termbox
type TermBox struct {
}

func NewTermBox() *TermBox {
	return &TermBox{}
}

func (t *TermBox) Init() error {
	return termbox.Init()
}

func (t *TermBox) Clear(fg, bg uint16) {
	termbox.Clear(termbox.Attribute(fg), termbox.Attribute(bg))
}

func (t *TermBox) Close() {
	termbox.Close()
}

func (t *TermBox) Flush() {
	termbox.Flush()
}

func (t *TermBox) SetExtendedColors(b bool) {
	termbox.SetExtendedColors(b)
}

func (t *TermBox) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (t *TermBox) Char(x, y int, c rune) {
	termbox.SetCell(x, y, c, termbox.Attribute(Ed.Fg.uint16), termbox.Attribute(Ed.Bg.uint16))
}

func (t *TermBox) Size() (int, int) {
	return termbox.Size()
}

func (t *TermBox) SetMouseMode(m termbox.MouseMode) {
	termbox.SetMouseMode(m)
}

func (t *TermBox) SetInputMode(m termbox.InputMode) {
	termbox.SetInputMode(m)
}

// ==================== Mock impl ===========================

// Mock  Terminal implementation for testing
type MockTerm struct {
	w, h             int
	cursorX, cursorY int
	text             [25][50]rune
}

func newMockTerm() *MockTerm {
	return &MockTerm{
		w:    50,
		h:    25,
		text: [25][50]rune{},
	}
}

func (t *MockTerm) Init() error {
	return nil
}

func (t *MockTerm) Close() {
}

func (t *MockTerm) Clear(fg, bg uint16) {
	t.text = [25][50]rune{}
}

func (t *MockTerm) Flush() {
}

func (t *MockTerm) SetExtendedColors(b bool) {
}

func (t *MockTerm) SetCursor(x, y int) {
	t.cursorX, t.cursorY = x, y
}

func (t *MockTerm) Char(x, y int, c rune) {
	if x < t.w && y < t.h {
		t.text[y][x] = c
	}
}

func (t *MockTerm) Size() (int, int) {
	return t.w, t.h
}

func (t *MockTerm) SetMouseMode(m termbox.MouseMode) {
}

func (t *MockTerm) SetInputMode(m termbox.InputMode) {
}

// for testing
func (t *MockTerm) charAt(x, y int) rune {
	return t.text[y][x]
}

//=================== Utilities =============================

// TermFB sets the "active" forground and backgrounds colors.
func (e *Editor) TermFB(fg, bg Style) {
	e.Fg = fg
	e.Bg = bg
}

func (e *Editor) TermChar(x, y int, c rune) {
	e.term.Char(x, y, c)
}

// TermStr draws an horizonttal string to the terminal
func (e *Editor) TermStr(x, y int, s string) {
	for _, c := range s {
		e.term.Char(x, y, c)
		x++
	}
}

// TermStrv draws a vertical string to the terminal
func (e *Editor) TermStrv(x, y int, s string) {
	for _, c := range s {
		e.term.Char(x, y, c)
		y++
	}
}

// TermFill fills an area of the terminal
func (e *Editor) TermFill(c rune, x1, y1, x2, y2 int) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			e.term.Char(x, y, c)
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
