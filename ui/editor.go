package ui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
)

// Editor is goed's main Editor pane (singleton)
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

// Editor with Mock terminal for testing
func NewMockEditor() *Editor {
	return &Editor{
		term:   core.NewMockTerm(),
		config: core.LoadDefaultConfig(),
	}
}

// Start starts-up the editor
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

	h, w := e.term.Size()
	e.Cmdbar = &Cmdbar{}
	e.Cmdbar.SetBounds(0, 0, 0, w)
	e.Statusbar = &Statusbar{}
	e.Statusbar.SetBounds(h-1, 0, h-1, w)
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

	go e.autoScroller()

	if !core.Testing {
		e.EventLoop()
	}
}

// Open opens a given location in the editor (in the given view)
func (e Editor) Open(loc string, view core.Viewable, rel string) (core.Viewable, error) {
	if len(rel) > 0 && !strings.HasPrefix(loc, string(os.PathSeparator)) {
		loc = path.Join(rel, loc)
	}
	stat, err := os.Stat(loc)
	if err != nil {
		return nil, fmt.Errorf("File not found %s", loc)
	}
	// make it absolute
	loc, err = filepath.Abs(loc)
	if err != nil {
		return nil, err
	}
	title := filepath.Base(loc)
	if stat.IsDir() {
		loc += string(os.PathSeparator)
		title += string(os.PathSeparator)
	}
	nv := false
	if view == nil {
		view = e.NewView()
		nv = true
	}
	view.Reset()
	view.SetTitle(title)
	if stat.IsDir() {
		err = e.openDir(loc, view)
	} else {
		err = e.openFile(loc, view)
	}
	view.SetWorkDir(filepath.Dir(loc))
	if nv {
		e.InsertView(view.(*View), e.CurView().(*View), 0.2)
		//		e.ActivateView(view.(*View), 0, 0)
	}
	return view, err
}

// OpenDir opens a directory listing
func (e *Editor) openDir(loc string, view core.Viewable) error {
	args := []string{"ls", "-a", "--color=no"}
	title := filepath.Base(loc) + "/"
	backend, err := backend.NewMemBackendCmd(args, loc, view.Id(), &title)
	if err != nil {
		return err
	}
	view.SetBackend(backend)
	e.SetStatus(fmt.Sprintf("%v", view.WorkDir()))
	view.SetDirty(false)
	return nil
}

// OpenFile opens a file in the editor
func (e *Editor) openFile(loc string, view core.Viewable) error {
	if !core.IsTextFile(loc) {
		return fmt.Errorf("Binary file ? %s", loc)
	}
	backend, err := backend.NewFileBackend(loc, view.Id())
	if err != nil {
		return err
	}
	view.SetBackend(backend)
	e.SetStatus(fmt.Sprintf("%v  [%d]", view.WorkDir(), view.Id()))
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

func (e Editor) SetCursor(y, x int) {
	e.term.SetCursor(y, x)
}

func (e Editor) CmdOn() bool {
	return e.cmdOn
}

func (e *Editor) SetCmdOn(v bool) {
	e.cmdOn = v
}

func (e *Editor) TermFlush() {
	e.term.Flush()
}

// Handle selection auto scrolling of views
func (e *Editor) autoScroller() {
	for {
		time.Sleep(200 * time.Millisecond)
		v := e.CurView().(*View)
		if v == nil {
			continue
		}
		x, y := v.autoScrollX, v.autoScrollY
		if x == 0 && y == 0 {
			continue
		}
		if len(v.selections) == 0 {
			continue
		}
		s := v.selections[0]
		ln := v.CurLine()
		v.offx += x
		v.offy += y
		if y > 0 {
			s.LineTo += y
		} else {
			s.LineFrom += y
		}
		if x > 0 {
			s.ColTo += x
		} else {
			s.ColFrom += x
		}
		// handle scroll / selection "overflows"
		lnLen := v.LineLen(v.Slice(), ln)
		if v.offy >= v.LineCount()-v.LastViewLine() {
			v.offy = v.LineCount() - v.LastViewLine()
		}
		if v.offy < 0 {
			v.offy = 0
		}
		if v.offx > lnLen-v.LastViewCol() {
			v.offx = lnLen - v.LastViewCol()
		}
		if v.offx < 0 {
			v.offx = 0
		}
		if s.LineFrom < 0 {
			s.LineFrom = 0
		} else if s.LineFrom > v.LineCount() {
			s.LineFrom = v.LineCount()
		}
		if s.LineTo < 0 {
			s.LineTo = 0
		} else if s.LineTo > v.LineCount() {
			s.LineTo = v.LineCount()
		}
		if s.ColFrom < 0 {
			s.ColFrom = 0
		} else if s.ColFrom > lnLen {
			s.ColFrom = lnLen
		}
		if s.ColTo < 0 {
			s.ColTo = 0
		} else if s.ColTo > lnLen {
			s.ColTo = lnLen
		}
		s.Normalize()
		v.selections = []core.Selection{
			s,
		}
		e.Render()
	}
}
