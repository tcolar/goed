package client

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/assert"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestEdActivateView(t *C) {
	err := Open(s.id, "test_data", "swapview.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	res, err := Action(s.id, []string{"ed_activate_view", vidStr(views[1])})
	assert.Nil(t, err)
	assert.Eq(t, actions.Ar.EdCurView(), views[1])
	res, err = Action(s.id, []string{"ed_activate_view", vidStr(views[0])})
	assert.Nil(t, err)
	assert.Eq(t, actions.Ar.EdCurView(), views[0])
	assert.Eq(t, len(res), 0)
}

func (s *ApiSuite) TestEdCurView(t *C) {
	res, err := Action(s.id, []string{"ed_cur_view"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, fmt.Sprintf("%d", actions.Ar.EdCurView()), res[0])
}

func (s *ApiSuite) TestEdDelCol(t *C) {
	v1 := actions.Ar.EdCurView()
	_, c1 := actions.Ar.EdViewIndex(v1)
	err := Open(s.id, "test_data", "delcol.txt")
	vid := actions.Ar.EdCurView()
	l, c2, _, _ := actions.Ar.ViewBounds(vid)
	actions.Ar.EdViewMove(vid, l, c2, 2, c2+5) // force view to it's own column
	_, col := actions.Ar.EdViewIndex(vid)
	assert.NotEq(t, col, c1)
	res, err := Action(s.id, []string{"ed_del_col", strconv.Itoa(col), "true"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	vid = actions.Ar.EdCurView()
	_, c3 := actions.Ar.EdViewIndex(vid)
	assert.NotEq(t, c3, col)
}

func (s *ApiSuite) TestEdDelView(t *C) {
	v0 := actions.Ar.EdCurView()
	err := Open(s.id, "test_data", "delview.txt")
	v1 := actions.Ar.EdCurView()
	assert.NotEq(t, v0, v1)
	actions.Ar.ViewSetDirty(v1, true)
	// view is dirty so first try should do nothing
	res, err := Action(s.id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.EdCurView(), v1)
	// but asking again, should force close v1 and send us back to v0
	res, err = Action(s.id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	assert.Eq(t, actions.Ar.EdCurView(), v0)
}

func (s *ApiSuite) TestEdOpen(t *C) {
	// open in a new view
	res, err := Action(s.id, []string{"ed_open", "theme.toml", "-1", "test_data", "true"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	vid := actions.Ar.EdCurView()
	assert.Eq(t, fmt.Sprintf("%d", vid), res[0])
	assert.Eq(t, actions.Ar.ViewTitle(vid), "theme.toml")

	// replace the view
	prevId := res[0]
	res, err = Action(s.id, []string{"ed_open", "testopen.txt", prevId, "test_data", "true"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	vid = actions.Ar.EdCurView()
	assert.Eq(t, fmt.Sprintf("%d", vid), res[0])
	assert.Eq(t, prevId, res[0])
	assert.Eq(t, actions.Ar.ViewTitle(vid), "testopen.txt")

	loc, _ := filepath.Abs("../test_data/theme.toml") // should no longer be found
	assert.Eq(t, actions.Ar.EdViewByLoc(loc), int64(-1))

	actions.Ar.EdDelView(actions.Ar.EdCurView(), true)
}

func (s *ApiSuite) TestEdQuitCheck(t *C) {
	views := actions.Ar.EdViews()
	actions.Ar.ViewSetDirty(views[0], true)
	res, err := Action(s.id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, "false", res[0])

	actions.Ar.ViewSetDirty(views[0], false)
	res, err = Action(s.id, []string{"ed_quit_check"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 1)
	assert.Eq(t, "true", res[0])
}

func (s *ApiSuite) TestEdResize(t *C) {
	r, c := actions.Ar.EdSize()
	res, err := Action(s.id, []string{"ed_resize", strconv.Itoa(r + 1), strconv.Itoa(c + 1)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
	r2, c2 := actions.Ar.EdSize()
	assert.Eq(t, r+1, r2)
	assert.Eq(t, c+1, c2)
	res, err = Action(s.id, []string{"ed_resize", strconv.Itoa(r), strconv.Itoa(c)})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 0)
}

func (s *ApiSuite) TestEdSize(t *C) {
	r, c := actions.Ar.EdSize()
	res, err := Action(s.id, []string{"ed_size"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 2)
	assert.Eq(t, strconv.Itoa(r), res[0])
	assert.Eq(t, strconv.Itoa(c), res[1])
}

func (s *ApiSuite) TestEdSwapViews(t *C) {
	boundStr := func(y1, x1, y2, x2 int) string {
		return fmt.Sprintf("%d %d %d %d", y1, x1, y2, x2)
	}
	err := Open(s.id, "test_data", "swapview.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	vid := views[0]
	v2id := views[1]
	b1 := boundStr(actions.Ar.ViewBounds(vid))
	b2 := boundStr(actions.Ar.ViewBounds(v2id))
	res, err := Action(s.id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	assert.Nil(t, err)
	assert.Eq(t, b1, boundStr(actions.Ar.ViewBounds(v2id)))
	assert.Eq(t, b2, boundStr(actions.Ar.ViewBounds(vid)))
	res, err = Action(s.id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	assert.Nil(t, err)
	assert.Eq(t, b1, boundStr(actions.Ar.ViewBounds(vid)))
	assert.Eq(t, b2, boundStr(actions.Ar.ViewBounds(v2id)))
	assert.Eq(t, len(res), 0)
}

func (s *ApiSuite) TestEdViewAt(t *C) {
	err := Open(s.id, "test_data", "viewat.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	res, err := Action(s.id, []string{"ed_view_at", "1", "1"})
	assert.Nil(t, err)
	assert.Eq(t, len(res), 3)
	assert.Eq(t, res[0], "-1")
	assert.Eq(t, res[1], "-1")
	assert.Eq(t, res[2], "-1")
	for _, v := range views {
		l, c, l2, c2 := actions.Ar.ViewBounds(v)
		cs := fmt.Sprintf("%d", c)
		c2s := fmt.Sprintf("%d", c2)
		ls := fmt.Sprintf("%d", l)
		l2s := fmt.Sprintf("%d", l2)
		res, err = Action(s.id, []string{"ed_view_at", ls, cs})
		assert.Nil(t, err)
		assert.Eq(t, len(res), 3)
		assert.Eq(t, res[0], vidStr(v))
		assert.Eq(t, res[1], "1")
		assert.Eq(t, res[2], "1")
		res, err = Action(s.id, []string{"ed_view_at", l2s, c2s})
		assert.Nil(t, err)
		assert.Eq(t, len(res), 3)
		assert.Eq(t, res[0], vidStr(v))
		assert.Eq(t, res[1], fmt.Sprintf("%d", l2-l+1))
		assert.Eq(t, res[2], c2s)
	}
}

func (s *ApiSuite) TestEdViewNavigate(t *C) {
	err := Open(s.id, "test_data", "viewnav.txt")
	assert.Nil(t, err)
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	vid := views[0]
	res, err := Action(s.id, []string{"ed_activate_view", vidStr(vid)})
	assert.Nil(t, err)
	actions.Ar.EdActivateView(vid)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtRight))})
	assert.Nil(t, err)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtDown))})
	assert.Nil(t, err)
	assert.NotEq(t, vid, actions.Ar.EdCurView())
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtUp))})
	assert.Nil(t, err)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtLeft))})
	assert.Nil(t, err)
	assert.Eq(t, vid, actions.Ar.EdCurView())
	assert.Eq(t, len(res), 0)
}

func (s *ApiSuite) TestEdViews(t *C) {
	err := Open(s.id, "test_data", "edviews.txt")
	views := actions.Ar.EdViews()
	assert.True(t, len(views) >= 2)
	res, err := Action(s.id, []string{"ed_views"})
	assert.Nil(t, err)
	assert.Eq(t, len(views), len(res))
	vs := []string{}
	for _, v := range views {
		vs = append(vs, vidStr(v))
	}
	assert.DeepEq(t, vs, res)
}
