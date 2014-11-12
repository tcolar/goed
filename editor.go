package main

import (
	"fmt"
	"os"

	"github.com/tcolar/termbox-go"
)

type Editor struct {
	Cmdbar    *Cmdbar
	Statusbar *Statusbar
	//Views      []*View
	Fg, Bg     Style
	Theme      *Theme
	Cols       []*Col
	CurView    *View
	CurCol     *Col
	CmdOn      bool
	pctw, pcth float64
	evtState   *EvtState
}

func (e *Editor) Start() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetExtendedColors(*colors == 256)
	e.evtState = &EvtState{}
	e.Theme = ReadTheme("themes/default.toml")
	e.Fg = e.Theme.Fg
	e.Bg = e.Theme.Bg

	w, h := e.Size()
	e.Cmdbar = &Cmdbar{}
	e.Cmdbar.SetBounds(0, 0, w, 0)
	e.Statusbar = &Statusbar{}
	e.Statusbar.SetBounds(0, h-1, w, h-1)
	view1 := e.NewFileView("view.go")
	view1.HeightRatio = 1.0
	view2 := e.NewFileView("themes/default.toml")
	view2.HeightRatio = 0.8
	view3 := e.NewView()
	view3.HeightRatio = 0.2
	c := e.NewCol(0.75, []*View{view1})
	c2 := e.NewCol(0.25, []*View{view2, view3})
	e.Cols = []*Col{
		c,
		c2,
	}

	e.CurView = view1
	e.CurCol = c

	e.Resize(e.Size())

	e.CurView.MoveCursor(0, 0)

	e.Render()

	e.SetStatus("Holla!")

	e.EventLoop()
}

func (e *Editor) OpenFile(loc string, view *View) error {
	if view == nil {
		return fmt.Errorf("No view selected !")
	}
	if _, err := os.Stat(loc); err != nil {
		return fmt.Errorf("File not found %s", loc)
	}
	view.Buffer = e.NewFileBuffer(loc)
	view.Dirty = false
	return nil
}

func (e *Editor) SetStatusErr(s string) {
	e.Statusbar.msg = s
	e.Statusbar.isErr = true
	e.Statusbar.Render()
}
func (e *Editor) SetStatus(s string) {
	e.Statusbar.msg = s
	e.Statusbar.isErr = false
	e.Statusbar.Render()
}
