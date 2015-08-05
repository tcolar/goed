package actions

import "github.com/tcolar/goed/core"

func ViewAddSelection(viewId int64, l1, c1, l2, c2 int) {
	d(viewAddSelection{viewId: viewId, l1: l1, c1: c1, l2: l2, c2: c2})
}

func ViewAutoScroll(viewId int64, y, x int, on bool) {
	d(viewAutoScroll{viewId: viewId, x: x, y: y, on: on})
}

func ViewBackspace(viewId int64) {
	d(viewBackspace{viewId: viewId})
}

func ViewClearSelections(viewId int64) {
	d(viewClearSelections{viewId: viewId})
}

func ViewCmdStop(viewId int64) {
	d(viewCmdStop{viewId: viewId})
}

func ViewCopy(viewId int64) {
	d(viewCopy{viewId: viewId})
}

func ViewCut(viewId int64) {
	d(viewCut{viewId: viewId})
}

func ViewCurPos(viewId int64) (ln, col int) {
	answer := make(chan (int), 2)
	d(viewCurPos{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

func ViewCursorMvmt(viewId int64, mvmt core.CursorMvmt) {
	d(viewCursorMvmt{viewId: viewId, mvmt: mvmt})
}

func ViewDelete(viewId int64, row1, col1, row2, col2 int, undoable bool) {
	d(viewDeleteAction{viewId: viewId, row1: row1, col1: col1, row2: row2, col2: col2, undoable: undoable})
}

func ViewDeleteCur(viewId int64) {
	d(viewDeleteCur{viewId: viewId})
}

func ViewInsert(viewId int64, row, col int, text string, undoable bool) {
	d(viewInsertAction{viewId: viewId, row: row, col: col, text: text, undoable: undoable})
}

func ViewInsertCur(viewId int64, text string) {
	d(viewInsertCur{viewId: viewId, text: text})
}

func ViewInsertNewLine(viewId int64) {
	d(viewInsertNewLine{viewId: viewId})
}

func ViewMoveCursor(viewId int64, y, x int) {
	d(viewMoveCursor{viewId: viewId, x: x, y: y})
}

func ViewMoveCursorRoll(viewId int64, y, x int) {
	d(viewMoveCursor{viewId: viewId, x: x, y: y, roll: true})
}

func ViewPaste(viewId int64) {
	d(viewPaste{viewId: viewId})
}

func ViewOpenSelection(viewId int64, newView bool) {
	d(viewOpenSelection{viewId: viewId, newView: newView})
}

func ViewRedo(viewId int64) {
	d(viewRedo{viewId: viewId})
}

func ViewReload(viewId int64) {
	d(viewReload{viewId: viewId})
}

func ViewRender(viewId int64) {
	d(viewRender{viewId: viewId})
}

func ViewSave(viewId int64) {
	d(viewSave{viewId: viewId})
}

func ViewSetDirty(viewId int64, on bool) {
	d(viewSetDirty{viewId: viewId, on: on})
}

func ViewSetTitle(viewId int64, title string) {
	d(viewSetTitle{viewId: viewId, title: title})
}

func ViewStretchSelection(viewId int64, prevLn, prevCol int) {
	d(viewStretchSelection{viewId: viewId, prevLn: prevLn, prevCol: prevCol})
}

func ViewSetWorkdir(viewId int64, workDir string) {
	d(viewSetWorkdir{viewId: viewId, workDir: workDir})
}

func ViewTrim(viewId int64, limit int) {
	d(viewTrim{viewId: viewId, limit: limit})
}

func ViewUndo(viewId int64) {
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
	answer chan (int)
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

type viewTrim struct {
	viewId int64
	limit  int
}

func (a viewTrim) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	if v.LineCount() > a.limit {
		v.Backend().Remove(1, 1, v.LineCount()-a.limit+1, 0)
	}
	return nil
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
