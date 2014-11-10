// History: Oct 02 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

const (
	Plain uint16 = iota + (1 << 8)
	Bold
	Underlined
)

func (e *Editor) WidgetAt(x, y int) Renderer {
	_, h := e.Size()
	if y == 1 {
		return e.Cmdbar
	}
	if y == h-1 {
		return e.Statusbar
	}
	for _, v := range e.Views {
		if x >= v.x1 && x <= v.x2 && y >= v.y1 && y <= v.y2 {
			return v
		}
	}
	return nil
}

func (e *Editor) Render() {
	e.FB(e.Theme.Fg, e.Theme.Bg)
	termbox.Clear(termbox.Attribute(e.Bg.uint16), termbox.Attribute(e.Bg.uint16))

	for _, v := range e.Views {
		v.Render()
	}

	e.Cmdbar.Render()
	e.Statusbar.Render()

	termbox.Flush()
}

type Renderer interface {
	Bounds() (int, int, int, int)
	Render()
	SetBounds(x1, y1, x2, y2 int)
	Event(*termbox.Event)
}

// Widget implements the base of UI widgets
type Widget struct {
	x1, x2, y1, y2 int
}

func (w *Widget) Bounds() (int, int, int, int) {
	return w.x1, w.y1, w.x2, w.y2
}

func (w *Widget) SetBounds(x1, y1, x2, y2 int) {
	w.x1 = x1
	w.x2 = x2
	w.y1 = y1
	w.y2 = y2
}

func (e *Editor) AddView(v *View) {
	e.Views = append(e.Views, v)
}

func (e *Editor) NewCol(pct int) *View {
	x1, y1, x2, y2 := e.CurView.Bounds()
	w := x2 - x1
	nc := int(float64(w) * float64(pct) / 100.0)
	oc := w - nc
	e.CurView.SetBounds(x1, y1, x1+oc, y2)
	nv := NewView()
	nv.SetBounds(x1+oc, y1, x2, y2) // todo : full height
	e.AddView(nv)
	// TODO: resize other views of same column
	return nv
}
