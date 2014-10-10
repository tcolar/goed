package main

import (
	"bytes"
	"io/ioutil"

	"github.com/tcolar/termbox-go"
)

type Editor struct {
	Menubar   *Menubar
	Statusbar *Statusbar
	Views     []View
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
	hs := w * 2 / 3
	vs := (h - 2) * 2 / 3
	view1 := View{
		Id:     1,
		Title:  "editor.go",
		Dirty:  true,
		Buffer: readFile("editor.go"),
	}
	view1.SetBounds(0, 1, hs, h-2)
	view2 := View{
		Id:     2,
		Title:  "themes/default.toml",
		Buffer: readFile("themes/default.toml"),
	}
	view2.SetBounds(hs+1, 1, w, vs)
	view3 := View{
		Id:    3,
		Title: "@scratch",
	}
	view3.SetBounds(hs+1, vs+1, w, h-2)

	e.Views = []View{view1, view2, view3}
	e.CurView = &view1
	//e.SetCursor(0, 0)

	e.Render()

	e.EventLoop()
}

// Temporary testing
func readFile(path string) [][]rune {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lines := bytes.Split(data, []byte("\n"))
	runes := [][]rune{}
	for _, l := range lines {
		runes = append(runes, bytes.Runes(l))
	}
	return runes
}
