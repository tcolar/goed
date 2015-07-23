/*
Set of actions that can be dispatched.
Actions are dispatched, and processed one at a time by the action bus for
concurency safety.
*/
package actions

import "github.com/tcolar/goed/core"

func d(action core.Action) {
	core.Bus.Dispatch(action)
}

func CmdbarRunAction() {
	d(cmdbarRunAction{})
}

func CmdbarToggleAction() {
	d(cmdbarToggleAction{})
}

func EdActivateViewAction(viewId int64, y, x int) {
	d(edActivateViewAction{viewId: viewId, y: y, x: x})
}

func EdDelViewCheckAction(viewId int64) {
	d(edDelViewCheckAction{viewId: viewId})
}

func EdRenderAction() {
	d(edRenderAction{})
}

func EdResizeAction(h, w int) {
	d(edResizeAction{h: h, w: w})
}

func EdSetStatusAction(status string) {
	d(edSetStatusAction{status: status, err: false})
}

func EdSetStatusErrAction(status string) {
	d(edSetStatusAction{status: status, err: true})
}

func EdViewMoveAction(viewId int64, y1, x1, y2, x2 int) {
	d(edViewMoveAction{viewId: viewId, y1: y1, x1: x1, y2: y2, x2: x2})
}

func EdTermFlushAction() {
	d(edTermFlushAction{})
}

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

// **** impls

type edResizeAction struct {
	h, w int
}

func (e edResizeAction) Run() error {
	core.Ed.Resize(e.h, e.w)
	return nil
}

type viewReloadAction struct{ viewId int64 }

func (e viewReloadAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v != nil {
		v.Reload()
	}
	return nil
}

type viewRenderAction struct {
	viewId int64
}

func (e viewRenderAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v != nil {
		v.Render()
	}
	return nil
}

type viewSetWorkdirAction struct {
	viewId  int64
	workDir string
}

func (e viewSetWorkdirAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v != nil {
		v.SetWorkDir(e.workDir)
	}
	return nil
}

type viewSetTitleAction struct {
	viewId int64
	title  string
}

func (e viewSetTitleAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v != nil {
		v.SetWorkDir(e.title)
	}
	return nil
}

type edSetStatusAction struct {
	status string
	err    bool
}

func (e edSetStatusAction) Run() error {
	if e.err {
		core.Ed.SetStatusErr(e.status)
	} else {
		core.Ed.SetStatus(e.status)
	}
	return nil
}

type edRenderAction struct{}

func (e edRenderAction) Run() error {
	core.Ed.Render()
	return nil
}

type edTermFlushAction struct{}

func (e edTermFlushAction) Run() error {
	core.Ed.TermFlush()
	return nil
}

type viewMoveCursorAction struct {
	viewId int64
	status string
	y, x   int
	roll   bool
}

func (e viewMoveCursorAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	if e.roll {
		v.MoveCursorRoll(e.y, e.x)
	} else {
		v.MoveCursor(e.y, e.x)
	}
	return nil
}

type cmdbarToggleAction struct{}

func (e cmdbarToggleAction) Run() error {
	core.Ed.CmdbarToggle()
	return nil
}

type viewOpenSelectionAction struct {
	viewId  int64
	newView bool
}

func (e viewOpenSelectionAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.OpenSelection(e.newView)
	return nil
}

type viewCursorMvmtAction struct {
	viewId int64
	mvmt   core.CursorMvmt
}

func (e viewCursorMvmtAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.CursorMvmt(e.mvmt)
	return nil
}

type viewAutoScrollAction struct {
	viewId int64
	y, x   int
	on     bool
}

func (e viewAutoScrollAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.SetAutoScroll(e.y, e.x, e.on)
	return nil
}

type viewInsertCurAction struct {
	viewId int64
	text   string
}

func (e viewInsertCurAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.InsertCur(e.text)
	return nil
}

type viewInsertNewLineAction struct {
	viewId int64
}

func (e viewInsertNewLineAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.InsertNewLineCur()
	return nil
}

type viewDeleteCurAction struct {
	viewId int64
}

func (e viewDeleteCurAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.DeleteCur()
	return nil
}

type viewBackspaceAction struct {
	viewId int64
}

func (e viewBackspaceAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.Backspace()
	return nil
}

type viewSaveAction struct {
	viewId int64
}

func (e viewSaveAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.Save()
	return nil
}

type edDelViewCheckAction struct {
	viewId int64
}

func (e edDelViewCheckAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	core.Ed.DelViewCheck(v)
	return nil
}

type viewCmdStopAction struct {
	viewId int64
}

func (e viewCmdStopAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
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

func (e viewCopyAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.Copy()
	return nil
}

type cmdbarRunAction struct {
}

func (e cmdbarRunAction) Run() error {
	core.Ed.CmdbarRun()
	return nil
}

type edActivateViewAction struct {
	viewId int64
	y, x   int
}

func (e edActivateViewAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ActivateView(v, e.y, e.x)
	return nil
}

type edViewMoveAction struct {
	viewId         int64
	y1, x1, y2, x2 int
}

func (e edViewMoveAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ViewMove(e.y1, e.x1, e.y2, e.x2)
	return nil
}

type viewPasteAction struct {
	viewId int64
}

func (e viewPasteAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.Paste()
	return nil
}

type viewCutAction struct {
	viewId int64
}

func (e viewCutAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
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

func (e viewClearSelectionsAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
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

func (e viewSetDirtyAction) Run() error {
	v := core.Ed.ViewById(e.viewId)
	if v == nil {
		return nil
	}
	v.SetDirty(e.on)
	return nil
}
