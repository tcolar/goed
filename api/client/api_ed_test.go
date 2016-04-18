package client

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
)

func TestEdActivateView(t *testing.T) {
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	res, err := Action(id, []string{"ed_activate_view", vidStr(views[1])})
	assert.Nil(t, err)
	assert.Equal(t, actions.Ar.EdCurView(), views[1])
	res, err = Action(id, []string{"ed_activate_view", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Equal(t, actions.Ar.EdCurView(), views[0])

	assert.Equal(t, len(res), 0)
}

func TestEdCurView(t *testing.T) {
	res, err := Action(id, []string{"ed_cur_view"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, fmt.Sprintf("%d", actions.Ar.EdCurView()), res[0])
}

func TestEdDelCol(t *testing.T) {
	v1 := actions.Ar.EdCurView()
	_, c1 := actions.Ar.EdViewIndex(v1)
	err := Open(id, "test_data", "delcol.txt")
	vid := actions.Ar.EdCurView()
	l, c, _, _ := actions.Ar.ViewBounds(vid)
	actions.Ar.EdViewMove(vid, l, c, 2, c+5) // force view to it's own column, if not already
	_, col := actions.Ar.EdViewIndex(vid)
	assert.NotEqual(t, col, c1)
	res, err := Action(id, []string{"ed_del_col", strconv.Itoa(col), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	vid = actions.Ar.EdCurView()
	_, c2 := actions.Ar.EdViewIndex(vid)
	assert.NotEqual(t, col, c2)
}

func TestEdDelView(t *testing.T) {
	v0 := actions.Ar.EdCurView()
	err := Open(id, "test_data", "delview.txt")
	v1 := actions.Ar.EdCurView()
	assert.NotEqual(t, v0, v1)
	actions.Ar.ViewSetDirty(v1, true)
	// view is dirty so first try should do nothing
	res, err := Action(id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, actions.Ar.EdCurView(), v1)
	// but asking again, should force close v1 and send us back to v0
	res, err = Action(id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	assert.Equal(t, actions.Ar.EdCurView(), v0)
}

func TestEdOpen(t *testing.T) {
	// open in a new view
	res, err := Action(id, []string{"ed_open", "theme.toml", "-1", "test_data", "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	vid := actions.Ar.EdCurView()
	assert.Equal(t, fmt.Sprintf("%d", vid), res[0])
	assert.Equal(t, actions.Ar.ViewTitle(vid), "theme.toml")

	// replace the view
	prevId := res[0]
	res, err = Action(id, []string{"ed_open", "testopen.txt", prevId, "test_data", "true"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	vid = actions.Ar.EdCurView()
	assert.Equal(t, fmt.Sprintf("%d", vid), res[0])
	assert.Equal(t, prevId, res[0])
	assert.Equal(t, actions.Ar.ViewTitle(vid), "testopen.txt")

	loc, _ := filepath.Abs("../test_data/theme.toml") // should no longer be found
	assert.Equal(t, actions.Ar.EdViewByLoc(loc), int64(-1))

	actions.Ar.EdDelView(actions.Ar.EdCurView(), true)
}

func TestEdQuitCheck(t *testing.T) {
	views := actions.Ar.EdViews()
	actions.Ar.ViewSetDirty(views[0], true)
	res, err := Action(id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "false", res[0])

	actions.Ar.ViewSetDirty(views[0], false)
	res, err = Action(id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)
	assert.Equal(t, "true", res[0])
}

func TestEdResize(t *testing.T) {
	r, c := actions.Ar.EdSize()
	res, err := Action(id, []string{"ed_resize", strconv.Itoa(r + 1), strconv.Itoa(c + 1)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
	r2, c2 := actions.Ar.EdSize()
	assert.Equal(t, r+1, r2)
	assert.Equal(t, c+1, c2)
	res, err = Action(id, []string{"ed_resize", strconv.Itoa(r), strconv.Itoa(c)})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
}

func TestEdSize(t *testing.T) {
	r, c := actions.Ar.EdSize()
	res, err := Action(id, []string{"ed_size"})
	assert.Nil(t, err)
	assert.Equal(t, len(res), 2)
	assert.Equal(t, strconv.Itoa(r), res[0])
	assert.Equal(t, strconv.Itoa(c), res[1])
}

func TestEdSwapViews(t *testing.T) {
	boundStr := func(y1, x1, y2, x2 int) string {
		return fmt.Sprintf("%d %d %d %d", y1, x1, y2, x2)
	}
	err := Open(id, "test_data", "swapview.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	vid := views[0]
	v2id := views[1]
	b1 := boundStr(actions.Ar.ViewBounds(vid))
	b2 := boundStr(actions.Ar.ViewBounds(v2id))
	res, err := Action(id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	assert.Nil(t, err)
	assert.Equal(t, b1, boundStr(actions.Ar.ViewBounds(v2id)))
	assert.Equal(t, b2, boundStr(actions.Ar.ViewBounds(vid)))
	res, err = Action(id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	assert.Nil(t, err)
	assert.Equal(t, b1, boundStr(actions.Ar.ViewBounds(vid)))
	assert.Equal(t, b2, boundStr(actions.Ar.ViewBounds(v2id)))
	assert.Equal(t, len(res), 0)
}

func TestEdViewAt(t *testing.T) {
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	for _, v := range views {
		l, c, l2, c2 := actions.Ar.ViewBounds(v)
		res, err := Action(id, []string{"ed_view_at", fmt.Sprintf("%d", l), fmt.Sprintf("%d", c)})
		assert.Nil(t, err)
		assert.Equal(t, len(res), 3)
		assert.Equal(t, res[0], vidStr(v))
		res, err = Action(id, []string{"ed_view_at", fmt.Sprintf("%d", l2), fmt.Sprintf("%d", c2)})
		assert.Nil(t, err)
		assert.Equal(t, len(res), 3)
		assert.Equal(t, res[0], vidStr(v))
	}
}

func TestEdViewNavigate(t *testing.T) {
	err := Open(id, "test_data", "viewnav.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	vid := views[0]
	res, err := Action(id, []string{"ed_activate_view", vidStr(vid)})
	assert.Nil(t, err)
	actions.Ar.EdActivateView(vid)
	res, err = Action(id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtRight))})
	assert.Nil(t, err)
	assert.NotEqual(t, vid, actions.Ar.EdCurView())
	res, err = Action(id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtLeft))})
	assert.Nil(t, err)
	assert.Equal(t, vid, actions.Ar.EdCurView())
	assert.Equal(t, len(res), 0)
}

func TestEdViews(t *testing.T) {
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	res, err := Action(id, []string{"ed_views"})
	assert.Nil(t, err)
	assert.Equal(t, len(views), len(res))
	vs := []string{}
	for _, v := range views {
		vs = append(vs, vidStr(v))
	}
	assert.Equal(t, vs, res)
}
