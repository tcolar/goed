package ui

import (
	"fmt"
	"time"

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
		//log.Printf("Event : %d %d", 0xFFFF-ev.Key, ev.Meta)
		switch ev.Type {
		case termbox.EventResize:
			e.Resize(ev.Height, ev.Width)
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
			w := e.WidgetAt(ev.MouseY, ev.MouseX)
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

	es := false                  //expand selection
	v.SetAutoScroll(0, 0, false) // any events stops autoscroll
	e.SetStatus(fmt.Sprintf("evt %d", ev.Type))
	switch ev.Type {
	case termbox.EventKey:
		ln, col := v.CurLine(), v.CurCol()
		e.evtState.InDrag = false
		// alt combo
		if ev.Mod == termbox.ModAlt {
			switch ev.Ch {
			case 'o':
				v.OpenSelection(false)
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
			v.MoveCursorRoll(0, offset)
			es = true
		case termbox.KeyArrowLeft:
			offset := 1
			c, _, _ := v.CursorChar(v.slice, ln, col-1)
			if c != nil {
				offset = v.runeSize(*c)
			}
			v.MoveCursorRoll(0, -offset)
			es = true
		case termbox.KeyArrowUp:
			v.MoveCursor(-1, 0)
			es = true
		case termbox.KeyArrowDown:
			v.MoveCursor(1, 0)
			es = true
		case termbox.KeyPgdn:
			dist := v.LastViewLine() + 1
			if v.LineCount()-ln < dist {
				dist = v.LineCount() - ln - 1
			}
			v.MoveCursor(dist, 0)
			es = true
		case termbox.KeyPgup:
			dist := v.LastViewLine() + 1
			if dist > ln {
				dist = ln
			}
			v.MoveCursor(-dist, 0)
			es = true
		case termbox.KeyEnd:
			v.MoveCursor(0, v.lineCols(v.slice, ln)-col)
			es = true
		case termbox.KeyHome:
			v.MoveCursor(0, -col)
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
			v.OpenSelection(true)
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
			if len(v.selections) == 0 {
				v.SelectLine(ln)
			}
			v.SelectionCopy(&v.selections[0])
		case termbox.KeyCtrlX:
			if len(v.selections) == 0 {
				v.SelectLine(ln)
			}
			v.SelectionCopy(&v.selections[0])
			v.SelectionDelete(&v.selections[0])
			v.ClearSelections()
		case termbox.KeyCtrlV:
			v.Paste()
			dirty = true
		case termbox.KeyCtrlQ:
			return
		case termbox.KeyCtrlR:
			v.Reload()
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
				v.InsertCur(string(ev.Ch))
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
			v.ClearSelections()
		}
	case termbox.EventMouse:
		col := ev.MouseX - v.x1 + v.offx - 2
		ln := ev.MouseY - v.y1 + v.offy - 2
		if isMouseUp(ev) && ev.MouseX == e.evtState.LastClickX &&
			ev.MouseY == e.evtState.LastClickY &&
			time.Now().Unix()-e.evtState.LastLeftClick < 2 {
			ev.Key = MouseLeftDbl
		}
		switch ev.Key {
		case MouseLeftDbl:
			if selection := v.ExpandSelectionWord(ln, col); selection != nil {
				v.selections = []core.Selection{
					*selection,
				}
			}
			return
		case termbox.MouseScrollUp:
			v.MoveCursor(-1, 0)
			return
		case termbox.MouseScrollDown:
			v.MoveCursor(1, 0)
			return
		case termbox.MouseRight:
			if isMouseUp(ev) {
				e.evtState.InDrag = false
				e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
				e.evtState.LastRightClick = time.Now().Unix()
				v.ClearSelections()
				v.MoveCursor(ev.MouseY-v.y1-2-v.CursorY, ev.MouseX-v.x1-2-v.CursorX)
				v.OpenSelection(true)
			}
			return
		case termbox.MouseLeft:
			if e.evtState.MovingView && isMouseUp(ev) {
				e.evtState.MovingView = false
				e.ViewMove(e.evtState.LastClickY, e.evtState.LastClickX, ev.MouseY, ev.MouseX)
				return
			}
			if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				e.DelViewCheck(v)
				return
			}
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 && isMouseUp(ev) {
				// handle
				e.evtState.MovingView = true
				e.evtState.LastClickX = ev.MouseX
				e.evtState.LastClickY = ev.MouseY
				e.SetStatusErr("Starting move, click new position.")
				return
			}
			if ev.MouseX <= v.x1 {
				return // scrollbar TBD
			}
			if ev.DragOn {
				if !e.evtState.InDrag {
					e.evtState.InDrag = true
					v.ClearSelections()
					e.ActivateView(v, e.evtState.DragCol, e.evtState.DragLn)
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
					v.SetAutoScroll(-v.LineCount()/10, 0, true)
				} else if ln >= v.offy+(v.y2-v.y1)-2 { // scroll down
					v.SetAutoScroll(v.LineCount()/10, 0, true)
				} else if col < v.offx { //scroll left
					v.SetAutoScroll(0, -5, true)
				} else if col >= v.offx+(v.x2-v.x1)-3 { // scroll right
					v.SetAutoScroll(0, 5, true)
				}
				return
			}

			if isMouseUp(ev) { // click
				if e.evtState.DragLn != ln || e.evtState.DragCol != col {
					e.evtState.InDrag = false
				}
				if !e.evtState.InDrag {
					v.ClearSelections()
					e.ActivateView(v, col, ln)
					e.evtState.LastLeftClick = time.Now().Unix()
					e.evtState.LastClickX, e.evtState.LastClickY = ev.MouseX, ev.MouseY
					e.SetStatus(fmt.Sprintf("%s  [%d]", v.WorkDir(), v.Id()))
				}
			}
			e.evtState.InDrag = false
			e.cmdOn = false
			e.evtState.DragLn = ln
			e.evtState.DragCol = col
		} // end switch
	}

	if dirty {
		v.SetDirty(true)
	}
}

func isMouseUp(ev *termbox.Event) bool {
	return ev.MouseBtnState == termbox.MouseBtnUp
}
