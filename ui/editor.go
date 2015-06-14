package ui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

type Editor struct {
	Cmdbar     *Cmdbar
	config     *core.Config
	Statusbar  *Statusbar
	Fg, Bg     core.Style
	theme      *core.Theme
	Cols       []*Col
	curView    *View
	CurCol     *Col
	cmdOn      bool
	pctw, pcth float64
	evtState   *EvtState
	term       core.Term
}

func NewEditor() *Editor {
	return &Editor{
		term:   core.NewTermBox(),
		config: core.LoadConfig(core.ConfFile),
	}
}

// Edior with Mock terminal for testing
func NewMockEditor() *Editor {
	return &Editor{
		term:   core.NewMockTerm(),
		config: core.LoadDefaultConfig(),
	}
}

func (e *Editor) Start(loc string) {
	err := e.term.Init()
	if err != nil {
		panic(err)
	}
	defer e.term.Close()
	e.term.SetExtendedColors(core.Colors == 256)
	e.evtState = &EvtState{}
	e.theme, err = core.ReadTheme(path.Join(core.Home, "themes", e.config.Theme))
	if err != nil {
		panic(err)
	}
	e.Fg = e.theme.Fg
	e.Bg = e.theme.Bg

	w, h := e.term.Size()
	e.Cmdbar = &Cmdbar{}
	e.Cmdbar.SetBounds(0, 0, w, 0)
	e.Statusbar = &Statusbar{}
	e.Statusbar.SetBounds(0, h-1, w, h-1)
	view1 := e.NewView()
	e.curView = view1
	e.Open(loc, view1, "")
	if view1.backend == nil { // new file that does not exist yet
		view1.backend, _ = backend.NewFileBackend("", view1.Id())
	}
	view1.HeightRatio = 1.0

	c := e.NewCol(1.0, []*View{view1})
	e.Cols = []*Col{c}

	e.CurCol = c

	// Add a directory listing view if we don't already have one
	if stat, err := os.Stat(loc); err == nil && !stat.IsDir() {
		view2 := e.NewView()
		e.Open(".", view2, "")
		view2.HeightRatio = 1.0
		c.WidthRatio = 0.75
		c2 := e.NewCol(0.25, []*View{view2})
		e.Cols = append([]*Col{c2}, c)
	}

	e.Resize(e.term.Size())

	e.Render()

	if !core.Testing {
		e.EventLoop()
	}
}

func (e Editor) Open(loc string, view core.Viewable, rel string) error {
	if view == nil {
		return fmt.Errorf("Invalid view !")
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
	view.SetTitle(title)
	if stat.IsDir() {
		err = e.openDir(loc, view)
	} else {
		err = e.openFile(loc, view)
	}
	view.SetWorkDir(filepath.Dir(loc))
	return nil
}

func (e *Editor) openDir(loc string, view core.Viewable) error {
	args := []string{"ls", "-a", "--color=no"}
	title := filepath.Base(loc) + "/"
	backend, err := backend.NewMemBackendCmd(args, loc, view.Id(), &title)
	if err != nil {
		return err
	}
	view.SetBackend(backend)
	e.SetStatus(fmt.Sprintf("[%d]%v", view.Id(), view.WorkDir()))
	view.SetDirty(false)
	return nil
}

func (e *Editor) openFile(loc string, view core.Viewable) error {
	backend, err := backend.NewFileBackend(loc, view.Id())
	if err != nil {
		return err
	}
	view.SetBackend(backend)
	e.SetStatus(fmt.Sprintf("[%d]%v", view.Id(), view.WorkDir()))
	view.SetDirty(false)
	return nil
}

func (e Editor) SetStatusErr(s string) {
	if e.Statusbar == nil {
		return
	}
	e.Statusbar.msg = s
	e.Statusbar.isErr = true
	e.Statusbar.Render()
}
func (e Editor) SetStatus(s string) {
	if e.Statusbar == nil {
		return
	}
	e.Statusbar.msg = s
	e.Statusbar.msg = s
	e.Statusbar.isErr = false
	e.Statusbar.Render()
}

// QuitCheck check if it's ok to quit
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

func (e Editor) Config() core.Config {
	return *e.config
}

func (e Editor) Theme() *core.Theme {
	return e.theme
}

func (e Editor) CurView() core.Viewable {
	return e.curView
}

func (e Editor) SetCursor(x, y int) {
	e.term.SetCursor(x, y)
}

func (e Editor) CmdOn() bool {
	return e.cmdOn
}

func (e *Editor) SetCmdOn(v bool) {
	e.cmdOn = v
}
