package core

type Viewable interface {
	Backend() Backend
	CurCol() int
	CurLine() int
	CursorTextPos(s *Slice, c, l int) (col, line int)
	Dirty() bool
	Id() int64
	LineCount() int
	MoveCursor(x, y int)
	Reload()
	Render()
	Reset()
	SetBackend(backend Backend)
	SetDirty(bool)
	Selections() *[]Selection
	SetAutoScroll(x, y int, isSelect bool)
	SetWorkDir(dir string)
	SetTitle(title string)
	Slice() *Slice
	Title() string
	WorkDir() string
}
