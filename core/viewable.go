package core

type Viewable interface {
	Backend() Backend
	CurCol() int
	CurLine() int
	CursorTextPos(s *Slice, c, l int) (col, ln int)
	Dirty() bool
	Id() int
	LineCount() int
	MoveCursor(x, y int)
	Render()
	Reset()
	SetBackend(backend Backend)
	SetDirty(bool)
	Selections() *[]Selection
	SetWorkDir(dir string)
	SetTitle(title string)
	Slice() *Slice
	Title() string
	WorkDir() string
}