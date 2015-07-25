package ui

import (
	"fmt"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/termbox-go"
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

// EventLoop is the main event loop that keeps waiting for events as long as
// the editor is running.
func (e *Editor) EventLoop() {

	e.term.SetMouseMode(termbox.MouseMotion)
	// Note: terminal might not support SGR mouse events, but trying anyway
	e.term.SetMouseMode(termbox.MouseSgr)

	e.term.SetInputMode(termbox.InputMouse)
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			actions.EdResize(ev.Height, ev.Width)
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				if !actions.EdQuitCheck() {
					actions.EdSetStatusErr("Unsaved changes. Save or request close again.")
				} else {
					return // that's all falks, quit
				}
			case termbox.KeyEsc:
				actions.CmdbarToggle()
			default:
				if e.cmdOn {
					e.Cmdbar.Event(e, &ev)
				} else if e.CurView != nil {
					e.curView.Event(e, &ev)
				}
			}
		case termbox.EventMouse:
			w := e.WidgetAt(ev.MouseY, ev.MouseX)
			if w != nil {
				w.Event(e, &ev)
			}
		}
		actions.EdRender()
	}
}

// ##################### CmdBar ########################################

// Event handler for Cmdbar
func (c *Cmdbar) Event(e *Editor, ev *termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
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

	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			if isMouseUp(ev) && !e.cmdOn {
				actions.CmdbarToggle()
			}
		}
	}
}

// ##################### StatusBar ########################################

// Event handler for Statusbar
func (s *Statusbar) Event(e *Editor, ev *termbox.Event) {
	// Anything ??
}

// ##################### View       ########################################

// Event handler for View
func (v *View) Event(e *Editor, ev *termbox.Event) {
	dirty := false
	es := false //expand selection
	vid := v.Id()
	actions.ViewAutoScroll(vid, 0, 0, false)
	switch ev.Type {
	case termbox.EventKey:
		ln, col := actions.ViewCurPos(vid)
		e.evtState.InDrag = false
		// alt combo
		if ev.Mod == termbox.ModAlt {
			switch ev.Ch {
			case 'o':
				actions.ViewOpenSelection(vid, false)
				return
			}
		}

		switch ev.Key {
		case termbox.KeyArrowRight:
			actions.ViewCursorMvmt(vid, core.CursorMvmtRight)
			es = true
		case termbox.KeyArrowLeft:
			actions.ViewCursorMvmt(vid, core.CursorMvmtLeft)
			es = true
		case termbox.KeyArrowUp:
			actions.ViewCursorMvmt(vid, core.CursorMvmtUp)
			es = true
		case termbox.KeyArrowDown:
			actions.ViewCursorMvmt(vid, core.CursorMvmtDown)
			es = true
		case termbox.KeyPgdn:
			actions.ViewCursorMvmt(vid, core.CursorMvmtPgDown)
			es = true
		case termbox.KeyPgup:
			actions.ViewCursorMvmt(vid, core.CursorMvmtPgUp)
			es = true
		case termbox.KeyEnd:
			actions.ViewCursorMvmt(vid, core.CursorMvmtEnd)
			es = true
		case termbox.KeyHome:
			actions.ViewCursorMvmt(vid, core.CursorMvmtHome)
			es = true
		case termbox.KeyTab:
			actions.ViewInsertCur(vid, "\t")
			dirty = true
			es = true
		case termbox.KeyEnter:
			actions.ViewInsertNewLine(vid)
			dirty = true
		case termbox.KeyDelete:
			actions.ViewDeleteCur(vid)
			dirty = true
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			actions.ViewBackspace(vid)
			dirty = true
		case termbox.KeyCtrlS:
			actions.ViewSave(vid)
		case termbox.KeyCtrlO:
			actions.ViewOpenSelection(vid, true)
		case termbox.KeyCtrlW:
			actions.EdDelViewCheck(e.curView.Id())
		case termbox.KeyCtrlC:
			switch v.backend.(type) {
			case *backend.BackendCmd:
				// CTRL+C process
				if v.backend.(*backend.BackendCmd).Running() {
					actions.ViewCmdStop(vid)
					return
				}
			}
			actions.ViewCopy(vid)
		case termbox.KeyCtrlX:
			actions.ViewCut(vid)
			dirty = true
		case termbox.KeyCtrlV:
			actions.ViewPaste(vid)
			dirty = true
		case termbox.KeyCtrlQ:
			return
		case termbox.KeyCtrlR:
			actions.ViewReload(vid)
		case termbox.KeyCtrlF:
			//			actions.ViewSearch(vid)
		default:
			// insert the key
			if ev.Ch != 0 && ev.Mod == 0 && ev.Meta == 0 { // otherwise some special key combo
				actions.ViewInsertCur(vid, string(ev.Ch))
				dirty = true
			}
		}
		// extend keyboard selection
		if es && ev.Meta == termbox.Shift {
			actions.ViewStretchSelection(vid, ln, col)
		} else {
			actions.ViewClearSelections(vid)
		}
	case termbox.EventMouse:
		col := ev.MouseX - v.x1 + v.offx - 2
		ln := ev.MouseY - v.y1 + v.offy - 2
		if isMouseUp(ev) && ev.MouseX == e.evtState.LastClickX &&
			ev.MouseY == e.evtState.LastClickY &&
			time.Now().Unix()-e.evtState.LastLeftClick <= 2 {
			ev.Key = MouseLeftDbl
			e.evtState.LastClickX = -1
		}
		switch ev.Key {
		case MouseLeftDbl:
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 {
				actions.EdSwapViews(e.CurView().Id(), vid)
				actions.EdActivateView(vid, v.CurLine(), v.CurCol())
				e.evtState.MovingView = false
				actions.EdSetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), vid))
				return
			}
			if selection := v.ExpandSelectionWord(ln, col); selection != nil {
				v.selections = []core.Selection{
					*selection,
				}
			}
			return
		case termbox.MouseScrollUp:
			actions.ViewMoveCursor(vid, -1, 0)
			return
		case termbox.MouseScrollDown:
			actions.ViewMoveCursor(vid, 1, 0)
			return
		case termbox.MouseRight:
			if isMouseUp(ev) {
				e.evtState.InDrag = false
				e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
				e.evtState.LastRightClick = time.Now().Unix()
				actions.ViewClearSelections(vid)
				actions.ViewMoveCursor(vid, ev.MouseY-v.y1-2-v.CursorY, ev.MouseX-v.x1-2-v.CursorX)
				actions.ViewOpenSelection(vid, true)
			}
			return
		case termbox.MouseLeft:
			if e.evtState.MovingView && isMouseUp(ev) {
				e.evtState.MovingView = false
				actions.EdViewMove(vid, e.evtState.LastClickY, e.evtState.LastClickX, ev.MouseY, ev.MouseX)
				actions.EdSetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), vid))
				return
			}
			if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				actions.EdDelViewCheck(vid)
				return
			}
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				// handle
				e.evtState.MovingView = true
				e.evtState.LastClickX = ev.MouseX
				e.evtState.LastClickY = ev.MouseY
				e.evtState.LastLeftClick = time.Now().Unix()
				actions.EdSetStatusErr("Starting move, click new position or dbl click to swap")
				return
			}
			if ev.MouseX <= v.x1 {
				return // scrollbar TBD
			}
			if ev.DragOn {
				if !e.evtState.InDrag {
					e.evtState.InDrag = true
					actions.ViewClearSelections(vid)
					actions.EdActivateView(vid, e.evtState.DragLn, e.evtState.DragCol)
				}
				// continued drag
				x1 := e.evtState.DragCol
				y1 := e.evtState.DragLn
				x2 := col
				y2 := ln

				actions.ViewClearSelections(vid)
				actions.ViewAddSelection(
					vid,
					y1,
					v.LineRunesTo(v.slice, y1, x1),
					y2,
					v.LineRunesTo(v.slice, y2, x2))

				// Handling scrolling while dragging
				if ln < v.offy { // scroll up
					actions.ViewAutoScroll(vid, -v.LineCount()/10, 0, true)
				} else if ln >= v.offy+(v.y2-v.y1)-2 { // scroll down
					actions.ViewAutoScroll(vid, v.LineCount()/10, 0, true)
				} else if col < v.offx { //scroll left
					actions.ViewAutoScroll(vid, 0, -5, true)
				} else if col >= v.offx+(v.x2-v.x1)-3 { // scroll right
					actions.ViewAutoScroll(vid, 0, 5, true)
				}
				return
			}

			if isMouseUp(ev) { // click
				if selected, _ := v.Selected(ln, col); selected {
					e.evtState.InDrag = false
					// otherwise it could be the mouseUp at the end of a drag.
				}
				if !e.evtState.InDrag {
					actions.ViewClearSelections(vid)
					actions.EdActivateView(vid, ln, col)
					e.evtState.LastLeftClick = time.Now().Unix()
					e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
					actions.EdSetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), vid))
				}
			}
			e.evtState.InDrag = false
			actions.CmdbarEnable(false)
			e.evtState.DragLn = ln
			e.evtState.DragCol = col
		} // end switch
	}

	if dirty {
		actions.ViewSetDirty(vid, true)
	}
}

func isMouseUp(ev *termbox.Event) bool {
	return ev.MouseBtnState == termbox.MouseBtnUp
}
