package core

// Viewable is the interface to a View
type Viewable interface {
	Widget
	Backspace()
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
	// Reset reinitializes the view to it's startup state.
	Reset()
	Save() // Save from buffer to src
	ScrollPos() (ln, col int)
	SetBackend(backend Backend)
	SetDirty(bool)
	SelectAll()
	SelectWord(ln, col int)
	Selections() *[]Selection
	// SetAutoScroll is used to make the view scroll contonuously in y,x increments
	// keeps scrolling until x and y are set to 0.
	SetAutoScroll(y, x int, isSelect bool)
	SetCursorPos(y, x int)
	SetScrollPct(ypct int)
	SetScrollPos(y, x int)
	SetTitle(title string)
	SetViewType(t ViewType)
	SetVtCols(cols int)
	// Sets the view work directory, commands and "open" actions will be relative
	// to this path.
	SetWorkDir(dir string)
	// Slice returns a view's text subset (rectangle)
	Slice() *Slice
	StretchSelection(prevl, prevc, ln, c int)
	SyncSlice()
	Title() string
	Text(ln1, col1, ln2, col2 int) [][]rune
	Type() ViewType
	WorkDir() string
}
