package actions

import (
	"log"

	"github.com/tcolar/goed/core"
)

// activate the given view, with the cursor at y,x
func (a *ar) EdActivateView(viewId int64, y, x int) {
	d(edActivateView{viewId: viewId, y: y, x: x})
}

// returns the currently active view
func (a *ar) EdCurView() int64 {
	vid := make(chan (int64), 1)
	d(edCurView{viewId: vid})
	return <-vid
}

// delete the given column (by index)
// if 'check' is true it will check if dirty first, in which case it will do nothing
// unless called twice in a row.
func (a *ar) EdDelCol(colIndex int, check bool) {
	d(edDelCol{colIndex: colIndex, check: check})
}

// delete the given view (by id)
// if 'check' is true it will check if dirty first, in which case it will do nothing
// unless called twice in a row.
func (a *ar) EdDelView(viewId int64, check bool) {
	d(edDelView{viewId: viewId, check: check})
}

// Open a file/dir(loc) in the editor
// rel is optionally the path to loc
// viewId is the viewId where to open into (or a new one if viewId<0)
// create indicates whether the file/dir needs to be created if it does not exist.
func (a *ar) EdOpen(loc string, viewId int64, rel string, create bool) int64 {
	vid := make(chan (int64), 1)
	d(edOpen{loc: loc, viewId: viewId, rel: rel, create: create, vid: vid})
	return <-vid
}

// Retuns whether the editor can be quit (ie: are any views "dirty")
func (a *ar) EdQuitCheck() bool {
	answer := make(chan (bool), 1)
	d(edQuitCheck{answer: answer})
	return <-answer
}

// Render/repaint the editor UI
func (a *ar) EdRender() {
	d(edRender{})
}

// resize the editor
func (a *ar) EdResize(h, w int) {
	d(edResize{h: h, w: w})
}

// Show a status message in the satus bar
func (a *ar) EdSetStatus(status string) {
	d(edSetStatus{status: status, err: false})
}

// Show an error message (red) in the status bar.
func (a *ar) EdSetStatusErr(status string) {
	d(edSetStatus{status: status, err: true})
}

// swap two views (their position in the UI)
func (a *ar) EdSwapViews(view1Id, view2Id int64) {
	d(edSwapViews{view1Id: view1Id, view2Id: view2Id})
}

// call flush on the underlying terminal (force sync)
func (a *ar) EdTermFlush() {
	d(edTermFlush{})
}

// returns the viewId of the view that holds a file/dir of the given path.
// or -1 if not found.
func (a *ar) EdViewByLoc(loc string) int64 {
	vid := make(chan (int64), 1)
	d(edViewByLoc{loc: loc, vid: vid})
	return <-vid
}

// move a view to the new coordinates (UI position)
func (a *ar) EdViewMove(viewId int64, y1, x1, y2, x2 int) {
	d(edViewMove{viewId: viewId, y1: y1, x1: x1, y2: y2, x2: x2})
}

// navigate between UI views given the CursorMvmt value (left,right,top,down)
func (a *ar) EdViewNavigate(mvmt core.CursorMvmt) {
	d(edViewNavigate{mvmt})
}

// ########  Impl ......

type edActivateView struct {
	viewId int64
	y, x   int
}

func (a edActivateView) Run() error {
	core.Ed.ViewActivate(a.viewId, a.y, a.x)
	return nil
}

type edCurView struct {
	viewId chan int64
}

func (a edCurView) Run() error {
	a.viewId <- core.Ed.CurViewId()
	return nil
}

type edDelCol struct {
	colIndex int
	check    bool
}

func (a edDelCol) Run() error {
	core.Ed.DelColByIndex(a.colIndex, a.check)
	return nil
}

type edDelView struct {
	viewId int64
	check  bool
}

func (a edDelView) Run() error {
	core.Ed.DelViewByIndex(a.viewId, a.check)
	return nil
}

type edOpen struct {
	loc, rel string
	viewId   int64
	create   bool
	vid      chan int64 // returned new viewid if viewId==-1
}

func (a edOpen) Run() error {
	vid, err := core.Ed.Open(a.loc, a.viewId, a.rel, a.create)
	a.vid <- vid
	if err != nil {
		log.Printf("EdOpen error : %s\n", err.Error())
	}
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
	core.Ed.SwapViews(a.view1Id, a.view2Id)
	return nil
}

type edTermFlush struct{}

func (a edTermFlush) Run() error {
	core.Ed.TermFlush()
	return nil
}

type edViewByLoc struct {
	loc string
	vid chan int64
}

func (a edViewByLoc) Run() error {
	vid := core.Ed.ViewByLoc(a.loc)
	a.vid <- vid
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

type edViewNavigate struct {
	mvmt core.CursorMvmt
}

func (a edViewNavigate) Run() error {
	core.Ed.ViewNavigate(a.mvmt)
	return nil
}
