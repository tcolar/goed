package core

// Editable provides editor features entry poins.
type Editable interface {
	Config() Config
	CurView() Viewable
	DelView(view Viewable, terminate bool)
	// CmdOn indicates whether the CommandBar is currently active
	CmdOn() bool
	// Open opens a file in the given view.
	Open(loc string, view Viewable, title string) (Viewable, error)
	// Render updates the whole editor UI
	Render()
	// SetStatusErr displays an error message in the status bar
	SetStatusErr(err string)
	// SetStatusErr displays a message in the status bar
	SetStatus(status string)
	SetCursor(y, x int)
	SetCurView(id int64) error
	// SetCmdOn activates or desactives the CommandBar
	SetCmdOn(v bool)
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
}
