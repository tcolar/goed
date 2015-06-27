package core

// Viewable is the interface to a View
type Viewable interface {
	Bounds() (y1, x1, y2, x2 int)
	Backend() Backend
	CurCol() int
	CurLine() int
	Dirty() bool
	Id() int64
	LineCount() int
	// LineRunesTo returns the number of raw runes to the given line column
	LineRunesTo(s *Slice, line, col int) int
	// MoveCursor moves the cursor by the y, x offsets (in runes)
	MoveCursor(y, x int)
	// Reload reloads the view data from it's source (backend)
	Reload()
	// Render forec re=-rendering the view UI.
	Render()
	// Reset reinitializes the view to it's startup state.
	Reset()
	SetBackend(backend Backend)
	SetDirty(bool)
	Selections() *[]Selection
	// SetAutoScroll is used to make the view scroll contonuously in y,x increments
	// keeps scrolling until x and y are set to 0.
	SetAutoScroll(y, x int, isSelect bool)
	// Sets the view work directory, commands and "open" actions will be relative
	// to this path.
	SetWorkDir(dir string)
	SetTitle(title string)
	// Slice returns a view's text subset (matrix)
	Slice() *Slice
	Title() string
	WorkDir() string
}
