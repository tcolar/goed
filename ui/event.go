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
			actions.EdResizeAction(ev.Height, ev.Width)
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				if !e.QuitCheck() {
					actions.EdSetStatusErrAction("Unsaved changes. Save or request close again.")
				} else {
					return // that's all falks, quit
				}
			case termbox.KeyEsc:
				actions.CmdbarToggleAction()
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
		actions.EdRenderAction()
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
			actions.CmdbarRunAction()
		default:
			if ev.Ch != 0 && ev.Mod == 0 { // otherwise special key combo
				c.Cmd = append(c.Cmd, ev.Ch)
			}
		}

	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			if isMouseUp(ev) && !e.cmdOn {
				actions.CmdbarToggleAction()
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
	actions.ViewAutoScrollAction(v.Id(), 0, 0, false)
	switch ev.Type {
	case termbox.EventKey:
		ln, col := v.CurLine(), v.CurCol()
		e.evtState.InDrag = false
		// alt combo
		if ev.Mod == termbox.ModAlt {
			switch ev.Ch {
			case 'o':
				actions.ViewOpenSelectionAction(v.Id(), false)
				return
			}
		}

		switch ev.Key {
		case termbox.KeyArrowRight:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtRight)
			es = true
		case termbox.KeyArrowLeft:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtLeft)
			es = true
		case termbox.KeyArrowUp:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtUp)
			es = true
		case termbox.KeyArrowDown:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtDown)
			es = true
		case termbox.KeyPgdn:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtPgDown)
			es = true
		case termbox.KeyPgup:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtPgUp)
			es = true
		case termbox.KeyEnd:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtEnd)
			es = true
		case termbox.KeyHome:
			actions.ViewCursorMvmtAction(v.Id(), core.CursorMvmtHome)
			es = true
		case termbox.KeyTab:
			actions.ViewInsertCurAction(v.Id(), "\t")
			dirty = true
			es = true
		case termbox.KeyEnter:
			actions.ViewInsertNewLineAction(v.Id())
			dirty = true
		case termbox.KeyDelete:
			actions.ViewDeleteCurAction(v.Id())
			dirty = true
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			actions.ViewBackspaceAction(v.Id())
			dirty = true
		case termbox.KeyCtrlS:
			actions.ViewSaveAction(v.Id())
		case termbox.KeyCtrlO:
			actions.ViewOpenSelectionAction(v.Id(), true)
		case termbox.KeyCtrlW:
			actions.EdDelViewCheckAction(e.curView.Id())
		case termbox.KeyCtrlC:
			switch v.backend.(type) {
			case *backend.BackendCmd:
				// CTRL+C process
				if v.backend.(*backend.BackendCmd).Running() {
					actions.ViewCmdStopAction(v.Id())
					return
				}
			}
			actions.ViewCopyAction(v.Id())
		case termbox.KeyCtrlX:
			actions.ViewCutAction(v.Id())
			dirty = true
		case termbox.KeyCtrlV:
			actions.ViewPasteAction(v.Id())
			dirty = true
		case termbox.KeyCtrlQ:
			return
		case termbox.KeyCtrlR:
			actions.ViewReloadAction(v.Id())
		case termbox.KeyCtrlF:
			if len(v.selections) == 0 {
				v.SelectWord(ln, col)
			}
			if len(v.selections) > 0 {
				text := core.RunesToString(v.SelectionText(&v.selections[0]))
				e.Cmdbar.Search(text)
			}
		default:
			// insert the key
			if ev.Ch != 0 && ev.Mod == 0 && ev.Meta == 0 { // otherwise some special key combo
				actions.ViewInsertCurAction(v.Id(), string(ev.Ch))
				dirty = true
			}
		}
		// extend keyboard selection
		if es && ev.Meta == termbox.Shift {
			v.StretchSelection(
				ln,
				v.LineRunesTo(v.slice, ln, col),
				v.CurLine(),
				v.LineRunesTo(v.slice, v.CurLine(), v.CurCol()),
			)
		} else {
			actions.ViewClearSelectionsAction(v.Id())
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
				e.SwapViews(e.CurView().(*View), v)
				actions.EdActivateViewAction(v.Id(), v.CurLine(), v.CurCol())
				e.evtState.MovingView = false
				actions.EdSetStatusAction(fmt.Sprintf("%s  [%d]", v.WorkDir(), v.Id()))
				return
			}
			if selection := v.ExpandSelectionWord(ln, col); selection != nil {
				v.selections = []core.Selection{
					*selection,
				}
			}
			return
		case termbox.MouseScrollUp:
			actions.ViewMoveCursorAction(v.Id(), -1, 0)
			return
		case termbox.MouseScrollDown:
			actions.ViewMoveCursorAction(v.Id(), 1, 0)
			return
		case termbox.MouseRight:
			if isMouseUp(ev) {
				e.evtState.InDrag = false
				e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
				e.evtState.LastRightClick = time.Now().Unix()
				actions.ViewClearSelectionsAction(v.Id())
				actions.ViewMoveCursorAction(v.Id(), ev.MouseY-v.y1-2-v.CursorY, ev.MouseX-v.x1-2-v.CursorX)
				actions.ViewOpenSelectionAction(v.Id(), true)
			}
			return
		case termbox.MouseLeft:
			if e.evtState.MovingView && isMouseUp(ev) {
				e.evtState.MovingView = false
				actions.EdViewMoveAction(v.Id(), e.evtState.LastClickY, e.evtState.LastClickX, ev.MouseY, ev.MouseX)
				actions.EdSetStatusAction(fmt.Sprintf("%s  [%d]", v.WorkDir(), v.Id()))
				return
			}
			if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				actions.EdDelViewCheckAction(v.Id())
				return
			}
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				// handle
				e.evtState.MovingView = true
				e.evtState.LastClickX = ev.MouseX
				e.evtState.LastClickY = ev.MouseY
				e.evtState.LastLeftClick = time.Now().Unix()
				actions.EdSetStatusErrAction("Starting move, click new position or dbl click to swap")
				return
			}
			if ev.MouseX <= v.x1 {
				return // scrollbar TBD
			}
			if ev.DragOn {
				if !e.evtState.InDrag {
					e.evtState.InDrag = true
					actions.ViewClearSelectionsAction(v.Id())
					actions.EdActivateViewAction(v.Id(), e.evtState.DragLn, e.evtState.DragCol)
				}
				// continued drag
				x1 := e.evtState.DragCol
				y1 := e.evtState.DragLn
				x2 := col
				y2 := ln

				s := core.NewSelection(
					y1,
					v.LineRunesTo(v.slice, y1, x1),
					y2,
					v.LineRunesTo(v.slice, y2, x2))
				v.selections = []core.Selection{
					*s,
				}
				// Handling scrolling while dragging
				if ln < v.offy { // scroll up
					actions.ViewAutoScrollAction(v.Id(), -v.LineCount()/10, 0, true)
				} else if ln >= v.offy+(v.y2-v.y1)-2 { // scroll down
					actions.ViewAutoScrollAction(v.Id(), v.LineCount()/10, 0, true)
				} else if col < v.offx { //scroll left
					actions.ViewAutoScrollAction(v.Id(), 0, -5, true)
				} else if col >= v.offx+(v.x2-v.x1)-3 { // scroll right
					actions.ViewAutoScrollAction(v.Id(), 0, 5, true)
				}
				return
			}

			if isMouseUp(ev) { // click
				if selected, _ := v.Selected(ln, col); selected {
					e.evtState.InDrag = false
					// otherwise it could be the mouseUp at the end of a drag.
				}
				if !e.evtState.InDrag {
					actions.ViewClearSelectionsAction(v.Id())
					actions.EdActivateViewAction(v.Id(), ln, col)
					e.evtState.LastLeftClick = time.Now().Unix()
					e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
					actions.EdSetStatusAction(fmt.Sprintf("%s  [%d]", v.WorkDir(), v.Id()))
				}
			}
			e.evtState.InDrag = false
			e.cmdOn = false
			e.evtState.DragLn = ln
			e.evtState.DragCol = col
		} // end switch
	}

	if dirty {
		actions.ViewSetDirtyAction(v.Id(), true)
	}
}

func isMouseUp(ev *termbox.Event) bool {
	return ev.MouseBtnState == termbox.MouseBtnUp
}
