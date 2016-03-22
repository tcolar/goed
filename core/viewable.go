package core

// Viewable is the interface to a View
type Viewable interface {
	Backspace()
	Bounds() (y1, x1, y2, x2 int)
	Backend() Backend
	ClearSelections()
	Copy()
	CurCol() int
	CurLine() int
	CursorMvmt(mvmt CursorMvmt)
	Cut()
	Delete(row1, col1, row, col2 int, undoable bool)
	DeleteCur()
	Dirty() bool
	Id() int64
	Insert(row, col int, text string, undoable bool)
	InsertCur(text string)
	InsertNewLineCur()
	LastViewCol() int
	LastViewLine() int
	LineCount() int
	// LineRunesTo returns the number of raw runes to the given line column
	LineRunesTo(s *Slice, line, col int) int
	// MoveCursor moves the cursor by the y, x offsets (in runes)
	MoveCursor(y, x int)
	MoveCursorRoll(y, x int)
	OpenSelection(newView bool)
	Paste()
	// Reload reloads the view data from it's source (backend)
	Reload()
	// Render forec re=-rendering the view UI.
	Render()
	// Reset reinitializes the view to it's startup state.
	Reset()
	Save() // Save from buffer to src
	ScrollPos() (ln, col int)
	SetBackend(backend Backend)
	SetDirty(bool)
	SelectAll()
	Selections() *[]Selection
	// SetAutoScroll is used to make the view scroll contonuously in y,x increments
	// keeps scrolling until x and y are set to 0.
	SetAutoScroll(y, x int, isSelect bool)
	SetViewType(t ViewType)
	// Sets the view work directory, commands and "open" actions will be relative
	// to this path.
	SetWorkDir(dir string)
	SetTitle(title string)
	SetVtCols(cols int)
	// Slice returns a view's text subset (matrix)
	Slice() *Slice
	StretchSelection(prevl, prevc, ln, c int)
	SyncSlice()
	Terminated() bool // view is nil or marked for termination
	Title() string
	WorkDir() string
}
