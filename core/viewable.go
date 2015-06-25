package core

// Viewable is the interface to a View
type Viewable interface {
	Backend() Backend
	CurCol() int
	CurLine() int
	// CusorTextPos return the text position (rune count) for a given Cursor
	// position
	CursorTextPos(s *Slice, c, l int) (col, line int)
	Dirty() bool
	Id() int64
	LineCount() int
	// MoveCusrosr moves the cursor by the x, y offsets
	MoveCursor(x, y int)
	// Reload reloads the view data from it's source (backend)
	Reload()
	// Render forec re=-rendering the view UI.
	Render()
	// Reset reinitializes the view to it's startup state.
	Reset()
	SetBackend(backend Backend)
	SetDirty(bool)
	Selections() *[]Selection
	// SetAutoScroll is used to make the view scroll contonuously in x,y increments
	// keeps scrolling until x and y are set to 0.
	SetAutoScroll(x, y int, isSelect bool)
	// Sets the view work directory, commands and "open" actions will be relative
	// to this path.
	SetWorkDir(dir string)
	SetTitle(title string)
	// Slice returns a view's text subset (matrix)
	Slice() *Slice
	Title() string
	WorkDir() string
}
