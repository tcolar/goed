// History: Oct 02 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

const (
	Plain uint16 = iota + (1 << 8)
	Bold
	Underlined
)

type Col struct {
	WidthRatio float64
	Views      []*View
}

func (e *Editor) NewCol(width float64, views []*View) *Col {
	if len(views) < 1 {
		panic("Column must have at least one view !")
	}
	return &Col{
		WidthRatio: width,
		Views:      views,
	}
}

func (e *Editor) WidgetAt(x, y int) Renderer {
	_, h := e.Size()
	if y == 1 {
		return e.Cmdbar
	}
	if y == h-1 {
		return e.Statusbar
	}
	for _, c := range e.Cols {
		for _, v := range c.Views {
			if x >= v.x1 && x <= v.x2 && y >= v.y1 && y <= v.y2 {
				return v
			}
		}
	}
	return nil
}

func (e *Editor) ViewColumn(v *View) *Col {
	for _, c := range e.Cols {
		for _, view := range c.Views {
			if v.Id == view.Id {
				return c
			}
		}
	}
	return nil
}

func (e *Editor) Render() {
	e.FB(e.Theme.Fg, e.Theme.Bg)
	termbox.Clear(termbox.Attribute(e.Bg.uint16), termbox.Attribute(e.Bg.uint16))

	for _, c := range e.Cols {
		for _, v := range c.Views {
			v.Render()
		}
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

// TODO: optimize, for example might only need to resize a single column
func (e *Editor) Resize(width, height int) {
	e.Cmdbar.SetBounds(0, 0, width, 0)
	e.Statusbar.SetBounds(0, height-1, width, height-1)
	wc := 0
	for i, c := range e.Cols {
		hc := 1
		w := int(float64(width) * c.WidthRatio)
		if i == len(e.Cols)-1 {
			w = width - wc // las column gets rest of width
		}
		for j, v := range c.Views {
			h := int(float64(height-2) * v.HeightRatio)
			if j == len(c.Views)-1 {
				h = height - hc - 2 // last view gets rest of height
			}
			v.SetBounds(wc, hc, wc+w, hc+h)
			hc += h
		}
		wc += w
	}
}

// AddCol adds a new column, space is "taken" from the current column
func (e *Editor) AddCol(pct int) *Col {
	/*_, eh := e.Size()
	x1, _, x2, _ := e.CurView.Bounds()
	w := x2 - x1
	nc := int(float64(w) * float6pct) / 100.0)
	oc := w - nc
	for _, v := range e.colViews(e.CurView.x1) {
		v.SetBounds(v.x1, v.y1, v.x1+oc, v.y2)
	}
	nv := NewView()
	nv.SetBounds(x1+oc, 1, x2, eh-2)
	e.Views = append(e.Views, nv)
	return nv*/
	return nil
}

// AddCol adds a new view in the current column, space is "taken" from the current view
func (e *Editor) AddView(ratio float64) *View {
	r := e.CurView.HeightRatio
	nv := e.NewView()
	nv.HeightRatio = r * ratio
	e.CurView.HeightRatio = r - (r * ratio)
	col := e.ViewColumn(e.CurView)
	// TODO: Need to insert it at the right index
	col.Views = append(col.Views, nv)
	e.Resize(e.Size())
	return nv
}
