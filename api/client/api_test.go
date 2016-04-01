package client

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/api"
	"github.com/tcolar/goed/core"
	"github.com/tcolar/goed/ui"
)

var id int64
var ed core.Editable

func init() {
	id = time.Now().Unix()
	core.Testing = true
	core.InitHome(id)
	core.Ed = ui.NewMockEditor()
	ed = core.Ed
	core.Bus = actions.NewActionBus()
	actions.RegisterActions()
	apiServer := api.Api{}
	apiServer.Start()
	go core.Bus.Start()
	core.Ed.Start([]string{"../test_data/file1.txt"})
}

func TestEdit(t *testing.T) {
	done := false
	completed := make(chan struct{})
	go func() {
		err := Edit(id, "test_data", "fooedit")
		done = true
		assert.Nil(t, err)
		close(completed)
	}()
	vid := int64(-1)
	// view should open up and stay open until the view is closed
	// at which time the open action should be completed
	loc, _ := filepath.Abs("./test_data/fooedit")
	for vid == -1 {
		vid = ed.ViewByLoc(loc)
		time.Sleep(100 * time.Millisecond)
	}
	assert.False(t, done)
	ed.DelView(vid, true)
	select {
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout waiting for edit to complete.")
	case <-completed: // good
	}
}

func TestNoSuchAction(t *testing.T) {
	res, err := Action(id, []string{"foobar"})
	assert.NotNil(t, err)
	assert.Equal(t, len(res), 0)
}

func TestCmdBarEnable(t *testing.T) {
	res, err := Action(id, []string{"foobar"})
	assert.NotNil(t, err)

	res, err = Action(id, []string{"cmdbar_enable", "true"})
	assert.Nil(t, err)
	assert.True(t, ed.CmdOn())

	res, err = Action(id, []string{"cmdbar_enable", "false"})
	assert.Nil(t, err)
	assert.False(t, ed.CmdOn())

	assert.Equal(t, len(res), 0)
}

func TestCmdBarToggle(t *testing.T) {
	res, err := Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.True(t, ed.CmdOn())

	res, err = Action(id, []string{"cmdbar_toggle"})
	assert.Nil(t, err)
	assert.False(t, ed.CmdOn())

	assert.Equal(t, len(res), 0)
}

func TestOpen(t *testing.T) {
	err := Open(id, "test_data", "empty.txt")
	assert.Nil(t, err)
	loc, _ := filepath.Abs("./test_data/empty.txt")
	vid := ed.ViewByLoc(loc)
	assert.NotEqual(t, vid, "-1")
	ed.DelView(vid, true)
}

func TestEdActivateView(t *testing.T) {
	views := ed.Views()
	assert.True(t, len(views) >= 2)
	res, err := Action(id, []string{"ed_activate_view", fmt.Sprintf("%d", views[1])})
	assert.Nil(t, err)

	assert.Equal(t, ed.CurViewId(), views[1])
	res, err = Action(id, []string{"ed_activate_view", fmt.Sprintf("%d", views[0])})
	assert.Nil(t, err)
	assert.Equal(t, ed.CurViewId(), views[0])

	assert.Equal(t, len(res), 0)
}

func TestEdCurView(t *testing.T) {
	res, err := Action(id, []string{"ed_cur_view"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, fmt.Sprintf("%d", ed.CurViewId()), res[0])
}

func TestEdDelCol(t *testing.T) {
	col0 := ed.CurColIndex()
	err := Open(id, "test_data", "delcol.txt")
	col1 := ed.CurColIndex()
	assert.NotEqual(t, col0, col1)
	res, err := Action(id, []string{"ed_del_col", strconv.Itoa(col1), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, ed.CurColIndex(), col0)
}

func TestEdDelView(t *testing.T) {
	v0 := ed.CurViewId()
	err := Open(id, "test_data", "delview.txt")
	v1 := ed.CurViewId()
	assert.NotEqual(t, v0, v1)
	ed.ViewById(v1).SetDirty(true)
	// view is dirty so frst try should do nothing
	res, err := Action(id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, ed.CurViewId(), v1)
	// but asking again, should force close v1 and send us back to v0
	res, err = Action(id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, ed.CurViewId(), v0)
}

func TestEdOpen(t *testing.T) {
	// open in a new view
	res, err := Action(id, []string{"ed_open", "theme.toml", "-1", "test_data", "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	vid := ed.CurViewId()
	assert.Equal(t, fmt.Sprintf("%d", vid), res[0])
	assert.Equal(t, ed.ViewById(vid).Title(), "theme.toml")

	// replace the view
	prevId := res[0]
	res, err = Action(id, []string{"ed_open", "testopen.txt", prevId, "test_data", "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	vid = ed.CurViewId()
	assert.Equal(t, fmt.Sprintf("%d", vid), res[0])
	assert.Equal(t, prevId, res[0])
	assert.Equal(t, ed.ViewById(vid).Title(), "testopen.txt")

	loc, _ := filepath.Abs("./test_data/theme.toml") // should no longer be found
	assert.Equal(t, ed.ViewByLoc(loc), int64(-1))

	ed.DelView(ed.CurViewId(), true)
}

func TestEdQuitCheck(t *testing.T) {
	v := ed.ViewById(ed.Views()[0])
	v.SetDirty(true)
	res, err := Action(id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "false", res[0])

	v.SetDirty(false)
	res, err = Action(id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "true", res[0])
}

func TestEdResize(t *testing.T) {
	r, c := ed.Size()
	res, err := Action(id, []string{"ed_resize", strconv.Itoa(r + 1), strconv.Itoa(c + 1)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	r2, c2 := ed.Size()
	assert.Equal(t, r+1, r2)
	assert.Equal(t, c+1, c2)
	res, err = Action(id, []string{"ed_resize", strconv.Itoa(r), strconv.Itoa(c)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
}

/*
ed_set_status(string)
ed_set_status_err(string)
ed_size
ed_swap_views(int64, int64)
ed_view_at(int, int) int64, int, int
ed_view_by_loc(string) int64
ed_view_move(int64, int, int, int, int)
ed_view_navigate(core.CursorMvmt)
ed_views() string
*/
