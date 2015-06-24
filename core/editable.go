package core

type Editable interface {
	Config() Config
	CurView() Viewable
	DelView(view Viewable, terminate bool)
	CmdOn() bool
	Open(loc string, view Viewable, title string) (Viewable, error)
	Render()
	SetStatusErr(err string)
	SetStatus(status string)
	SetCursor(x, y int)
	SetCurView(id int64) error
	SetCmdOn(v bool)
	Start(loc string)
	TermChar(x, y int, c rune)
	TermFB(fg, bg Style)
	TermFill(c rune, x1, y1, x2, y2 int)
	TermFlush()
	TermStr(x, y int, s string)
	TermStrv(x, y int, s string)
	Theme() *Theme
	ViewByLoc(loc string) Viewable
	ViewById(id int64) Viewable
}
