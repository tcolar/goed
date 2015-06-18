package syntax

import (
	"path/filepath"
	"strings"
	"unicode"
)

type Highlight struct {
	Style          StyleId
	ColFrom, ColTo int
}

func NewHighlight(s StyleId, colFrom, colTo int) Highlight {
	return Highlight{
		Style:   s,
		ColFrom: colFrom, ColTo: colTo,
	}
}

type Highlights struct {
	Lines           [][]Highlight
	ln, col         int // internal use
	curLn, curIndex int
}

// Update updates the highlights for the given text (current slice)
func (h *Highlights) Update(text [][]rune, file string) {
	h.ln = 0
	h.col = 0
	h.curLn = 0
	h.curIndex = 0
	ext := strings.ToLower(filepath.Ext(file))
	base := strings.ToLower(filepath.Base(file))
	h.Lines = make([][]Highlight, len(text))
	syntax, found := Syntaxes[ext]
	if !found {
		syntax, found = Syntaxes[base]
	}
	if !found {
		syntax = Syntaxes["_"]
	}
	h.consumeLeftovers(syntax.Patterns, text)
	for h.ln < len(text) {
		consumed := h.consumePatterns(syntax.Patterns, text) ||
			h.consume(syntax.Symbols, text, false) ||
			h.consume(syntax.Keywords, text, true)
		if !consumed {
			h.col++
		}
		if h.col >= len(text[h.ln]) {
			h.ln++
			h.col = 0
		}
	}
}

// StyleAt returnsthe highlight at a given location,
// Expects forward sweep from start to finish !
func (h *Highlights) StyleAt(ln, col int) StyleId {
	if ln >= len(h.Lines) {
		return StyleNone
	}
	if ln != h.curLn {
		h.curIndex = 0
	}
	h.curLn = ln
	line := h.Lines[ln]
	if line == nil || h.curIndex >= len(line) {
		return StyleNone
	}
	item := line[h.curIndex]
	if col < item.ColFrom {
		return StyleNone
	}
	if col >= item.ColFrom && col <= item.ColTo {
		if col == item.ColTo {
			h.curIndex++
		}
		return item.Style
	}
	h.curIndex++
	return StyleNone
}

// consumeLeftovers tries to consume multline patterns leftovers at the top of
// the screen. ie: "end" lines of a partial '/*' ...... '*/' comment.
func (h *Highlights) consumeLeftovers(patterns []SyntaxPattern, text [][]rune) {
	for _, p := range patterns {
		if !p.MultiLine {
			continue
		}
		h.col = 0
		h.ln = 0
		for h.ln < len(text) {
			if h.peek(p.Start, text) {
				if p.Start != p.End || h.col+len(p.Start) < len(text[h.ln]) {
					// if p.Start comes before p.End, it's not a leftover.
					// When the opening and ending of a multiline comment is the
					// same there is no "easy" way to tell which it is other than
					// counting from the beginning of file, which is costly.
					// for now will make an "often correct" assumption that the
					// end one is usally (not always) toward EOL
					break
				}
			}
			if h.peek(p.End, text) {
				// Found a close first, we have a multiline leftover
				for i := 0; i < h.ln; i++ {
					hl := NewHighlight(p.StyleId, 0, len(text[i])-1)
					h.Lines[i] = append(h.Lines[i], hl)
				}
				h.col += len(p.End)
				hl := NewHighlight(p.StyleId, 0, h.col-1)
				h.Lines[h.ln] = append(h.Lines[h.ln], hl)
				return
			}
			//continue
			h.col++
			if h.col >= len(text[h.ln]) {
				h.ln++
				h.col = 0
			}
		}
	}
	h.col = 0
	h.ln = 0
}

// consumePatterns consumes the text patterns
func (h *Highlights) consumePatterns(patterns []SyntaxPattern, text [][]rune) bool {
	var hl Highlight
	var p *SyntaxPattern
	for _, pat := range patterns {
		if pat.MustStartLine && h.col > 0 {
			continue
		}
		if h.peek(pat.Start, text) {
			hl = NewHighlight(pat.StyleId, h.col, h.col)
			h.col += len(pat.Start)
			p = &pat
			break
		}
	}
	if p != nil {
		if len(p.End) > 0 { // find ending
			for {
				prev := "\u0000"
				found := true
				for !h.peek(p.End, text) || prev == p.Escape {
					h.col++
					if h.col >= len(text[h.ln]) {
						found = false
						break
					}
					prev = string(text[h.ln][h.col-1])
				}
				if found {
					h.col += len(p.End)
				}
				hl.ColTo = h.col - 1
				h.Lines[h.ln] = append(h.Lines[h.ln], hl)
				if found || !p.MultiLine || h.ln >= len(text)-1 {
					break
				}
				// add whole line, and keep looking
				h.ln++
				h.col = 0
				hl = NewHighlight(p.StyleId, h.col, h.col)
			}
		} else { // To EOL
			h.col = len(text[h.ln])
			hl.ColTo = h.col
			h.Lines[h.ln] = append(h.Lines[h.ln], hl)
		}
		return true
	}
	return false
}

// consumes normal items (exact matches)
func (h *Highlights) consume(items SyntaxItems, text [][]rune, isKw bool) bool {
	for _, si := range items {
		if isKw && h.col > 0 {
			// check keyword is not part of longer "word" in which case it
			// should not be highlighted (ie: "go" in "pogo")
			r := text[h.ln][h.col-1]
			if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
				continue
			}
		}
		if h.peek(si.Text, text) {
			matchLen := len(si.Text)
			if isKw && h.col+matchLen < len(text[h.ln]) {
				// check keyword is not part of longer "word" in which case it
				// should not be highlighted (ie: "go" in "gopher")
				r := text[h.ln][h.col+matchLen]
				if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
					continue
				}
			}
			h.Lines[h.ln] = append(h.Lines[h.ln], NewHighlight(si.Id, h.col, h.col+len(si.Text)-1))
			h.col += matchLen
			return true
		}
	}
	return false
}

// peek sees if the given string (s) is found at the given location
// Does not advance h.Ln or h.Col
func (h *Highlights) peek(s string, text [][]rune) bool {
	offCol := 0
	if h.ln >= len(text) {
		return false
	}
	ln := text[h.ln]
	if h.col+len(s) > len(ln) {
		return false
	}
	for _, r := range s {
		if r != ln[h.col+offCol] {
			return false
		}
		offCol++
	}
	return true
}
