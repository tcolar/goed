package widgets

import "github.com/tcolar/goed/core"

// Widget represents a UI component of the editor.
type BaseWidget struct {
	x1, x2, y1, y2 int
	Bg, Fg         core.Style
	parent         core.Widget
}

func (w *BaseWidget) Bounds() (y1, x1, y2, x2 int) {
	return w.y1, w.x1, w.y2, w.x2
}

func (w *BaseWidget) X1() int { return w.x1 }
func (w *BaseWidget) X2() int { return w.x2 }
func (w *BaseWidget) Y1() int { return w.y1 }
func (w *BaseWidget) Y2() int { return w.y2 }

func (w *BaseWidget) SetBounds(y1, x1, y2, x2 int) {
	w.x1 = x1
	w.x2 = x2
	w.y1 = y1
	w.y2 = y2
}

func GetTermWidget(w core.Widget) *TermWidget {
	ww := w
	for {
		if ww == nil {
			return nil
		}
		if wdg, ok := ww.(*TermWidget); ok {
			return wdg
		}
		ww = ww.GetParent()
	}
}

func (w *BaseWidget) GetParent() core.Widget {
	return w.parent
}

func (w *BaseWidget) SetParent(parent core.Widget) {
	w.parent = parent
}

func (w *BaseWidget) Move(y, x int) {
	w.x1 += x
	w.x2 += x
	w.y1 += y
	w.y2 += y
}
