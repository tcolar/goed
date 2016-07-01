package ui

import (
	"log"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
	termbox "github.com/tcolar/termbox-go"
)

// Termbox : Term implementation using termbox
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
	// TODO: should set it to original value, but how to read it ??
	t.SetMouseMode(termbox.MouseClick)
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
		event.Queue(*es)
	}
}

// parses a termbox event into the 'es' goed event (event.Event)
func (t *TermBox) parseEvent(e termbox.Event, es *event.Event) {
	if e.Ch > 0 {
		es.Glyph = string(e.Ch)
	} else {
		es.Glyph = ""
	}
	m := e.Meta
	es.Combo = event.Combo{
		LAlt:   m == termbox.Alt || m == termbox.AltCtrl || m == termbox.AltCtrlShift || m == termbox.AltShift,
		LCtrl:  m == termbox.Ctrl || m == termbox.CtrlShift || m == termbox.AltCtrl || m == termbox.AltCtrlShift,
		LShift: m == termbox.Shift || m == termbox.AltShift || m == termbox.CtrlShift || m == termbox.AltCtrlShift,
		LSuper: m == termbox.Meta,
	}
	es.Keys = []string{}
	es.MouseBtns[8] = false  // reset wheel down, no separate "up" event
	es.MouseBtns[16] = false // reset wheel up

	es.Type = event.Evt_None
	switch e.Type {
	case termbox.EventResize:
		actions.Ar.EdResize(e.Height, e.Width)
		es.Type = event.EvtWinResize
		return
	}

	if len(es.Glyph) > 0 {
		es.KeyDown(es.Glyph)
		return
	}

	ctrl := func(k string) {
		es.KeyDown(event.KeyLeftControl)
		es.KeyDown(k)
	}

	k := e.Key

	switch k {
	case termbox.MouseLeft:
		if e.MouseBtnState != termbox.MouseBtnUp {
			es.MouseDown(event.MouseLeft, e.MouseY, e.MouseX)
		} else {
			es.MouseUp(event.MouseLeft, e.MouseY, e.MouseX)
		}
	case termbox.MouseRight:
		if e.MouseBtnState != termbox.MouseBtnUp {
			es.MouseDown(event.MouseRight, e.MouseY, e.MouseX)
		} else {
			es.MouseUp(event.MouseRight, e.MouseY, e.MouseX)
		}
	case termbox.MouseMiddle:
		if e.MouseBtnState != termbox.MouseBtnUp {
			es.MouseDown(event.MouseMiddle, e.MouseY, e.MouseX)
		} else {
			es.MouseUp(event.MouseMiddle, e.MouseY, e.MouseX)
		}
	case termbox.MouseScrollDown:
		es.MouseDown(event.MouseWheelDown, e.MouseY, e.MouseX)
	case termbox.MouseScrollUp:
		es.MouseDown(event.MouseWheelUp, e.MouseY, e.MouseX)
	case termbox.KeyBackspace, termbox.KeyBackspace2:
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
	log.Printf("es : %#v", es)
}
