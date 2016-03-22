package core

// Editable provides editor features entry poins.
type Editable interface {
	CmdbarToggle()
	Config() Config
	CurColIndex() int
	CurViewId() int64
	DelColByIndex(col int, check bool)
	DelView(viewId int64, terminate bool)
	DelViewByIndex(viewId int64, check bool)
	Dispatch(action Action)
	// CmdOn indicates whether the CommandBar is currently active
	CmdOn() bool
	// Open opens a file in the given view (new view if viewid<0)
	// create -> create file at loc if does not exist yet
	Open(loc string, viewId int64, rel string, create bool) (int64, error)
	QuitCheck() bool
	// Render updates the whole editor UI
	Render()
	Resize(h, w int)
	StartTermView(args []string) int64
	// SetStatusErr displays an error message in the status bar
	SetStatusErr(err string)
	// SetStatusErr displays a message in the status bar
	SetStatus(status string)
	SetCursor(y, x int)
	SetCurView(id int64) error
	// SetCmdOn activates or desactives the CommandBar
	SetCmdOn(v bool)
	SwapViews(v1, v2 int64)
	Start(locs []string)
	TermChar(y, x int, c rune)
	TermFB(fg, bg Style)
	TermFill(c rune, y1, x1, y2, x2 int)
	TermFlush()
	TermStr(y, x int, s string)
	TermStrv(y, x int, s string)
	Theme() *Theme
	ViewActivate(v int64, y, x int)
	ViewAt(ln, col int) int64
	// ViewByLoc finds if there is an existing view for the given file (loc)
	ViewByLoc(loc string) int64
	ViewById(id int64) Viewable
	// Move a view
	ViewMove(y1, x1, y2, x2 int)
	// Navigate from a view to another
	ViewNavigate(mvmt CursorMvmt)
}
