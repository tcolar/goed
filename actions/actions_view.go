package actions

import "github.com/tcolar/goed/core"

func ViewAutoScrollAction(viewId int64, y, x int, on bool) {
	d(viewAutoScrollAction{viewId: viewId, x: x, y: y, on: on})
}

func ViewBackspaceAction(viewId int64) {
	d(viewBackspaceAction{viewId: viewId})
}

func ViewClearSelectionsAction(viewId int64) {
	d(viewClearSelectionsAction{viewId: viewId})
}

func ViewCmdStopAction(viewId int64) {
	d(viewCmdStopAction{viewId: viewId})
}

func ViewCopyAction(viewId int64) {
	d(viewCopyAction{viewId: viewId})
}

func ViewCutAction(viewId int64) {
	d(viewCutAction{viewId: viewId})
}

func ViewCurPosAction(viewId int64) (ln, col int) {
	answer := make(chan (int), 2)
	d(viewCurPosAction{viewId: viewId, answer: answer})
	return <-answer, <-answer
}

func ViewCursorMvmtAction(viewId int64, mvmt core.CursorMvmt) {
	d(viewCursorMvmtAction{viewId: viewId, mvmt: mvmt})
}

func ViewDeleteCurAction(viewId int64) {
	d(viewDeleteCurAction{viewId: viewId})
}

func ViewInsertCurAction(viewId int64, text string) {
	d(viewInsertCurAction{viewId: viewId, text: text})
}

func ViewInsertNewLineAction(viewId int64) {
	d(viewInsertNewLineAction{viewId: viewId})
}

func ViewMoveCursorAction(viewId int64, y, x int) {
	d(viewMoveCursorAction{viewId: viewId, x: x, y: y})
}

func ViewMoveCursorRollAction(viewId int64, y, x int) {
	d(viewMoveCursorAction{viewId: viewId, x: x, y: y, roll: true})
}

func ViewPasteAction(viewId int64) {
	d(viewPasteAction{viewId: viewId})
}

func ViewOpenSelectionAction(viewId int64, newView bool) {
	d(viewOpenSelectionAction{viewId: viewId, newView: newView})
}
func ViewReloadAction(viewId int64) {
	d(viewReloadAction{viewId: viewId})
}

func ViewRenderAction(viewId int64) {
	d(viewRenderAction{viewId: viewId})
}
func ViewSaveAction(viewId int64) {
	d(viewSaveAction{viewId: viewId})
}

func ViewSetDirtyAction(viewId int64, on bool) {
	d(viewSetDirtyAction{viewId: viewId, on: on})
}
func ViewSetWorkdirAction(viewId int64, workDir string) {
	d(viewSetWorkdirAction{viewId: viewId, workDir: workDir})
}

func ViewSetTitleAction(viewId int64, title string) {
	d(viewSetTitleAction{viewId: viewId, title: title})
}

func ViewStretchSelection(viewId int64, prevLn, prevCol int) {
	d(viewStretchSelection{viewId: viewId, prevLn: prevLn, prevCol: prevCol})
}

// ########  Impl ......

type viewReloadAction struct{ viewId int64 }

func (a viewReloadAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Reload()
	}
	return nil
}

type viewRenderAction struct {
	viewId int64
}

func (a viewRenderAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.Render()
	}
	return nil
}

type viewSetWorkdirAction struct {
	viewId  int64
	workDir string
}

func (a viewSetWorkdirAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetWorkDir(a.workDir)
	}
	return nil
}

type viewSetTitleAction struct {
	viewId int64
	title  string
}

func (a viewSetTitleAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v != nil {
		v.SetWorkDir(a.title)
	}
	return nil
}

type viewMoveCursorAction struct {
	viewId int64
	status string
	y, x   int
	roll   bool
}

func (a viewMoveCursorAction) Run() error {
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

type viewOpenSelectionAction struct {
	viewId  int64
	newView bool
}

func (a viewOpenSelectionAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.OpenSelection(a.newView)
	return nil
}

type viewCursorMvmtAction struct {
	viewId int64
	mvmt   core.CursorMvmt
}

func (a viewCursorMvmtAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.CursorMvmt(a.mvmt)
	return nil
}

type viewAutoScrollAction struct {
	viewId int64
	y, x   int
	on     bool
}

func (a viewAutoScrollAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SetAutoScroll(a.y, a.x, a.on)
	return nil
}

type viewInsertCurAction struct {
	viewId int64
	text   string
}

func (a viewInsertCurAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.InsertCur(a.text)
	return nil
}

type viewInsertNewLineAction struct {
	viewId int64
}

func (a viewInsertNewLineAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.InsertNewLineCur()
	return nil
}

type viewDeleteCurAction struct {
	viewId int64
}

func (a viewDeleteCurAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.DeleteCur()
	return nil
}

type viewBackspaceAction struct {
	viewId int64
}

func (a viewBackspaceAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Backspace()
	return nil
}

type viewSaveAction struct {
	viewId int64
}

func (a viewSaveAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Save()
	return nil
}

type viewCmdStopAction struct {
	viewId int64
}

func (a viewCmdStopAction) Run() error {
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

type viewCopyAction struct {
	viewId int64
}

func (a viewCopyAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Copy()
	return nil
}

type viewPasteAction struct {
	viewId int64
}

func (a viewPasteAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Paste()
	return nil
}

type viewCutAction struct {
	viewId int64
}

func (a viewCutAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.Copy()
	v.Delete()
	return nil
}

type viewClearSelectionsAction struct {
	viewId int64
}

func (a viewClearSelectionsAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.ClearSelections()
	return nil
}

type viewSetDirtyAction struct {
	viewId int64
	on     bool
}

func (a viewSetDirtyAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	v.SetDirty(a.on)
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

type viewCurPosAction struct {
	answer chan (int)
	viewId int64
}

func (a viewCurPosAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		a.answer <- 0
		a.answer <- 0
	}
	a.answer <- v.CurLine()
	a.answer <- v.CurCol()
	return nil
}
