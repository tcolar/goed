package actions

import (
	"fmt"

	"github.com/tcolar/goed/core"
)

// Add a text selection to the view. from l1,c1 to l2,c2
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

// return the current view location in the ui
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

// return the current cursor position in the view
func (a *ar) ViewCurPos(viewId int64) (ln, col int) {
	answer := make(chan int, 2)
	d(viewCurPos{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

// send a movement event to the view. (ie: down, up, left, right, etc...)
func (a *ar) ViewCursorMvmt(viewId int64, mvmt core.CursorMvmt) {
	d(viewCursorMvmt{viewId: viewId, mvmt: mvmt})
}

// delete text from the view (from row1,col1 to row2,col2)
func (a *ar) ViewDelete(viewId int64, row1, col1, row2, col2 int, undoable bool) {
	d(viewDeleteAction{viewId: viewId, row1: row1, col1: col1, row2: row2, col2: col2, undoable: undoable})
}

// delete text from the view (current selection, if none, current line)
func (a *ar) ViewDeleteCur(viewId int64) {
	d(viewDeleteCur{viewId: viewId})
}

// insert text into the view at the row,col location
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

// move the cursor by y,x (relative)
func (a *ar) ViewMoveCursor(viewId int64, y, x int) {
	d(viewMoveCursor{viewId: viewId, x: x, y: y})
}

// move the cursor by y,x (relative) but also scroll te view as needed to keep
// the cursor in view and in place.
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

// extends the current text selection toward the given location
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

// return the current scrolling position (0,0 is top, left)
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

// undo
func (a *ar) ViewUndo(viewId int64) {
	d(viewUndo{viewId: viewId})
}

// ########  Impl ......

type viewAddSelection struct {
	viewId         int64
	l1, c1, l2, c2 int
}

func (a viewAddSelection) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	s := core.NewSelection(a.l1, a.c1, a.l2, a.c2)
	selections := v.Selections()
	*selections = append(*selections, *s)
	return nil
}

type viewAutoScroll struct {
	viewId int64
	y, x   int
	on     bool
}

func (a viewAutoScroll) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SetAutoScroll(a.y, a.x, a.on)
	return nil
}

type viewBackspace struct {
	viewId int64
}

func (a viewBackspace) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Backspace()
	return nil
}

type viewBounds struct {
	answer chan int
	viewId int64
}

func (a viewBounds) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
		a.answer <- 0
		a.answer <- 0
	}
	l1, c1, l2, c2 := v.Bounds()
	a.answer <- l1
	a.answer <- c1
	a.answer <- l2
	a.answer <- c2
	return nil
}

type viewClearSelections struct {
	viewId int64
}

func (a viewClearSelections) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.ClearSelections()
	return nil
}

type viewCmdStop struct {
	viewId int64
}

func (a viewCmdStop) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	b := v.Backend()
	if b != nil {
		b.Close()
	}
	return nil
}

type viewCols struct {
	answer chan int
	viewId int64
}

func (a viewCols) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
	}
	a.answer <- v.LastViewCol()
	return nil
}

type viewCopy struct {
	viewId int64
}

func (a viewCopy) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Copy()
	return nil
}

type viewCurPos struct {
	answer chan int
	viewId int64
}

func (a viewCurPos) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
	}
	a.answer <- v.CurLine()
	a.answer <- v.CurCol()
	return nil
}

type viewCursorMvmt struct {
	viewId int64
	mvmt   core.CursorMvmt
}

func (a viewCursorMvmt) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.CursorMvmt(a.mvmt)
	return nil
}

type viewCut struct {
	viewId int64
}

func (a viewCut) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Cut()
	return nil
}

type viewDeleteAction struct {
	viewId                 int64
	row1, col1, row2, col2 int
	undoable               bool
}

func (a viewDeleteAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Delete(a.row1, a.col1, a.row2, a.col2, a.undoable)
	return nil
}

type viewDeleteCur struct {
	viewId int64
}

func (a viewDeleteCur) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.DeleteCur()
	return nil
}

type viewInsertAction struct {
	viewId   int64
	row, col int
	text     string
	undoable bool
}

func (a viewInsertAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Insert(a.row, a.col, a.text, a.undoable)
	return nil
}

type viewInsertCur struct {
	viewId int64
	text   string
}

func (a viewInsertCur) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.InsertCur(a.text)
	return nil
}

type viewInsertNewLine struct {
	viewId int64
}

func (a viewInsertNewLine) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.InsertNewLineCur()
	return nil
}

type viewMoveCursor struct {
	viewId int64
	status string
	y, x   int
	roll   bool
}

func (a viewMoveCursor) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	if a.roll {
		v.MoveCursorRoll(a.y, a.x)
	} else {
		v.MoveCursor(a.y, a.x)
	}
	return nil
}

type viewOpenSelection struct {
	viewId  int64
	newView bool
}

func (a viewOpenSelection) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.OpenSelection(a.newView)
	return nil
}

type viewPaste struct {
	viewId int64
}

func (a viewPaste) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Paste()
	return nil
}

type viewRedo struct {
	viewId int64
}

func (a viewRedo) Run() error {
	return Redo(a.viewId)
}

type viewReload struct{ viewId int64 }

func (a viewReload) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Reload()
	}
	return nil
}

type viewRender struct {
	viewId int64
}

func (a viewRender) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Render()
	}
	return nil
}

type viewRows struct {
	answer chan int
	viewId int64
}

func (a viewRows) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
	}
	a.answer <- v.LastViewLine()
	return nil
}

type viewSave struct {
	viewId int64
}

func (a viewSave) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Save()
	return nil
}

type viewScrollPos struct {
	answer chan int
	viewId int64
}

func (a viewScrollPos) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
	}
	a.answer <- v.CurLine()
	a.answer <- v.CurCol()
	return nil
}

type viewSelectAll struct {
	viewId int64
}

func (a viewSelectAll) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SelectAll()
	return nil
}

type viewSetDirty struct {
	viewId int64
	on     bool
}

func (a viewSetDirty) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SetDirty(a.on)
	return nil
}

type viewSetTitle struct {
	viewId int64
	title  string
}

func (a viewSetTitle) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetTitle(a.title)
	}
	return nil
}

type viewSetWorkdir struct {
	viewId  int64
	workDir string
}

func (a viewSetWorkdir) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetWorkDir(a.workDir)
	}
	return nil
}

type viewSetVtCols struct {
	viewId int64
	cols   int
}

func (a viewSetVtCols) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetVtCols(a.cols)
	}
	return nil
}

type viewSrcLoc struct {
	viewId int64
	answer chan string
}

func (a viewSrcLoc) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil || v.Id() == 0 {
		a.answer <- fmt.Sprintf("No such view %d", a.viewId)
		return fmt.Errorf("No such view : %d", a.viewId)
	}
	a.answer <- v.Backend().SrcLoc()
	return nil
}

type viewStretchSelection struct {
	viewId          int64
	prevLn, prevCol int
}

func (a viewStretchSelection) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.StretchSelection(
		a.prevLn,
		v.LineRunesTo(v.Slice(), a.prevLn, a.prevCol),
		v.CurLine(),
		v.LineRunesTo(v.Slice(), v.CurLine(), v.CurCol()),
	)
	return nil
}

type viewSyncSlice struct {
	viewId int64
}

func (a viewSyncSlice) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SyncSlice()
	return nil
}

type viewTrim struct {
	viewId int64
	limit  int
}

type viewUndo struct {
	viewId int64
}

func (a viewUndo) Run() error {
	return Undo(a.viewId)
}

func NewViewInsertAction(viewId int64, row, col int, text string, undoable bool) core.Action {
	return viewInsertAction{viewId: viewId, row: row, col: col, text: text, undoable: undoable}
}

func NewViewDeleteAction(viewId int64, row1, col1, row2, col2 int, undoable bool) core.Action {
	return viewDeleteAction{viewId: viewId, row1: row1, col1: col1, row2: row2, col2: col2, undoable: undoable}
}
