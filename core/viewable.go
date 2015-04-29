package core

type Viewable interface {
	CurCol() int
	CurLine() int
	CursorTextPos(s *Slice, c, l int) (col, ln int)
	Id() int
	LineCount() int
	MoveCursor(x, y int)
	Render()
	Reset()
	SetBackend(backend Backend)
	SetWorkDir(dir string)
	SetTitle(title string)
	SetDirty(dirty bool)
	Slice() *Slice
	WorkDir() string
}
