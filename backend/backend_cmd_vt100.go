package backend

import (
	"log"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui/style"
)

// TODO : Handle VT100 codes
// http://www.ccs.neu.edu/research/gpc/VonaUtils/vona/terminal/VT100_Escape_Codes.html
// http://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-VT100-Mode
func (b *backendAppender) vt100(data []byte) (int, error) {
	// handle VT100 codes
	from := 0
	log.Printf("VT100 raw data: %v %q\n", data, string(data))
	for i := 0; i < len(data); i++ {
		if b.consumeVt100(data, from, &i) {
			from = i
			i--
		}
	}
	// flush leftover
	err := b.flush(data[from:len(data)])
	return len(data), err
}

func (b *backendAppender) setCol(col int) {
	if col < 0 {
		col = 0
	}
	b.backend.MemBackend.lock.Lock()
	defer b.backend.MemBackend.lock.Unlock()
	b.col = col
}

func (b *backendAppender) setLine(line int) {
	if line < 0 {
		line = 0
	}
	b.backend.MemBackend.lock.Lock()
	defer b.backend.MemBackend.lock.Unlock()
	b.line = line
}

func (b *backendAppender) consumeVt100(data []byte, from int, i *int) bool {
	t := core.Ed.Theme()

	start := *i
	// bell
	if b.consume(data, i, 7) {
		b.flush(data[from:start])
		actions.Ar.EdSetStatusErr("Beep !!")
		return true
	}
	*i = start
	// "backspace" (move left)
	if b.consume(data, i, 8) {
		b.flush(data[from:start])
		b.setCol(b.col - 1)
		return true
	}
	// horizontal tab
	if b.consume(data, i, 9) {
		spaces := []byte{}
		b.flush(data[from:start])
		for i := 0; i < 8-(b.col%8); i++ {
			spaces = append(spaces, ' ')
		}
		b.flush(spaces)
		return true
	}
	// delete
	*i = start
	if b.consume(data, i, 127) {
		b.flush(data[from:start])
		actions.Ar.ViewDeleteCur(b.viewId)
		return true
	}
	*i = start // \r\n
	if b.consume(data, i, 13) && b.consume(data, i, 10) {
		b.flush(data[from:start])
		b.setLine(b.line + 1)
		b.setCol(0)
		return true
	}
	*i = start
	if b.consume(data, i, 10) { // \n
		b.flush(data[from:start])
		b.setLine(b.line + 1)
		b.setCol(0)
		return true
	}
	*i = start
	if b.consume(data, i, 13) { // \r
		b.flush(data[from:start])
		//b.line++
		b.setCol(0)
		return true
	}
	// Other unhandled control characters, skip them
	*i = start
	if data[*i] != 27 && data[*i] < 32 {
		b.flush(data[from:start])
		*i++
		return true
	}

	// various ignored escape sequences
	*i = start
	if b.consume(data, i, 27) {
		if b.consume(data, i, '=') || b.consume(data, i, '>') ||
			(b.consume(data, i, '(') && b.consume(data, i, 'B')) ||
			(b.consume(data, i, ')') && b.consume(data, i, '0')) {
			return true
		}
	}

	// set title (xterm-ish)
	// http://www.xfree86.org/4.7.0/ctlseqs.html
	*i = start
	if b.consume(data, i, 27) && b.consume(data, i, 93) &&
		// 0 is (icon&window title), 2 is window only
		(b.consume(data, i, '0') || b.consume(data, i, '2')) &&
		b.consume(data, i, ';') {
		b.consumeUntil(data, i, 7)
		actions.Ar.ViewSetTitle(b.viewId, string(data[start+4:*i]))
		return true
	}

	// Other unhandled xterm codes -> ignore
	*i = start
	if b.consume(data, i, 27) && b.consume(data, i, 93) {
		b.consumeUntil(data, i, 7)
		return true
	}

	// ###### Start "real" VT100 codes
	*i = start
	if !(b.consume(data, i, 27) && b.consume(data, i, 91)) { // ^[
		return false
	}

	*i = start + 2
	// cursor up (A)
	if b.consume(data, i, 'A') || (b.consumeNb(data, i) && b.consume(data, i, 'A')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setLine(b.line - nb)
		if b.line-nb < 0 {
			b.setCol(0)
		}
		return true
	}
	*i = start + 2
	// cursor down (B)
	if b.consume(data, i, 'B') || (b.consumeNb(data, i) && b.consume(data, i, 'B')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setLine(b.line + nb)
		return true
	}
	*i = start + 2
	// cursor right (C)
	if b.consume(data, i, 'C') || (b.consumeNb(data, i) && b.consume(data, i, 'C')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setCol(b.col + nb)
		return true
	}
	*i = start + 2
	// cursor left (D)
	if b.consume(data, i, 'D') || (b.consumeNb(data, i) && b.consume(data, i, 'D')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setCol(b.col - nb)
		return true
	}
	*i = start + 2
	// cursor next line (E)
	if b.consume(data, i, 'E') || (b.consumeNb(data, i) && b.consume(data, i, 'E')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setLine(b.line + nb)
		b.setCol(0)
		return true
	}
	*i = start + 2
	// cursor prev line (F)
	if b.consume(data, i, 'F') || (b.consumeNb(data, i) && b.consume(data, i, 'F')) {
		b.flush(data[from:start])
		nb, _ := b.readNb(data, start+2, 1)
		b.setLine(b.line - nb)
		b.setCol(0)
		return true
	}
	*i = start + 2
	// set cursor (H)
	if b.consume(data, i, 'H') || (b.consumeNbTuple(data, i) && b.consume(data, i, 'H')) {
		b.flush(data[from:start])
		row, col, _ := b.readNbTuple(data, start+2, 1, 1)
		b.setLine(row - 1)
		b.setCol(col - 1)
		return true
	}
	*i = start + 2
	// clear screen (J)
	if b.consume(data, i, 'J') || (b.consumeNb(data, i) && b.consume(data, i, 'J')) {
		b.flush(data[from:start])
		ps, _ := b.readNb(data, start+2, 0)
		switch ps {
		case 0: // 0
			b.backend.clearScreen(b.line, b.col)
		case 2:
			b.backend.clearScreen(0, 0)
		default: // 1 is above, 3 is "saved lines"
			log.Printf("TODO Unsupported clear mode (J) %d", ps)
		}
		return true
	}
	*i = start + 2
	// clear line (K)
	if b.consume(data, i, 'K') {
		b.flush(data[from:start])
		ps, _ := b.readNb(data, start+2, 0)
		switch ps {
		case 0: // 0
			b.backend.clearLn(b.line, b.col)
		default: // 1 is left, 2 is whole line
			log.Printf("TODO Unsupported clear line mode (K) %d", ps)
		}
		return true
	}

	*i = start + 2
	// Color attribute + fg color  + bg color
	if b.consumeNbTriple(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start]) // flush what was before escape seq first.
		x, y, z := b.readNbTriple(data, start+2, 0, 0, 0)
		b.applyColors(x, y, z)
		return true
	}
	*i = start + 2
	// Color attribute + fg color
	if b.consumeNbTuple(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		x, y, _ := b.readNbTuple(data, start+2, 0, 0)
		b.applyColors(x, y)
		return true
	}
	*i = start + 2
	// color attr alone
	if b.consumeNb(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		x, _ := b.readNb(data, start+2, 0)
		b.applyColors(x)
		return true
	}
	*i = start + 2
	// reset char/color attributes
	if b.consume(data, i, 'm') {
		b.flush(data[from:start])
		b.curFg, b.curBg = t.Fg, t.Bg
		return true
	}

	// ###### Start unhandled ! ######################################

	*i = start + 2
	// set scrolling region (r)
	if b.consumeNbTuple(data, i) && b.consume(data, i, 'r') {
		b.flush(data[from:start])
		// TODO VT100 scrolling region
		log.Printf("TODO Vt100 scrolling region: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// Various Set comand, ignore for now
	if b.consume(data, i, '?') && b.consumeNb(data, i) && b.consume(data, i, 'h') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 set command: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// Various Set comand, ignore for now
	if b.consume(data, i, '?') && b.consumeNb(data, i) && b.consume(data, i, 'l') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 set command: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// Set alternate keypad mode
	if b.consume(data, i, '=') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 alternate keypad")
		// ignore for now
		return true
	}
	*i = start + 2
	// Set numeric keypad mode
	if b.consume(data, i, '>') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 alternate keypad")
		// ignore for now
		return true
	}

	// Debug : write other potential sequences to log for now
	to := start + 10
	if to > len(data) {
		to = len(data)
	}
	log.Printf("Unhandled escape sequence ? %v | %+q |\n", data[start:to], string(data[start+2:to]))
	// end debug

	// no match
	*i = start
	return false
}

// apply VT100 color attributes
// http://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Functions-using-CSI-_-ordered-by-the-final-character_s_
func (b *backendAppender) applyColors(colors ...int) {
	t := core.Ed.Theme()
	for _, color := range colors {
		switch {
		case color == 0: // normal
			b.curFg = t.Fg
			b.curBg = t.Bg
		case color == 1: // bold
			b.curFg = b.curFg.WithAttr(style.Bold)
		case color == 39: // reset fg
			b.curFg = t.Fg
		case color == 49: // reset bg
			b.curBg = t.Bg
		case color >= 30 && color <= 37: // set fg
			b.curFg = style.NewStyle(uint16(color - 30 + 8))
		case color >= 40 && color <= 47: // set bg
			b.curBg = style.NewStyle(uint16(color - 40 + 8))
		case color >= 90 && color <= 97:
			b.curFg = style.NewStyle(uint16(color - 90 + 8))
		case color >= 100 && color <= 107:
			b.curBg = style.NewStyle(uint16(color - 100 + 8))
		}
	}
}

func (b *backendAppender) consumeNb(data []byte, i *int) bool {
	found := false
	for *i < len(data) && data[*i] >= '0' && data[*i] <= '9' {
		*i++
		found = true
	}
	return found
}

func (b *backendAppender) consumeNbTuple(data []byte, i *int) bool {
	return b.consumeNb(data, i) && b.consume(data, i, ';') &&
		b.consumeNb(data, i)
}

func (b *backendAppender) consumeNbTriple(data []byte, i *int) bool {
	return b.consumeNbTuple(data, i) && b.consume(data, i, ';') &&
		b.consumeNb(data, i)
}

func (b *backendAppender) consume(data []byte, i *int, c byte) bool {
	if *i >= len(data) {
		return false
	}
	if data[*i] == c {
		*i++
		return true
	}
	return false
}

func (b *backendAppender) consumeUntil(data []byte, i *int, c byte) {
	for *i < len(data) && data[*i] != c {
		*i++
	}
	return
}

func (b *backendAppender) readNb(data []byte, i int, defVal int) (nb, readTo int) {
	if i >= len(data) || data[i] < '0' || data[i] > '9' {
		return defVal, i
	}
	n := 0
	for data[i] >= '0' && data[i] <= '9' {
		n = 10*n + int(data[i]-'0')
		i++
	}
	return n, i
}

func (b *backendAppender) readNbTuple(data []byte, i int, defVal1, defVal2 int) (nb1, nb2, readTo int) {
	n1, readTo := b.readNb(data, i, defVal1)
	readTo++ //;
	n2, readTo := b.readNb(data, readTo, defVal2)
	return n1, n2, readTo
}

func (b *backendAppender) readNbTriple(data []byte, i int, defVal1, defVal2, defVal3 int) (int, int, int) {
	n1, n2, readTo := b.readNbTuple(data, i, defVal1, defVal2)
	readTo++ // ;
	n3, readTo := b.readNb(data, readTo, defVal3)
	return n1, n2, n3
}
