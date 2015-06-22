package core

type Editable interface {
	Config() Config
	CurView() Viewable
	CmdOn() bool
	ViewById(id int64) Viewable
	Render()
	SetStatusErr(err string)
	SetStatus(status string)
	Open(loc string, view Viewable, title string) error
	SetCursor(x, y int)
	SetCurView(id int64) error
	SetCmdOn(v bool)
	Start(loc string)
	TermFB(fg, bg Style)
	TermChar(x, y int, c rune)
	TermStr(x, y int, s string)
	TermStrv(x, y int, s string)
	TermFill(c rune, x1, y1, x2, y2 int)
	TermFlush()
	Theme() *Theme
}
