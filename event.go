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
		/*switch ev.Key {
		case termbox.KeyArrowRight:
			e.MoveCursor(1, 0)
		case termbox.KeyArrowLeft:
			e.MoveCursor(-1, 0)
		case termbox.KeyArrowUp:
			e.MoveCursor(0, -1)
		case termbox.KeyArrowDown:
			e.MoveCursor(0, 1)
		case termbox.KeyEsc:
			return
		}*/
	case termbox.EventMouse:
		switch ev.Key {
		case termbox.MouseLeft:
			termbox.SetCursor(ev.MouseX, ev.MouseY)
			v.CursorX = ev.MouseX - v.x1
			v.CursorY = ev.MouseY - v.y1
			Ed.CurView = v
		}
	}
}
