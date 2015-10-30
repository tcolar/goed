package core

type Highlighter interface {
	UpdateHighlights(v Viewable)
	ApplyHighlight(v Viewable, lnOffset, ln, col int)
}
