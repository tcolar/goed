// History: Oct 07 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

// Evtstate stores some state about kb/mouse events
type EvtState struct {
	MovingView                     bool
	X, Y                           int
	DragX1, DragY1, DragX2, DragY2 int
}

func (e *Editor) EventLoop() {

	termbox.SetMouseMode(termbox.MouseMotion)
	termbox.SetInputMode(termbox.InputMouse)

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			e.Resize(ev.Width, ev.Height)
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				if !e.QuitCheck() {
					Ed.SetStatusErr("Unsaved changes. Save or request close again.")
				} else {
					return // the end
				}
			case termbox.KeyEsc:
				if !Ed.CmdOn {
					Ed.Cmdbar.Cmd = []rune{}
				}
				Ed.CmdOn = !Ed.CmdOn
			default:
				if e.CmdOn {
					e.Cmdbar.Event(&ev)
				} else if e.CurView != nil {
					e.CurView.Event(&ev)
				}
			}
		case termbox.EventMouse:
			w := e.WidgetAt(ev.MouseX, ev.MouseY)
			if w != nil {
				w.Event(&ev)
			}
		}
		e.Render()
	}
}

// ##################### CmdBar ########################################

// Event handler for Cmdbar
func (m *Cmdbar) Event(ev *termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		//case termbox.KeyDelete:
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
			Ed.CmdOn = true
		}
	}
}

// ##################### StatusBar ########################################

// Event handler for Statusbar
func (s *Statusbar) Event(ev *termbox.Event) {
	// Anyhting ??
}

// ##################### View       ########################################

// Event handler for View
func (v *View) Event(ev *termbox.Event) {
	dirty := false
	switch ev.Type {
	case termbox.EventKey:
		// alt combo
		if ev.Mod == termbox.ModAlt {
			switch ev.Ch {
			case 'o':
				Ed.Cmdbar.OpenSelection(v, false)
			}
			return
		}
		// Not alt
		switch ev.Key {
		case termbox.KeyArrowRight:
			offset := 1
			c, _, _ := v.CurChar()
			if c != nil {
				offset = v.runeSize(*c)
			}
			v.MoveCursor(offset, 0)
		case termbox.KeyArrowLeft:
			offset := 1
			c, _, _ := v.CursorChar(v.CurCol()-1, v.CurLine())
			if c != nil {
				offset = v.runeSize(*c)
			}
			v.MoveCursor(-offset, 0)
		case termbox.KeyArrowUp:
			v.MoveCursor(0, -1)
		case termbox.KeyArrowDown:
			v.MoveCursor(0, 1)
		case termbox.KeyPgdn:
			dist := v.LastViewLine() + 1
			if v.LineCount()-v.CurLine() < dist {
				dist = v.LineCount() - v.CurLine() - 1
			}
			v.MoveCursor(0, dist)
		case termbox.KeyPgup:
			dist := v.LastViewLine() + 1
			if dist > v.CurLine() {
				dist = v.CurLine()
			}
			v.MoveCursor(0, -dist)
		case termbox.KeyEnd:
			v.MoveCursor(v.lineCols(v.CurLine())-v.CurCol(), 0)
		case termbox.KeyHome:
			v.MoveCursor(-v.CurCol(), 0)
		case termbox.KeyTab:
			v.Insert('\t')
			dirty = true
		case termbox.KeyEnter:
			v.InsertNewLine()
			dirty = true
		case termbox.KeyDelete:
			v.Delete()
			dirty = true
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			v.Backspace()
			dirty = true
		case termbox.KeyCtrlS:
			v.Save()
		case termbox.KeyCtrlO:
			Ed.Cmdbar.OpenSelection(v, true)
		case termbox.KeyCtrlW:
			Ed.DelViewCheck(Ed.CurView)
		case termbox.KeyCtrlC:
			if len(v.Selections) > 0 {
				v.Copy(v.Selections[0])
			}
		case termbox.KeyCtrlV:
			v.Paste()
			dirty = true
		case termbox.KeyCtrlQ:
			return
		default:
			// insert the key
			if ev.Ch != 0 && ev.Mod == 0 { // otherwise special key combo
				v.Insert(ev.Ch)
				dirty = true
			}
		}
	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseScrollUp:
			v.MoveCursor(0, -1)
		case termbox.MouseScrollDown:
			v.MoveCursor(0, 1)
		case termbox.MouseLeft:
			if Ed.evtState.MovingView {
				Ed.evtState.MovingView = false
				Ed.ViewMove(Ed.evtState.X, Ed.evtState.Y, ev.MouseX, ev.MouseY)
				return
			}
			if ev.MouseX == v.x2-1 && ev.MouseY == v.y1 {
				Ed.DelViewCheck(v)
				return
			}
			if ev.MouseX == v.x1 && ev.MouseY == v.y1 {
				// handle
				Ed.evtState.MovingView = true
				Ed.evtState.X = ev.MouseX
				Ed.evtState.Y = ev.MouseY
				Ed.SetStatusErr("Starting move, click new position.")
				return
			}
			if ev.DragOn {
				if Ed.evtState.DragX1 == 0 {
					// start drag
					Ed.evtState.DragX1, Ed.evtState.DragY1 = ev.MouseX, ev.MouseY
				}
				// continued drag
				Ed.evtState.DragX2, Ed.evtState.DragY2 = ev.MouseX, ev.MouseY
				x1 := Ed.evtState.DragX1 - v.x1 + v.offx - 2
				x2 := Ed.evtState.DragX2 - v.x1 + v.offx - 2
				y1 := Ed.evtState.DragY1 - v.y1 + v.offy - 2
				y2 := Ed.evtState.DragY2 - v.y1 + v.offy - 2

				s := Selection{
					LineFrom: y1,
					LineTo:   y2,
					ColFrom:  v.lineRunesTo(y1, x1),
					ColTo:    v.lineRunesTo(y2, x2),
				}
				// Deal with "reverse" selection
				reverse := false
				if s.LineFrom == s.LineTo && s.ColFrom > s.ColTo {
					reverse = true
					s.ColFrom, s.ColTo = s.ColTo, s.ColFrom
				} else if s.LineFrom > s.LineTo {
					reverse = true
					s.LineFrom, s.LineTo = s.LineTo, s.LineFrom
					s.ColFrom, s.ColTo = s.ColTo, s.ColFrom
				}
				// Because we only receive the event after a "move", we need to add the start location
				if reverse {
					s.ColTo++
				} else {
					s.ColFrom--
				}
				// set the selection
				v.Selections = []Selection{
					s,
				}
				return
			} else {
				// reset drag
				Ed.evtState.DragX1, Ed.evtState.DragY1 = 0, 0
				Ed.evtState.DragX2, Ed.evtState.DragY2 = 0, 0
				v.Selections = []Selection{}
			}
			Ed.CmdOn = false
			// MoveCursor use text coordinates which starts at offset 2,2
			v.MoveCursor(ev.MouseX-v.x1-2-v.CursorX, ev.MouseY-v.y1-2-v.CursorY)
			// Make the clicked view active
			Ed.ActivateView(v, 0, 0)
			Ed.SetStatus(v.WorkDir)
		}
	}
	if dirty {
		v.Dirty = true
	}
}
