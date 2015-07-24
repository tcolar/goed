package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/tcolar/goed/core"
)

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

// ########  Impl ......

type edResizeAction struct {
	h, w int
}

func (a edResizeAction) Run() error {
	core.Ed.Resize(a.h, a.w)
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

type edQuitCheckAction struct {
	answer chan (bool)
}

func (a edQuitCheckAction) Run() error {
	a.answer <- core.Ed.QuitCheck()
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
