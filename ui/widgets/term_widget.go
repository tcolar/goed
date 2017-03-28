package widgets

import (
	"github.com/tcolar/goed/core/term"
	"github.com/tcolar/goed/ui/style"
)

type TermWidget struct {
	term    term.Term
	widgets []Widget
}

func NewTermWidget(term term.Term) *TermWidget {
	return &TermWidget{
		term: term,
	}
}

func (w *TermWidget) GetParent() Widget {
	return nil // the terminal widget is always the root
}

func (w *TermWidget) Bounds() (y1, x1, y2, x2 int) {
	y, x := w.term.Size()
	return 0, 0, y - 1, x - 1
}

func (w *TermWidget) Char(y, x int, r rune, fg, bg style.Style) {
	w.term.Char(y, x, r, fg, bg)
}

func (w *TermWidget) Flush() {
	w.term.Flush()
}

func (w *TermWidget) Render() {
	for _, child := range w.widgets {
		child.Render()
	}
}

func (w *TermWidget) SetBounds(y1, x1, y2, x2 int) {
	// NOOP, use the term bounds
}

func (w *TermWidget) AddWidget(ww Widget) {
	ww.SetParent(w)
	w.widgets = append(w.widgets, ww)
}

func (w *TermWidget) SetParent(_ Widget) {
	// alwasy nil
}
