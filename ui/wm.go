package ui

import (
	"fmt"
	"log"
	"runtime"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	termbox "github.com/tcolar/termbox-go"
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

func (e *Editor) ViewAt(y, x int) int64 {
	vid := int64(-1)
	w := e.WidgetAt(y, x)
	if w != nil {
		if v, ok := w.(core.Viewable); ok {
			vid = v.Id()
		}
	}
	return vid
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

func (e *Editor) Render() {
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
	MouseEvent(e *Editor, ev *termbox.Event)
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
			w = width - wc // last column gets rest of width
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
			y1, _, y2, _ := v.Bounds()
			if y1 != hc || y2 != hc+h-1 { // height changed
				if v.CursorY > h-4 {
					// When the view is shrank, scroll if needed to keep
					// the current line visible.
					off := v.CursorY - h + 4
					v.offy += off
					v.CursorY -= off
				}
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

// Narrowest column
func (e *Editor) ColNarrowest() *Col {
	x := 0
	for i, c := range e.Cols {
		if c.WidthRatio < e.Cols[x].WidthRatio {
			x = i
		}
	}
	return e.Cols[x]
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
	e.ViewActivate(tv.Id())
}

func (e *Editor) CurColIndex() int {
	return e.ColIndex(e.CurCol)
}

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

	e.ViewActivate(nv.Id())
	e.Resize(e.term.Size())
	return c
}

func (e *Editor) AddViewSmart(v *View) *View {
	col := e.tryNewCol(v)
	if col != nil {
		// new full column
		v = e.views[col.Views[0]]
	} else {
		// if we did not have space to create a new column
		// will find the one with the most empty space, or least views
		col = e.emptiestCol()
		v = e.addToCol(col, v)
	}

	e.Resize(e.term.Size())
	return v
}

func (e *Editor) addToCol(c *Col, nv *View) *View {
	a := make([]bool, len(c.Views))
	h, _ := e.term.Size()
	var used float64
	tbd := len(c.Views) + 1 // +1 for nv to be added
	p := h / tbd
	// squeeze empty space
	for i, vid := range c.Views {
		v, _ := e.views[vid]
		if v.LineCount() < p { // view can be shrunk
			a[i] = true
			v.HeightRatio = float64(v.LineCount()+5) / float64(h)
			used += v.HeightRatio
			tbd--
		}
	}
	ratio := (1.0 - used) / float64(tbd)
	for i, vid := range c.Views {
		if a[i] {
			continue
		}
		v, _ := e.views[vid]
		v.HeightRatio = ratio
		used += ratio
	}
	if nv == nil {
		nv = e.NewView("")
	}
	nv.HeightRatio = 1.0 - used
	c.Views = append(c.Views, nv.Id())
	return nv
}

func (e *Editor) emptiestCol() *Col {
	var mf, mfv int
	var lc, lcv int
	dirs := e.ColNarrowest() // assume dir listings column

	for i, c := range e.Cols {
		if c == dirs {
			continue
		}
		f := e.colFreeSpace(c)
		if f > mfv {
			mfv = f
			mf = i
		}
		if lcv == 0 || len(c.Views) <= lcv {
			lcv = len(c.Views)
			lc = i
		}
	}
	// if we have a column with at least x free lines, return that
	if mfv > 10 {
		return e.Cols[mf]
	}
	// else column with least views and right most
	return e.Cols[lc]
}

func (e *Editor) colFreeSpace(c *Col) int {
	h, _ := e.term.Size()
	h -= 4
	for _, vid := range c.Views {
		v, _ := e.views[vid]
		h -= v.LineCount()
	}
	if h < 0 {
		h = 0
	}
	return h
}

// Do we have room for a new column ?
// if so, create it and add the view to it, compress the other columns
func (e *Editor) tryNewCol(v *View) *Col {
	dirs := e.ColNarrowest() // assume dir listings column
	w := 0
	for _, c := range e.Cols {
		if c == dirs {
			w += 23
		} else {
			w += 83
		}
	}
	_, tw := e.term.Size()
	// if so, create it and compact the other columns
	if tw-w < 80 {
		return nil
	}
	r := 0.0
	for _, c := range e.Cols {
		if c == dirs {
			c.WidthRatio = 23.0 / float64(tw)
		} else {
			c.WidthRatio = 83.0 / float64(tw)
		}
		r += c.WidthRatio
	}
	if v == nil {
		v = e.NewView("")
	}
	v.HeightRatio = 1.0
	nc := e.NewCol(1.0-r, []int64{v.Id()})
	e.Cols = append(e.Cols, nc)
	return nc
}

func (e *Editor) AddDirViewSmart(view *View) {
	e.addToCol(e.ColNarrowest(), view)
	e.ViewActivate(view.Id())
	e.Resize(e.term.Size())
}

func (e *Editor) InsertViewSmart(view *View) {
	e.AddViewSmart(view)
}

// AddCol adds a new view in the current column, space is "taken" from toView
func (e *Editor) AddView(toView *View, ratio float64) *View {
	nv := e.NewView("")
	e.InsertView(nv, toView, ratio)
	e.ViewActivate(nv.Id())
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

func (e *Editor) DelColByIndex(index int, check bool) {
	if check {
		e.DelColCheck(e.Cols[index])
	} else {
		e.DelCol(e.Cols[index], true)
	}
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
			e.ViewActivate(v.Id())
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
			e.ViewActivate(cv.Id())
			break
		}
		prev, _ = e.views[c.Views[i]]
	}
	if terminate {
		e.TerminateView(viewId)
	}
	e.Resize(e.term.Size())
}

func (e *Editor) TerminateView(vid int64) {
	v, found := e.views[vid]
	if !found {
		return
	}
	delete(e.views, vid)
	// This probably way overkill, but without nugging the GC it tends to not
	// be very agressive and leave the memory allocated quite a while.
	if v.backend != nil {
		v.backend.Close()
	}
	v.backend = nil
	runtime.GC()
	actions.UndoClear(vid)
}

// Delete (close) a view, with dirty check
func (e *Editor) DelViewCheck(viewId int64, terminate bool) {
	view := e.ViewById(viewId).(*View)
	if view == nil {
		return
	}
	if !view.canClose() {
		e.SetStatusErr("Unsaved changes. Save or request close again.")
		return
	}
	e.DelView(view.Id(), terminate)
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

func (e *Editor) ViewActivate(viewId int64) {
	v := e.ViewById(viewId).(*View)
	if v == nil {
		return
	}
	e.curViewId = viewId
	e.CurCol = e.ViewColumn(e.curViewId)
	v.updateCursor()
	e.SetStatus(fmt.Sprintf("%s [%d]", v.WorkDir(), viewId))
}

func (e *Editor) ViewById(id int64) core.Viewable {
	v, found := e.views[id]
	if !found || v.Id() < 0 {
		log.Printf("View not found %v, \n", id)
		return nil
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
		if v != nil && v.backend != nil && v.backend.SrcLoc() == loc {
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
