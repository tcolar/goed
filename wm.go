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

// ViewIndex returns the index of a view in the column
func (e *Editor) ViewIndex(col *Col, view *View) int {
	for i, v := range col.Views {
		if v == view {
			return i
		}
	}
	return -1
}

// ColIndex returns the index of a column in the editor
func (e *Editor) ColIndex(col *Col) int {
	for i, c := range e.Cols {
		if c == col {
			return i
		}
	}
	return -1
}

// AddCol adds a new column, space is "taken" from the current column
func (e *Editor) AddCol(ratio float64) *Col {
	r := e.CurCol.WidthRatio
	nv := e.NewView()
	nv.HeightRatio = 1.0
	c := e.NewCol(r*ratio, []*View{nv})
	e.CurCol.WidthRatio = r - (r * ratio)
	// Insert it after curcol
	i := e.ColIndex(e.CurCol) + 1
	e.Cols = append(e.Cols, nil)
	copy(e.Cols[i+1:], e.Cols[i:])
	e.Cols[i] = c

	e.CurCol = c
	e.CurView = nv
	e.Resize(e.Size())
	return c
}

// AddCol adds a new view in the current column, space is "taken" from the current view
func (e *Editor) AddView(ratio float64) *View {
	r := e.CurView.HeightRatio
	nv := e.NewView()
	nv.HeightRatio = r * ratio
	e.CurView.HeightRatio = r - (r * ratio)
	col := e.CurCol
	// Insert it at after curView
	i := e.ViewIndex(col, e.CurView) + 1
	col.Views = append(col.Views, nil)
	copy(col.Views[i+1:], col.Views[i:])
	col.Views[i] = nv
	e.CurView = nv
	e.Resize(e.Size())
	return nv
}

func (e *Editor) DelCol() {
	// TODO: check and warn if any view is dirty ??
	if len(e.Cols) <= 1 {
		e.SetStatusErr("Only one column left !")
		return
	}
	var prev *Col
	for i, c := range e.Cols {
		if c == e.CurCol {
			if prev != nil {
				prev.WidthRatio += c.WidthRatio
				e.CurCol = prev
			} else {
				e.Cols[i+1].WidthRatio += c.WidthRatio
				e.CurCol = e.Cols[i+1]
			}
			e.CurView = e.CurCol.Views[0]
			e.Cols = append(e.Cols[:i], e.Cols[i+1:]...)
			break
		}
		prev = e.Cols[i]
	}
	e.Resize(e.Size())
}

func (e *Editor) DelView() {
	// TODO: check and warn if dirty ??
	c := e.CurCol
	if len(e.Cols) <= 1 && len(c.Views) <= 1 {
		e.SetStatusErr("Only one view left !")
		return
	}
	// only one view left in col, delcol
	if len(c.Views) <= 1 {
		e.DelCol()
		return
	}
	// otherwise remove curview and reassign space
	var prev *View
	for i, v := range c.Views {
		if v == e.CurView {
			if prev != nil {
				prev.HeightRatio += v.HeightRatio
				e.CurView = prev
			} else {
				c.Views[i+1].HeightRatio += v.HeightRatio
				e.CurView = c.Views[i+1]
			}
			c.Views = append(c.Views[:i], c.Views[i+1:]...)
			break
		}
		prev = c.Views[i]
	}
	e.Resize(e.Size())
}
