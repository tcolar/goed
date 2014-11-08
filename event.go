// History: Oct 07 14 tcolar Creation

package main

import "github.com/tcolar/termbox-go"

func (e *Editor) EventLoop() {

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventResize:
			e.Render()
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ:
				return
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
		case termbox.KeySpace:
			m.Cmd = append(m.Cmd, ' ') // hum why is ev.Ch not space when pressing space ???
		default:
			m.Cmd = append(m.Cmd, ev.Ch)
		}

	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			Ed.CmdOn = true
		}
	}
}

// Event handler for Statusbar
func (s *Statusbar) Event(ev *termbox.Event) {
	// TBD
}

// Event handler for View
func (v *View) Event(ev *termbox.Event) {
	dirty := false
	switch ev.Type {
	case termbox.EventKey:
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
		case termbox.KeyCtrlQ:
			return
		default:
			// insert the key
			v.Insert(ev.Ch)
			dirty = true
		}
	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			Ed.CmdOn = false
			// MoveCursor use text coordinates which starts at offset 2,2
			v.MoveCursor(ev.MouseX-v.x1-2-v.CursorX, ev.MouseY-v.y1-2-v.CursorY)
			// Make the clicked view active
			Ed.CurView = v
		}
	}
	if dirty {
		v.Dirty = true
	}
}
