package ui

// Widget represents a UI component of the editor.
type Widget struct {
	x1, x2, y1, y2 int
}

func (w *Widget) Bounds() (y1, x1, y2, x2 int) {
	return w.y1, w.x1, w.y2, w.x2
}

func (w *Widget) SetBounds(y1, x1, y2, x2 int) {
	w.x1 = x1
	w.x2 = x2
	w.y1 = y1
	w.y2 = y2
}
