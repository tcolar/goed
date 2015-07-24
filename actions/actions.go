/*
Set of actions that can be dispatched.
Actions are dispatched, and processed one at a time by the action bus for
concurency safety.
*/
package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/tcolar/goed/core"
)

func d(action core.Action) {
	core.Bus.Dispatch(action)
}

func CmdbarEnableAction(on bool) {
	d(cmdbarEnableAction{on: on})
}

func CmdbarToggleAction() {
	d(cmdbarToggleAction{})
}

func EdActivateViewAction(viewId int64, y, x int) {
	d(edActivateViewAction{viewId: viewId, y: y, x: x})
}

func EdDelColCheckAction(colIndex int) {
	d(edDelColCheckAction{colIndex: colIndex})
}

func EdDelViewCheckAction(viewId int64) {
	d(edDelViewCheckAction{viewId: viewId})
}

func EdExternalAction(name string) {
	d(edExternalAction{name: name})
}

func EdOpenAction(loc string, view core.Viewable, rel string) {
	d(edOpenAction{loc: loc, view: view, rel: rel})
}

// Retuns whether the editor can be quit.
func EdQuitCheck() bool {
	answer := make(chan (bool), 1)
	d(edQuitCheckAction{answer: answer})
	return <-answer
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

func EdSwapViewsAction(view1Id, view2Id int64) {
	d(edSwapViewsAction{view1Id: view1Id, view2Id: view2Id})
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

// **** impls

type edResizeAction struct {
	h, w int
}

func (a edResizeAction) Run() error {
	core.Ed.Resize(a.h, a.w)
	return nil
}

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

type edSetStatusAction struct {
	status string
	err    bool
}

func (a edSetStatusAction) Run() error {
	if a.err {
		core.Ed.SetStatusErr(a.status)
	} else {
		core.Ed.SetStatus(a.status)
	}
	return nil
}

type edRenderAction struct{}

func (a edRenderAction) Run() error {
	core.Ed.Render()
	return nil
}

type edTermFlushAction struct{}

func (a edTermFlushAction) Run() error {
	core.Ed.TermFlush()
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

type cmdbarToggleAction struct{}

func (a cmdbarToggleAction) Run() error {
	core.Ed.CmdbarToggle()
	return nil
}

type cmdbarEnableAction struct {
	on bool
}

func (a cmdbarEnableAction) Run() error {
	core.Ed.SetCmdOn(a.on)
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

type edDelViewCheckAction struct {
	viewId int64
}

func (a edDelViewCheckAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.DelViewCheck(v)
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

type edActivateViewAction struct {
	viewId int64
	y, x   int
}

func (a edActivateViewAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ActivateView(v, a.y, a.x)
	return nil
}

type edViewMoveAction struct {
	viewId         int64
	y1, x1, y2, x2 int
}

func (a edViewMoveAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ViewMove(a.y1, a.x1, a.y2, a.x2)
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

type edDelColCheckAction struct {
	colIndex int
}

func (a edDelColCheckAction) Run() error {
	core.Ed.DelColCheckByIndex(a.colIndex)
	return nil
}

type edSwapViewsAction struct {
	view1Id, view2Id int64
}

func (a edSwapViewsAction) Run() error {
	v := core.Ed.ViewById(a.view1Id)
	v2 := core.Ed.ViewById(a.view2Id)
	if v == nil || v2 == nil {
		return nil
	}
	core.Ed.SwapViews(v, v2)
	return nil
}

type edOpenAction struct {
	loc, rel string
	view     core.Viewable
}

func (a edOpenAction) Run() error {
	core.Ed.Open(a.loc, a.view, a.rel)
	return nil
}

type edQuitCheckAction struct {
	answer chan (bool)
}

func (a edQuitCheckAction) Run() error {
	a.answer <- core.Ed.QuitCheck()
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

type edExternalAction struct {
	name string
}

func (a edExternalAction) Run() error {
	e := core.Ed
	v := e.CurView()
	loc := core.FindResource(path.Join("actions", a.name))
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		return fmt.Errorf("Action not found : %s", a.name)
	}
	env := os.Environ()
	env = append(env, fmt.Sprintf("GOED_INSTANCE=%d", core.InstanceId))
	env = append(env, fmt.Sprintf("GOED_VIEW=%d", v.Id()))
	cmd := exec.Command(loc)
	cmd.Env = env
	out, err := cmd.CombinedOutput()
	fp := path.Join(core.Home, "errors.txt")
	if err != nil {
		file, _ := os.Create(fp)
		file.Write([]byte(err.Error()))
		file.Write([]byte{'\n'})
		file.Write(out)
		file.Close()
		errv := e.ViewByLoc(fp)
		errv, err = e.Open(fp, errv, "Errors")
		if err != nil {
			e.SetStatusErr(err.Error())
		}
		return fmt.Errorf("%s failed", a.name)
	}
	errv := e.ViewByLoc(fp)
	if errv != nil {
		e.DelView(errv, true)
	}
	return nil
}

/*
type viewSearchAction struct {
	viewId int64
}

func (a viewSearchAction) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	if len(v.Selections()) == 0 {
		v.SelectWord(ln, col)
	}
	if len(v.Selections()) > 0 {
		text := core.RunesToString(v.SelectionText(&v.selections[0]))
		e.Cmdbar.Search(text)
	}
}*/
