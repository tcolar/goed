package backend

import (
	"log"

	"github.com/tcolar/goed/actions"
)

// TODO : Handle VT100 codes
// http://www.ccs.neu.edu/research/gpc/VonaUtils/vona/terminal/VT100_Escape_Codes.html
// http://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-VT100-Mode
func (b *backendAppender) vt100(data []byte) (int, error) {
	// handle VT100 codes
	from := 0
	log.Printf("VT100 raw data: %v\n", data)
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

func (b *backendAppender) consumeVt100(data []byte, from int, i *int) bool {

	start := *i
	// bell
	if b.consume(data, i, 7) {
		b.flush(data[from:start])
		actions.EdSetStatusErr("Beep !!")
		return true
	}
	*i = start
	// "backspace" (move left)
	if b.consume(data, i, 8) {
		b.flush(data[from:start])
		b.col--
		if b.col < 0 {
			b.col = 0
		}
		return true
	}
	// delete
	*i = start
	if b.consume(data, i, 127) {
		b.flush(data[from:start])
		actions.ViewDeleteCur(b.viewId)
		return true
	}
	*i = start // \r\n
	if b.consume(data, i, 13) && b.consume(data, i, 10) {
		b.flush(data[from:start])
		b.line++
		b.col = 0
		return true
	}
	*i = start
	if b.consume(data, i, 10) { // \n
		b.flush(data[from:start])
		b.line++
		b.col = 0
		return true
	}
	*i = start
	if b.consume(data, i, 13) { // \r
		b.flush(data[from:start])
		//b.line++
		b.col = 0
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
	*i = start
	if b.consume(data, i, 27) && b.consume(data, i, 93) &&
		b.consume(data, i, '0') && b.consume(data, i, ';') {
		b.consumeUntil(data, i, 7)
		actions.ViewSetTitle(b.viewId, string(data[start+4:*i]))
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
		nb := b.readNb(data, start+2, 1)
		b.line -= nb
		if b.line < 0 {
			b.line = 0
		}
		return true
	}
	*i = start + 2
	// cursor down (B)
	if b.consume(data, i, 'B') || (b.consumeNb(data, i) && b.consume(data, i, 'B')) {
		b.flush(data[from:start])
		nb := b.readNb(data, start+2, 1)
		b.line += nb
		return true
	}
	*i = start + 2
	// cursor right (C)
	if b.consume(data, i, 'C') || (b.consumeNb(data, i) && b.consume(data, i, 'C')) {
		b.flush(data[from:start])
		nb := b.readNb(data, start+2, 1)
		b.col += nb
		return true
	}
	*i = start + 2
	// cursor left (D)
	if b.consume(data, i, 'D') || (b.consumeNb(data, i) && b.consume(data, i, 'D')) {
		b.flush(data[from:start])
		nb := b.readNb(data, start+2, 1)
		b.col -= nb
		if b.col < 0 {
			b.col = 0
		}
		return true
	}
	*i = start + 2
	// cursor next line (E)
	if b.consume(data, i, 'E') || (b.consumeNb(data, i) && b.consume(data, i, 'E')) {
		b.flush(data[from:start])
		nb := b.readNb(data, start+2, 1)
		b.line += nb
		b.col = 0
		return true
	}
	*i = start + 2
	// cursor prev line (F)
	if b.consume(data, i, 'F') || (b.consumeNb(data, i) && b.consume(data, i, 'F')) {
		b.flush(data[from:start])
		nb := b.readNb(data, start+2, 1)
		b.line -= nb
		if b.line < 0 {
			b.line = 0
		}
		b.col = 0
		return true
	}
	*i = start + 2
	// set cursor (H)
	if b.consume(data, i, 'H') || (b.consumeNbTuple(data, i) && b.consume(data, i, 'H')) {
		b.flush(data[from:start])
		row, col := b.readNbTuple(data, start+2, 1, 1)
		b.line, b.col = row-1, col-1
		if b.line < 0 {
			b.line = 0
		}
		if b.col < 0 {
			b.col = 0
		}
		return true
	}
	*i = start + 2
	// clear screen (J)
	if b.consume(data, i, 'J') || (b.consumeNb(data, i) && b.consume(data, i, 'J')) {
		b.flush(data[from:start])
		ps := b.readNb(data, start+2, 0)
		switch ps {
		case 0: // 0
			b.backend.(*MemBackend).ClearScreen(b.line, b.col)
		case 2:
			b.backend.Wipe()
		default: // 1 is above, 3 is "saved lines"
			log.Printf("TODO Unsupported clear mode (J) %d", ps)
		}
		return true
	}
	*i = start + 2
	// clear line (K)
	if b.consume(data, i, 'K') {
		b.flush(data[from:start])
		ps := b.readNb(data, start+2, 0)
		switch ps {
		case 0: // 0
			b.backend.(*MemBackend).ClearLn(b.line, b.col)
		default: // 1 is left, 2 is whole line
			log.Printf("TODO Unsupported clear line mode (K) %d", ps)
		}
		return true
	}

	// ###### Start unhandled ! ######################################

	*i = start + 2
	// Color attribute + fg color  + bg color
	if b.consumeNbTriple(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start]) // flush what was before escape seq first.
		log.Printf("TODO Vt100 color: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// Color attribute + fg color
	if b.consumeNbTuple(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 color: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// color attr alone
	if b.consumeNb(data, i) && b.consume(data, i, 'm') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 color: %v", data[start+2:*i])
		return true
	}
	*i = start + 2
	// reset char/color attributes
	if b.consume(data, i, 'm') {
		b.flush(data[from:start])
		log.Printf("TODO Vt100 reset colors")
		return true
	}
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

func (b *backendAppender) readNb(data []byte, i int, defVal int) int {
	if data[i] <= '0' || data[i] >= '9' {
		return defVal
	}
	n := 0
	for data[i] >= '0' && data[i] <= '9' {
		n = 10*n + int(data[i]-'0')
		i++
	}
	return n
}

func (b *backendAppender) readNbTuple(data []byte, i int, defVal1, defVal2 int) (int, int) {
	if data[i] <= '0' || data[i] >= '9' {
		return defVal1, defVal2
	}
	n1, n2 := 0, 0
	for data[i] >= '0' && data[i] <= '9' {
		n1 = 10*n1 + int(data[i]-'0')
		i++
	}
	i++ // ';'
	for data[i] >= '0' && data[i] <= '9' {
		n2 = 10*n2 + int(data[i]-'0')
		i++
	}
	return n1, n2
}
