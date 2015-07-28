package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/tcolar/goed/core"
)

func EdActivateView(viewId int64, y, x int) {
	d(edActivateView{viewId: viewId, y: y, x: x})
}

func EdDelColCheck(colIndex int) {
	d(edDelColCheck{colIndex: colIndex})
}

func EdDelViewCheck(viewId int64) {
	d(edDelViewCheck{viewId: viewId})
}

func EdExternal(name string) {
	d(edExternal{name: name})
}

func EdOpen(loc string, view core.Viewable, rel string, create bool) {
	d(edOpen{loc: loc, view: view, rel: rel, create: create})
}

// Retuns whether the editor can be quit.
func EdQuitCheck() bool {
	answer := make(chan (bool), 1)
	d(edQuitCheck{answer: answer})
	return <-answer
}

func EdRender() {
	d(edRender{})
}

func EdResize(h, w int) {
	d(edResize{h: h, w: w})
}

func EdSetStatus(status string) {
	d(edSetStatus{status: status, err: false})
}

func EdSetStatusErr(status string) {
	d(edSetStatus{status: status, err: true})
}

func EdSwapViews(view1Id, view2Id int64) {
	d(edSwapViews{view1Id: view1Id, view2Id: view2Id})
}

func EdTermFlush() {
	d(edTermFlush{})
}

func EdViewMove(viewId int64, y1, x1, y2, x2 int) {
	d(edViewMove{viewId: viewId, y1: y1, x1: x1, y2: y2, x2: x2})
}

// ########  Impl ......

type edActivateView struct {
	viewId int64
	y, x   int
}

func (a edActivateView) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ActivateView(v, a.y, a.x)
	return nil
}

type edDelColCheck struct {
	colIndex int
}

func (a edDelColCheck) Run() error {
	core.Ed.DelColCheckByIndex(a.colIndex)
	return nil
}

type edDelViewCheck struct {
	viewId int64
}

func (a edDelViewCheck) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.DelViewCheck(v)
	return nil
}

func (a edExternal) Run() error {
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
		errv, err = e.Open(fp, errv, "", true)
		if err != nil {
			e.SetStatusErr(err.Error())
		}
		e.Render()
		return fmt.Errorf("%s failed", a.name)
	}
	errv := e.ViewByLoc(fp)
	if errv != nil {
		e.DelView(errv, true)
	}
	return nil
}

type edOpen struct {
	loc, rel string
	view     core.Viewable
	create   bool
}

func (a edOpen) Run() error {
	_, err := core.Ed.Open(a.loc, a.view, a.rel, a.create)
	return err
}

type edQuitCheck struct {
	answer chan (bool)
}

func (a edQuitCheck) Run() error {
	a.answer <- core.Ed.QuitCheck()
	return nil
}

type edRender struct{}

func (a edRender) Run() error {
	core.Ed.Render()
	return nil
}

type edResize struct {
	h, w int
}

func (a edResize) Run() error {
	core.Ed.Resize(a.h, a.w)
	return nil
}

type edSetStatus struct {
	status string
	err    bool
}

func (a edSetStatus) Run() error {
	if a.err {
		core.Ed.SetStatusErr(a.status)
	} else {
		core.Ed.SetStatus(a.status)
	}
	return nil
}

type edSwapViews struct {
	view1Id, view2Id int64
}

func (a edSwapViews) Run() error {
	v := core.Ed.ViewById(a.view1Id)
	v2 := core.Ed.ViewById(a.view2Id)
	if v == nil || v2 == nil {
		return nil
	}
	core.Ed.SwapViews(v, v2)
	return nil
}

type edTermFlush struct{}

func (a edTermFlush) Run() error {
	core.Ed.TermFlush()
	return nil
}

type edViewMove struct {
	viewId         int64
	y1, x1, y2, x2 int
}

func (a edViewMove) Run() error {
	v := core.Ed.ViewById(a.viewId)
	if v == nil {
		return nil
	}
	core.Ed.ViewMove(a.y1, a.x1, a.y2, a.x2)
	return nil
}

type edExternal struct {
	name string
}

/*
type viewSearch struct {
	viewId int64
}

func (a viewSearch) Run() error {
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
