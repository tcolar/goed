package ui

import (
	"fmt"
	"log"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/termbox-go"
)

// Col represent a column of the editor (a set of views)
type Col struct {
	WidthRatio float64
	Views      []int64
}

func (e *Editor) NewCol(width float64, views []int64) *Col {
	if len(views) < 1 {
		panic("Column must have at least one view !")
	}
	return &Col{
		WidthRatio: width,
		Views:      views,
	}
}

// WidgetAt returns the widget at a given editor location
func (e *Editor) WidgetAt(y, x int) Renderer {
	h, _ := e.term.Size()
	if y == 0 {
		return e.Cmdbar
	}
	if y == h-1 {
		return e.Statusbar
	}
	for _, c := range e.Cols {
		for _, vid := range c.Views {
			v, found := e.views[vid]
			if found && x >= v.x1 && x <= v.x2 && y >= v.y1 && y <= v.y2 {
				return v
			}
		}
	}
	return nil
}

func (e Editor) Render() {
	e.TermFB(e.theme.Fg, e.theme.Bg)
	e.term.Clear(e.Bg.Uint16(), e.Bg.Uint16())

	for _, c := range e.Cols {
		for _, v := range c.Views {
			e.ViewById(v).Render()
		}
	}

	// cursor
	v := e.CurView().(*View)
	cc, cl := v.CurCol(), v.CurLine()
	c, _, _ := v.CurChar()
	// With some terminals & color schemes the cursor might be "invisible" if we are at a
	// location with no text (ie: end of line)
	// so in that case put as space there to cause the cursor to appear.
	var car = ' '
	if c != nil {
		car = *c
	}
	// Note theterminal inverse the colors where the cursor is
	// this is why this statement might appear "backward"
	e.TermFB(e.theme.BgCursor, e.theme.FgCursor)
	e.TermChar(cl+v.y1-v.offy+2, cc+v.x1-v.offx+2, car)
	e.TermFB(e.theme.Fg, e.theme.Bg)

	e.Cmdbar.Render()
	e.Statusbar.Render()

	e.TermFlush()
}

// Renderer is the interface for a renderable UI component.
type Renderer interface {
	Bounds() (y1, x1, y2, x2 int)
	Render()
	SetBounds(y1, x1, y2, x2 int)
	Event(e *Editor, ev *termbox.Event)
}

// TODO: optimize, for example might only need to resize a single column
func (e *Editor) Resize(height, width int) {
	e.Cmdbar.SetBounds(0, 0, 0, width-1)
	e.Statusbar.SetBounds(height-1, 0, height-1, width-1)
	wc := 0
	wr := 0.0
	for i, c := range e.Cols {
		hc := 1
		hr := 0.0
		w := int(float64(width) * c.WidthRatio)
		if i == len(e.Cols)-1 {
			w = width - wc // las column gets rest of width
			c.WidthRatio = 1.0 - wr
		}
		for j, vi := range c.Views {
			v, found := e.views[vi]
			if !found {
				continue
			}
			h := int(float64(height-2) * v.HeightRatio)
			if h < 1 {
				h = 1
			}
			if j == len(c.Views)-1 {
				h = height - hc - 1 // last view gets rest of height
				v.HeightRatio = 1.0 - hr
			}
			v.SetBounds(hc, wc, hc+h-1, wc+w-1)
			hc += h
			hr += v.HeightRatio
		}
		wc += w
		wr += c.WidthRatio
	}
}

// ViewMove handles moving & resizing views/columns, typically using the mouse
func (e *Editor) ViewMove(y1, x1, y2, x2 int) {
	h, w := e.term.Size()
	v1 := e.WidgetAt(y1, x1).(*View)
	v2 := e.WidgetAt(y2, x2).(*View)
	c1 := e.ViewColumn(v1.Id())
	c2 := e.ViewColumn(v2.Id())
	c1i := e.ColIndex(c1)
	v1i := e.ViewIndex(c1, v1.Id())
	v2i := e.ViewIndex(c2, v2.Id())
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
				e.views[c1.Views[v1i-1]].HeightRatio += ratio // giving space to prev view
			}
		} else if c1i == c2i && v2i == v1i-1 && y2 != v2.y1 {
			// Expanding the view
			ratio := float64(y1-y2) / float64(h-2)
			v1.HeightRatio += ratio
			e.views[c1.Views[v1i-1]].HeightRatio -= ratio // taking space from prev view
		} else if c1i == c2i {
			// moved within same column
			i1, i2 := v1i, v2i
			if i1 < i2 { // down
				copy(c1.Views[i1:i2], c2.Views[i1+1:i2+1])
				c1.Views[v2i] = v1.Id()
			} else { // up
				copy(c1.Views[i2+1:i1+1], c1.Views[i2:i1])
				c1.Views[v2i+1] = v1.Id()
			}
		} else {
			// moved to a different column
			// taking space out of target view
			ratio := float64(y2-v2.y1) / float64(v2.y2-v2.y1)
			e.DelView(v1.Id(), false)
			v1.HeightRatio = v2.HeightRatio * (1.0 - ratio)
			v2.HeightRatio *= ratio
			c2.Views = append(c2.Views, 0)
			copy(c2.Views[v2i+2:], c2.Views[v2i+1:])
			c2.Views[v2i+1] = v1.Id()
		}
	} else {
		// Moving a whole column or a view to it's own column
		if y2 == 1 && v1i > 0 && x2 > v2.x1 && x2 < v2.x2 {
			// Moving a view to it's own column
			ratio := float64(x2-v2.x1) / float64(v2.x2-v2.x1)
			nc := e.AddCol(c2, 1.0-ratio)
			e.DelView(v1.Id(), false)
			nc.Views[0] = v1.Id()
		} else if c1i == c2i {
			// Reducing a column
			if c1i > 0 {
				ratio := float64(x2-x1) / float64(w)
				c1.WidthRatio -= ratio
				e.Cols[c1i-1].WidthRatio += ratio // giving space to prev col
			}
		} else if c2i == c1i-1 {
			// Expanding the column
			ratio := float64(x1-x2) / float64(w)
			c1.WidthRatio += ratio
			e.Cols[c1i-1].WidthRatio -= ratio // taking space from prev col
		} else {
			//reorder
			i1, i2 := c1i, c2i
			if i1 < i2 { // right
				copy(e.Cols[i1:i2], e.Cols[i1+1:i2+1])
				e.Cols[c2i] = c1
			} else { // left
				copy(e.Cols[i2+1:i1+1], e.Cols[i2:i1])
				e.Cols[c2i+1] = c1
			}

		}
	}
	e.SetStatus("")
	e.Resize(e.term.Size())
}

// ViewColumn returns the column that is holding a given view
func (e *Editor) ViewColumn(vid int64) *Col {
	for _, c := range e.Cols {
		for _, view := range c.Views {
			if vid == view {
				return c
			}
		}
	}
	return nil
}

// ViewIndex returns the index of a view in the column
func (e *Editor) ViewIndex(col *Col, vid int64) int {
	for i, v := range col.Views {
		if v == vid {
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

// ViewNavigate navigates from he current view (left, right, up, down)
func (e *Editor) ViewNavigate(mvmt core.CursorMvmt) {
	v := e.CurView().(*View)
	if v == nil {
		return
	}
	c := e.ViewColumn(v.Id())
	if c == nil {
		return
	}
	col := e.ColIndex(c)
	view := e.ViewIndex(e.Cols[col], v.Id())
	if col < 0 || view < 0 {
		return
	}
	switch mvmt {
	case core.CursorMvmtLeft:
		col--
	case core.CursorMvmtRight:
		col++
	case core.CursorMvmtUp:
		view--
	case core.CursorMvmtDown:
		view++
	}
	if col < 0 || col >= len(e.Cols) || view < 0 {
		return
	}
	if view >= len(e.Cols[col].Views) {
		view = len(e.Cols[col].Views) - 1
	}
	tv, _ := e.views[e.Cols[col].Views[view]]
	e.ViewActivate(tv.Id(), tv.CurLine(), tv.CurCol())
}

func (e *Editor) CurColIndex() int {
	return e.ColIndex(e.CurCol)
}

// AddCol adds a new column, space is "taken" from toCol
func (e *Editor) AddCol(toCol *Col, ratio float64) *Col {
	r := toCol.WidthRatio
	nv := e.NewView("")
	nv.HeightRatio = 1.0
	c := e.NewCol(r*ratio, []int64{nv.Id()})
	toCol.WidthRatio = r - (r * ratio)
	// Insert it after curcol
	i := e.ColIndex(toCol) + 1
	e.Cols = append(e.Cols, nil)
	copy(e.Cols[i+1:], e.Cols[i:])
	e.Cols[i] = c

	e.ViewActivate(nv.Id(), 0, 0)
	e.Resize(e.term.Size())
	return c
}

// TODO: Not all that smart yet
func (e *Editor) AddViewSmart() *View {
	var nv *View
	var emptiestCol *Col
	for _, c := range e.Cols {
		// if we have room for a new column,favor that
		v, _ := e.views[c.Views[0]]
		if v.x2-v.x1 > 120 {
			nc := e.AddCol(c, 0.5)
			nv = e.views[nc.Views[0]]
			break
		}
		// otherwise favor emptiest, most to the right col
		if emptiestCol == nil || len(c.Views) <= len(emptiestCol.Views) {
			emptiestCol = c
		}
	}
	if nv == nil {
		// will take some of emptiest col tallest view
		var emptiestView *View
		for _, vid := range emptiestCol.Views {
			v, _ := e.views[vid]
			if emptiestView == nil || emptiestView.HeightRatio <= v.HeightRatio {
				emptiestView = v
			}
		}
		// TODO : consider buffer text length ?
		nv = e.AddView(emptiestView, 0.5)
	}
	e.ViewActivate(nv.Id(), 0, 0)
	e.Resize(e.term.Size())
	return nv
}

func (e *Editor) InsertViewSmart(view *View) {
	nv := e.AddViewSmart()
	e.ReplaceView(nv, view)
}

// AddCol adds a new view in the current column, space is "taken" from toView
func (e *Editor) AddView(toView *View, ratio float64) *View {
	nv := e.NewView("")
	e.InsertView(nv, toView, ratio)
	e.ViewActivate(nv.Id(), 0, 0)
	return nv
}

func (e *Editor) InsertView(view, toView *View, ratio float64) {
	if ratio > 1.0 {
		ratio = 1.0
	}
	r := toView.HeightRatio
	view.HeightRatio = r * ratio
	toView.HeightRatio = r - (r * ratio)
	col := e.ViewColumn(toView.Id())
	// Insert it at after toView
	i := e.ViewIndex(col, toView.Id()) + 1
	col.Views = append(col.Views, 0)
	copy(col.Views[i+1:], col.Views[i:])
	col.Views[i] = view.Id()
	e.Resize(e.term.Size())
}

func (e *Editor) ReplaceView(oldView, newView *View) {
	newView.x1 = oldView.x1
	newView.x2 = oldView.x2
	newView.y1 = oldView.y1
	newView.y2 = oldView.y2
	newView.HeightRatio = oldView.HeightRatio
	col := e.ViewColumn(oldView.Id())
	i := e.ViewIndex(col, oldView.Id())
	col.Views[i] = newView.Id()
	e.TerminateView(oldView.Id())
}

func (e *Editor) DelColCheckByIndex(index int) {
	e.DelColCheck(e.Cols[index])
}

func (e *Editor) DelCol(col *Col, terminateViews bool) {
	if len(e.Cols) <= 1 {
		e.SetStatusErr("Only one column left !")
		return
	}
	var prev *Col
	for i, c := range e.Cols {
		if c == col {
			if prev != nil {
				prev.WidthRatio += c.WidthRatio
				e.CurCol = prev
			} else {
				e.Cols[i+1].WidthRatio += c.WidthRatio
				e.CurCol = e.Cols[i+1]
			}
			v := e.views[e.CurCol.Views[0]]
			e.ViewActivate(v.Id(), v.CurLine(), v.CurCol())
			e.Cols = append(e.Cols[:i], e.Cols[i+1:]...)
			break
		}
		prev = e.Cols[i]
	}
	if terminateViews {
		for _, v := range col.Views {
			e.TerminateView(v)
		}
	}
	col = nil
	e.Resize(e.term.Size())
}

func (e *Editor) DelView(viewId int64, terminate bool) {
	c := e.ViewColumn(viewId)
	if c == nil {
		e.TerminateView(viewId)
		return
	}
	if len(e.Cols) == 1 && len(c.Views) <= 1 {
		e.SetStatusErr("Only one view left !")
		return
	}
	// only one view left in col, delcol
	if len(c.Views) <= 1 {
		e.DelCol(c, terminate)
		return
	}
	// otherwise remove curview and reassign space
	var prev *View
	for i, vid := range c.Views {
		if vid == viewId {
			v, _ := e.views[vid]
			if prev != nil {
				prev.HeightRatio += v.HeightRatio
				e.curViewId = prev.Id()
			} else {
				e.views[c.Views[i+1]].HeightRatio += v.HeightRatio
				e.curViewId = c.Views[i+1]
			}
			c.Views = append(c.Views[:i], c.Views[i+1:]...)
			cv, _ := e.views[e.curViewId]
			e.ViewActivate(cv.Id(), cv.CurLine(), cv.CurCol())
			if terminate {
				e.TerminateView(vid)
			}
			break
		}
		prev, _ = e.views[c.Views[i]]
	}
	e.Resize(e.term.Size())
}

func (e *Editor) TerminateView(vid int64) {
	v, found := e.views[vid]
	if !found {
		return
	}
	if v.backend != nil {
		v.backend.Close()
	}
	actions.UndoClear(vid)
	delete(e.views, vid)
}

// Delete (close) a view, but with dirty check
func (e *Editor) DelViewCheck(viewId int64) {
	view := e.ViewById(viewId).(*View)
	if view == nil {
		return
	}
	if !view.canClose() {
		e.SetStatusErr("Unsaved changes. Save or request close again.")
		return
	}
	e.DelView(view.Id(), true)
}

// Delete (close) a col, but with dirty check
func (e *Editor) DelColCheck(c *Col) {
	ok := true
	for _, vid := range c.Views {
		v, found := e.views[vid]
		if !found {
			continue
		}
		ok = ok && v.canClose()
	}
	if !ok {
		e.SetStatusErr("Unsaved changes. Save or request close again.")
		return
	}
	e.DelCol(c, true)
}

func (e *Editor) ViewActivate(viewId int64, cursory, cursorx int) {
	v := e.ViewById(viewId).(*View)
	if v == nil {
		return
	}
	e.curViewId = viewId
	e.CurCol = e.ViewColumn(e.curViewId)
	v.MoveCursor(cursory-v.CurLine(), cursorx-v.CurCol())
}

func (e *Editor) SetCurView(id int64) error {
	v := e.ViewById(id)
	if v == nil {
		return fmt.Errorf("No such view %d", id)
	}
	e.ViewActivate(v.Id(), 0, 0)
	return nil
}

func (e Editor) ViewById(id int64) core.Viewable {
	v, found := e.views[id]
	if !found {
		log.Printf("View not found %v, \n", id)
	}
	return v
}

// ViewByLoc returns the view matching the given location
// or -1 if no match
func (e *Editor) ViewByLoc(loc string) int64 {
	if len(loc) == 0 {
		return -1
	}
	for _, v := range e.views {
		if v.backend.SrcLoc() == loc {
			return v.Id()
		}
	}
	return -1
}

// SwapView swaps 2 views (UI wise)
func (e *Editor) SwapViews(vv1, vv2 int64) {
	v1 := e.ViewById(vv1).(*View)
	v2 := e.ViewById(vv2).(*View)
	if v1 == nil || v2 == nil {
		return
	}
	c1 := e.ViewColumn(v1.Id())
	i1 := e.ViewIndex(c1, v1.Id())
	c2 := e.ViewColumn(v2.Id())
	i2 := e.ViewIndex(c2, v2.Id())
	c1.Views[i1], c2.Views[i2] = vv2, vv1
	v1.HeightRatio, v2.HeightRatio = v2.HeightRatio, v1.HeightRatio
	v1.y1, v1.x1, v2.y1, v2.x1 = v2.y1, v2.x1, v1.y1, v1.x1
	v1.y2, v1.x2, v2.y2, v2.x2 = v2.y2, v2.x2, v1.y2, v1.x2
}
