package core

// Editable provides editor features entry poins.
type Editable interface {
	ActivateView(v Viewable, y, x int)
	CmdbarToggle()
	Config() Config
	CurColIndex() int
	CurView() Viewable
	DelColCheckByIndex(col int)
	DelView(view Viewable, terminate bool)
	DelViewCheck(view Viewable)
	Dispatch(action Action)
	// CmdOn indicates whether the CommandBar is currently active
	CmdOn() bool
	// Open opens a file in the given view.
	Open(loc string, view Viewable, rel string, create bool) (Viewable, error)
	QuitCheck() bool
	// Render updates the whole editor UI
	Render()
	Resize(h, w int)
	// SetStatusErr displays an error message in the status bar
	SetStatusErr(err string)
	// SetStatusErr displays a message in the status bar
	SetStatus(status string)
	SetCursor(y, x int)
	SetCurView(id int64) error
	// SetCmdOn activates or desactives the CommandBar
	SetCmdOn(v bool)
	SwapViews(v1, v2 Viewable)
	Start(locs []string)
	TermChar(y, x int, c rune)
	TermFB(fg, bg Style)
	TermFill(c rune, y1, x1, y2, x2 int)
	TermFlush()
	TermStr(y, x int, s string)
	TermStrv(y, x int, s string)
	Theme() *Theme
	// ViewByLoc finds if there is an existing view for the given file (loc)
	ViewByLoc(loc string) Viewable
	ViewById(id int64) Viewable
	// Move a view
	ViewMove(y1, x1, y2, x2 int)
	// Navigate from a view to another
	ViewNavigate(mvmt CursorMvmt)
}
