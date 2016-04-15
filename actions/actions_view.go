package actions

import "github.com/tcolar/goed/core"

// Add a text selection to the view. from l1,c1 to l2,c2 (1 indexed)
func (a *ar) ViewAddSelection(viewId int64, l1, c1, l2, c2 int) {
	d(viewAddSelection{viewId: viewId, l1: l1, c1: c1, l2: l2, c2: c2})
}

// Enable/disable a view autoscrolling (by y,x increments)
func (a *ar) ViewAutoScroll(viewId int64, y, x int, on bool) {
	d(viewAutoScroll{viewId: viewId, x: x, y: y, on: on})
}

// send 'backspace' to the view
func (a *ar) ViewBackspace(viewId int64) {
	d(viewBackspace{viewId: viewId})
}

// return the current view location in the ui (1 indexed)
func (a *ar) ViewBounds(viewId int64) (ln, col, ln2, col2 int) {
	answer := make(chan int, 4)
	d(viewBounds{viewId: viewId, answer: answer})
	return <-answer, <-answer, <-answer, <-answer
}

// remove all the view selections.
func (a *ar) ViewClearSelections(viewId int64) {
	d(viewClearSelections{viewId: viewId})
}

// stop the command currenty running in the view (for exec views.)
func (a *ar) ViewCmdStop(viewId int64) {
	d(viewCmdStop{viewId: viewId})
}

// return the nuber of columns (width) of the view.
func (a *ar) ViewCols(viewId int64) (cols int) {
	answer := make(chan int, 1)
	d(viewCols{viewId: viewId, answer: answer})
	return <-answer
}

// copy text from the view (current selection, if none, current line)
func (a *ar) ViewCopy(viewId int64) {
	d(viewCopy{viewId: viewId})
}

// cut text from the view (current selection, if none, current line)
func (a *ar) ViewCut(viewId int64) {
	d(viewCut{viewId: viewId})
}

// return the current cursor UI position in the view (1 indexed)
func (a *ar) ViewCursorCoords(viewId int64) (y, x int) {
	answer := make(chan int, 2)
	d(viewCursorCoords{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

// return the current cursor text position in the view (1 indexed)
func (a *ar) ViewCursorPos(viewId int64) (y, x int) {
	answer := make(chan int, 2)
	d(viewCursorPos{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

// send a movement event to the view. (ie: down, up, left, right, etc...)
func (a *ar) ViewCursorMvmt(viewId int64, mvmt core.CursorMvmt) {
	d(viewCursorMvmt{viewId: viewId, mvmt: mvmt})
}

// delete text from the view (from row1,col1 to row2,col2). 1 indexed
func (a *ar) ViewDelete(viewId int64, row1, col1, row2, col2 int, undoable bool) {
	d(viewDeleteAction{viewId: viewId, row1: row1, col1: col1, row2: row2, col2: col2, undoable: undoable})
}

// delete text from the view (current selection, if none, current line)
func (a *ar) ViewDeleteCur(viewId int64) {
	d(viewDeleteCur{viewId: viewId})
}

// insert text into the view at the row,col location. 1 indexed
func (a *ar) ViewInsert(viewId int64, row, col int, text string, undoable bool) {
	d(viewInsertAction{viewId: viewId, row: row, col: col, text: text, undoable: undoable})
}

// insert text into the view at the current cursor location
func (a *ar) ViewInsertCur(viewId int64, text string) {
	d(viewInsertCur{viewId: viewId, text: text})
}

// insert a newLine at the current cursor location
func (a *ar) ViewInsertNewLine(viewId int64) {
	d(viewInsertNewLine{viewId: viewId})
}

// move the cursor by ln, col runes (relative), scroll as needed
func (a *ar) ViewMoveCursor(viewId int64, y, x int) {
	d(viewMoveCursor{viewId: viewId, x: x, y: y})
}

// move the cursor by y,x runes (relative) but also "roll" to rev/next line as needed.
func (a *ar) ViewMoveCursorRoll(viewId int64, y, x int) {
	d(viewMoveCursor{viewId: viewId, x: x, y: y, roll: true})
}

// paste text into the view at the curent location
// if in a selection, paste over it.
func (a *ar) ViewPaste(viewId int64) {
	d(viewPaste{viewId: viewId})
}

// try to "open" the current selection into a view (ie: expect a file path)
func (a *ar) ViewOpenSelection(viewId int64, newView bool) {
	d(viewOpenSelection{viewId: viewId, newView: newView})
}

// redo
func (a *ar) ViewRedo(viewId int64) {
	d(viewRedo{viewId: viewId})
}

// reload the view from it's source file, discard all unsaved buffer changes
func (a *ar) ViewReload(viewId int64) {
	d(viewReload{viewId: viewId})
}

// render/repaint the view
func (a *ar) ViewRender(viewId int64) {
	d(viewRender{viewId: viewId})
}

// return the number of rows (lines) in the view
func (a *ar) ViewRows(viewId int64) (rows int) {
	answer := make(chan int, 1)
	d(viewRows{viewId: viewId, answer: answer})
	return <-answer
}

// save the view content to the backing file
func (a *ar) ViewSave(viewId int64) {
	d(viewSave{viewId: viewId})
}

// select all
func (a *ar) ViewSelectAll(viewId int64) {
	d(viewSelectAll{viewId: viewId})
}

// Return a list of view selctions (one per line), 1 indexed, ie:
// 2 1 2 6
// 3 2 4 7
func (a *ar) ViewSelections(viewId int64) []core.Selection {
	answer := make(chan []core.Selection, 1)
	d(viewSelections{answer: answer, viewId: viewId})
	return <-answer
}

// move the cursor to the given text position(1 indexed), scroll as needed
func (a *ar) ViewSetCursorPos(viewId int64, y, x int) {
	d(viewSetCursorPos{viewId: viewId, y: y, x: x})
}

// mark the view "dirty" or not (ie: modified, unsaved)
func (a *ar) ViewSetDirty(viewId int64, on bool) {
	d(viewSetDirty{viewId: viewId, on: on})
}

// se the view title (typically file path)
func (a *ar) ViewSetTitle(viewId int64, title string) {
	d(viewSetTitle{viewId: viewId, title: title})
}

// set the number of vt100 columns, this is useful so that tty programs that
// can use the full view wisth properly
func (a *ar) ViewSetVtCols(viewId int64, cols int) {
	d(viewSetVtCols{viewId: viewId, cols: cols})
}

// extends the current text selection toward the given location(1 indexed)
// this maybe in any direction
func (a *ar) ViewStretchSelection(viewId int64, prevLn, prevCol int) {
	d(viewStretchSelection{viewId: viewId, prevLn: prevLn, prevCol: prevCol})
}

// set the current working dir of the view, especially usefull for terminal views.
// this is used when "opening" relative locations, among other things.
func (a *ar) ViewSetWorkdir(viewId int64, workDir string) {
	d(viewSetWorkdir{viewId: viewId, workDir: workDir})
}

// return the absolute path of the file backing the view (if any)
func (a *ar) ViewSrcLoc(viewId int64) string {
	answer := make(chan string, 1)
	d(viewSrcLoc{viewId: viewId, answer: answer})
	return <-answer
}

// return the current scrolling position (1,1 is top, left)
func (a *ar) ViewScrollPos(viewId int64) (ln, col int) {
	answer := make(chan int, 2)
	d(viewScrollPos{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

// this force sync of the in memory slice representing the part of the content
// that is currently visible in the view (performance optimization)
func (a *ar) ViewSyncSlice(viewId int64) {
	d(viewSyncSlice{viewId: viewId})
}

// return the vew title
func (a *ar) ViewTitle(viewId int64) string {
	answer := make(chan string, 1)
	d(viewTitle{viewId: viewId, answer: answer})
	return <-answer
}

// undo
func (a *ar) ViewUndo(viewId int64) {
	d(viewUndo{viewId: viewId})
}

// ########  Impl ......

type viewAddSelection struct {
	viewId         int64
	l1, c1, l2, c2 int
}

func (a viewAddSelection) Run() {
	v := core.Ed.ViewById(a.viewId)
	if a.l2 != -1 {
		a.l2--
	}
	if a.c2 != -1 {
		a.c2--
	}
	if v != nil {
		s := core.NewSelection(a.l1-1, a.c1-1, a.l2, a.c2)
		selections := v.Selections()
		*selections = append(*selections, *s)
	}
}

type viewAutoScroll struct {
	viewId int64
	y, x   int
	on     bool
}

func (a viewAutoScroll) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetAutoScroll(a.y, a.x, a.on)
	}
}

type viewBackspace struct {
	viewId int64
}

func (a viewBackspace) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Backspace()
	}
}

type viewBounds struct {
	answer chan int
	viewId int64
}

func (a viewBounds) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
		a.answer <- 0
		a.answer <- 0
		return
	}
	l1, c1, l2, c2 := v.Bounds()
	a.answer <- l1 + 1
	a.answer <- c1 + 1
	a.answer <- l2 + 1
	a.answer <- c2 + 1
}

type viewClearSelections struct {
	viewId int64
}

func (a viewClearSelections) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.ClearSelections()
	}
}

type viewCmdStop struct {
	viewId int64
}

func (a viewCmdStop) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return
	}
	b := v.Backend()
	if b != nil {
		b.Close()
	}
	return
}

type viewCols struct {
	answer chan int
	viewId int64
}

func (a viewCols) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		return
	}
	a.answer <- v.LastViewCol()
}

type viewCopy struct {
	viewId int64
}

func (a viewCopy) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Copy()
	}
}

type viewCursorCoords struct {
	answer chan int
	viewId int64
}

func (a viewCursorCoords) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
		return
	}
	a.answer <- v.CurLine() + 1
	a.answer <- v.CurCol() + 1
}

type viewCursorPos struct {
	answer chan int
	viewId int64
}

func (a viewCursorPos) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
		return
	}
	a.answer <- v.CurLine() + 1
	a.answer <- v.LineRunesTo(v.Slice(), v.CurLine(), v.CurCol()) + 1
}

type viewCursorMvmt struct {
	viewId int64
	mvmt   core.CursorMvmt
}

func (a viewCursorMvmt) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.CursorMvmt(a.mvmt)
	}
}

type viewCut struct {
	viewId int64
}

func (a viewCut) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Cut()
	}
}

type viewDeleteAction struct {
	viewId                 int64
	row1, col1, row2, col2 int
	undoable               bool
}

func (a viewDeleteAction) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		if a.row2 != -1 {
			a.row2--
		}
		if a.col2 != -1 {
			a.col2--
		}
		v.Delete(a.row1-1, a.col1-1, a.row2, a.col2, a.undoable)
	}
}

type viewDeleteCur struct {
	viewId int64
}

func (a viewDeleteCur) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.DeleteCur()
	}
}

type viewInsertAction struct {
	viewId   int64
	row, col int
	text     string
	undoable bool
}

func (a viewInsertAction) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Insert(a.row-1, a.col-1, a.text, a.undoable)
	}
}

type viewInsertCur struct {
	viewId int64
	text   string
}

func (a viewInsertCur) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.InsertCur(a.text)
	}
}

type viewInsertNewLine struct {
	viewId int64
}

func (a viewInsertNewLine) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.InsertNewLineCur()
	}
}

type viewMoveCursor struct {
	viewId int64
	y, x   int
	roll   bool
}

func (a viewMoveCursor) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return
	}
	if a.roll {
		v.MoveCursorRoll(a.y, a.x)
	} else {
		v.MoveCursor(a.y, a.x)
	}
}

type viewOpenSelection struct {
	viewId  int64
	newView bool
}

func (a viewOpenSelection) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.OpenSelection(a.newView)
	}
}

type viewPaste struct {
	viewId int64
}

func (a viewPaste) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Paste()
	}
}

type viewRedo struct {
	viewId int64
}

func (a viewRedo) Run() {
	if viewExists(a.viewId) {
		Redo(a.viewId)
	}
}

type viewReload struct{ viewId int64 }

func (a viewReload) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Reload()
	}
}

type viewRender struct {
	viewId int64
}

func (a viewRender) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Render()
	}
}

type viewRows struct {
	answer chan int
	viewId int64
}

func (a viewRows) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		return
	}
	a.answer <- v.LastViewLine()
}

type viewSave struct {
	viewId int64
}

func (a viewSave) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Save()
	}
}

type viewScrollPos struct {
	answer chan int
	viewId int64
}

func (a viewScrollPos) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
		return
	}
	y, x := v.ScrollPos()
	a.answer <- y + 1
	a.answer <- x + 1
}

type viewSelectAll struct {
	viewId int64
}

func (a viewSelectAll) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SelectAll()
	}
}

type viewSelections struct {
	answer chan []core.Selection
	viewId int64
}

func (a viewSelections) Run() {
	v := core.Ed.ViewById(a.viewId)
	result := []core.Selection{}
	if v == nil {
		a.answer <- result
		return
	}
	for _, s := range *v.Selections() {
		ct := 1
		lt := 1
		if s.ColTo == -1 {
			ct = 0
		}
		if s.LineTo == -1 {
			lt = 0
		}
		result = append(result, *core.NewSelection(s.LineFrom+1, s.ColFrom+1,
			s.LineTo+lt, s.ColTo+ct))
	}
	a.answer <- result
}

type viewSetCursorPos struct {
	viewId int64
	y, x   int
}

func (a viewSetCursorPos) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.MoveCursor(a.y-v.CurLine()-1, a.x-v.CurCol()-1)
	}
}

type viewSetDirty struct {
	viewId int64
	on     bool
}

func (a viewSetDirty) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetDirty(a.on)
	}
}

type viewSetTitle struct {
	viewId int64
	title  string
}

func (a viewSetTitle) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetTitle(a.title)
	}
}

type viewSetWorkdir struct {
	viewId  int64
	workDir string
}

func (a viewSetWorkdir) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetWorkDir(a.workDir)
	}
}

type viewSetVtCols struct {
	viewId int64
	cols   int
}

func (a viewSetVtCols) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetVtCols(a.cols)
	}
}

type viewSrcLoc struct {
	viewId int64
	answer chan string
}

func (a viewSrcLoc) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil || v.Id() == 0 {
		a.answer <- ""
		return
	}
	a.answer <- v.Backend().SrcLoc()
}

type viewStretchSelection struct {
	viewId          int64
	prevLn, prevCol int
}

func (a viewStretchSelection) Run() {
	v := core.Ed.ViewById(a.viewId)
	ln, col := v.CurLine(), v.CurCol()
	if v != nil {
		v.StretchSelection(
			a.prevLn,
			v.LineRunesTo(v.Slice(), a.prevLn-1, a.prevCol-1),
			ln,
			v.LineRunesTo(v.Slice(), ln, col),
		)
	}
}

type viewSyncSlice struct {
	viewId int64
}

func (a viewSyncSlice) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SyncSlice()
	}
}

type viewTitle struct {
	viewId int64
	answer chan string
}

func (a viewTitle) Run() {
	v := core.Ed.ViewById(a.viewId)
	if v == nil || v.Id() == 0 {
		a.answer <- ""
		return
	}
	a.answer <- v.Title()
}

type viewUndo struct {
	viewId int64
}

func (a viewUndo) Run() {
	if viewExists(a.viewId) {
		Undo(a.viewId)
	}
}

func NewViewInsertAction(viewId int64, row, col int, text string, undoable bool) core.Action {
	return viewInsertAction{viewId: viewId, row: row + 1, col: col + 1,
		text: text, undoable: undoable}
}

func NewViewDeleteAction(viewId int64, row1, col1, row2, col2 int, undoable bool) core.Action {
	return viewDeleteAction{viewId: viewId, row1: row1 + 1, col1: col1 + 1,
		row2: row2 + 1, col2: col2 + 1, undoable: undoable}
}
