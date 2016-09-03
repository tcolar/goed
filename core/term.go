package core

import "fmt"

// Terminal interface
type Term interface {
	Close()
	Clear(fg, bg Style)
	Char(y, x int, c rune, fg, bg Style)
	Flush()
	Init() error
	Listen()
	SetExtendedColors(bool)
	SetCursor(y, x int)
	Size() (y, x int)
}

// ==================== Mock impl ===========================

// Mock  Terminal implementation for testing
type MockTerm struct {
	w, h             int
	cursorX, cursorY int
	text             [25][50]rune
}

func NewMockTerm() *MockTerm {
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

func (t *MockTerm) Clear(fg, bg Style) {
	t.text = [25][50]rune{}
}

func (t *MockTerm) Flush() {
}

func (t *MockTerm) Listen() {
}

func (t *MockTerm) SetExtendedColors(b bool) {
}

func (t *MockTerm) SetCursor(y, x int) {
	t.cursorX, t.cursorY = x, y
}

func (t *MockTerm) Char(y, x int, c rune, fg, bg Style) {
	if x < 0 || y < 0 {
		return
	}
	if x < t.w && y < t.h {
		t.text[y][x] = c
	}
}

func (t *MockTerm) Size() (h, w int) {
	return t.h, t.w
}

// for testing
func (t *MockTerm) CharAt(y, x int) rune {
	return t.text[y][x]
}

//=================== Utilities =============================

// Print colors to terminal to try it.
func TermColors() {
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

func DetectColors() int {
	// TBD
	return 256
}
