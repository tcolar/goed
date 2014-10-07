package main

import "github.com/tcolar/termbox-go"

type Editor struct {
	Cols   []Col
	Fg, Bg Style
	Theme  *Theme
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
	//termbox.SetCursor(ed.cursor_position())
	w, h := e.Size()
	vs := h * 2 / 3
	hs := w * 2 / 3
	col1 := Col{
		X1: 0,
		X2: hs - 1,
	}
	view1 := View{
		Y1:    0,
		Y2:    hs + 1,
		Title: "/home/tcolar/DEV/mesa/Makefile",
		Dirty: true,
	}
	col1.Views = append(col1.Views, view1)
	col2 := Col{
		X1: hs,
		X2: w,
	}
	view2 := View{
		Y1:    0,
		Y2:    vs - 1,
		Title: "Test.txt",
	}
	view3 := View{
		Y1:    vs,
		Y2:    h,
		Title: "dummy_test.go",
	}
	col2.Views = append(col2.Views, view2)
	col2.Views = append(col2.Views, view3)
	e.Cols = []Col{
		col1,
		col2,
	}

	e.draw()
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			}
		case termbox.EventResize:
			e.draw()
		}
	}
}

func (e *Editor) draw() {

	e.Render()
	e.FB(e.Theme.Fg, e.Theme.Bg)
	e.Str(2, 3, "Hello World")

	e.FB(e.Theme.String, e.Theme.Bg)
	e.Str(1, 9, "Kw: What's up -\u0E5Bಠﭛಠ! \u2611 \u2612 ➩")

	termbox.Flush()
}
