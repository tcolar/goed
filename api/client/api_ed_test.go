package client

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/tcolar/goed/actions"
	"github.com/tcolar/goed/core"
	. "gopkg.in/check.v1"
)

func (s *ApiSuite) TestEdActivateView(c *C) {
	err := Open(s.id, "test_data", "swapview.txt")
	c.Assert(err, IsNil)
	views := actions.Ar.EdViews()
	c.Assert(len(views) >= 2, Equals, true)
	res, err := Action(s.id, []string{"ed_activate_view", vidStr(views[1])})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.EdCurView(), Equals, views[1])
	res, err = Action(s.id, []string{"ed_activate_view", vidStr(views[0])})
	c.Assert(err, IsNil)
	c.Assert(actions.Ar.EdCurView(), Equals, views[0])
	c.Assert(len(res), Equals, 0)
}

func (s *ApiSuite) TestEdCurView(c *C) {
	res, err := Action(s.id, []string{"ed_cur_view"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert(fmt.Sprintf("%d", actions.Ar.EdCurView()), Equals, res[0])
}

func (s *ApiSuite) TestEdDelCol(c *C) {
	v1 := actions.Ar.EdCurView()
	_, c1 := actions.Ar.EdViewIndex(v1)
	err := Open(s.id, "test_data", "delcol.txt")
	vid := actions.Ar.EdCurView()
	l, c2, _, _ := actions.Ar.ViewBounds(vid)
	actions.Ar.EdViewMove(vid, l, c2, 2, c2+5) // force view to it's own column, if not already
	_, col := actions.Ar.EdViewIndex(vid)
	c.Assert(col, Not(Equals), c1)
	res, err := Action(s.id, []string{"ed_del_col", strconv.Itoa(col), "true"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	vid = actions.Ar.EdCurView()
	_, c3 := actions.Ar.EdViewIndex(vid)
	c.Assert(col, Not(Equals), c3)
}

func (s *ApiSuite) TestEdDelView(c *C) {
	v0 := actions.Ar.EdCurView()
	err := Open(s.id, "test_data", "delview.txt")
	v1 := actions.Ar.EdCurView()
	c.Assert(v0, Not(Equals), v1)
	actions.Ar.ViewSetDirty(v1, true)
	// view is dirty so first try should do nothing
	res, err := Action(s.id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.EdCurView(), Equals, v1)
	// but asking again, should force close v1 and send us back to v0
	res, err = Action(s.id, []string{"ed_del_view", fmt.Sprintf("%d", v1), "true"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	c.Assert(actions.Ar.EdCurView(), Equals, v0)
}

func (s *ApiSuite) TestEdOpen(c *C) {
	// open in a new view
	res, err := Action(s.id, []string{"ed_open", "theme.toml", "-1", "test_data", "true"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	vid := actions.Ar.EdCurView()
	c.Assert(fmt.Sprintf("%d", vid), Equals, res[0])
	c.Assert(actions.Ar.ViewTitle(vid), Equals, "theme.toml")

	// replace the view
	prevId := res[0]
	res, err = Action(s.id, []string{"ed_open", "testopen.txt", prevId, "test_data", "true"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	vid = actions.Ar.EdCurView()
	c.Assert(fmt.Sprintf("%d", vid), Equals, res[0])
	c.Assert(prevId, Equals, res[0])
	c.Assert(actions.Ar.ViewTitle(vid), Equals, "testopen.txt")

	loc, _ := filepath.Abs("../test_data/theme.toml") // should no longer be found
	c.Assert(actions.Ar.EdViewByLoc(loc), Equals, int64(-1))

	actions.Ar.EdDelView(actions.Ar.EdCurView(), true)
}

func (s *ApiSuite) TestEdQuitCheck(c *C) {
	views := actions.Ar.EdViews()
	actions.Ar.ViewSetDirty(views[0], true)
	res, err := Action(s.id, []string{"ed_quit_check"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert("false", Equals, res[0])

	actions.Ar.ViewSetDirty(views[0], false)
	res, err = Action(s.id, []string{"ed_quit_check"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 1)
	c.Assert("true", Equals, res[0])
}

func (s *ApiSuite) TestEdResize(c *C) {
	r, cc := actions.Ar.EdSize()
	res, err := Action(s.id, []string{"ed_resize", strconv.Itoa(r + 1), strconv.Itoa(cc + 1)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
	r2, c2 := actions.Ar.EdSize()
	c.Assert(r+1, Equals, r2)
	c.Assert(cc+1, Equals, c2)
	res, err = Action(s.id, []string{"ed_resize", strconv.Itoa(r), strconv.Itoa(cc)})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 0)
}

func (s *ApiSuite) TestEdSize(c *C) {
	r, cc := actions.Ar.EdSize()
	res, err := Action(s.id, []string{"ed_size"})
	c.Assert(err, IsNil)
	c.Assert(len(res), Equals, 2)
	c.Assert(strconv.Itoa(r), Equals, res[0])
	c.Assert(strconv.Itoa(cc), Equals, res[1])
}

func (s *ApiSuite) TestEdSwapViews(c *C) {
	boundStr := func(y1, x1, y2, x2 int) string {
		return fmt.Sprintf("%d %d %d %d", y1, x1, y2, x2)
	}
	err := Open(s.id, "test_data", "swapview.txt")
	c.Assert(err, IsNil)
	views := actions.Ar.EdViews()
	c.Assert(len(views) >= 2, Equals, true)
	vid := views[0]
	v2id := views[1]
	b1 := boundStr(actions.Ar.ViewBounds(vid))
	b2 := boundStr(actions.Ar.ViewBounds(v2id))
	res, err := Action(s.id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	c.Assert(err, IsNil)
	c.Assert(b1, Equals, boundStr(actions.Ar.ViewBounds(v2id)))
	c.Assert(b2, Equals, boundStr(actions.Ar.ViewBounds(vid)))
	res, err = Action(s.id, []string{"ed_swap_views", vidStr(vid), vidStr(v2id)})
	c.Assert(err, IsNil)
	c.Assert(b1, Equals, boundStr(actions.Ar.ViewBounds(vid)))
	c.Assert(b2, Equals, boundStr(actions.Ar.ViewBounds(v2id)))
	c.Assert(len(res), Equals, 0)
}

func (s *ApiSuite) TestEdViewAt(c *C) {
	err := Open(s.id, "test_data", "viewat.txt")
	c.Assert(err, IsNil)
	views := actions.Ar.EdViews()
	c.Assert(len(views) >= 2, Equals, true)
	for _, v := range views {
		l, cc, l2, c2 := actions.Ar.ViewBounds(v)
		res, err := Action(s.id, []string{"ed_view_at", fmt.Sprintf("%d", l), fmt.Sprintf("%d", cc)})
		c.Assert(err, IsNil)
		c.Assert(len(res), Equals, 3)
		c.Assert(res[0], Equals, vidStr(v))
		res, err = Action(s.id, []string{"ed_view_at", fmt.Sprintf("%d", l2), fmt.Sprintf("%d", c2)})
		c.Assert(err, IsNil)
		c.Assert(len(res), Equals, 3)
		c.Assert(res[0], Equals, vidStr(v))
	}
}

func (s *ApiSuite) TestEdViewNavigate(c *C) {
	err := Open(s.id, "test_data", "viewnav.txt")
	c.Assert(err, IsNil)
	views := actions.Ar.EdViews()
	c.Assert(len(views) >= 2, Equals, true)
	vid := views[0]
	res, err := Action(s.id, []string{"ed_activate_view", vidStr(vid)})
	c.Assert(err, IsNil)
	actions.Ar.EdActivateView(vid)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtRight))})
	c.Assert(err, IsNil)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtDown))})
	c.Assert(err, IsNil)
	c.Assert(vid, Not(Equals), actions.Ar.EdCurView())
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtUp))})
	c.Assert(err, IsNil)
	res, err = Action(s.id, []string{"ed_view_navigate", strconv.Itoa(int(core.CursorMvmtLeft))})
	c.Assert(err, IsNil)
	c.Assert(vid, Equals, actions.Ar.EdCurView())
	c.Assert(len(res), Equals, 0)
}

func (s *ApiSuite) TestEdViews(c *C) {
	err := Open(s.id, "test_data", "edviews.txt")
	views := actions.Ar.EdViews()
	c.Assert(len(views) >= 2, Equals, true)
	res, err := Action(s.id, []string{"ed_views"})
	c.Assert(err, IsNil)
	c.Assert(len(views), Equals, len(res))
	vs := []string{}
	for _, v := range views {
		vs = append(vs, vidStr(v))
	}
	c.Assert(vs, DeepEquals, res)
}
