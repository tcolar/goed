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
	if y == 0 {
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

func (e *Editor) Render() {
	e.FB(e.Theme.Fg, e.Theme.Bg)
	termbox.Clear(termbox.Attribute(e.Bg.uint16), termbox.Attribute(e.Bg.uint16))

	for _, c := range e.Cols {
		for _, v := range c.Views {
			v.Render()
		}
	}

	// cursor
	v := Ed.CurView
	cc, cl := v.CurCol(), v.CurLine()
	c, _, _ := v.CursorChar(cc, cl)
	// With some terminals & color schemes the cursor might be "invisible" if we are at a
	// location with no text (ie: end of line)
	// so in that case put as space there to cause the cursor to appear.
	var car = ' '
	if c != nil {
		car = *c
	}
	// Note theterminal inverse the colors where the cursor is
	// this is why this statement might appear "backward"
	Ed.FB(Ed.Theme.BgCursor, Ed.Theme.FgCursor)
	Ed.Char(cc+v.x1-v.offx+2, cl+v.y1-v.offy+2, car)
	Ed.FB(Ed.Theme.Fg, Ed.Theme.Bg)

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
				h = height - hc - 1 // last view gets rest of height
				// TODO: maybe adjust the ratio so it always adds up to ~100%
			}
			v.SetBounds(wc, hc, wc+w-1, hc+h-1)
			hc += h
		}
		wc += w
	}
}

// ViewMove handles moving & resizing views/columns, typically using the mouse
func (e *Editor) ViewMove(x1, y1, x2, y2 int) {
	w, h := Ed.Size()
	v1 := e.WidgetAt(x1, y1).(*View)
	v2 := e.WidgetAt(x2, y2).(*View)
	c1 := e.ViewColumn(v1)
	c2 := e.ViewColumn(v2)
	c1i := e.ColIndex(c1)
	v1i := e.ViewIndex(c1, v1)
	v2i := e.ViewIndex(c2, v2)
	c2i := e.ColIndex(c2)
	onSep := x2 == v2.x1 // dropped on a column "scrollbar"
	if x1 == x2 && y1 == y2 {
		// noop
		e.SetStatus("")
		return
	}
	if onSep {
		// We are moving a single view
		if v1 == v2 {
			if v1i > 0 {
				// Reducing a view
				ratio := float64(y2-y1) / float64(h-2)
				v1.HeightRatio -= ratio
				c1.Views[v1i-1].HeightRatio += ratio // giving space to prev view
			}
		} else if c1i == c2i && v2i == v1i-1 {
			// Expanding the view
			ratio := float64(y1-y2) / float64(h-2)
			v1.HeightRatio += ratio
			c1.Views[v1i-1].HeightRatio -= ratio // taking space from prev view
		} else if c1i == c2i {
			// moved within same column
			i1, i2 := v1i, v2i
			if i1 < i2 { // down
				copy(c1.Views[i1:i2], c2.Views[i1+1:i2+1])
				c1.Views[v2i] = v1
			} else { // up
				copy(c1.Views[i2+1:i1+1], c1.Views[i2:i1])
				c1.Views[v2i+1] = v1
			}
		} else {
			// moved to a different column
			// TODO: Try to be samrt like acme as far as minimizing whitespace
			// For now will just take 50% of the target view
			v2.HeightRatio /= 2
			v1.HeightRatio = v2.HeightRatio
			c2.Views = append(c2.Views, nil)
			copy(c2.Views[v2i+2:], c2.Views[v2i+1:])
			c2.Views[v2i+1] = v1
			ov := e.CurView
			oc := e.CurCol
			e.CurView = v1
			e.CurCol = c1
			if len(c1.Views) == 0 {
				e.DelCol()
			} else {
				e.DelView()
			}
			e.CurView = ov
			e.CurCol = oc
		}
	} else {
		// Moving a whole column
		if c1i == c2i {
			if c1i > 0 {
				// Reducing a column
				ratio := float64(x2-x1) / float64(w)
				c1.WidthRatio -= ratio
				Ed.Cols[c1i-1].WidthRatio += ratio // giving space to prev col
			}
		} else if c2i == c1i-1 {
			// Expanding the column
			ratio := float64(x1-x2) / float64(w)
			c1.WidthRatio += ratio
			Ed.Cols[c1i-1].WidthRatio -= ratio // taking space from prev col
		} else {
			//reorder
			i1, i2 := c1i, c2i
			if i1 < i2 { // right
				copy(Ed.Cols[i1:i2], Ed.Cols[i1+1:i2+1])
				e.Cols[c2i] = c1
			} else { // left
				copy(Ed.Cols[i2+1:i1+1], Ed.Cols[i2:i1])
				e.Cols[c2i+1] = c1
			}

		}
	}
	e.SetStatus("")
	e.Resize(e.Size())
}

// ViewColumn returns the column that is holding a given view
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
	if ratio > 1.0 {
		ratio = 1.0
	}
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
