package ui

import (
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/syntax"
)

// CodeHighlighter is used to highlight(color) source code
type CodeHighlighter struct {
	highlights syntax.Highlights
}

func (h *CodeHighlighter) UpdateHighlights(v core.Viewable) {
	if len(v.Backend().SrcLoc()) > 0 {
		h.highlights.Update(*v.Slice().Text(), v.Backend().SrcLoc())
	}
}

func (h *CodeHighlighter) ApplyHighlight(v core.Viewable, lnOffset, ln, col int) {
	style := h.highlights.StyleAt(ln, col)
	e := core.Ed
	t := e.Theme()
	var s core.Style
	switch style {
	case syntax.StyleComment:
		s = t.Comment
	case syntax.StyleString:
		s = t.String
	case syntax.StyleKw1:
		s = t.Keyword1
	case syntax.StyleKw2:
		s = t.Keyword2
	case syntax.StyleKw3:
		s = t.Keyword3
	case syntax.StyleSep1:
		s = t.Separator1
	case syntax.StyleSep2:
		s = t.Separator2
	case syntax.StyleSep3:
		s = t.Separator3
	case syntax.StyleSymb1:
		s = t.Symbol1
	case syntax.StyleSymb2:
		s = t.Symbol2
	case syntax.StyleSymb3:
		s = t.Symbol3
	default:
		s = t.Fg
	}
	e.TermFB(s, t.Bg)
}

// TermHighlighter is used to highlight(color) terminal output
type TermHighlighter struct {
}

func (h *TermHighlighter) UpdateHighlights(v core.Viewable) {
}

func (h *TermHighlighter) ApplyHighlight(v core.Viewable, lnOffset, ln, col int) {
	core.Ed.TermFB(v.Backend().ColorAt(ln+lnOffset, col))
}
