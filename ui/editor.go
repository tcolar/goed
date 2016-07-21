// package ui provides the UI components of Goed.
package ui

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/backend"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/event"
)

var _ core.Editable = (*Editor)(nil)

// Editor is goed's main Editor pane (singleton)
type Editor struct {
	Cmdbar    *Cmdbar
	config    *core.Config
	Statusbar *Statusbar
	Fg, Bg    core.Style
	theme     *core.Theme
	Cols      []*Col
	curViewId int64
	CurCol    *Col
	cmdOn     bool
	term      core.Term
	views     map[int64]*View
}

func NewEditor(gui bool) *Editor {
	var term core.Term
	if gui {
		term = NewGuiTerm(1200, 800)
	} else {
		term = NewTermBox()
	}
	return &Editor{
		term:   term,
		config: core.LoadConfig(core.ConfFile),
		views:  map[int64]*View{},
	}
}

// Editor with Mock terminal for testing
func NewMockEditor() *Editor {
	return &Editor{
		term:   core.NewMockTerm(),
		config: core.LoadConfig("config.toml"),
		views:  map[int64]*View{},
	}
}

func (e *Editor) Dispatch(action core.Action) {
	core.Bus.Dispatch(action)
}

func (e *Editor) Commandbar() core.Commander {
	return e.Cmdbar
}

func (e *Editor) Quit() {
	event.Shutdown()
	e.term.Close()
	os.Exit(0)
}

// Start starts-up the editor
func (e *Editor) Start(locs []string) {
	err := e.term.Init()
	if err != nil {
		panic(err)
	}

	defer func() {
		e.term.Close()
	}()
	e.term.SetExtendedColors(core.Colors == 256)
	e.theme, err = core.ReadTheme(core.FindResource(path.Join("themes", e.config.Theme)))
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
	dirs := []string{}
	files := []string{}
	for _, loc := range locs {
		if stat, err := os.Stat(loc); err == nil && stat.IsDir() {
			dirs = append(dirs, loc)
		} else {
			files = append(files, loc)
		}
	}
	if len(dirs) == 0 {
		if len(files) > 0 {
			dirs = []string{path.Dir(locs[0])}
		} else {
			dirs = []string{"."}
		}
	}
	e.Cols = append(e.Cols, &Col{WidthRatio: 1.0})
	ratio := 1.0 / float64(len(dirs))
	for _, dir := range dirs {
		view := e.NewView(dir)
		view.HeightRatio = ratio
		e.Cols[0].Views = append(e.Cols[0].Views, view.Id())
		e.Open(dir, view.Id(), "", true)
	}
	e.CurCol = e.Cols[0]
	e.curViewId = e.CurCol.Views[0]
	if len(files) > 0 {
		e.CurCol.WidthRatio = 0.2
		c := &Col{WidthRatio: 0.8}
		ratio := 1.0 / float64(len(files))
		for _, f := range files {
			view := e.NewView(f)
			view.HeightRatio = ratio
			c.Views = append(c.Views, view.Id())
			e.Open(f, view.Id(), "", true)
		}
		e.Cols = append(e.Cols, c)
		e.CurCol = c
		e.curViewId = c.Views[0]
	}

	actions.Ar.EdResize(e.term.Size())

	actions.Ar.EdRender()

	go core.Bus.Start()

	go e.autoScroller()

	go event.Listen()

	e.term.Listen()
}

// Open opens a given location in the editor (in the given view)
// or new view if viewId < 0
func (e *Editor) Open(loc string, viewId int64, rel string, create bool) (int64, error) {
	loc = strings.TrimSpace(loc)
	rel = strings.TrimSpace(rel)
	if len(rel) > 0 && !strings.HasPrefix(loc, string(os.PathSeparator)) {
		loc = path.Join(rel, loc)
	}
	// make it absolute
	loc, err := filepath.Abs(loc)
	if err != nil {
		return -1, err
	}
	stat, err := os.Stat(loc)
	newFile := false
	if os.IsNotExist(err) {
		if !create {
			return -1, err
		}
		newFile = true
	}
	title := filepath.Base(loc)
	if !newFile && stat.IsDir() {
		loc += string(os.PathSeparator)
		title += string(os.PathSeparator)
	}
	nv := false
	var view core.Viewable
	if viewId < 0 {
		view = e.NewFileView(loc)
		nv = true
	} else {
		view = e.ViewById(viewId)
	}
	view.Reset()
	view.SetTitle(title)
	if newFile || !stat.IsDir() {
		err = e.openFile(loc, view)
	} else {
		err = e.openDir(loc, view)
	}
	if err != nil {
		return -1, err
	}
	if nv {
		if stat != nil && stat.IsDir() {
			e.AddDirViewSmart(viewCast(view))
		} else {
			e.InsertViewSmart(viewCast(view))
		}
	}
	view.Reset()
	view.SetWorkDir(filepath.Dir(loc))
	return view.Id(), nil
}

// OpenDir opens a directory listing
func (e *Editor) openDir(loc string, view core.Viewable) error {
	v := viewCast(view)
	if v == nil {
		return fmt.Errorf("No such view")
	}
	v.highlighter = &TermHighlighter{}
	args := append([]string{"ls"}, core.OsLsArgs...)
	title := filepath.Base(loc) + "/"
	backend, err := backend.NewMemBackendCmd(args, loc, view.Id(), &title, true)
	if err != nil {
		return err
	}
	view.SetBackend(backend)
	e.SetStatus(fmt.Sprintf("%v", view.WorkDir()))
	return nil
}

// OpenFile opens a file in the editor
func (e *Editor) openFile(loc string, view core.Viewable) error {
	if !core.IsTextFile(loc) {
		return fmt.Errorf("Binary file ? %s", loc)
	}
	b, err := backend.NewFileBackend(loc, view.Id())
	if err != nil {
		return err
	}
	view.SetBackend(b)
	e.SetStatus(fmt.Sprintf("%v  [%d]", view.WorkDir(), view.Id()))
	view.SetDirty(false)
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

func (e *Editor) Config() core.Config {
	return *e.config
}

func (e *Editor) Theme() *core.Theme {
	return e.theme
}

func (e *Editor) CurView() core.Viewable {
	v, found := e.views[e.curViewId]
	if !found {
		return e.views[e.Cols[0].Views[0]]
	}
	return v
}

func (e *Editor) CurViewId() int64 {
	return e.curViewId
}

func (e *Editor) SetCursor(y, x int) {
	e.term.SetCursor(x, y)
}

func (e *Editor) CmdOn() bool {
	return e.cmdOn
}

func (e *Editor) SetCmdOn(v bool) {
	e.cmdOn = v
}

func (e *Editor) TermFlush() {
	e.term.Flush()
}

func (e *Editor) QuitCheck() bool {
	for _, c := range e.Cols {
		for _, vi := range c.Views {
			v, found := e.views[vi]
			if found && !v.canClose() {
				return false
			}
		}
	}
	return true
}

func (e *Editor) StartTermView(args []string) int64 {
	vid := exec(args, true)
	v := viewCast(core.Ed.ViewById(vid))
	if v == nil || v.backend == nil {
		return -1
	}
	b := v.backend.(*backend.BackendCmd)
	time.Sleep(500 * time.Millisecond)
	ext := ".sh"
	if os.Getenv("SHELL") == "rc" {
		ext = ".rc"
	}
	cmd := ". $HOME/.goed/default/actions/goed" +
		fmt.Sprintf("%s %d %d\n", ext, core.InstanceId, v.Id())
	b.SendBytes([]byte(cmd))
	return vid
}

// Handle selection auto scrolling of views
func (e *Editor) autoScroller() {
	action := autoScrollAction{}
	for {
		time.Sleep(200 * time.Millisecond)
		core.Bus.Dispatch(action)
	}
}

type autoScrollAction struct {
}

func (e autoScrollAction) Run() {
	v := viewCast(core.Ed.ViewById(core.Ed.CurViewId()))
	if v == nil {
		return
	}
	x, y := v.autoScrollX, v.autoScrollY
	if x == 0 && y == 0 {
		return
	}
	if len(v.selections) == 0 {
		return
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
	core.Ed.Render()
}

// TODO: Do away with those ugly assertions
func viewCast(v core.Viewable) *View {
	if v == nil {
		return nil
	}
	return v.(*View)
}

func widgetCast(w Renderer) *View {
	if w == nil {
		return nil
	}
	if v, ok := w.(*View); ok {
		return v
	}
	return nil
}
