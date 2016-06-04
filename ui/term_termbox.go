package ui

import (
	"fmt"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
	termbox "github.com/tcolar/termbox-go"
)

const (
	// Double clicks
	MouseLeftDbl termbox.Key = 0xFF00 + iota
	MouseRightDbl
)

// Evtstate stores some state about kb/mouse events
type EvtState struct {
	MovingView                    bool
	LastClickX, LastClickY        int
	LastLeftClick, LastRightClick int64 // timestamp
	DragLn, DragCol               int
	InDrag                        bool
}

// ==================== Termbox impl ===========================

// Real Terinal implementation using termbox
type TermBox struct {
}

func NewTermBox() *TermBox {
	return &TermBox{}
}

func (t *TermBox) Init() error {
	return termbox.Init()
}

func (t *TermBox) Clear(fg, bg uint16) {
	termbox.Clear(termbox.Attribute(fg), termbox.Attribute(bg))
}

func (t *TermBox) Close() {
	termbox.Close()
}

func (t *TermBox) Flush() {
	termbox.Flush()
}

func (t *TermBox) SetExtendedColors(b bool) {
	termbox.SetExtendedColors(b)
}

func (t *TermBox) SetCursor(y, x int) {
	termbox.SetCursor(y, x)
}

func (t *TermBox) Char(y, x int, c rune, fg, bg core.Style) {
	termbox.SetCell(x, y, c, termbox.Attribute(fg.Uint16()), termbox.Attribute(bg.Uint16()))
}

func (t *TermBox) Size() (h, w int) {
	w, h = termbox.Size()
	return h, w
}

func (t *TermBox) SetMouseMode(m termbox.MouseMode) {
	termbox.SetMouseMode(m)
}

func (t *TermBox) SetInputMode(m termbox.InputMode) {
	termbox.SetInputMode(m)
}

func (t *TermBox) Listen() {
	t.SetMouseMode(termbox.MouseMotion)
	// Note: terminal might not support SGR mouse events, but trying anyway
	t.SetMouseMode(termbox.MouseSgr)

	t.SetInputMode(termbox.InputMouse)
	es := event.NewEvent()
	for {
		ev := termbox.PollEvent()
		t.parseEvent(ev, es)
		if es.Type == event.EvtQuit {
			return
		}
		event.Queue(es)
	}
}

// parses a termbox event into the 'es' goed event (event.Event)
func (t *TermBox) parseEvent(e termbox.Event, es *event.Event) {
	es.Type = event.Evt_None
	if e.Ch > 0 {
		es.Glyph = string(e.Ch)
	} else {
		es.Glyph = ""
	}
	es.MouseBtns = map[int]bool{}
	es.Combo = event.Combo{}
	es.Keys = []string{}

	if len(es.Glyph) > 0 {
		es.KeyDown(es.Glyph)
	}

	ctrl := func(k string) {
		es.KeyDown(event.KeyLeftControl)
		es.KeyDown(k)
	}

	k := e.Key
	switch k {
	case termbox.MouseLeft:
		es.MouseDown(event.MouseLeft, e.MouseY, e.MouseX)
	case termbox.MouseRight:
		es.MouseDown(event.MouseRight, e.MouseY, e.MouseX)
	case termbox.MouseMiddle:
		es.MouseDown(event.MouseMiddle, e.MouseY, e.MouseX)
	case termbox.MouseScrollDown:
		es.MouseDown(event.MouseWheelDown, e.MouseY, e.MouseX)
	case termbox.MouseScrollUp:
		es.MouseDown(event.MouseWheelUp, e.MouseY, e.MouseX)
	case termbox.KeyBackspace:
		es.KeyDown(event.KeyBackspace)
	case termbox.KeyTab:
		es.KeyDown(event.KeyTab)
	case termbox.KeyEnter:
		es.KeyDown(event.KeyReturn)
	case termbox.KeySpace:
		es.KeyDown(event.KeySpace)
	case termbox.KeyF1:
		es.KeyDown(event.KeyF1)
	case termbox.KeyF2:
		es.KeyDown(event.KeyF2)
	case termbox.KeyF3:
		es.KeyDown(event.KeyF3)
	case termbox.KeyF4:
		es.KeyDown(event.KeyF4)
	case termbox.KeyF5:
		es.KeyDown(event.KeyF5)
	case termbox.KeyF6:
		es.KeyDown(event.KeyF7)
	case termbox.KeyF7:
		es.KeyDown(event.KeyF7)
	case termbox.KeyF8:
		es.KeyDown(event.KeyF8)
	case termbox.KeyF9:
		es.KeyDown(event.KeyF9)
	case termbox.KeyF10:
		es.KeyDown(event.KeyF10)
	case termbox.KeyF11:
		es.KeyDown(event.KeyF11)
	case termbox.KeyF12:
		es.KeyDown(event.KeyF12)
	case termbox.KeyInsert:
		es.KeyDown(event.KeyInsert)
	case termbox.KeyDelete:
		es.KeyDown(event.KeyDelete)
	case termbox.KeyHome:
		es.KeyDown(event.KeyHome)
	case termbox.KeyEnd:
		es.KeyDown(event.KeyEnd)
	case termbox.KeyPgup:
		es.KeyDown(event.KeyPrior)
	case termbox.KeyPgdn:
		es.KeyDown(event.KeyNext)
	case termbox.KeyArrowUp:
		es.KeyDown(event.KeyUpArrow)
	case termbox.KeyArrowDown:
		es.KeyDown(event.KeyDownArrow)
	case termbox.KeyArrowLeft:
		es.KeyDown(event.KeyLeftArrow)
	case termbox.KeyArrowRight:
		es.KeyDown(event.KeyRightArrow)

	// Termbox list of supported ctrl characters is weird ....
	case termbox.KeyCtrl2:
		ctrl("2")
	case termbox.KeyCtrl3:
		ctrl("3")
	case termbox.KeyCtrl4:
		ctrl("4")
	case termbox.KeyCtrl5:
		ctrl("5")
	case termbox.KeyCtrl6:
		ctrl("6")
	case termbox.KeyCtrl7:
		ctrl("7")
	case termbox.KeyCtrl8:
		ctrl("8")
	case termbox.KeyCtrlA:
		ctrl("a")
	case termbox.KeyCtrlB:
		ctrl("b")
	case termbox.KeyCtrlC:
		ctrl("c")
	case termbox.KeyCtrlD:
		ctrl("d")
	case termbox.KeyCtrlE:
		ctrl("e")
	case termbox.KeyCtrlF:
		ctrl("f")
	case termbox.KeyCtrlG:
		ctrl("g")
	case termbox.KeyCtrlJ:
		ctrl("j")
	case termbox.KeyCtrlK:
		ctrl("k")
	case termbox.KeyCtrlL:
		ctrl("l")
	case termbox.KeyCtrlN:
		ctrl("n")
	case termbox.KeyCtrlO:
		ctrl("o")
	case termbox.KeyCtrlP:
		ctrl("p")
	case termbox.KeyCtrlQ:
		ctrl("q")
	case termbox.KeyCtrlR:
		ctrl("r")
	case termbox.KeyCtrlS:
		ctrl("s")
	case termbox.KeyCtrlT:
		ctrl("t")
	case termbox.KeyCtrlU:
		ctrl("u")
	case termbox.KeyCtrlV:
		ctrl("v")
	case termbox.KeyCtrlW:
		ctrl("w")
	case termbox.KeyCtrlX:
		ctrl("x")
	case termbox.KeyCtrlY:
		ctrl("y")
	case termbox.KeyCtrlZ:
		ctrl("z")

		// hu ?? all those are duplicated values in termbox .....
		//case termbox.KeyCtrlH:
		//	ctrl("h")
		//case termbox.KeyCtrlI:
		//ctrl("i")
		//case termbox.KeyCtrlM:
		//ctrl("m")
		//case termbox.KeyCtrlSpace:
		//	ctrl(" ")
		//case termbox.KeyCtrlTilde:
		//	ctrl("~")
		//case termbox.KeyCtrlLsqBracket:
		//	ctrl("[")
		//case termbox.KeyCtrlRsqBracket:
		//	ctrl("]")
		//case termbox.KeyCtrlBackslash:
		//	ctrl("\\")
		//case termbox.KeyCtrlSlash:
		//	ctrl("/")
		//case termbox.KeyCtrlUnderscore:
		//	ctrl("_")
	}
}

// ############# TO BE OBSOLETED

func (t *TermBox) old() {
	t.SetMouseMode(termbox.MouseMotion)
	// Note: terminal might not support SGR mouse events, but trying anyway
	t.SetMouseMode(termbox.MouseSgr)

	t.SetInputMode(termbox.InputMouse)

	e := core.Ed.(*Editor)
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyCtrlQ {
			if !actions.Ar.EdQuitCheck() {
				actions.Ar.EdSetStatusErr("Unsaved changes. Save or request close again.")
			} else {
				return // that's all falks, quit
			}
		}
		switch ev.Type {
		case termbox.EventResize:
			actions.Ar.EdResize(ev.Height, ev.Width)
			// TODO: for terminal view, tell terminal about new size
		case termbox.EventKey:
			// Special case for terminal/command views, pass most keystrokes through
			v := viewCast(e.CurView())
			if v != nil && v.viewType == core.ViewTypeShell {
				v.TermEvent(e, &ev)
				break
			}

			switch ev.Key {
			// "Global" keys
			case termbox.KeyEsc:
				actions.Ar.CmdbarToggle()
			default:
				if e.cmdOn {
					e.Cmdbar.Event(e, &ev)
				} else if e.CurView != nil {
					viewCast(e.CurView()).Event(e, &ev)
				}
			}
		case termbox.EventMouse:
			w := e.WidgetAt(ev.MouseY, ev.MouseX)
			if w != nil {
				w.MouseEvent(e, &ev)
			}
		}
		actions.Ar.EdRender()
	}
}

// ##################### CmdBar ########################################

// Event handler for Cmdbar
func (c *Cmdbar) Event(e *Editor, ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if len(c.Cmd) > 0 {
			c.Cmd = c.Cmd[:len(c.Cmd)-1]
		}
	case termbox.KeyEnter:
		c.RunCmd()
	default:
		if ev.Ch != 0 && ev.Mod == 0 { // otherwise special key combo
			c.Cmd = append(c.Cmd, ev.Ch)
		}
	}
}

func (s *Cmdbar) MouseEvent(e *Editor, ev *termbox.Event) {
	switch ev.Key {
	case termbox.MouseLeft:
		if isMouseDown(ev) && !e.cmdOn {
			actions.Ar.CmdbarToggle()
		}
	}
}

// ##################### StatusBar ########################################

func (s *Statusbar) Event(e *Editor, ev *termbox.Event) {
}

func (s *Statusbar) MouseEvent(e *Editor, ev *termbox.Event) {
}

// ##################### View       ########################################

// Event handler for View
func (v *View) Event(e *Editor, ev *termbox.Event) {
	dirty := false
	es := false //expand selection
	vid := v.Id()
	actions.Ar.ViewAutoScroll(vid, 0, 0)
	//ln, col := actions.Ar.ViewCursorCoords(vid)
	e.evtState.InDrag = false

	if v.viewCommonEvent(e, ev) {
		return
	}

	switch ev.Key {
	// Ctrl combos
	case termbox.KeyCtrlA:
		actions.Ar.ViewSelectAll(vid)
		return
	case termbox.KeyCtrlC:
		actions.Ar.ViewCopy(vid)
	//case termbox.KeyCtrlF:
	//	actions.External("search.ank")
	case termbox.KeyCtrlQ:
		return
	case termbox.KeyCtrlR:
		actions.Ar.ViewReload(vid)
	case termbox.KeyCtrlS:
		actions.Ar.ViewSave(vid)
	case termbox.KeyCtrlV:
		actions.Ar.ViewPaste(vid)
		dirty = true
	case termbox.KeyCtrlW:
		actions.Ar.EdDelView(e.curViewId, true)
		return
	case termbox.KeyCtrlX:
		actions.Ar.ViewCut(vid)
		dirty = true
	case termbox.KeyCtrlY:
		actions.Ar.ViewRedo(vid)
	case termbox.KeyCtrlZ:
		actions.Ar.ViewUndo(vid)
	// "Regular" keys
	case termbox.KeyArrowRight:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtRight)
		es = true
	case termbox.KeyArrowLeft:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtLeft)
		es = true
	case termbox.KeyArrowUp:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtUp)
		es = true
	case termbox.KeyArrowDown:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtDown)
		es = true
	case termbox.KeyPgdn:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgDown)
		es = true
	case termbox.KeyPgup:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgUp)
		es = true
	case termbox.KeyEnd:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtEnd)
		es = true
	case termbox.KeyHome:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtHome)
		es = true
	case termbox.KeyTab:
		actions.Ar.ViewInsertCur(vid, "\t")
		dirty = true
		es = true
	case termbox.KeyEnter:
		actions.Ar.ViewInsertNewLine(vid)
		dirty = true
	case termbox.KeyDelete:
		actions.Ar.ViewDeleteCur(vid)
		dirty = true
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		actions.Ar.ViewBackspace(vid)
		dirty = true
	default:
		// insert the key
		if ev.Ch != 0 && ev.Mod == 0 && ev.Meta == 0 { // otherwise some special key combo
			actions.Ar.ViewInsertCur(vid, string(ev.Ch))
			dirty = true
		}
	}
	// extend keyboard selection
	if es && ev.Meta == termbox.Shift {
		//		actions.Ar.ViewStretchSelection(vid, ln, col)
	} else {
		actions.Ar.ViewClearSelections(vid)
	}

	if dirty {
		actions.Ar.ViewSetDirty(vid, true)
	}
}

func (v *View) viewCommonEvent(e *Editor, ev *termbox.Event) bool {
	vid := v.Id()
	// alt combos
	if ev.Mod == termbox.ModAlt {
		switch ev.Ch {
		case 'o':
			actions.Ar.ViewOpenSelection(vid, false)
			return true
		}
		return true
	}

	// Combos not supported directly by termbox
	if ev.Meta == termbox.Ctrl {
		switch ev.Key {
		case termbox.KeyArrowDown:
			actions.Ar.EdViewNavigate(core.CursorMvmtDown)
			return true
		case termbox.KeyArrowUp:
			actions.Ar.EdViewNavigate(core.CursorMvmtUp)
			return true
		case termbox.KeyArrowLeft:
			actions.Ar.EdViewNavigate(core.CursorMvmtLeft)
			return true
		case termbox.KeyArrowRight:
			actions.Ar.EdViewNavigate(core.CursorMvmtRight)
			return true
		}
	}
	switch ev.Key {
	case termbox.KeyCtrlO:
		actions.Ar.ViewOpenSelection(vid, true)
		return true
	case termbox.KeyCtrlT:
		v := actions.Ar.EdOpenTerm([]string{core.Terminal})
		actions.Ar.EdActivateView(v)
		return true
	}
	return false
}

// Events for terminal/command views
func (v *View) TermEvent(e *Editor, ev *termbox.Event) {
	vid := v.Id()
	actions.Ar.ViewAutoScroll(vid, 0, 0)
	e.evtState.InDrag = false
	bc := v.backend.(*backend.BackendCmd)
	es := false
	//ln, col := actions.Ar.ViewCursorCoords(vid)

	if v.viewCommonEvent(e, ev) {
		return
	}
	// Handle termbox special keys to VT100
	switch ev.Key {
	case 0: // "normal" character
		actions.Ar.ViewInsertCur(vid, string(ev.Ch))
	case termbox.KeyCtrlC:
		if len(v.selections) > 0 {
			actions.Ar.ViewCopy(vid)
		} else {
			bc.SendBytes([]byte{byte(ev.Key)})
		}
	case termbox.KeyCtrlV:
		actions.Ar.ViewPaste(vid)
	// "special"/navigation keys
	case termbox.KeyDelete:
		bc.SendBytes([]byte{27, 'O', 'C'}) // move right
		bc.SendBytes([]byte{127})          // delete (~ backspace)
	case termbox.KeyArrowUp:
		bc.SendBytes([]byte{27, 'O', 'A'})
	case termbox.KeyArrowDown:
		bc.SendBytes([]byte{27, 'O', 'B'})
	case termbox.KeyArrowRight:
		bc.SendBytes([]byte{27, 'O', 'C'})
	case termbox.KeyArrowLeft:
		bc.SendBytes([]byte{27, 'O', 'D'})
	case termbox.KeyPgdn:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgDown)
		es = true
	case termbox.KeyPgup:
		actions.Ar.ViewCursorMvmt(vid, core.CursorMvmtPgUp)
		es = true
	case termbox.KeyEnd:
		bc.SendBytes([]byte{byte(termbox.KeyCtrlE)})
		es = true
	case termbox.KeyHome:
		bc.SendBytes([]byte{byte(termbox.KeyCtrlA)})
		es = true
	// function keys
	case termbox.KeyF1:
		bc.SendBytes([]byte{27, 'O', 'P'})
	case termbox.KeyF2:
		bc.SendBytes([]byte{27, 'O', 'Q'})
	case termbox.KeyF3:
		bc.SendBytes([]byte{27, 'O', 'R'})
	case termbox.KeyF4:
		bc.SendBytes([]byte{27, 'O', 'S'})
	case termbox.KeyF5:
		bc.SendBytes([]byte{27, '[', '1', '5', '~'})
	case termbox.KeyF6:
		bc.SendBytes([]byte{27, '[', '1', '7', '~'})
	case termbox.KeyF7:
		bc.SendBytes([]byte{27, '[', '1', '8', '~'})
	case termbox.KeyF8:
		bc.SendBytes([]byte{27, '[', '1', '9', '~'})
	case termbox.KeyF9:
		bc.SendBytes([]byte{27, '[', '2', '0', '~'})
	case termbox.KeyF10:
		bc.SendBytes([]byte{27, '[', '2', '1', '~'})
	case termbox.KeyF11:
		bc.SendBytes([]byte{27, '[', '2', '3', '~'})
	case termbox.KeyF12:
		bc.SendBytes([]byte{27, '[', '2', '4', '~'})
	default: // some other special character? pass through for now
		b := []byte{}
		if ev.Key > 256 {
			b = append(b, byte(ev.Key/256))
		}
		b = append(b, byte(ev.Key%256))
		bc.SendBytes(b)
	}

	// extend keyboard selection
	if es && ev.Meta == termbox.Shift {
		//		actions.Ar.ViewStretchSelection(vid, ln, col)
	} else {
		actions.Ar.ViewClearSelections(vid)
	}
}

func (v *View) MouseEvent(e *Editor, ev *termbox.Event) {
	vid := v.Id()
	_, y, x := actions.Ar.EdViewAt(ev.MouseY+1, ev.MouseX+1)
	ln, col := actions.Ar.ViewTextPos(vid, y, x)
	if isMouseDown(ev) && ev.MouseX == e.evtState.LastClickX &&
		ev.MouseY == e.evtState.LastClickY &&
		time.Now().Unix()-e.evtState.LastLeftClick <= 2 {
		ev.Key = MouseLeftDbl
		e.evtState.LastClickX = -1
	}
	switch ev.Key {
	case MouseLeftDbl:
		if ev.MouseX == v.x1 && ev.MouseY == v.y1 {
			// view swap
			actions.Ar.EdSwapViews(e.CurViewId(), vid)
			actions.Ar.EdActivateView(vid)
			e.evtState.MovingView = false
			actions.Ar.EdSetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), vid))
			return
		}
		if ev.MouseY == v.y1 && ev.MouseX > v.x1 {
			// TODO : collapse / expand view
			return
		}
		if selection := v.ExpandSelectionWord(ln, col); selection != nil {
			v.selections = []core.Selection{
				*selection,
			}
		}
		return
	case termbox.MouseScrollUp:
		actions.Ar.ViewMoveCursor(vid, -1, 0, false)
		return
	case termbox.MouseScrollDown:
		actions.Ar.ViewMoveCursor(vid, 1, 0, false)
		return
	case termbox.MouseRight:
		if isMouseDown(ev) {
			e.evtState.InDrag = false
			e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
			e.evtState.LastRightClick = time.Now().Unix()
			actions.Ar.ViewClearSelections(vid)
			actions.Ar.ViewSetCursorPos(vid, ln, col)
			actions.Ar.ViewOpenSelection(vid, true)
		}
		return
	case termbox.MouseLeft:
		// end view move
		if e.evtState.MovingView && isMouseDown(ev) {
			e.evtState.MovingView = false
			actions.Ar.EdViewMove(e.evtState.LastClickY+1, e.evtState.LastClickX+1, ev.MouseY+1, ev.MouseX+1)
			actions.Ar.EdSetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), vid))
			return
		}
		// 'X' button (close view)
		if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 && isMouseDown(ev) {
			actions.Ar.EdDelView(vid, true)
			return
		}
		// start view move
		if ev.MouseX == v.x1 && ev.MouseY == v.y1 && isMouseDown(ev) {
			// handle
			e.evtState.MovingView = true
			e.evtState.LastClickX = ev.MouseX
			e.evtState.LastClickY = ev.MouseY
			e.evtState.LastLeftClick = time.Now().Unix()
			actions.Ar.EdSetStatusErr("Starting move, click new position or dbl click to swap")
			return
		}
		if ev.MouseX <= v.x1 {
			return // scrollbar TBD
		}
		if ev.DragOn {
			if !e.evtState.InDrag {
				e.evtState.InDrag = true
				actions.Ar.ViewClearSelections(vid)
				actions.Ar.ViewSetCursorPos(vid, e.evtState.DragLn+1, e.evtState.DragCol+1)
				actions.Ar.EdActivateView(vid)
			}
			// continued drag
			x1 := e.evtState.DragCol
			y1 := e.evtState.DragLn
			x2 := col - 1
			y2 := ln - 1

			actions.Ar.ViewClearSelections(vid)
			actions.Ar.ViewAddSelection(
				vid,
				y1,
				v.LineRunesTo(v.slice, y1, x1)+1,
				y2,
				v.LineRunesTo(v.slice, y2, x2)+1)

			// Handling scrolling while dragging
			if ln < v.offy { // scroll up
				actions.Ar.ViewAutoScroll(vid, -v.LineCount()/10, 0)
			} else if ln >= v.offy+(v.y2-v.y1)-2 { // scroll down
				actions.Ar.ViewAutoScroll(vid, v.LineCount()/10, 0)
			} else if col < v.offx { //scroll left
				actions.Ar.ViewAutoScroll(vid, 0, -5)
			} else if col >= v.offx+(v.x2-v.x1)-3 { // scroll right
				actions.Ar.ViewAutoScroll(vid, 0, 5)
			}
			return
		}

		if isMouseDown(ev) { // click
			if selected, _ := v.Selected(ln, col); selected {
				e.evtState.InDrag = false
				// otherwise it could be the MouseDown at the end of a drag.
			}
			if !e.evtState.InDrag {
				actions.Ar.ViewClearSelections(vid)
				actions.Ar.ViewSetCursorPos(vid, ln, col)
				actions.Ar.EdActivateView(vid)
				e.evtState.LastLeftClick = time.Now().Unix()
				e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
				actions.Ar.EdSetStatus(fmt.Sprintf("%s  [%d] %d:%d %d:%d", v.WorkDir(), vid, y, x, ln, col))
			}
		}
		e.evtState.InDrag = false
		actions.Ar.CmdbarEnable(false)
		e.evtState.DragLn = ln
		e.evtState.DragCol = col
	}
}

func isMouseDown(ev *termbox.Event) bool {
	return ev.MouseBtnState == termbox.MouseBtnUp
}
