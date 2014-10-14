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
			case termbox.KeyEsc:
				return
			default:
				if e.CurView != nil {
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

// Event handler for Menubar
func (m *Menubar) Event(ev *termbox.Event) {
	// TBD
}

// Event handler for Statusbar
func (s *Statusbar) Event(ev *termbox.Event) {
	// TBD
}

// Event handler for View
func (v *View) Event(ev *termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowRight:
			v.MoveCursor(1, 0)
		case termbox.KeyArrowLeft:
			v.MoveCursor(-1, 0)
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
		case termbox.KeyEsc:
			return
		case termbox.KeyEnter:
			// special TBD
		case termbox.KeyDelete:
			// special TBD
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			// special TBD
		default:
			// insert the key
			v.Insert(ev.Ch)
		}
	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			// MoveCursor use text coordinates which starts at offset 2,2
			v.MoveCursor(ev.MouseX-v.x1-2-v.CursorX, ev.MouseY-v.y1-2-v.CursorY)
			// Make the clicked view active
			Ed.CurView = v
		}
	}
}
