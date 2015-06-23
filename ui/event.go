package ui

import (
	"fmt"

	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/termbox-go"
)

// Evtstate stores some state about kb/mouse events
type EvtState struct {
	MovingView      bool
	X, Y            int
	DragLn, DragCol int
}

func (e *Editor) EventLoop() {

	e.term.SetMouseMode(termbox.MouseMotion)
	// Note: terminal might not support SGR mouse events, but trying anyway
	e.term.SetMouseMode(termbox.MouseSgr)

	e.term.SetInputMode(termbox.InputMouse)
	for {
		ev := termbox.PollEvent()
		//log.Printf("Event : %d %d", 0xFFFF-ev.Key, ev.Meta)
		switch ev.Type {
		case termbox.EventResize:
			e.Resize(ev.Width, ev.Height)
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				if !e.QuitCheck() {
					e.SetStatusErr("Unsaved changes. Save or request close again.")
				} else {
					return // the end
				}
			case termbox.KeyEsc:
				if !e.cmdOn {
					e.Cmdbar.Cmd = []rune{}
				}
				e.cmdOn = !e.cmdOn
			default:
				if e.cmdOn {
					e.Cmdbar.Event(e, &ev)
				} else if e.CurView != nil {
					e.curView.Event(e, &ev)
				}
			}
		case termbox.EventMouse:
			w := e.WidgetAt(ev.MouseX, ev.MouseY)
			if w != nil {
				w.Event(e, &ev)
			}
		}
		e.Render()
	}
}

// ##################### CmdBar ########################################

// Event handler for Cmdbar
func (m *Cmdbar) Event(e *Editor, ev *termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			if len(m.Cmd) > 0 {
				m.Cmd = m.Cmd[:len(m.Cmd)-1]
			}
		case termbox.KeyEnter:
			m.RunCmd()
		default:
			if ev.Ch != 0 && ev.Mod == 0 { // otherwise special key combo
				m.Cmd = append(m.Cmd, ev.Ch)
			}
		}

	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			if isMouseUp(ev) {
				e.cmdOn = true
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
	ln, col := v.CurLine(), v.CurCol()
	es := false                  //expand selection
	v.SetAutoScroll(0, 0, false) // any events stops autoscroll
	switch ev.Type {
	case termbox.EventKey:
		// alt combo
		if ev.Mod == termbox.ModAlt {
			switch ev.Ch {
			case 'o':
				e.Cmdbar.OpenSelection(v, false)
				return
			}
		}

		switch ev.Key {
		case termbox.KeyArrowRight:
			offset := 1
			c, _, _ := v.CurChar()
			if c != nil {
				offset = v.runeSize(*c)
			}
			v.MoveCursorRoll(offset, 0)
			es = true
		case termbox.KeyArrowLeft:
			offset := 1
			c, _, _ := v.CursorChar(v.slice, v.CurCol()-1, v.CurLine())
			if c != nil {
				offset = v.runeSize(*c)
			}
			v.MoveCursorRoll(-offset, 0)
			es = true
		case termbox.KeyArrowUp:
			v.MoveCursor(0, -1)
			es = true
		case termbox.KeyArrowDown:
			v.MoveCursor(0, 1)
			es = true
		case termbox.KeyPgdn:
			dist := v.LastViewLine() + 1
			if v.LineCount()-v.CurLine() < dist {
				dist = v.LineCount() - v.CurLine() - 1
			}
			v.MoveCursor(0, dist)
			es = true
		case termbox.KeyPgup:
			dist := v.LastViewLine() + 1
			if dist > v.CurLine() {
				dist = v.CurLine()
			}
			v.MoveCursor(0, -dist)
			es = true
		case termbox.KeyEnd:
			v.MoveCursor(v.lineCols(v.slice, v.CurLine())-v.CurCol(), 0)
			es = true
		case termbox.KeyHome:
			v.MoveCursor(-v.CurCol(), 0)
			es = true
		case termbox.KeyTab:
			v.InsertCur("\t")
			dirty = true
			es = true
		case termbox.KeyEnter:
			v.InsertNewLineCur()
			dirty = true
		case termbox.KeyDelete:
			v.DeleteCur()
			dirty = true
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			v.Backspace()
			dirty = true
		case termbox.KeyCtrlS:
			v.Save()
		case termbox.KeyCtrlO:
			e.Cmdbar.OpenSelection(v, true)
		case termbox.KeyCtrlW:
			e.DelViewCheck(e.curView)
		case termbox.KeyCtrlC:
			switch v.backend.(type) {
			case *backend.BackendCmd:
				// CTRL+C process
				if v.backend.(*backend.BackendCmd).Running() {
					v.backend.Close()
					return
				}
			}
			// copy
			if len(v.selections) > 0 {
				v.SelectionCopy(&v.selections[0])
			}
		case termbox.KeyCtrlX:
			if len(v.selections) > 0 {
				v.SelectionCopy(&v.selections[0])
				v.SelectionDelete(&v.selections[0])
				v.ClearSelections()
			}
		case termbox.KeyCtrlV:
			v.Paste()
			dirty = true
		case termbox.KeyCtrlQ:
			return
		case termbox.KeyCtrlR:
			v.Reload()
		case termbox.KeyCtrlF:
			if len(v.selections) > 0 {
				text := core.RunesToString(v.SelectionText(&v.selections[0]))
				e.Cmdbar.Search(text)
			}
		default:
			// insert the key
			if ev.Ch != 0 && ev.Mod == 0 && ev.Meta == 0 { // otherwise some special key combo
				v.InsertCur(string(ev.Ch))
				dirty = true
			}
		}
		// extend keyboard slection
		if es && ev.Meta == termbox.Shift {
			v.ExpandSelection(
				ln+1,
				v.lineRunesTo(v.slice, ln, col)+1,
				v.CurLine()+1,
				v.lineRunesTo(v.slice, v.CurLine(), v.CurCol())+1,
			)
		} else {
			v.ClearSelections()
		}
	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseScrollUp:
			v.MoveCursor(0, -1)
		case termbox.MouseScrollDown:
			v.MoveCursor(0, 1)
		case termbox.MouseRight:
			if isMouseUp(ev) {
				v.ClearSelections()
				v.MoveCursor(ev.MouseX-v.x1-2-v.CursorX, ev.MouseY-v.y1-2-v.CursorY)
				e.Cmdbar.OpenSelection(v, true)
			}
		case termbox.MouseLeft:
			if e.evtState.MovingView && isMouseUp(ev) {
				e.evtState.MovingView = false
				e.ViewMove(e.evtState.X, e.evtState.Y, ev.MouseX, ev.MouseY)
				return
			}
			if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				e.DelViewCheck(v)
				return
			}
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				// handle
				e.evtState.MovingView = true
				e.evtState.X = ev.MouseX
				e.evtState.Y = ev.MouseY
				e.SetStatusErr("Starting move, click new position.")
				return
			}
			if ev.MouseX <= v.x1 {
				return // scrollbar TBD
			}
			col := ev.MouseX - v.x1 + v.offx - 1
			ln := ev.MouseY - v.y1 + v.offy - 1
			if ev.DragOn {
				// continued drag
				x1 := e.evtState.DragCol
				y1 := e.evtState.DragLn
				x2 := col
				y2 := ln

				s := core.NewSelection(
					y1,
					v.lineRunesTo(v.slice, y1-1, x1-1)+1,
					y2,
					v.lineRunesTo(v.slice, y2-1, x2-1)+1)
				v.selections = []core.Selection{
					*s,
				}
				// Handling scrolling while dragging
				if ln <= v.offy { // scroll up
					v.SetAutoScroll(0, -v.LineCount()/10, true)
				} else if ln >= v.offy+(v.y2-v.y1)-1 { // scroll down
					v.SetAutoScroll(0, v.LineCount()/10, true)
				} else if col <= v.offx { //scroll left
					v.SetAutoScroll(-5, 0, true)
				} else if col >= v.offx+(v.x2-v.x1)-2 { // scroll right
					v.SetAutoScroll(5, 0, true)
				}
				return
			} else {
				if !isMouseUp(ev) { // reset drag
					v.ClearSelections()
				}
				e.evtState.DragLn = ln
				e.evtState.DragCol = col
			}
			if isMouseUp(ev) {
				e.cmdOn = false
				e.ActivateView(v, ev.MouseX-v.x1-2+v.offx, ev.MouseY-v.y1-2+v.offy)
				e.SetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), v.Id()))
			}
		}
	}

	if dirty {
		v.SetDirty(true)
	}
}

func isMouseUp(ev *termbox.Event) bool {
	return ev.MouseBtnState == termbox.MouseBtnUp
}
