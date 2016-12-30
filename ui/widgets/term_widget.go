package widgets

import "github.com/tcolar/goed/core"

type TermWidget struct {
	term    core.Term
	widgets []core.Widget
}

func NewTermWidget(term core.Term) *TermWidget {
	return &TermWidget{
		term: term,
	}
}

func (w *TermWidget) GetParent() core.Widget {
	return nil // the terminal widget is always the root
}

func (w *TermWidget) Bounds() (y1, x1, y2, x2 int) {
	y, x := w.term.Size()
	return 0, 0, y - 1, x - 1
}

func (w *TermWidget) Char(y, x int, r rune, fg, bg core.Style) {
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

func (w *TermWidget) AddWidget(ww core.Widget) {
	ww.SetParent(w)
	w.widgets = append(w.widgets, ww)
}

func (w *TermWidget) SetParent(_ core.Widget) {
	// alwasy nil
}
