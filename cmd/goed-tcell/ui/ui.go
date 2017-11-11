package ui

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/tcolar/goed"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
)

var _ core.Term = (*Tcell)(nil)

func Main() {
	config := goed.Initialize()
	if config == nil {
		return
	}
	term := NewTcell()
	defer goed.Terminate(term)
	goed.Start(term, config)
}

// Tcell : Term implementation using tcell
type Tcell struct {
	screen tcell.Screen
}

func NewTcell() *Tcell {
	return &Tcell{}
}

func (t *Tcell) Init() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	encoding.Register()

	if err = screen.Init(); err != nil {
		return err
	}
	t.screen = screen
	return nil
}

func (t *Tcell) Clear(fg, bg core.Style) {
	st := t.toTcellStyle(fg, bg)
	w, h := t.screen.Size()
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			t.screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

func (t *Tcell) Close() {
	t.screen.DisableMouse()
	t.screen.Fini()
}

func (t *Tcell) Flush() {
	t.screen.Show()
}

func (t *Tcell) SetExtendedColors(b bool) {
}

func (t *Tcell) SetCursor(y, x int) {
	t.screen.ShowCursor(y, x)
}

func (t *Tcell) Char(y, x int, c rune, fg, bg core.Style) {
	st := t.toTcellStyle(fg, bg)
	t.screen.SetContent(x, y, c, nil, st)
}

func (t *Tcell) toTcellStyle(fg, bg core.Style) tcell.Style {
	f := tcell.Color(fg.Color())
	b := tcell.Color(bg.Color())
	st := tcell.StyleDefault.Foreground(f).Background(b)
	if fg.IsBold() {
		st = st.Bold(true)
	}
	if fg.IsUnderlined() {
		st = st.Underline(true)
	}
	return st
}

func (t *Tcell) Size() (h, w int) {
	w, h = t.screen.Size()
	return h, w
}

func (t *Tcell) Listen() {
	t.screen.EnableMouse()
	es := event.NewEvent()
	for {
		e := t.screen.PollEvent()
		t.parseEvent(e, es)
		if es.Type == event.EvtQuit {
			return
		}
		event.Queue(es.Clone())
	}
}

func (t *Tcell) updateMouseState(es *event.Event, down bool, button event.MouseButton, y, x int) {
	if down {
		es.MouseDown(button, y, x)
	} else if es.MouseBtns[button] {
		es.MouseUp(button, y, x)
	}
}

func (t *Tcell) parseEvent(e tcell.Event, es *event.Event) {
	es.Glyph = ""
	es.Keys = []string{}
	es.MouseBtns[event.MouseWheelDown] = false // reset wheel down, no separate "up" event
	es.MouseBtns[event.MouseWheelUp] = false   // reset wheel up
	es.Type = event.Evt_None

	ctrl := func(k string) {
		es.KeyDown(event.KeyLeftControl)
		es.KeyDown(k)
	}

	switch te := e.(type) {
	case *tcell.EventResize:
		w, h := te.Size()
		actions.Ar.EdResize(h, w)
		return
	case *tcell.EventMouse:
		m := te.Modifiers()
		es.Combo = event.Combo{
			LAlt:   m&tcell.ModAlt != 0,
			LCtrl:  m&tcell.ModCtrl != 0,
			LShift: m&tcell.ModShift != 0,
			LSuper: m&tcell.ModMeta != 0,
		}
		b := te.Buttons()
		x, y := te.Position()
		t.updateMouseState(es, b&tcell.Button1 != 0, event.MouseLeft, y, x)
		t.updateMouseState(es, b&tcell.Button2 != 0, event.MouseMiddle, y, x)
		t.updateMouseState(es, b&tcell.Button3 != 0, event.MouseRight, y, x)
		t.updateMouseState(es, b&tcell.WheelUp != 0, event.MouseWheelUp, y, x)
		t.updateMouseState(es, b&tcell.WheelDown != 0, event.MouseWheelDown, y, x)
		t.updateMouseState(es, b&tcell.WheelLeft != 0, event.MouseWheelLeft, y, x)
		t.updateMouseState(es, b&tcell.WheelRight != 0, event.MouseWheelRight, y, x)

	case *tcell.EventKey:
		m := te.Modifiers()
		es.Combo = event.Combo{
			LAlt:   m&tcell.ModAlt != 0,
			LCtrl:  m&tcell.ModCtrl != 0,
			LShift: m&tcell.ModShift != 0,
			LSuper: m&tcell.ModMeta != 0,
		}
		if te.Key() == tcell.KeyRune {
			es.Glyph = string(te.Rune())
			es.KeyDown(es.Glyph)
			return
		}
		switch te.Key() {
		case tcell.KeyEsc:
			es.KeyDown(event.KeyEscape)
		case tcell.KeyBackspace2:
			es.KeyDown(event.KeyBackspace)
		case tcell.KeyTab:
			es.KeyDown(event.KeyTab)
		case tcell.KeyEnter:
			es.KeyDown(event.KeyReturn)
		case tcell.KeyF1:
			es.KeyDown(event.KeyF1)
		case tcell.KeyF2:
			es.KeyDown(event.KeyF2)
		case tcell.KeyF3:
			es.KeyDown(event.KeyF3)
		case tcell.KeyF4:
			es.KeyDown(event.KeyF4)
		case tcell.KeyF5:
			es.KeyDown(event.KeyF5)
		case tcell.KeyF6:
			es.KeyDown(event.KeyF7)
		case tcell.KeyF7:
			es.KeyDown(event.KeyF7)
		case tcell.KeyF8:
			es.KeyDown(event.KeyF8)
		case tcell.KeyF9:
			es.KeyDown(event.KeyF9)
		case tcell.KeyF10:
			es.KeyDown(event.KeyF10)
		case tcell.KeyF11:
			es.KeyDown(event.KeyF11)
		case tcell.KeyF12:
			es.KeyDown(event.KeyF12)
		case tcell.KeyInsert:
			es.KeyDown(event.KeyInsert)
		case tcell.KeyDelete:
			es.KeyDown(event.KeyDelete)
		case tcell.KeyHome:
			es.KeyDown(event.KeyHome)
		case tcell.KeyEnd:
			es.KeyDown(event.KeyEnd)
		case tcell.KeyPgUp:
			es.KeyDown(event.KeyPrior)
		case tcell.KeyPgDn:
			es.KeyDown(event.KeyNext)
		case tcell.KeyUp:
			es.KeyDown(event.KeyUpArrow)
		case tcell.KeyDown:
			es.KeyDown(event.KeyDownArrow)
		case tcell.KeyLeft:
			es.KeyDown(event.KeyLeftArrow)
		case tcell.KeyRight:
			es.KeyDown(event.KeyRightArrow)

		case tcell.KeyCtrlA:
			ctrl("a")
		case tcell.KeyCtrlB:
			ctrl("b")
		case tcell.KeyCtrlC:
			ctrl("c")
		case tcell.KeyCtrlD:
			ctrl("d")
		case tcell.KeyCtrlE:
			ctrl("e")
		case tcell.KeyCtrlF:
			ctrl("f")
		case tcell.KeyCtrlG:
			ctrl("g")
		case tcell.KeyCtrlH:
			ctrl("h")
		case tcell.KeyCtrlJ:
			ctrl("j")
		case tcell.KeyCtrlK:
			ctrl("k")
		case tcell.KeyCtrlL:
			ctrl("l")
		case tcell.KeyCtrlN:
			ctrl("n")
		case tcell.KeyCtrlO:
			ctrl("o")
		case tcell.KeyCtrlP:
			ctrl("p")
		case tcell.KeyCtrlQ:
			ctrl("q")
		case tcell.KeyCtrlR:
			ctrl("r")
		case tcell.KeyCtrlS:
			ctrl("s")
		case tcell.KeyCtrlT:
			ctrl("t")
		case tcell.KeyCtrlU:
			ctrl("u")
		case tcell.KeyCtrlV:
			ctrl("v")
		case tcell.KeyCtrlW:
			ctrl("w")
		case tcell.KeyCtrlX:
			ctrl("x")
		case tcell.KeyCtrlY:
			ctrl("y")
		case tcell.KeyCtrlZ:
			ctrl("z")
		}
	}
}
