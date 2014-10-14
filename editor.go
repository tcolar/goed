package main

import "github.com/tcolar/termbox-go"

type Editor struct {
	Menubar   *Menubar
	Statusbar *Statusbar
	Views     []*View
	Fg, Bg    Style
	Theme     *Theme
	CurView   *View
}

func (e *Editor) Start() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetExtendedColors(*colors == 256)
	e.Theme = ReadTheme("themes/default.toml")
	e.Fg = e.Theme.Fg
	e.Bg = e.Theme.Bg

	w, h := e.Size()
	e.Menubar = &Menubar{}
	e.Menubar.SetBounds(0, 0, w, 0)
	e.Statusbar = &Statusbar{}
	e.Statusbar.SetBounds(0, h-1, w, h-1)
	hs := w*2/3 - 1
	vs := (h - 2) * 2 / 3
	view1 := View{
		Id:     1,
		Title:  "view.go",
		Dirty:  true,
		Buffer: NewFileBuffer("view.go"),
	}
	view1.SetBounds(0, 1, hs, h-2)
	view2 := View{
		Id:     2,
		Title:  "themes/default.toml",
		Buffer: NewFileBuffer("themes/default.toml"),
	}
	view2.SetBounds(hs+1, 1, w, vs)
	view3 := View{
		Id:     3,
		Title:  "@scratch",
		Buffer: NewFileBuffer("cam.props"),
	}
	view3.SetBounds(hs+1, vs+1, w, h-2)

	e.Views = []*View{&view1, &view2, &view3}
	e.CurView = &view1
	e.CurView.MoveCursor(0, 0)

	e.Render()
	e.SetStatus("Holla!")

	e.EventLoop()
}

func (e *Editor) SetStatus(s string) {
	e.Statusbar.msg = s
	e.Statusbar.Render()
}
