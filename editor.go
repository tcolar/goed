package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tcolar/termbox-go"
)

var Ed *Editor

type Editor struct {
	Cmdbar     *Cmdbar
	Statusbar  *Statusbar
	Fg, Bg     Style
	Theme      *Theme
	Cols       []*Col
	CurView    *View
	CurCol     *Col
	CmdOn      bool
	pctw, pcth float64
	evtState   *EvtState
	Home       string
}

func (e *Editor) Start(loc string) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	e.initHome()
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
	view1 := &View{}
	e.CurView = view1
	e.Open(loc, view1, "")
	if view1.Buffer == nil { // new file that does not exist yet
		view1.Buffer = &Buffer{
			file: loc,
		}
	}
	view1.HeightRatio = 1.0

	c := e.NewCol(1.0, []*View{view1})
	e.Cols = []*Col{c}

	e.CurCol = c

	// Add a directory listing view if we don't already have one
	if stat, err := os.Stat(loc); err == nil && !stat.IsDir() {
		view2 := &View{}
		e.Open(".", view2, "")
		view2.HeightRatio = 1.0
		c.WidthRatio = 0.75
		c2 := e.NewCol(0.25, []*View{view2})
		e.Cols = append([]*Col{c2}, c)
	}

	e.Resize(e.Size())

	e.CurView.MoveCursor(0, 0)

	e.Render()

	e.EventLoop()
}

func (e *Editor) Open(loc string, view *View, rel string) error {
	if view == nil {
		return fmt.Errorf("No view selected !")
	}
	if len(rel) > 0 && !strings.HasPrefix(loc, string(os.PathSeparator)) {
		loc = path.Join(rel, loc)
	}
	stat, err := os.Stat(loc)
	if err != nil {
		return fmt.Errorf("File not found %s", loc)
	}
	if stat.Size() > 500000 { // TODO : check if utf8 / executable
		return fmt.Errorf("File too large %s", loc)
	}
	// make it absolute
	loc, err = filepath.Abs(loc)
	if err != nil {
		return err
	}
	title := filepath.Base(loc)
	if stat.IsDir() {
		loc += string(os.PathSeparator)
		title += string(os.PathSeparator)
	}
	view.Reset()
	view.title = title
	if stat.IsDir() {
		err = e.openDir(loc, view)
	} else {
		err = e.openFile(loc, view)
	}
	view.WorkDir = filepath.Dir(loc)
	e.SetStatus(view.WorkDir)
	return nil
}

func (e *Editor) openDir(loc string, view *View) error {
	buffer, err := e.NewDirBuffer(loc)
	if err != nil {
		return err
	}
	view.Buffer = buffer
	view.Dirty = false
	return nil
}

func (e *Editor) openFile(loc string, view *View) error {
	view.Buffer = e.NewFileBuffer(loc)
	view.Dirty = false
	return nil
}

func (e *Editor) SetStatusErr(s string) {
	if e.Statusbar == nil {
		return
	}
	e.Statusbar.msg = s
	e.Statusbar.isErr = true
	e.Statusbar.Render()
}
func (e *Editor) SetStatus(s string) {
	if e.Statusbar == nil {
		return
	}
	e.Statusbar.msg = s
	e.Statusbar.msg = s
	e.Statusbar.isErr = false
	e.Statusbar.Render()
}

// RunesToString returns a rune section as a srting.
func (e Editor) RunesToString(runes [][]rune) string {
	r := []rune{}
	for i, line := range runes {
		if i != 0 && i != len(runes) {
			r = append(r, '\n')
		}
		r = append(r, line...)
	}
	return string(r)
}

func (e Editor) StringToRunes(s []byte) [][]rune {
	lines := bytes.Split(s, []byte("\n"))
	runes := [][]rune{}
	for i, l := range lines {
		// Ignore last line if empty
		if i != len(lines)-1 || len(l) != 0 {
			runes = append(runes, bytes.Runes(l))
		}
	}
	return runes
}

// QuitCeck check if it's ok to quit
// if there are no dirty buffer
// or if requested twice in a row
func (e *Editor) QuitCheck() bool {
	ok := true
	for _, c := range e.Cols {
		for _, v := range c.Views {
			ok = ok && v.canClose()
		}
	}
	return ok
}
